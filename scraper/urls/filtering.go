package urls

import (
	"fmt"
	"regexp"
	"strings"
)

// genericPatterns represent text that usually appears in links and that is not specific enough to
// be a valid battle or participant
var genericPatterns = strings.Join([]string{
	"campaign_against",
	"history_of",
	"list_of",
	"timeline_of",
}, "|")

// confusingPatterns represent text that usually appears in links and that could potentially be
// confused to be a battle or participant
var confusingPatterns = strings.Join([]string{
	"(category|file|help|portal|talk|wikipedia):",
	"[\\w-]*(advisor|chief_of|division|force|marines|politics|rangers|regiment)[\\w-]*",
	"army", "auxiliaries",
	"caliphate", "cia", "commandery", "conscription", "crusades",
	"empire",
	"flag",
	"in_absentia", "islam", "islamism",
	"jewish", "jews",
	"killed_in_action",
	"muslim_conquests",
	"offensive_jihad",
	"participants", "pow", "prisoner_of_war",
	"roman_emperor",
	"sicherheitsdienst", "surrender_\\(military\\)",
}, "|")

var genericURLs = regexp.MustCompile(fmt.Sprintf(`(?i)\w*%s\w*`, genericPatterns))
var redLinks = regexp.MustCompile(`(?i)[\?\&]redlink=1`)
var confusingURLs = regexp.MustCompile(fmt.Sprintf("(?i)/wiki/(%s)", confusingPatterns))
var fragmentOrNestedWikis = regexp.MustCompile(`(?i)/wiki/([\w-]*#|[\w-]+/)[\w-]+`)

// NotSpecific returns true when the URL refers to a Wikipedia article that is not specific enough
// to be a battle or participant. Example: https://en.wikipedia.org/wiki/History_of_Norway
func NotSpecific(url string) bool {
	return genericURLs.Match([]byte(url))
}

// ShouldSkip returns true when the URL refers to a non-wiki URL, or to a Wikipedia article that is
// not really a battle or participant, but that usually appears in places where they could be
// mistaken for being one. Example: https://en.wikipedia.org/wiki/Prisoner_of_war
func ShouldSkip(url string) bool {
	bytesURL := []byte(url)
	return !strings.Contains(url, "/wiki/") ||
		redLinks.Match(bytesURL) ||
		confusingURLs.Match(bytesURL) ||
		fragmentOrNestedWikis.Match(bytesURL)
}
