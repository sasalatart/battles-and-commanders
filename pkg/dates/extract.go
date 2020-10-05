package dates

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// extractYear extracts the year of a date in YYYY-MM-DD format
var extractYear = createExtractor(yearMatcher)

// extractMonth extracts the month of a date in YYYY-MM-DD format
var extractMonth = createExtractor(monthMatcher)

// extractDay extracts the day of a date in YYYY-MM-DD format
var extractDay = createExtractor(dayMatcher)

func createExtractor(matcher *regexp.Regexp) func(date string) string {
	return func(date string) string {
		if !matcher.MatchString(date) {
			return ""
		}
		return matcher.ReplaceAllString(matcher.FindString(date), `$1`)
	}
}

func extract(s string) (Historic, error) {
	if !historicStringFormat.MatchString(s) {
		return Historic{}, ErrNotDate
	}

	safeConv := func(value string) (int, error) {
		if value == "" {
			return 0, nil
		}
		return strconv.Atoi(value)
	}
	isBCE := strings.Contains(s, "BC")
	year, err := safeConv(extractYear(s))
	if err != nil {
		return Historic{}, errors.New("Invalid year")
	}
	month, err := safeConv(extractMonth(s))
	if err != nil || month < 0 || month > 12 {
		return Historic{}, errors.New("Invalid month")
	}
	day, err := safeConv(extractDay(s))
	if err != nil || day < 0 || day > daysPerMonth[month] {
		return Historic{}, errors.New("Invalid day")
	}
	return Historic{IsBCE: isBCE, Year: year, Month: month, Day: day}, nil
}
