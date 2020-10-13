package strclean

import (
	"regexp"
	"strings"
)

// Apply applies a series of transformations to its input in order to transform and remove unwanted
// text that resulted from scraping activities
func Apply(text string) string {
	result := strings.Trim(text, " ,\n")
	for _, op := range pipeline {
		result = op.regex.ReplaceAllString(result, op.replaceWith)
	}
	return result
}

var pipeline = []struct {
	regex       *regexp.Regexp
	replaceWith string
}{
	{
		// Remove irregular whitespaces
		regex:       regexp.MustCompile(`[\xa0  ]`),
		replaceWith: " ",
	},
	{
		// Remove references (examples: [Note 1], [2], [better source needed], [3]:44)
		regex:       regexp.MustCompile(`\[[\w\s]+\](:\d*)?`),
		replaceWith: "",
	},
	{
		// Remove line breaks after colons
		regex:       regexp.MustCompile(`:\n+`),
		replaceWith: ": ",
	},
	{
		// Replace line breaks with periods
		regex:       regexp.MustCompile(`\n+`),
		replaceWith: ". ",
	},
	{
		// Remove spaces before periods
		regex:       regexp.MustCompile(`\s+\.`),
		replaceWith: ".",
	},
	{
		// Ensure there are spaces between word-number pairs
		regex:       regexp.MustCompile(`([a-zA-Z])(\d)`),
		replaceWith: "$1 $2",
	},
	{
		// Ensure there are spaces after commas for non-numeric text
		regex:       regexp.MustCompile(`(\D),(\w)`),
		replaceWith: "$1, $2",
	},
	{
		// Ensure there are spaces after colons
		regex:       regexp.MustCompile(`:([^\s])`),
		replaceWith: ": $1",
	},
	{
		// Ensure there are spaces before opening brackes
		regex:       regexp.MustCompile(`(\w)([\(\[])`),
		replaceWith: "$1 $2",
	},
	{
		// Ensure there are spaces after closing brackes
		regex:       regexp.MustCompile(`([\)\]])(\w)`),
		replaceWith: "$1 $2",
	},
	{
		// Remove excess spaces
		regex:       regexp.MustCompile(`\s{2,}`),
		replaceWith: " ",
	},
	{
		// Ensure numbers are not mixed between themselves (example: 1,129,619478,741 -> 1,129,619. 478,741)
		regex:       regexp.MustCompile(`(,\d{3})(\d{1,})`),
		replaceWith: "$1. $2",
	},
	{
		// Ensure there are no spaces before and after time colons
		regex:       regexp.MustCompile(`(?i)(\d)\s*:\s*(\d{1,2})\s*([a|p]\.?m\.?|hrs)`),
		replaceWith: "$1:$2$3",
	},
}
