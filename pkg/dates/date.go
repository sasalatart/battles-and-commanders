package dates

import (
	"fmt"
	"regexp"
)

var historicStringFormat = regexp.MustCompile(`^\d{1,4}(-\d{1,2}(-\d{1,2})?)?(\sBC)?$`)

// Historic represents a specific date from history. It stores the year, month, day and if whether
// or not the underlying date was BCE. Both month and day are optional, because some dates in
// history do not have that level of precision
type Historic struct {
	Year  int  `json:"year"`
	Month int  `json:"month,omitempty"`
	Day   int  `json:"day,omitempty"`
	IsBCE bool `json:"isBCE"`
}

func (h Historic) String() string {
	result := fmt.Sprint(h.Year)
	if h.Month != 0 {
		result += fmt.Sprintf("-%02d", h.Month)
	}
	if h.Day != 0 {
		result += fmt.Sprintf("-%02d", h.Day)
	}
	return result
}

// ToNum converts the Historic to a floating-point number, where the magnitude of the result
// will be directly proportional to how recent the underlying date was. For this purpose, if the
// Month or Day values of the Historic are zero, they will not be considered in this calculaiton.
// Following this logic, "1769-08-15" will be greater than "1769-08", which will in turn be greater
// than "1769". Whether or not the date is BCE is also considered, so "1769" will be greater than
// "1769 BC", in the same way that "31-09-02 BC" will be greater than "52-09 BC"
func (h Historic) ToNum() float64 {
	monthFraction := (float64(h.Month) / (monthsInYear + 1))
	dayFraction := (float64(h.Day) / ((monthsInYear + 1) * (maxDaysInMonth + 1)))
	decimalPart := monthFraction + dayFraction
	if h.IsBCE {
		return float64((-1 * h.Year)) + decimalPart
	}
	return float64(h.Year) + decimalPart
}

// ToBeginning fills the missing month and day of a partial Historic by setting them to 1 if
// they are missing. For example, the underlying date "1769" is converted into "1769-01-01",
// "1769-08" is converted into "1769-08-01", and "1769-08-15" stays the same
func (h Historic) ToBeginning() Historic {
	if h.Month == 0 {
		h.Month = 1
	}
	if h.Day == 0 {
		h.Day = 1
	}
	return h
}

// ToEnd fills the missing month and day of a partial Historic by setting them to the last month
// of the year (unless specified), and to the last day of that month (unless specified). For
// example, the underlying date "1769" is converted into "1769-12-31", "1769-08" is converted into
// "1769-08-31", and "1769-08-15" stays the same
func (h Historic) ToEnd() Historic {
	if h.Month == 0 {
		h.Month = monthsInYear
	}
	if h.Day == 0 {
		h.Day = daysPerMonth[h.Month]
	}
	return h
}

// IsValid applies the IsValid algorithm to check if the underlying date is valid or not
func (h Historic) IsValid() bool {
	return IsValid(h.String())
}

// New attempts to create a Historic date from its input string. At the moment this string has to
// be in YYYY-MM-DD [BC] format. For example, the following string dates are all valid inputs for
// this function: "1769-08-15", "1769-8-15", "1769-8", "1769" and "1457-04-16 BC"
func New(s string) (Historic, error) {
	return extract(s)
}

const monthsInYear = 12
const maxDaysInMonth = 31

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
