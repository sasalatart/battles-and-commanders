package dates

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type formatMatcher struct {
	regex     *regexp.Regexp
	formatter func(src string) []string
}

var bcMatcher = regexp.MustCompile(`[\s,]*(B(\.)?C(\.)?E?(\.)?)[\s,]*`)
var monthNameMatcher = withMonths(`(?i).*(%s).*`)
var yearMatcher = regexp.MustCompile(`^(\d{1,4})-?`)
var monthMatcher = regexp.MustCompile(`^\d{1,4}-(\d{1,2})`)
var dayMatcher = regexp.MustCompile(`^\d{1,4}-\d{1,2}-(\d{1,2})`)

var formatMatchers = []formatMatcher{
	buildFormatMatcher("YMD|YMD", `$1-$2-$3|$4-$5-$6`), // 1935 Jan 19 — 1935 March 22
	buildFormatMatcher("DMY|DMY", `$3-$2-$1|$6-$5-$4`), // 18 May 1803 – 20 November 1815
	buildFormatMatcher("MDY|MDY", `$3-$1-$2|$6-$4-$5`), // May 13, 1867 to May 24, 1867
	buildFormatMatcher("DMY|MDY", `$3-$2-$1|$6-$4-$5`), // 4 December 2009 – December 12, 2009
	buildFormatMatcher("MDY|MY", `$3-$1-$2|$5-$4`),     // November 18, 1918 – March, 1919
	buildFormatMatcher("DMY|MY", `$3-$2-$1|$5-$4`),     // 13 October 1945 – April 1946
	buildFormatMatcher("MY|DMY", `$2-$1|$5-$4-$3`),     // May 1913 - 25 September 1920
	buildFormatMatcher("DM|MDY", `$5-$2-$1|$5-$3-$4`),  // 8 August – November 11 1918
	buildFormatMatcher("DM|DMY", `$5-$2-$1|$5-$4-$3`),  // 21 August – 2 September 1644
	buildFormatMatcher("MD|DY", `$4-$1-$2|$4-$1-$3`),   // December 11-15, 1862
	buildFormatMatcher("MD|MDY", `$5-$1-$2|$5-$3-$4`),  // April 30 – May 6, 1863
	buildFormatMatcher("MY|MDY", `$2-$1|$5-$3-$4`),     // August 1769 – 5 May 1821
	buildFormatMatcher("MY|DMY", `$2-$1|$5-$4-$3`),     // May 1913 - 25 September 1920
	buildFormatMatcher("MY|MY", `$2-$1|$4-$3`),         // May 640 – December 640
	buildFormatMatcher("Y|MDY", `$1|$4-$2-$3`),         // 455 BC – May 8, 453 BC
	buildFormatMatcher("M|MDY", `$4-$1|$4-$2-$3`),      // January - April 7, 1337
	buildFormatMatcher("M|DMY", `$4-$1|$4-$3-$2`),      // June – 29 November 1855
	buildFormatMatcher("Y|DMY", `$1|$4-$3-$2`),         // 1267–14 March 1273
	buildFormatMatcher("D|DMY", `$4-$3-$1|$4-$3-$2`),   // 1-3 August 1798
	buildFormatMatcher("M|MY", `$3-$1|$3-$2`),          // January-March 309
	buildFormatMatcher("Y|MY", `$1|$3-$2`),             // 996 – May 998
	buildFormatMatcher("Y|Y", `$1|$2`),                 // 1769 – 1821
	buildFormatMatcher("MDY", `$3-$1-$2`),              // April 19, 1775
	buildFormatMatcher("DMY", `$3-$2-$1`),              // 9 August 48
	buildFormatMatcher("MY", `$2-$1`),                  // August 490
	buildFormatMatcher("Y", `$1`),                      // 1769
}

func buildFormatMatcher(format string, expandedFormat string) formatMatcher {
	p := strings.ReplaceAll(format, "|", `[\s,]*(?:—|-|–|−|to|and|&)[\s,]*`)
	p = strings.ReplaceAll(p, "D", `[\s,]+(\d{1,2})(?:th)?(?:\.)?[\s,]+`)
	p = strings.ReplaceAll(p, "M", fmt.Sprintf(`[\s,]+(%s)(?:\.)?[\s,]+`, strings.Join(append(months, shortMonths...), "|")))
	p = strings.ReplaceAll(p, "Y", `[\s,]+(\d{1,4})(?:\.)?[\s,]+`)
	p = strings.ReplaceAll(p, `[\s,]+[\s,]+`, `[\s,]+`)
	p = strings.ReplaceAll(p, `[\s,]+[\s,]*`, `[\s,]*`)
	p = strings.ReplaceAll(p, `[\s,]*[\s,]+`, `[\s,]*`)
	p = strings.Trim(p, `[\s,]+`)
	r := regexp.MustCompile(`(?i)` + p + `\.?`)

	formatter := func(src string) []string {
		dates := strings.Split(r.ReplaceAllString(src, expandedFormat), "|")
		var res []string
		for _, date := range dates {
			parsed := strings.Trim(date, " ")
			if monthNameMatcher.MatchString(parsed) {
				monthName := monthNameMatcher.ReplaceAllString(parsed, `$1`)
				parsed = strings.ReplaceAll(parsed, monthName, mToi[strings.ToLower(monthName)])
			}
			parsed = padZeroes(parsed)
			res = append(res, parsed)
		}
		return res
	}

	return formatMatcher{
		regex:     r,
		formatter: formatter,
	}
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

func withMonths(format string) *regexp.Regexp {
	joined := strings.Join(append(months, shortMonths...), "|")
	amount := len(regexp.MustCompile(`\(%s\)`).FindAllStringIndex(format, -1))

	var joins []interface{}
	for i := 0; i < amount; i++ {
		joins = append(joins, joined)
	}
	return regexp.MustCompile(fmt.Sprintf(format, joins...))
}

var months = []string{
	"january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december",
}

var shortMonths = []string{
	"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sept", "oct", "nov", "dec",
}

var mToi = make(map[string]string)

func init() {
	for n, month := range months {
		num := strconv.Itoa(n + 1)
		if len(num) != 2 {
			num = "0" + num
		}
		mToi[month] = num
		mToi[shortMonths[n]] = num
	}
}
