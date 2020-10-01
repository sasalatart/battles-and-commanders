package dates_test

import (
	"testing"

	"github.com/sasalatart/batcoms/pkg/dates"
	"github.com/stretchr/testify/assert"
)

func TestDateToNum(t *testing.T) {
	t.Run("InputsValidation", func(t *testing.T) {
		t.Parallel()
		_, err := dates.ToNum("not-a-date")
		assert.Error(t, err, "Returns an error for non-date inputs")
		_, err = dates.ToNum("1769-08-15")
		assert.NoError(t, err, "Does not return error for valid dates")
	})

	t.Run("RelativeValues", func(t *testing.T) {
		t.Parallel()
		assertNumerical := func(date1, date2 string) (float32, float32) {
			value1, err := dates.ToNum(date1)
			assert.NoErrorf(t, err, "Converting %q to number", date1)
			value2, err := dates.ToNum(date2)
			assert.NoErrorf(t, err, "Converting %q to number", date2)
			return value1, value2
		}
		assertGreaterThan := func(date1, date2 string) {
			value1, value2 := assertNumerical(date1, date2)
			assert.Greaterf(t, value1, value2, "Comparing date %q with %q", date1, date2)
		}
		assertEqual := func(date1, date2 string) {
			value1, value2 := assertNumerical(date1, date2)
			assert.Equalf(t, value1, value2, "Comparing date %q with %q", date1, date2)
		}
		assertGreaterThan("1821-05-05", "1769-08-15")
		assertGreaterThan("1821-05-05", "1821-05-04")
		assertGreaterThan("1821-05-05", "1821-04-05")
		assertGreaterThan("1821-05-05", "1820-05-05")
		assertGreaterThan("1821-05-05", "1821-05")
		assertGreaterThan("1821-05-05", "1821")
		assertGreaterThan("1821-05", "1821")
		assertGreaterThan("1821-05-05", "1821-05-05 BC")
		assertGreaterThan("1", "1 BC")
		assertGreaterThan("31-09-02 BC", "31-09-01 BC")
		assertGreaterThan("31-09-02 BC", "32-09-02 BC")
		assertGreaterThan("31-09-02 BC", "31-09 BC")
		assertGreaterThan("31-09-02 BC", "31 BC")
		assertGreaterThan("31-09 BC", "31 BC")
		assertGreaterThan("1769", "1768-12-31")
		assertGreaterThan("1768-01-01 BC", "1769 BC")
		assertGreaterThan("1768-01-01 BC", "1769-12-31 BC")
		assertEqual("1821-05-05", "1821-5-5")
		assertEqual("31-09-02 BC", "31-09-2 BC")
		assertEqual("31-09-02 BC", "31-9-02 BC")
	})
}
