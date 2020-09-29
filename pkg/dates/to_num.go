package dates

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// ToNum converts a date in YYYY-MM-DD [BC] format to a floating-point number, where dates
// chronologically after other dates will always be greater than them. For this purpose, if the MM
// and/or DD portion of the date are missing, they are assumed to be "zero". Hence, "1769-08-15"
// will be greater than "1769-08", which will in turn be greater than "1769"
func ToNum(date string) (float32, error) {
	if !IsValid(date) {
		return 0, errors.Wrapf(ErrNotDate, "Validating date %q", date)
	}

	fmtErr := func(unit string) error {
		return fmt.Errorf("Unable to obtain numerical %s from %q", unit, date)
	}
	safeConv := func(value string) (int, error) {
		if value == "" {
			return 0, nil
		}
		return strconv.Atoi(value)
	}
	isBC := strings.Contains(date, "BC")
	year, err := safeConv(extractYear(date))
	if err != nil {
		return 0, fmtErr("year")
	}
	month, err := safeConv(extractMonth(date))
	if err != nil {
		return 0, fmtErr("month")
	}
	day, err := safeConv(extractDay(date))
	if err != nil {
		return 0, fmtErr("day")
	}

	decimalPart := (float32(month) / 13) + (float32(day) / 32)
	if isBC {
		return float32((-1 * year)) + decimalPart, nil
	}
	return float32(year) + decimalPart, nil
}
