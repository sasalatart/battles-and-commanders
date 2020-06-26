package urls

import (
	"fmt"
	"regexp"
	"strings"
)

// forbiddenKeywords represents text that usually appears in links and that could potentially be
// confused to be a battle or participant.
var forbiddenKeywords = strings.Join([]string{
	"category:", "help:", "portal:", "talk:", "wikipedia:",
	"army", "auxiliaries", "auxiliary_division",
	"caliphate", "chief_of_police", "cia", "commandery", "conscription", "crusades",
	"delta_force",
	"empire",
	"flag",
	"islam", "islamism",
	"jewish", "jews",
	"killed_in_action",
	"left-wing_politics",
	"military_advisor", "muslim_conquests",
	"offensive_jihad",
	"pow", "prisoner_of_war",
	"right-wing_politics", "roman_emperor",
	"sicherheitsdienst", "surrender_\\(military\\)",
	"united_states_army_rangers",
}, "|")

var genericURLs = regexp.MustCompile(`(?i)(history_of|list_of|campaign_against|timeline_of)`)
var forbiddenQS = regexp.MustCompile(`(?i)[\?\&](redlink=1)`)
var forbiddenURLs = regexp.MustCompile(fmt.Sprintf("(?i)/wiki/(%s)", forbiddenKeywords))

// NotSpecific returns true when the URL refers to a Wikipedia article that is not specific enough
// to be a battle or participant. Example: https://en.wikipedia.org/wiki/History_of_Norway
func NotSpecific(url string) bool {
	return genericURLs.Match([]byte(url))
}

// ShouldSkip returns true when the URL refers to a Wikipedia article that is not really a battle or
// participant, but that usually appears in places where they could be mistaken for being one.
// Example: https://en.wikipedia.org/wiki/Prisoner_of_war
func ShouldSkip(url string) bool {
	bytesURL := []byte(url)
	return forbiddenQS.Match(bytesURL) || forbiddenURLs.Match(bytesURL)
}

// IsExternal returns true when the URL does not point to a Wiki page.
func IsExternal(url string) bool {
	return !strings.Contains(url, "/wiki/")
}
