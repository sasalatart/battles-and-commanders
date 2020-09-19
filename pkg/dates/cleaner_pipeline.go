package dates

import "regexp"

var cleanerPipeline = []struct {
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
