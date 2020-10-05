package dates_test

import (
	"testing"

	"github.com/sasalatart/batcoms/pkg/dates"
	"github.com/stretchr/testify/assert"
)

func TestDate(t *testing.T) {
	t.Run("ToNum", func(t *testing.T) {
		t.Parallel()
		cases := []struct {
			upper dates.Historic
			lower dates.Historic
		}{
			{
				upper: dates.Historic{Year: 1821, Month: 5, Day: 5},
				lower: dates.Historic{Year: 1769, Month: 8, Day: 15},
			},
			{
				upper: dates.Historic{Year: 1821, Month: 5, Day: 5},
				lower: dates.Historic{Year: 1821, Month: 5, Day: 4},
			},
			{
				upper: dates.Historic{Year: 1821, Month: 5, Day: 5},
				lower: dates.Historic{Year: 1821, Month: 4, Day: 5},
			},
			{
				upper: dates.Historic{Year: 1821, Month: 5, Day: 5},
				lower: dates.Historic{Year: 1820, Month: 5, Day: 5},
			},
			{
				upper: dates.Historic{Year: 1821, Month: 5, Day: 5},
				lower: dates.Historic{Year: 1821, Month: 5},
			},
			{
				upper: dates.Historic{Year: 1821, Month: 5, Day: 5},
				lower: dates.Historic{Year: 1821},
			},
			{
				upper: dates.Historic{Year: 1821, Month: 5},
				lower: dates.Historic{Year: 1821},
			},
			{
				upper: dates.Historic{Year: 1821, Month: 5, Day: 5},
				lower: dates.Historic{Year: 1821, Month: 5, Day: 5, IsBCE: true},
			},
			{
				upper: dates.Historic{Year: 1},
				lower: dates.Historic{Year: 1, IsBCE: true},
			},
			{
				upper: dates.Historic{Year: 31, Month: 9, Day: 2, IsBCE: true},
				lower: dates.Historic{Year: 31, Month: 9, Day: 1, IsBCE: true},
			},
			{
				upper: dates.Historic{Year: 31, Month: 9, Day: 2, IsBCE: true},
				lower: dates.Historic{Year: 32, Month: 9, Day: 2, IsBCE: true},
			},
			{
				upper: dates.Historic{Year: 31, Month: 9, Day: 2, IsBCE: true},
				lower: dates.Historic{Year: 31, Month: 9, IsBCE: true},
			},
			{
				upper: dates.Historic{Year: 31, Month: 9, Day: 2, IsBCE: true},
				lower: dates.Historic{Year: 31, IsBCE: true},
			},
			{
				upper: dates.Historic{Year: 31, Month: 9, IsBCE: true},
				lower: dates.Historic{Year: 31, IsBCE: true},
			},
			{
				upper: dates.Historic{Year: 1769},
				lower: dates.Historic{Year: 1768, Month: 12, Day: 31},
			},
			{
				upper: dates.Historic{Year: 1768, Month: 1, Day: 1, IsBCE: true},
				lower: dates.Historic{Year: 1769, IsBCE: true},
			},
			{
				upper: dates.Historic{Year: 1768, Month: 1, IsBCE: true},
				lower: dates.Historic{Year: 1769, IsBCE: true},
			},
			{
				upper: dates.Historic{Year: 1768, IsBCE: true},
				lower: dates.Historic{Year: 1769, IsBCE: true},
			},
			{
				upper: dates.Historic{Year: 1768, Month: 1, Day: 1, IsBCE: true},
				lower: dates.Historic{Year: 1769, Month: 12, Day: 31, IsBCE: true},
			},
		}
		for _, c := range cases {
			assert.Greaterf(
				t,
				c.upper.ToNum(),
				c.lower.ToNum(),
				"Comparing date %s with %s", c.upper, c.upper,
			)
		}
	})

	t.Run("ToBeginning", func(t *testing.T) {
		t.Parallel()
		cases := []struct {
			date     dates.Historic
			expected dates.Historic
		}{
			{
				date:     dates.Historic{Year: 1769, Month: 8, Day: 15},
				expected: dates.Historic{Year: 1769, Month: 8, Day: 15},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 8, Day: 15, IsBCE: true},
				expected: dates.Historic{Year: 1769, Month: 8, Day: 15, IsBCE: true},
			},
			{
				date:     dates.Historic{Year: 1769},
				expected: dates.Historic{Year: 1769, Month: 1, Day: 1},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 1},
				expected: dates.Historic{Year: 1769, Month: 1, Day: 1},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 2},
				expected: dates.Historic{Year: 1769, Month: 2, Day: 1},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 3},
				expected: dates.Historic{Year: 1769, Month: 3, Day: 1},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 4},
				expected: dates.Historic{Year: 1769, Month: 4, Day: 1},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 5},
				expected: dates.Historic{Year: 1769, Month: 5, Day: 1},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 6},
				expected: dates.Historic{Year: 1769, Month: 6, Day: 1},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 7},
				expected: dates.Historic{Year: 1769, Month: 7, Day: 1},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 8},
				expected: dates.Historic{Year: 1769, Month: 8, Day: 1},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 9},
				expected: dates.Historic{Year: 1769, Month: 9, Day: 1},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 10},
				expected: dates.Historic{Year: 1769, Month: 10, Day: 1},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 11},
				expected: dates.Historic{Year: 1769, Month: 11, Day: 1},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 12},
				expected: dates.Historic{Year: 1769, Month: 12, Day: 1},
			},
		}
		for _, c := range cases {
			assert.Equal(t, c.expected, c.date.ToBeginning())
		}
	})

	t.Run("ToEnd", func(t *testing.T) {
		t.Parallel()
		cases := []struct {
			date     dates.Historic
			expected dates.Historic
		}{
			{
				date:     dates.Historic{Year: 1769, Month: 8, Day: 15},
				expected: dates.Historic{Year: 1769, Month: 8, Day: 15},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 8, Day: 15, IsBCE: true},
				expected: dates.Historic{Year: 1769, Month: 8, Day: 15, IsBCE: true},
			},
			{
				date:     dates.Historic{Year: 1769},
				expected: dates.Historic{Year: 1769, Month: 12, Day: 31},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 1},
				expected: dates.Historic{Year: 1769, Month: 1, Day: 31},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 2},
				expected: dates.Historic{Year: 1769, Month: 2, Day: 29},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 3},
				expected: dates.Historic{Year: 1769, Month: 3, Day: 31},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 4},
				expected: dates.Historic{Year: 1769, Month: 4, Day: 30},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 5},
				expected: dates.Historic{Year: 1769, Month: 5, Day: 31},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 6},
				expected: dates.Historic{Year: 1769, Month: 6, Day: 30},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 7},
				expected: dates.Historic{Year: 1769, Month: 7, Day: 31},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 8},
				expected: dates.Historic{Year: 1769, Month: 8, Day: 31},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 9},
				expected: dates.Historic{Year: 1769, Month: 9, Day: 30},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 10},
				expected: dates.Historic{Year: 1769, Month: 10, Day: 31},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 11},
				expected: dates.Historic{Year: 1769, Month: 11, Day: 30},
			},
			{
				date:     dates.Historic{Year: 1769, Month: 12},
				expected: dates.Historic{Year: 1769, Month: 12, Day: 31},
			},
		}
		for _, c := range cases {
			assert.Equal(t, c.expected, c.date.ToEnd())
		}
	})
}
