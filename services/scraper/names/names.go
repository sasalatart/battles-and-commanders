package names

import (
	"fmt"
	"regexp"
	"strings"
)

var keywords = strings.Join([]string{
	"action", "ambush", "assault", "attack",
	"battle", "blockade", "bloodbath", "bombardment", "bombing", "burning",
	"campaign", "capture", "clash", "clashes", "combat", "conflict", "confrontation", "conquest", "crisis", "crossing", "crusade",
	"defense",
	"engagement", "expedition",
	"fall", "fight",
	"incident", "insurgency", "intervention", "invasion",
	"landing", "liberation",
	"massacre", "march", "mutiny",
	"occupation", "offensive", "operatie", "operation",
	"prison break",
	"raid", "rebellion", "recovery", "relief", "revolt", "rising",
	"sack", "siege", "sinking", "skirmish", "stand", "standoff", "strike",
	"takeover",
	"uprising",
	"war",
}, "|")

func buildRegex(format string) *regexp.Regexp {
	return regexp.MustCompile("(?i)" + fmt.Sprintf(format, keywords))
}

var inside = buildRegex(`(%s)s?\s(at|for|in|of|on|to)`)
var suffix = buildRegex(`\s(%s)s?$`)
var prefix = buildRegex(`^(%s)\s`)

// IsBattle returns true if the given name probably corresponds to a battle, and false if not
func IsBattle(name string) bool {
	return inside.MatchString(name) || suffix.MatchString(name) || prefix.MatchString(name)
}
