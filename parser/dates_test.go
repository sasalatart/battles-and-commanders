package parser_test

import (
	"reflect"
	"testing"

	"github.com/sasalatart/batcoms/parser"
)

func TestDate(t *testing.T) {
	cases := []struct {
		raw      string
		expected []string
	}{
		{
			"~2500 BC",
			[]string{"2500 BC"},
		},
		{
			"211 BC",
			[]string{"211 BC"},
		},
		{
			"311 B.C.",
			[]string{"311 BC"},
		},
		{
			"272 BCE",
			[]string{"272 BC"},
		},
		{
			"Late 686",
			[]string{"686"},
		},
		{
			"344 CE",
			[]string{"344"},
		},
		{
			"404 BC or 403 BC",
			[]string{"404 BC"},
		},
		{
			"268 or 269 CE",
			[]string{"268"},
		},
		{
			"Between 1480, or 1483",
			[]string{"1480"},
		},
		{
			"April or May 1521",
			[]string{"1521-04"},
		},
		{
			"June, July or August 251 CE",
			[]string{"251-07"},
		},
		{
			"circa September 9 CE",
			[]string{"9-09"},
		},
		{
			"circa. 380 BC",
			[]string{"380 BC"},
		},
		{
			"circ. 360 BC",
			[]string{"360 BC"},
		},
		{
			"c. 870 AD",
			[]string{"870"},
		},
		{
			"April 19, 1775",
			[]string{"1775-04-19"},
		},
		{
			"July 6, 1950; 69 years ago (1950-07-06)",
			[]string{"1950-07-06"},
		},
		{
			"Sunday, March 8, 1722",
			[]string{"1722-03-08"},
		},
		{
			"9 August 48 BC",
			[]string{"48-08-09 BC"},
		},
		{
			"25 October 1415 (Saint Crispin's Day)",
			[]string{"1415-10-25"},
		},
		{
			"March of 1287",
			[]string{"1287-03"},
		},
		{
			"August/September (Metageitnion), 490 BC",
			[]string{"490-08 BC"},
		},
		{
			"Winter solstice, December 218 BC",
			[]string{"218-12 BC"},
		},
		{
			"First half of 1263",
			[]string{"1263"},
		},
		{
			"Mid to Late November 2001",
			[]string{"2001-11"},
		},
		{
			"February 11, 1700 (O.S.) February 12, 1700 (Swedish calendar) February 22, 1700 (N.S.)",
			[]string{"1700-02-22"},
		},
		{
			"July 8, 1702 (O.S.). July 9, 1702 (Swedish calendar). July 19, 1702 (1702-07-19) (N.S.)",
			[]string{"1702-07-19"},
		},
		{
			"1769 – 1821",
			[]string{"1769", "1821"},
		},
		{
			"1821 – 1769 BC",
			[]string{"1821 BC", "1769 BC"},
		},
		{
			"Spring 218 – 201 BC (17 years)",
			[]string{"218 BC", "201 BC"},
		},
		{
			"550 – spring of 551",
			[]string{"550", "551"},
		},
		{
			"996 – May 998",
			[]string{"996", "998-05"},
		},
		{
			"May 640 – December 640",
			[]string{"640-05", "640-12"},
		},
		{
			"April 19, 1775 – September 3, 1783 (8 years, 4 months and 15 days)",
			[]string{"1775-04-19", "1783-09-03"},
		},
		{
			"January-March 309 B.C.",
			[]string{"309-01 BC", "309-03 BC"},
		},
		{
			"April 30 (1863-04-30) – May 6, 1863 (1863-05-06)",
			[]string{"1863-04-30", "1863-05-06"},
		},
		{
			"July 1–August 1, 30 BC",
			[]string{"30-07-01 BC", "30-08-01 BC"},
		},
		{
			"July 25 to August 1, 1626",
			[]string{"1626-07-25", "1626-08-01"},
		},
		{
			"September 19 and October 7, 1777",
			[]string{"1777-09-19", "1777-10-07"},
		},
		{
			"December 11–15, 1862",
			[]string{"1862-12-11", "1862-12-15"},
		},
		{
			"July 15 & 16, 1839",
			[]string{"1839-07-15", "1839-07-16"},
		},
		{
			"May 1913 - 25 September 1920",
			[]string{"1913-05", "1920-09-25"},
		},
		{
			"May 13, 1867 to May 24, 1867",
			[]string{"1867-05-13", "1867-05-24"},
		},
		{
			"February 1482 – January 2, 1492",
			[]string{"1482-02", "1492-01-02"},
		},
		{
			"January - April 7, 1337",
			[]string{"1337-01", "1337-04-07"},
		},
		{
			"June – 29 November 1855",
			[]string{"1855-06", "1855-11-29"},
		},
		{
			"November 18, 1918 – March, 1919",
			[]string{"1918-11-18", "1919-03"},
		},
		{
			"18 May 1803 – 20 November 1815 (12 years, 5 months and 4 weeks)",
			[]string{"1803-05-18", "1815-11-20"},
		},
		{
			"8 February 1904. – 5 September 1905",
			[]string{"1904-02-08", "1905-09-05"},
		},
		{
			"21 August – 2 September 1644",
			[]string{"1644-08-21", "1644-09-02"},
		},
		{
			"8 August – November 11 1918",
			[]string{"1918-08-08", "1918-11-11"},
		},
		{
			"1-3 August 1798",
			[]string{"1798-08-01", "1798-08-03"},
		},
		{
			"9–11 February 1586.",
			[]string{"1586-02-09", "1586-02-11"},
		},
		{
			"Night of 3–4 August 1327",
			[]string{"1327-08-03", "1327-08-04"},
		},
		{
			"24 and 25 September 1812",
			[]string{"1812-09-24", "1812-09-25"},
		},
		{
			"13 October 1945 – April 1946",
			[]string{"1945-10-13", "1946-04"},
		},
		{
			"4 December 2009 – December 12, 2009",
			[]string{"2009-12-04", "2009-12-12"},
		},
		{
			"455 BC – May 8, 453 BC",
			[]string{"455 BC", "453-05-08 BC"},
		},
		{
			"1267–14 March 1273",
			[]string{"1267", "1273-03-14"},
		},
		{
			"1935 Jan 19 — 1935 March 22",
			[]string{"1935-01-19", "1935-03-22"},
		},
		{
			"August 1927 – 22 December 1936 (9 years, 4 months and 3 weeks). 10 August 1945 – 7 August 1950 (4 years, 4 months and 1 week)",
			[]string{"1927-08", "1950-08-07"},
		},
		{
			"18-19 July 1913 (N.S.)(5-6 July in O.S.)",
			[]string{"1913-07-18", "1913-07-19"},
		},
		{
			"May 23, 1592 – December 16, 1598 (Gregorian Calendar); April 13, 1592 – November 19, 1598 (Lunar calendar)",
			[]string{"1592-05-23", "1598-12-16"},
		},
		{
			"29–30 November 1612 (Julian calendar); 9-10 December 1612 (Gregorian calendar)",
			[]string{"1612-12-09", "1612-12-10"},
		},
	}

	for _, c := range cases {
		got, err := parser.Date(c.raw)

		if err != nil {
			t.Errorf("Expected %q to be parsed as %q, but instead errored with %s", c.raw, c.expected, err)
			continue
		}

		if !reflect.DeepEqual(got, c.expected) {
			t.Errorf("Expected %q to be parsed as %q, but instead got %q", c.raw, c.expected, got)
		}
	}
}
