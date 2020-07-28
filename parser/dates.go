package parser

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type dateMatcher struct {
	regex *regexp.Regexp
	build func(src string) []string
}

var bcMatcher = regexp.MustCompile(`[\s,]*(B(\.)?C(\.)?E?(\.)?)[\s,]*`)
var monthNameMatcher = withMonths(`(?i).*(%s).*`)
var yearMatcher = regexp.MustCompile(`^(\d{1,4})-?`)
var monthMatcher = regexp.MustCompile(`^\d{1,4}-(\d{1,2})`)
var dayMatcher = regexp.MustCompile(`^\d{1,4}-\d{1,2}-(\d{1,2})`)

var months = []string{
	"january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december",
}

var shortMonths = []string{
	"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sept", "oct", "nov", "dec",
}

var cleanDatePipeline = []struct {
	regex       *regexp.Regexp
	replaceWith string
}{
	{
		// Prefer Gregorian Calendar and New System (N.S.) over other kinds
		// Example: 29–30 November 1612 (Julian calendar); 9-10 December 1612 (Gregorian calendar) -> 9-10 December 1612
		regex:       regexp.MustCompile(`(?i).*?[\)\.;]([^\(\.;]*)(\s\(\d{1,4}-\d{1,2}-\d{1,2}\)\s?)?\((gregorian|N\.?S\.?).*`),
		replaceWith: `$1`,
	},
	{
		// Remove parenthes
		regex:       regexp.MustCompile(`\([^(]*\)`),
		replaceWith: "",
	},
	{
		// Remove brackets
		regex:       regexp.MustCompile(`\[[^[]*\]`),
		replaceWith: "",
	},
	{
		// Decide on numerical uncertainties
		// Example: 404 or 403 -> 404
		regex:       regexp.MustCompile(`(\d{1,4})[\s,]*(?:BC)?[\s,]*(/|or)[\s,]*(\d{1,4})[\s,]*(?:BC)?`),
		replaceWith: `$1`,
	},
	{
		// Decide on monthly uncertainties
		// Example: April or May 1521 -> April 1521
		regex:       withMonths(`(?i)(%s)[\s,]*(/|or)[\s,]*(%s)`),
		replaceWith: `$1`,
	},
	{
		// Remove days of week
		regex:       regexp.MustCompile(`(?i)(Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Sunday)\s*`),
		replaceWith: "",
	},
	{
		// Remove text after semi-colons
		regex:       regexp.MustCompile(`;.*$`),
		replaceWith: "",
	},
	{
		// Remove era
		regex:       regexp.MustCompile(`[\s,]*(B(\.)?C(\.)?E?(\.)?|C(\.)?E(\.)?|A(\.)?D(\.)?)[\s,]*`),
		replaceWith: "",
	},
	{
		// Remove circa variants
		regex:       regexp.MustCompile(`(?i)[\s,]*(c\.|circ(a)?(\.)?)[\s,]*`),
		replaceWith: "",
	},
	{
		// Remove other kinds of approximations
		regex:       regexp.MustCompile(`(?i)[\s,]*(~|Between|Spring|Summer|Fall|Autumn|Winter|Early to (Mid|Late)|Mid to Late|Early (days)?|Mid|Late|Solstice|Morning|Noon|Night|(First|Second) half|\?)(\sof)?[\s,]*`),
		replaceWith: "",
	},
	{
		// Remove useless connectors
		regex:       regexp.MustCompile(`(?i)[\s]+of[\s+]`),
		replaceWith: " ",
	},
}

var mToIMap = make(map[string]string)

func init() {
	for n, month := range months {
		num := strconv.Itoa(n + 1)
		if len(num) != 2 {
			num = "0" + num
		}
		mToIMap[month] = num
		mToIMap[shortMonths[n]] = num
	}
}

func mToI(m string) string {
	return mToIMap[strings.ToLower(m)]
}

func withMonths(format string) *regexp.Regexp {
	joined := strings.Join(append(months, shortMonths...), "|")
	amount := len(regexp.MustCompile(`\(%s\)`).FindAllStringIndex(format, -1))

	var joins []interface{}
	for i := 0; i < amount; i++ {
		joins = append(joins, joined)
	}
	return regexp.MustCompile(fmt.Sprintf(format, joins...))
}

func padZeroes(d string) string {
	var z0 = regexp.MustCompile(`(.*)-(\d)-(\d)`)
	var z1 = regexp.MustCompile(`(.*)-(\d)-(.*)`)
	var z2 = regexp.MustCompile(`(.*)-(\d)$`)

	parsed := d
	if z0.MatchString(parsed) {
		parsed = z0.ReplaceAllString(parsed, `$1-0$2-0$3`)
	} else if z1.MatchString(parsed) {
		parsed = z1.ReplaceAllString(parsed, `$1-0$2-$3`)
	} else if z2.MatchString(parsed) {
		parsed = z2.ReplaceAllString(parsed, `$1-0$2`)
	}
	return parsed
}

func newDateMatcher(format string, expandedFormat string) dateMatcher {
	p := strings.ReplaceAll(format, "|", `[\s,]*(?:—|-|–|−|to|and|&)[\s,]*`)
	p = strings.ReplaceAll(p, "D", `[\s,]+(\d{1,2})(?:th)?(?:\.)?[\s,]+`)
	p = strings.ReplaceAll(p, "M", fmt.Sprintf(`[\s,]+(%s)(?:\.)?[\s,]+`, strings.Join(append(months, shortMonths...), "|")))
	p = strings.ReplaceAll(p, "Y", `[\s,]+(\d{1,4})(?:\.)?[\s,]+`)
	p = strings.ReplaceAll(p, `[\s,]+[\s,]+`, `[\s,]+`)
	p = strings.ReplaceAll(p, `[\s,]+[\s,]*`, `[\s,]*`)
	p = strings.ReplaceAll(p, `[\s,]*[\s,]+`, `[\s,]*`)
	p = strings.Trim(p, `[\s,]+`)
	r := regexp.MustCompile(`(?i)` + p + `\.?`)

	build := func(src string) []string {
		dates := strings.Split(r.ReplaceAllString(src, expandedFormat), "|")
		var res []string
		for _, date := range dates {
			parsed := strings.Trim(date, " ")
			if monthNameMatcher.MatchString(parsed) {
				monthName := monthNameMatcher.ReplaceAllString(parsed, `$1`)
				parsed = strings.ReplaceAll(parsed, monthName, mToI(monthName))
			}
			parsed = padZeroes(parsed)
			res = append(res, parsed)
		}
		return res
	}

	return dateMatcher{
		regex: r,
		build: build,
	}
}

var datesMatchers = []dateMatcher{
	newDateMatcher("YMD|YMD", `$1-$2-$3|$4-$5-$6`), // 1935 Jan 19 — 1935 March 22
	newDateMatcher("DMY|DMY", `$3-$2-$1|$6-$5-$4`), // 18 May 1803 – 20 November 1815
	newDateMatcher("MDY|MDY", `$3-$1-$2|$6-$4-$5`), // May 13, 1867 to May 24, 1867
	newDateMatcher("DMY|MDY", `$3-$2-$1|$6-$4-$5`), // 4 December 2009 – December 12, 2009
	newDateMatcher("MDY|MY", `$3-$1-$2|$5-$4`),     // November 18, 1918 – March, 1919
	newDateMatcher("DMY|MY", `$3-$2-$1|$5-$4`),     // 13 October 1945 – April 1946
	newDateMatcher("MY|DMY", `$2-$1|$5-$4-$3`),     // May 1913 - 25 September 1920
	newDateMatcher("DM|MDY", `$5-$2-$1|$5-$3-$4`),  // 8 August – November 11 1918
	newDateMatcher("DM|DMY", `$5-$2-$1|$5-$4-$3`),  // 21 August – 2 September 1644
	newDateMatcher("MD|DY", `$4-$1-$2|$4-$1-$3`),   // December 11-15, 1862
	newDateMatcher("MD|MDY", `$5-$1-$2|$5-$3-$4`),  // April 30 – May 6, 1863
	newDateMatcher("MY|MDY", `$2-$1|$5-$3-$4`),     // August 1769 – 5 May 1821
	newDateMatcher("MY|DMY", `$2-$1|$5-$4-$3`),     // May 1913 - 25 September 1920
	newDateMatcher("MY|MY", `$2-$1|$4-$3`),         // May 640 – December 640
	newDateMatcher("Y|MDY", `$1|$4-$2-$3`),         // 455 BC – May 8, 453 BC
	newDateMatcher("M|MDY", `$4-$1|$4-$2-$3`),      // January - April 7, 1337
	newDateMatcher("M|DMY", `$4-$1|$4-$3-$2`),      // June – 29 November 1855
	newDateMatcher("Y|DMY", `$1|$4-$3-$2`),         // 1267–14 March 1273
	newDateMatcher("D|DMY", `$4-$3-$1|$4-$3-$2`),   // 1-3 August 1798
	newDateMatcher("M|MY", `$3-$1|$3-$2`),          // January-March 309
	newDateMatcher("Y|MY", `$1|$3-$2`),             // 996 – May 998
	newDateMatcher("Y|Y", `$1|$2`),                 // 1769 – 1821
	newDateMatcher("MDY", `$3-$1-$2`),              // April 19, 1775
	newDateMatcher("DMY", `$3-$2-$1`),              // 9 August 48
	newDateMatcher("MY", `$2-$1`),                  // August 490
	newDateMatcher("Y", `$1`),                      // 1769
}

// minMax finds the minimum and maximum dates within an array of formatted, non-BC dates (YYYY-MM-DD)
func minMax(ss []string) (string, string) {
	maxYearLength := 0
	for _, s := range ss {
		yearLen := len(year(s))
		if yearLen > maxYearLength {
			maxYearLength = yearLen
		}
	}
	for i, s := range ss {
		lenDiff := maxYearLength - len(year(s))
		if lenDiff > 0 {
			ss[i] = strings.Repeat("0", lenDiff) + ss[i]
		}
	}
	sort.Strings(ss)
	return strings.TrimLeft(ss[0], "0"), strings.TrimLeft(ss[len(ss)-1], "0")
}

func extractFactory(matcher *regexp.Regexp) func(date string) string {
	return func(date string) string {
		if !matcher.MatchString(date) {
			return ""
		}
		return matcher.ReplaceAllString(matcher.FindString(date), `$1`)
	}
}

// year extracts the year of a date in YYYY-MM-DD format
var year = extractFactory(yearMatcher)

// month extracts the month of a date in YYYY-MM-DD format
var month = extractFactory(monthMatcher)

// day extracts the day of a date in YYYY-MM-DD format
var day = extractFactory(dayMatcher)

// discardMatches runs a series of heuristics to check if a set of potential dates should be
// discarded or not. Sometimes a regex may give a false positive.
func discardMatches(dates []string) bool {
	isValidDatePart := func(s string, max int) bool {
		n := strings.TrimLeft(s, "0")
		nInt, err := strconv.Atoi(n)
		if err != nil || nInt > max || nInt < 1 {
			return false
		}
		return true
	}

	for _, date := range dates {
		if m := month(date); m != "" && !isValidDatePart(m, 12) {
			fmt.Println("INVALID M", m)
			return false
		}
		if d := day(date); d != "" && !isValidDatePart(d, 31) {
			fmt.Println("INVALID D", d)
			return false
		}
	}

	min, max := minMax(dates)
	minYear := year(min)
	maxYear := year(max)
	if len(dates) > 1 && len(maxYear)-len(minYear) >= 2 {
		return true
	}
	return false
}

// sentenceDates find date pairs (start-finish) in YYYY-MM-DD format given a sentence
func sentenceDates(s string) []string {
	var res []string
	for _, dateMatcher := range datesMatchers {
		if !dateMatcher.regex.MatchString(s) {
			continue
		}

		discardMatcher := false
		matches := dateMatcher.regex.FindAllString(s, -1)
		for _, match := range matches {
			subMatches := dateMatcher.build(match)
			discardMatcher = discardMatcher || discardMatches(subMatches)
			if !discardMatcher {
				res = append(res, subMatches...)
			}
		}

		if !discardMatcher {
			break
		}
		res = []string{}
	}

	if len(res) == 0 {
		return []string{}
	}

	min, max := minMax(res)
	if min == max {
		return []string{min}
	}
	return []string{min, max}
}

// Date parses text containing a potential date or range of dates, and translates them into
// YYYY-MM-DD format, returning a slice of each found date, sorted chronologically
func Date(text string) ([]string, error) {
	t := text

	isBC := bcMatcher.MatchString(t)
	for _, p := range cleanDatePipeline {
		t = p.regex.ReplaceAllString(t, p.replaceWith)
	}

	var dates []string
	for _, s := range strings.Split(t, ". ") {
		dates = append(dates, sentenceDates(s)...)
	}

	if len(dates) == 0 {
		return dates, errors.Errorf("Unable to parse date %q", t)
	}

	min, max := minMax(dates)
	minYear := year(min)
	maxYear := year(max)
	if isBC {
		min = fmt.Sprintf("%s BC", min)
		max = fmt.Sprintf("%s BC", max)
	}
	if min == max {
		return []string{min}, nil
	}
	if isBC && minYear != maxYear {
		min, max = max, min
	}
	return []string{min, max}, nil
}
