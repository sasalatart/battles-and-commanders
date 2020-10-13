package dates

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Parse receives text containing a potential date or range of dates, and translates them into a
// slice of Historic, where the first element is the earliest date detected in the text, and the
// second element is the latest date detected in the text. If there is just one date detected, then
// the slice will just have one element
func Parse(t string) ([]Historic, error) {
	isBCE := bcMatcher.MatchString(t)
	for _, p := range cleanerPipeline {
		t = p.regex.ReplaceAllString(t, p.replaceWith)
	}

	var dates []Historic
	for _, s := range strings.Split(t, ". ") {
		dates = append(dates, fromPhrase(s)...)
	}
	if len(dates) == 0 {
		return []Historic{}, errors.Errorf("No valid dates found in %q", t)
	}

	min, max := minMax(dates...)
	min.IsBCE = isBCE
	max.IsBCE = isBCE
	if min == max {
		return []Historic{min}, nil
	}
	if isBCE && min.Year != max.Year {
		min, max = max, min
	}
	return []Historic{min, max}, nil
}

// fromPhrase finds historic date pairs (start, finish) given a phrase
func fromPhrase(s string) []Historic {
	var res []Historic
	for _, formatMatcher := range formatMatchers {
		if !formatMatcher.regex.MatchString(s) {
			continue
		}
		discardMatcher := false
		matches := formatMatcher.regex.FindAllString(s, -1)
		for _, match := range matches {
			subMatches := formatMatcher.formatter(match)
			discardMatcher = shouldDiscard(subMatches)
			if discardMatcher {
				break
			}
			for _, sm := range subMatches {
				h, err := New(sm)
				if err != nil {
					continue
				}
				res = append(res, h)
			}
		}
		if !discardMatcher {
			break
		}
		res = []Historic{}
	}
	if len(res) == 0 {
		return []Historic{}
	}

	min, max := minMax(res...)
	if min == max {
		return []Historic{min}
	}
	return []Historic{min, max}
}

// shouldDiscard runs a series of heuristics to check if a set of potential string dates should be
// discarded or not. Sometimes a regex may give a false positive
func shouldDiscard(dates []string) bool {
	for _, date := range dates {
		if !IsValid(date) {
			return true
		}
	}

	var minYear, maxYear int
	for _, d := range dates {
		year := extractYear(d)
		yearNum, err := strconv.Atoi(year)
		if err != nil {
			return true
		}
		if yearNum < minYear || minYear == 0 {
			minYear = yearNum
		}
		if yearNum > maxYear || maxYear == 0 {
			maxYear = yearNum
		}
	}
	if minYear != maxYear && maxYear-minYear >= 100 {
		return true
	}
	return false
}

// minMax finds the minimum and maximum dates from a set of HistoricDates
func minMax(dates ...Historic) (Historic, Historic) {
	var min, max Historic
	var minNum, maxNum float64
	for _, h := range dates {
		num := h.ToNum()
		if num < minNum || minNum == 0 {
			min = h
			minNum = num
		}
		if num > maxNum || maxNum == 0 {
			max = h
			maxNum = num
		}
	}
	return min, max
}

var cleanerPipeline = []struct {
	regex       *regexp.Regexp
	replaceWith string
}{
	{
		// Prefer Gregorian Calendar, New System (N.S.) and "probable" over other kinds
		// Example: 29–30 November 1612 (Julian calendar); 9-10 December 1612 (Gregorian calendar) -> 9-10 December 1612
		regex:       regexp.MustCompile(`(?i).*?[\)\.;]([^\(\.;]*)(\s\(\d{1,4}-\d{1,2}-\d{1,2}\)\s?)?\((gregorian|N\.?S\.?|probable).*`),
		replaceWith: `$1`,
	},
	{
		// Remove parenthes
		regex:       regexp.MustCompile(`\([^(]*\)`),
		replaceWith: "",
	},
	{
		// Remove brackets
		regex:       regexp.MustCompile(`\[[^[]*\]`),
		replaceWith: "",
	},
	{
		// Decide on numerical uncertainties
		// Example: 404 or 403 -> 404
		regex:       regexp.MustCompile(`(\d{1,4})[\s,]*(?:BC)?[\s,]*(/|or)[\s,]*(\d{1,4})[\s,]*(?:BC)?`),
		replaceWith: `$1`,
	},
	{
		// Decide on monthly uncertainties
		// Example: April or May 1521 -> April 1521
		regex:       withMonths(`(?i)(%s)[\s,]*(/|or)[\s,]*(%s)`),
		replaceWith: `$1`,
	},
	{
		// Remove days of week
		regex:       regexp.MustCompile(`(?i)(Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Sunday)\s*`),
		replaceWith: "",
	},
	{
		// Remove text after semi-colons
		regex:       regexp.MustCompile(`;.*$`),
		replaceWith: "",
	},
	{
		// Remove era
		regex:       regexp.MustCompile(`[\s,]*(B(\.)?C(\.)?E?(\.)?|C(\.)?E(\.)?|A(\.)?D(\.)?)[\s,]*`),
		replaceWith: "",
	},
	{
		// Remove circa variants
		regex:       regexp.MustCompile(`(?i)[\s,]*(c\.|circ(a)?(\.)?)[\s,]*`),
		replaceWith: "",
	},
	{
		// Remove other kinds of approximations
		regex:       regexp.MustCompile(`(?i)[\s,]*(~|Between|Spring|Summer|Fall|Autumn|Winter|Early to (Mid|Late)|Mid to Late|Early (days)?|Mid|Late|Solstice|Morning|Noon|Night|(First|Second) half|\?)(\sof)?`),
		replaceWith: "",
	},
	{
		// Remove useless connectors
		regex:       regexp.MustCompile(`(?i)[\s]+of[\s+]`),
		replaceWith: " ",
	},
}
