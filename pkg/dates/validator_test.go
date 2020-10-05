package dates_test

import (
	"testing"

	"github.com/sasalatart/batcoms/pkg/dates"
	"github.com/stretchr/testify/assert"
)

func TestDatesValidator(t *testing.T) {
	cases := []struct {
		value    string
		expected bool
	}{
		{value: "1769-08-15", expected: true},
		{value: "1769-8-15", expected: true},
		{value: "1769-15-08", expected: false},
		{value: "1769-15-8", expected: false},
		{value: "15-08-1769", expected: false},
		{value: "08-15-1769", expected: false},
		{value: "08-1769-15", expected: false},
		{value: "8-1769-15", expected: false},
		{value: "15-1769-08", expected: false},
		{value: "15-1769-8", expected: false},
		{value: "1769-08", expected: true},
		{value: "1769-8", expected: true},
		{value: "1769-15", expected: false},
		{value: "08-1769", expected: false},
		{value: "8-1769", expected: false},
		{value: "08-15", expected: false},
		{value: "8-15", expected: false},
		{value: "15-08", expected: true},
		{value: "15-8", expected: true},
		{value: "1769", expected: true},
		{value: "08", expected: true},
		{value: "8", expected: true},
		{value: "15", expected: true},
		{value: "1769-08-15 something else", expected: false},
		{value: "some text 1769-08-15", expected: false},
		{value: "1769-08-15 BC", expected: true},
		{value: "1769-08 BC", expected: true},
		{value: "1769 BC", expected: true},
		{value: "just-text", expected: false},
		{value: "1769-xx-yy", expected: false},
		{value: "1769-xx", expected: false},
		{value: "1769-", expected: false},
	}
	for _, c := range cases {
		assert.Equalf(t, c.expected, dates.IsValid(c.value), "Validating %s", c.value)
	}
}
