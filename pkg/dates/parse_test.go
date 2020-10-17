package dates_test

import (
	"testing"

	"github.com/sasalatart/batcoms/pkg/dates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatesParse(t *testing.T) {
	cases := []struct {
		raw      string
		expected []dates.Historic
	}{
		{
			"~2500 BC",
			[]dates.Historic{{Year: 2500, IsBCE: true}},
		},
		{
			"211 BC",
			[]dates.Historic{{Year: 211, IsBCE: true}},
		},
		{
			"311 B.C.",
			[]dates.Historic{{Year: 311, IsBCE: true}},
		},
		{
			"272 BCE",
			[]dates.Historic{{Year: 272, IsBCE: true}},
		},
		{
			"344 CE",
			[]dates.Historic{{Year: 344}},
		},
		{
			"First half of 1263",
			[]dates.Historic{{Year: 1263}},
		},
		{
			"Late 686",
			[]dates.Historic{{Year: 686}},
		},
		{
			"404 BC or 403 BC",
			[]dates.Historic{{Year: 404, IsBCE: true}},
		},
		{
			"268 or 269 CE",
			[]dates.Historic{{Year: 268}},
		},
		{
			"Between 1480, or 1483",
			[]dates.Historic{{Year: 1480}},
		},
		{
			"April or May 1521",
			[]dates.Historic{{Year: 1521, Month: 4}},
		},
		{
			"June, July or August 251 CE",
			[]dates.Historic{{Year: 251, Month: 7}},
		},
		{
			"circa September 9 CE",
			[]dates.Historic{{Year: 9, Month: 9}},
		},
		{
			"circa. 380 BC",
			[]dates.Historic{{Year: 380, IsBCE: true}},
		},
		{
			"circ. 360 BC",
			[]dates.Historic{{Year: 360, IsBCE: true}},
		},
		{
			"c. 870 AD",
			[]dates.Historic{{Year: 870}},
		},
		{
			"April 19, 1775",
			[]dates.Historic{{Year: 1775, Month: 4, Day: 19}},
		},
		{
			"July 6, 1950; 69 years ago (1950-07-06)",
			[]dates.Historic{{Year: 1950, Month: 7, Day: 6}},
		},
		{
			"Sunday, March 8, 1722",
			[]dates.Historic{{Year: 1722, Month: 3, Day: 8}},
		},
		{
			"9 August 48 BC",
			[]dates.Historic{{Year: 48, Month: 8, Day: 9, IsBCE: true}},
		},
		{
			"25 October 1415 (Saint Crispin's Day)",
			[]dates.Historic{{Year: 1415, Month: 10, Day: 25}},
		},
		{
			"March of 1287",
			[]dates.Historic{{Year: 1287, Month: 3}},
		},
		{
			"August/September (Metageitnion), 490 BC",
			[]dates.Historic{{Year: 490, Month: 8, IsBCE: true}},
		},
		{
			"Winter solstice, December 218 BC",
			[]dates.Historic{{Year: 218, Month: 12, IsBCE: true}},
		},
		{
			"Mid to Late November 2001",
			[]dates.Historic{{Year: 2001, Month: 11}},
		},
		{
			"February 11, 1700 (O.S.) February 12, 1700 (Swedish calendar) February 22, 1700 (N.S.)",
			[]dates.Historic{{Year: 1700, Month: 2, Day: 22}},
		},
		{
			"July 8, 1702 (O.S.). July 9, 1702 (Swedish calendar). July 19, 1702 (1702-07-19) (N.S.)",
			[]dates.Historic{{Year: 1702, Month: 7, Day: 19}},
		},
		{
			"1769 – 1821",
			[]dates.Historic{
				{Year: 1769},
				{Year: 1821},
			},
		},
		{
			"1821 – 1769 BC",
			[]dates.Historic{
				{Year: 1821, IsBCE: true},
				{Year: 1769, IsBCE: true},
			},
		},
		{
			"Spring 218 – 201 BC (17 years)",
			[]dates.Historic{
				{Year: 218, IsBCE: true},
				{Year: 201, IsBCE: true},
			},
		},
		{
			"550 – spring of 551",
			[]dates.Historic{
				{Year: 550},
				{Year: 551},
			},
		},
		{
			"996 – May 998",
			[]dates.Historic{
				{Year: 996},
				{Year: 998, Month: 5},
			},
		},
		{
			"May 640 – December 640",
			[]dates.Historic{
				{Year: 640, Month: 5},
				{Year: 640, Month: 12},
			},
		},
		{
			"April 19, 1775 – September 3, 1783 (8 years, 4 months and 15 days)",
			[]dates.Historic{
				{Year: 1775, Month: 4, Day: 19},
				{Year: 1783, Month: 9, Day: 3},
			},
		},
		{
			"May 13, 1867 to May 24, 1867",
			[]dates.Historic{
				{Year: 1867, Month: 5, Day: 13},
				{Year: 1867, Month: 5, Day: 24},
			},
		},
		{
			"January-March 309 B.C.",
			[]dates.Historic{
				{Year: 309, Month: 1, IsBCE: true},
				{Year: 309, Month: 3, IsBCE: true},
			},
		},
		{
			"April 30 (1863-04-30) – May 6, 1863 (1863-05-06)",
			[]dates.Historic{
				{Year: 1863, Month: 4, Day: 30},
				{Year: 1863, Month: 5, Day: 6},
			},
		},
		{
			"July 1–August 1, 30 BC",
			[]dates.Historic{
				{Year: 30, Month: 7, Day: 1, IsBCE: true},
				{Year: 30, Month: 8, Day: 1, IsBCE: true},
			},
		},
		{
			"July 25 to August 1, 1626",
			[]dates.Historic{
				{Year: 1626, Month: 7, Day: 25},
				{Year: 1626, Month: 8, Day: 1},
			},
		},
		{
			"September 19 and October 7, 1777",
			[]dates.Historic{
				{Year: 1777, Month: 9, Day: 19},
				{Year: 1777, Month: 10, Day: 7},
			},
		},
		{
			"December 11–15, 1862",
			[]dates.Historic{
				{Year: 1862, Month: 12, Day: 11},
				{Year: 1862, Month: 12, Day: 15},
			},
		},
		{
			"July 15 & 16, 1839",
			[]dates.Historic{
				{Year: 1839, Month: 7, Day: 15},
				{Year: 1839, Month: 7, Day: 16},
			},
		},
		{
			"May 1913 - 25 September 1920",
			[]dates.Historic{
				{Year: 1913, Month: 5},
				{Year: 1920, Month: 9, Day: 25},
			},
		},
		{
			"February 1482 – January 2, 1492",
			[]dates.Historic{
				{Year: 1482, Month: 2},
				{Year: 1492, Month: 1, Day: 2},
			},
		},
		{
			"January - April 7, 1337",
			[]dates.Historic{
				{Year: 1337, Month: 1},
				{Year: 1337, Month: 4, Day: 7},
			},
		},
		{
			"June – 29 November 1855",
			[]dates.Historic{
				{Year: 1855, Month: 6},
				{Year: 1855, Month: 11, Day: 29},
			},
		},
		{
			"November 18, 1918 – March, 1919",
			[]dates.Historic{
				{Year: 1918, Month: 11, Day: 18},
				{Year: 1919, Month: 3},
			},
		},
		{
			"18 May 1803 – 20 November 1815 (12 years, 5 months and 4 weeks)",
			[]dates.Historic{
				{Year: 1803, Month: 5, Day: 18},
				{Year: 1815, Month: 11, Day: 20},
			},
		},
		{
			"8 February 1904. – 5 September 1905",
			[]dates.Historic{
				{Year: 1904, Month: 2, Day: 8},
				{Year: 1905, Month: 9, Day: 5},
			},
		},
		{
			"21 August – 2 September 1644",
			[]dates.Historic{
				{Year: 1644, Month: 8, Day: 21},
				{Year: 1644, Month: 9, Day: 2},
			},
		},
		{
			"8 August – November 11 1918",
			[]dates.Historic{
				{Year: 1918, Month: 8, Day: 8},
				{Year: 1918, Month: 11, Day: 11},
			},
		},
		{
			"1-3 August 1798",
			[]dates.Historic{
				{Year: 1798, Month: 8, Day: 1},
				{Year: 1798, Month: 8, Day: 3},
			},
		},
		{
			"9–11 February 1586.",
			[]dates.Historic{
				{Year: 1586, Month: 2, Day: 9},
				{Year: 1586, Month: 2, Day: 11},
			},
		},
		{
			"Night of 3–4 August 1327",
			[]dates.Historic{
				{Year: 1327, Month: 8, Day: 3},
				{Year: 1327, Month: 8, Day: 4},
			},
		},
		{
			"24 and 25 September 1812",
			[]dates.Historic{
				{Year: 1812, Month: 9, Day: 24},
				{Year: 1812, Month: 9, Day: 25},
			},
		},
		{
			"13 October 1945 – April 1946",
			[]dates.Historic{
				{Year: 1945, Month: 10, Day: 13},
				{Year: 1946, Month: 4},
			},
		},
		{
			"4 December 2009 – December 12, 2009",
			[]dates.Historic{
				{Year: 2009, Month: 12, Day: 4},
				{Year: 2009, Month: 12, Day: 12},
			},
		},
		{
			"455 BC – May 8, 453 BC",
			[]dates.Historic{
				{Year: 455, IsBCE: true},
				{Year: 453, Month: 5, Day: 8, IsBCE: true},
			},
		},
		{
			"1267–14 March 1273",
			[]dates.Historic{
				{Year: 1267},
				{Year: 1273, Month: 3, Day: 14},
			},
		},
		{
			"1935 Jan 19 — 1935 March 22",
			[]dates.Historic{
				{Year: 1935, Month: 1, Day: 19},
				{Year: 1935, Month: 3, Day: 22},
			},
		},
		{
			"August 1927 – 22 December 1936 (9 years, 4 months and 3 weeks). 10 August 1945 – 7 August 1950 (4 years, 4 months and 1 week)",
			[]dates.Historic{
				{Year: 1927, Month: 8},
				{Year: 1950, Month: 8, Day: 7},
			},
		},
		{
			"18-19 July 1913 (N.S.)(5-6 July in O.S.)",
			[]dates.Historic{
				{Year: 1913, Month: 7, Day: 18},
				{Year: 1913, Month: 7, Day: 19},
			},
		},
		{
			"May 23, 1592 – December 16, 1598 (Gregorian Calendar); April 13, 1592 – November 19, 1598 (Lunar calendar)",
			[]dates.Historic{
				{Year: 1592, Month: 5, Day: 23},
				{Year: 1598, Month: 12, Day: 16},
			},
		},
		{
			"29–30 November 1612 (Julian calendar); 9-10 December 1612 (Gregorian calendar)",
			[]dates.Historic{
				{Year: 1612, Month: 12, Day: 9},
				{Year: 1612, Month: 12, Day: 10},
			},
		},
		{
			"18 July 390 BC (traditional), 387 BC (probable)",
			[]dates.Historic{{Year: 387, IsBCE: true}},
		},
		{
			"September 25 – September 28?, 539 BC",
			[]dates.Historic{
				{Year: 539, Month: 9, Day: 25, IsBCE: true},
				{Year: 539, Month: 9, Day: 28, IsBCE: true},
			},
		},
	}
	for _, c := range cases {
		got, err := dates.Parse(c.raw)
		require.NoErrorf(t, err, "Parsing dates in text %q", c.raw)
		assert.Equal(t, c.expected, got, "Error parsing date %q", c.raw)
	}
}
