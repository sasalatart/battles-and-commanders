package dates

import "regexp"

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
