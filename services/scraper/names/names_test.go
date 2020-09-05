package names_test

import (
	"testing"

	"github.com/sasalatart/batcoms/services/scraper/names"
	"github.com/stretchr/testify/assert"
)

func TestNames(t *testing.T) {
	cases := []struct {
		name     string
		expected bool
	}{
		{"Battle of Austerlitz", true},
		{"Hadong Ambush", true},
		{"Assault on Copenhagen", true},
		{"Attack at Fromelles", true},
		{"Union blockade", true},
		{"Bombardment of Algiers", true},
		{"Bombing of Singapore", true},
		{"Sinai and Palestine campaign", true},
		{"Capture of Fort Erie", true},
		{"2013 Sidon clash", true},
		{"2015 Kumanovo clashes", true},
		{"Combat of the Thirty", true},
		{"Nagorno-Karabakh conflict", true},
		{"Lungi Lol confrontation", true},
		{"Mongol conquest of Western Xia", true},
		{"Suez Crisis. Tripartite aggression. Sinai War", true},
		{"Crossing of the Düna", true},
		{"Defense Of Kozelsk", true},
		{"Engagement near Carthage", true},
		{"Wolseley expedition", true},
		{"Fall of Tenochtitlan", true},
		{"Grass Fight", true},
		{"Jingkang Incident", true},
		{"Lahij insurgency", true},
		{"North Russia intervention", true},
		{"Japanese invasion of Thailand", true},
		{"Åndalsnes landings", true},
		{"Liberation of Kuwait", true},
		{"Sand Creek massacre", true},
		{"Long March", true},
		{"Occupation of German Samoa", true},
		{"Uman-Botoshany Offensive", true},
		{"Operation Astonia", true},
		{"Ndop prison break", true},
		{"Raid on Chester", true},
		{"Kett's Rebellion", true},
		{"Recovery of Ré Island", true},
		{"Relief of Goes", true},
		{"Ionian Revolt", true},
		{"The Jacobite Rising of 1715", true},
		{"Sack of Rome (410)", true},
		{"Siege of Acre (1291)", true},
		{"Sinking of Prince of Wales and Repulse", true},
		{"Skirmish at Chalk Bluff", true},
		{"Great Stand on the Ugra River", true},
		{"Scarborough Shoal standoff", true},
		{"UAE takeover of Socotra", true},
		{"Wuchang Uprising", true},
		{"King William's War", true},
		{"Napoleon", false},
		{"Erwin Rommel", false},
		{"British Empire", false},
		{"GitHub", false},
	}
	for _, c := range cases {
		got := names.IsBattle(c.name)
		assert.Equal(t, c.expected, got, "Expected names.IsBattle(%q) to be %t", c.name, c.expected)
	}
}
