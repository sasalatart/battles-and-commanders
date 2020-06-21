package urls_test

import (
	"testing"

	"github.com/sasalatart/batcoms/scraper/urls"
)

func TestURLs(t *testing.T) {
	t.Run("NotSpecific", func(t *testing.T) {
		cases := []struct {
			url      string
			expected bool
		}{
			{"https://en.wikipedia.org/wiki/History_of_Norway", true},
			{"https://en.wikipedia.org/wiki/Norway", false},
			{"/wiki/History_of_Norway", true},
			{"/wiki/Norway", false},
			{"https://en.wikipedia.org/wiki/history_of_bavaria", true},
			{"https://en.wikipedia.org/wiki/bavaria", false},
			{"/wiki/history_of_bavaria", true},
			{"/wiki/bavaria", false},
			{"/wiki/Military_history_of_Australia_during_World_War_I", true},
			{"/wiki/List_of_kings_of_Leinster", true},
			{"/wiki/Campaign_against_Dong_Zhuo", true},
			{"/wiki/Dong_Zhuo", false},
		}

		for _, c := range cases {
			got := urls.NotSpecific(c.url)
			if got != c.expected {
				t.Errorf(
					"Expected urls.NotSpecific(%s) to be %t, but got %t",
					c.url,
					c.expected,
					got,
				)
			}
		}
	})

	t.Run("ShouldSkip", func(t *testing.T) {
		cases := []struct {
			url      string
			expected bool
		}{
			{"example.org?redlink=1", true},
			{"example.org?Redlink=1", true},
			{"example.org?foo=bar&redlink=1", true},
			{"example.org?foo=bar", false},
			{"example.org?foo=bar&baz=quux", false},
			{"https://en.wikipedia.org/wiki/Talk:Battle_of_Vyazma", true},
			{"https://en.wikipedia.org/wiki/Battle_of_Vyazma", false},
			{"/wiki/Talk:Battle_of_Vyazma", true},
			{"/wiki/Battle_of_Vyazma", false},
			{"/wiki/Wikipedia:Citation_needed", true},
			{"/wiki/Killed_in_action", true},
			{"/wiki/POW", true},
			{"/wiki/Prisoner_of_war", true},
			{"/wiki/Surrender_(military)", true},
			{"/wiki/Army", true},
			{"/wiki/Auxiliaries", true},
			{"/wiki/Auxiliary_Division", true},
			{"/wiki/Caliphate", true},
			{"/wiki/CIA", true},
			{"/wiki/Commandery", true},
			{"/wiki/Conscription", true},
			{"/wiki/Crusades", true},
			{"/wiki/Delta_Force", true},
			{"/wiki/Empire", true},
			{"/wiki/Flag", true},
			{"/wiki/Islam", true},
			{"/wiki/Islamism", true},
			{"/wiki/Jewish", true},
			{"/wiki/Jews", true},
			{"/wiki/Left-wing_politics", true},
			{"/wiki/Right-wing_politics", true},
			{"/wiki/Military_advisor", true},
			{"/wiki/Muslim_conquests", true},
			{"/wiki/Offensive_jihad", true},
			{"/wiki/Roman_Emperor", true},
			{"/wiki/United_States_Army_Rangers", true},
			{"/wiki/Napoleon", false},
		}

		for _, c := range cases {
			got := urls.ShouldSkip(c.url)
			if got != c.expected {
				t.Errorf(
					"Expected urls.ShouldSkip(%s) to be %t, but got %t",
					c.url,
					c.expected,
					got,
				)
			}
		}
	})
}
