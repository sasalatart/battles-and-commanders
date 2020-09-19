package dates

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Parse parses text containing a potential date or range of dates, and translates them into
// YYYY-MM-DD format, returning a slice of each found date, sorted chronologically
func Parse(text string) ([]string, error) {
	t := text

	isBC := bcMatcher.MatchString(t)
	for _, p := range cleanerPipeline {
		t = p.regex.ReplaceAllString(t, p.replaceWith)
	}

	var dates []string
	for _, s := range strings.Split(t, ". ") {
		dates = append(dates, fromPhrase(s)...)
	}

	if len(dates) == 0 {
		return dates, errors.Errorf("Unable to parse date %q", t)
	}

	min, max := minMax(dates)
	minYear := extractYear(min)
	maxYear := extractYear(max)
	if isBC {
		min = fmt.Sprintf("%s BC", min)
		max = fmt.Sprintf("%s BC", max)
	}
	if min == max {
		return []string{min}, nil
	}
	if isBC && minYear != maxYear {
		min, max = max, min
	}
	return []string{min, max}, nil
}

// fromPhrase finds date pairs (start, finish) in YYYY-MM-DD format given a phrase
func fromPhrase(s string) []string {
	var res []string
	for _, formatMatcher := range formatMatchers {
		if !formatMatcher.regex.MatchString(s) {
			continue
		}

		discardMatcher := false
		matches := formatMatcher.regex.FindAllString(s, -1)
		for _, match := range matches {
			subMatches := formatMatcher.formatter(match)
			discardMatcher = discardMatcher || shouldDiscard(subMatches)
			if !discardMatcher {
				res = append(res, subMatches...)
			}
		}

		if !discardMatcher {
			break
		}
		res = []string{}
	}

	if len(res) == 0 {
		return []string{}
	}

	min, max := minMax(res)
	if min == max {
		return []string{min}
	}
	return []string{min, max}
}

// shouldDiscard runs a series of heuristics to check if a set of potential dates should be
// discarded or not. Sometimes a regex may give a false positive
func shouldDiscard(dates []string) bool {
	isValidDatePart := func(s string, max int) bool {
		n := strings.TrimLeft(s, "0")
		nInt, err := strconv.Atoi(n)
		if err != nil || nInt > max || nInt < 1 {
			return false
		}
		return true
	}

	for _, date := range dates {
		if m := extractMonth(date); m != "" && !isValidDatePart(m, 12) {
			return false
		}
		if d := extractDay(date); d != "" && !isValidDatePart(d, 31) {
			return false
		}
	}

	min, max := minMax(dates)
	minYear := extractYear(min)
	maxYear := extractYear(max)
	if len(dates) > 1 && len(maxYear)-len(minYear) >= 2 {
		return true
	}
	return false
}

// minMax finds the minimum and maximum dates within an array of formatted, non-BC dates (YYYY-MM-DD)
func minMax(ss []string) (string, string) {
	maxYearLength := 0
	for _, s := range ss {
		yearLen := len(extractYear(s))
		if yearLen > maxYearLength {
			maxYearLength = yearLen
		}
	}
	for i, s := range ss {
		lenDiff := maxYearLength - len(extractYear(s))
		if lenDiff > 0 {
			ss[i] = strings.Repeat("0", lenDiff) + ss[i]
		}
	}
	sort.Strings(ss)
	return strings.TrimLeft(ss[0], "0"), strings.TrimLeft(ss[len(ss)-1], "0")
}
