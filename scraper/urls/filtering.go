package urls

import "regexp"

var genericURLs = regexp.MustCompile(`(?i)(history_of|list_of|campaign_against)`)
var forbiddenQS = regexp.MustCompile(`(?i)[\?\&](redlink=1)`)
var forbiddenURLs = regexp.MustCompile(`(?i)/wiki/(talk:|wikipedia:|pow|prisoner_of_war|killed_in_action|surrender_\(military\)|army|auxiliaries|auxiliary_division|caliphate|chief_of_police|cia|commandery|conscription|crusades|delta_force|empire|flag|islam|islamism|jewish|jews|left-wing_politics|right-wing_politics|military_advisor|muslim_conquests|offensive_jihad|roman_emperor|sicherheitsdienst|united_states_army_rangers)`)

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
