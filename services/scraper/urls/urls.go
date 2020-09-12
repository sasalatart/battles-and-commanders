package urls

import (
	"fmt"
	"regexp"
	"strings"
)

// notSpecificMatcher represents text that usually appears in links and that is not specific enough
// to be a valid battle or participant
var notSpecificMatcher = regexp.MustCompile(`(?i)` + strings.Join([]string{
	"campaign_against",
	"history_of",
	"list_of",
	"timeline_of",
}, "|"))

// falsePositivesMatcher represents text that usually appears in links and that could potentially be
// confused to be a battle or participant
var falsePositivesMatcher = regexp.MustCompile(fmt.Sprintf(`(?i)/wiki/(%s)`, strings.Join([]string{
	`(category|file|help|portal|talk|wikipedia):`,
	`(flag|in_absentia|killed_in_action|pow|prisoner_of_war|surrender_\(military\)|wia|wounded_in_action)$`,
	`[\w-,]*(advisor|chief_of|division|force|marines|participants|politics|rangers|regiment)[\w-,]*`,
	`(army|auxiliaries|caliphate|cia|commandery|conscription|crusades|empire|islam|islamism|jewish|jews|muslim_conquests|offensive_jihad|roman_emperor|sicherheitsdienst)$`,
}, "|")))

// redLinksMatcher matches links that do not yet exist (Wikipedia adds a "redlink" query param)
var redLinksMatcher = regexp.MustCompile(`(?i)[\?\&]redlink=1`)

// fragmentsMatcher matches fragments in URLs (foo.bar/baz#fragment) and nested resources
var fragmentsMatcher = regexp.MustCompile(`(?i)/wiki/([\w-,]*#|[\w-,]+/)[\w-,]+`)

// NotSpecific returns true when the URL refers to a Wikipedia article that is not specific enough
// to be a battle or participant. Example: https://en.wikipedia.org/wiki/History_of_Norway
func NotSpecific(url string) bool {
	return notSpecificMatcher.MatchString(url)
}

// ShouldSkip returns true when the URL refers to a non-wiki URL, or to a Wikipedia article that is
// not really a battle or participant, but that usually appears in places where they could be
// mistaken for being one. Example: https://en.wikipedia.org/wiki/Prisoner_of_war
func ShouldSkip(url string) bool {
	return !strings.Contains(url, "/wiki/") ||
		redLinksMatcher.MatchString(url) ||
		fragmentsMatcher.MatchString(url) ||
		falsePositivesMatcher.MatchString(url)
}
