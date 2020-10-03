package dates

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var daysPerMonth = map[int]int{
	1:  31,
	2:  29,
	3:  31,
	4:  30,
	5:  31,
	6:  30,
	7:  31,
	8:  31,
	9:  30,
	10: 31,
	11: 30,
	12: 31,
}

// Parse receives text containing a potential date or range of dates, and translates them into
// YYYY-MM-DD format, returning a chronologically sorted slice of each one of those dates
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
	for _, date := range dates {
		if !IsValid(date) {
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

// ToBeginning fills the missing month and day of a partial date by setting them to "1" if they are
// missing. For example, "1769" is converted into "1769-01-01", "1769-08" is converted into
// "1769-08-01", and "1769-08-15" stays the same
func ToBeginning(date string) (string, error) {
	if !IsValid(date) {
		return date, errors.Wrapf(ErrNotDate, "Validating date %q", date)
	}

	isBC := strings.Contains(date, "BC")
	year := extractYear(date)
	month := extractMonth(date)
	day := extractDay(date)
	if month == "" {
		month = "01"
	}
	if day == "" {
		day = "01"
	}
	result := fmt.Sprintf("%s-%s-%s", year, month, day)
	if !isBC {
		return result, nil
	}
	return result + " BC", nil
}

// ToEnd fills the missing month and day of a partial date by setting them to the last month of the
// year (unless specified), and to the last day of that month (unless specified). For example,
// "1769" is converted into "1769-12-31", "1769-08" is converted into "1769-08-31", and "1769-08-15"
// stays the same
func ToEnd(date string) (string, error) {
	if !IsValid(date) {
		return date, errors.Wrapf(ErrNotDate, "Validating date %q", date)
	}

	isBC := strings.Contains(date, "BC")
	year := extractYear(date)
	month := extractMonth(date)
	day := extractDay(date)
	if month == "" {
		month = "12"
	}
	if day == "" {
		monthNumber, err := strconv.Atoi(month)
		if err != nil {
			return "", errors.Wrapf(err, "Converting month %q to int", month)
		}
		day = fmt.Sprint(daysPerMonth[monthNumber])
	}
	result := fmt.Sprintf("%s-%s-%s", year, month, day)
	if !isBC {
		return result, nil
	}
	return result + " BC", nil
}
