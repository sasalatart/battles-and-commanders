package parser

import (
	"regexp"
	"strings"
)

var pipeline = []struct {
	regex       *regexp.Regexp
	replaceWith string
}{
	{
		// Remove references (example: [Note 1], [2], [3]:44)
		regex:       regexp.MustCompile(`\[\w*\s*\w\](:\d*)?`),
		replaceWith: "",
	},
	{
		// Remove line breaks after colons
		regex:       regexp.MustCompile(`:[\n]+`),
		replaceWith: ": ",
	},
	{
		// Replace line breaks with periods
		regex:       regexp.MustCompile(`[\n]+`),
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
		regex:       regexp.MustCompile(`(:)([^\s])`),
		replaceWith: "$1 $2",
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
}

// Clean applies a series of transformations to its input in order to transform and remove unwanted
// text that resulted from scraping.
func Clean(text string) string {
	result := strings.Trim(text, "\n")
	for _, op := range pipeline {
		result = op.regex.ReplaceAllString(result, op.replaceWith)
	}
	return result
}