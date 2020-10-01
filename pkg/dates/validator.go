package dates

import (
	"strconv"
	"strings"
)

const monthsInYear = 12
const maxDaysInMonth = 31

// IsValid checks whether or not the input date is in "YYYY-MM-DD [BC]" format or not. For example,
// the following dates are all valid: 1769-08-15, 1769-8-15, 1769-8, 1769, 1769-08-15 BC, etc.
func IsValid(date string) bool {
	isValidDatePart := func(s string, max int) bool {
		n := strings.TrimLeft(s, "0")
		nInt, err := strconv.Atoi(n)
		if err != nil || nInt > max || nInt < 1 {
			return false
		}
		return true
	}

	y := extractYear(date)
	m := extractMonth(date)
	d := extractDay(date)

	if y == "" {
		return false
	}
	if m != "" && !isValidDatePart(m, monthsInYear) {
		return false
	}
	if d != "" && !isValidDatePart(d, maxDaysInMonth) {
		return false
	}

	reconstructed := y
	if m != "" {
		reconstructed = reconstructed + "-" + m
	}
	if d != "" {
		reconstructed = reconstructed + "-" + d
	}
	return reconstructed == date || reconstructed+" BC" == date
}
