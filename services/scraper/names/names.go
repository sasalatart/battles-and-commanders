package names

import (
	"fmt"
	"regexp"
	"strings"
)

var keywords = strings.Join([]string{
	"action", "affair", "ambush", "assault", "attack",
	"battle", "blockade", "bloodbath", "bombardment", "bombing", "burning",
	"campaign", "capture", "clash", "clashes", "combat", "conflict", "confrontation", "conquest", "crisis", "crossing", "crusade",
	"defeat", "defense",
	"engagement", "expedition",
	"fall", "fight",
	"incident", "insurgency", "intervention", "invasion",
	"landing", "liberation",
	"massacre", "march", "mutiny",
	"occupation", "offensive", "operatie", "operation",
	"prison break",
	"raid", "rebellion", "recovery", "relief", "revolt", "revolution", "rising",
	"sack", "siege", "singeing", "sinking", "skirmish", "stand", "standoff", "strike",
	"takeover",
	"uprising",
}, "|")

var matcher = regexp.MustCompile(fmt.Sprintf(`(?i)(%s)s?`, keywords))

// IsBattle returns true if the given name probably corresponds to a battle, and false if not
func IsBattle(name string) bool {
	return matcher.MatchString(name)
}
