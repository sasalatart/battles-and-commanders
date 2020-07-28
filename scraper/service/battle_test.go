package service_test

import (
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/sasalatart/batcoms/scraper/domain"
	"github.com/sasalatart/batcoms/scraper/service"
	"github.com/sasalatart/batcoms/scraper/store"
)

func TestBattle(t *testing.T) {
	battlesStore := store.NewBattlesMem()
	participantsStore := store.NewParticipantsMem()
	service := service.NewScraper(battlesStore, participantsStore, ioutil.Discard)

	assertBattle := func(t *testing.T, url string) domain.Battle {
		t.Helper()
		battle, err := service.Battle(url)
		if err != nil {
			t.Fatalf("Unexpected service.Battle error: %s", err)
		}
		return battle
	}

	assertStruct := func(t *testing.T, name string, got, expected interface{}) {
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Expected %s to be %+v, but instead got %+v", name, expected, got)
		}
	}

	t.Run("UsualStructure", func(t *testing.T) {
		t.Parallel()
		const battleURL = "https://en.wikipedia.org/wiki/Battle_of_Austerlitz"
		battle := assertBattle(t, battleURL)

		stringAttrCases := []struct {
			attr     string
			got      string
			expected string
		}{
			{
				attr:     "URL",
				got:      battle.URL,
				expected: battleURL,
			},
			{
				attr:     "Name",
				got:      battle.Name,
				expected: "Battle of Austerlitz",
			},
			{
				attr:     "PartOf",
				got:      battle.PartOf,
				expected: "Part of the War of the Third Coalition",
			},
			{
				attr:     "Description",
				got:      battle.Description,
				expected: "Battle of the Napoleonic Wars",
			},
			{
				attr:     "Extract",
				got:      battle.Extract,
				expected: "The Battle of Austerlitz, also known as the Battle of the Three Emperors, was one of the most important and decisive engagements of the Napoleonic Wars. In what is widely regarded as the greatest victory achieved by Napoleon, the Grande Armée of France defeated a larger Russian and Austrian army led by Emperor Alexander I and Holy Roman Emperor Francis II. The battle occurred near the town of Austerlitz in the Austrian Empire. Austerlitz brought the War of the Third Coalition to a rapid end, with the Treaty of Pressburg signed by the Austrians later in the month. The battle is often cited as a tactical masterpiece, in the same league as other historic engagements like Cannae or Gaugamela.",
			},
			{
				attr:     "Date",
				got:      battle.Date,
				expected: "2 December 1805",
			},
			{
				attr:     "Place",
				got:      battle.Location.Place,
				expected: "Austerlitz, Moravia, Austria",
			},
			{
				attr:     "Latitude",
				got:      battle.Location.Latitude,
				expected: "49°8′N",
			},
			{
				attr:     "Longitude",
				got:      battle.Location.Longitude,
				expected: "16°46′E",
			},
			{
				attr:     "Result",
				got:      battle.Result,
				expected: "Decisive French victory. Treaty of Pressburg. Effective end of the Third Coalition",
			},
			{
				attr:     "TerritorialChanges",
				got:      battle.TerritorialChanges,
				expected: "Dissolution of the Holy Roman Empire and creation of the Confederation of the Rhine",
			},
		}
		for _, sc := range stringAttrCases {
			if sc.got != sc.expected {
				t.Errorf("Expected battle %s to be %q, but instead got %q", sc.attr, sc.expected, sc.got)
			}
		}

		assertStruct(t, "factions", battle.Factions, domain.SideParticipants{
			A: []int{21418258},
			B: []int{20611504, 266894},
		})
		assertStruct(t, "commanders", battle.Commanders, domain.SideParticipants{
			A: []int{69880},
			B: []int{27126603, 251000, 11551, 14092123},
		})
		assertStruct(t, "strength", battle.Strength, domain.SideNumbers{
			A:  "65,000–75,000",
			B:  "84,000–95,000",
			AB: "",
		})
		assertStruct(t, "casualties", battle.Casualties, domain.SideNumbers{
			A:  "1,305 killed 6,991 wounded 573 captured",
			B:  "16,000 killed and wounded 20,000 captured",
			AB: "",
		})
		assertStruct(t, "commanders grouping", battle.CommandersByFaction, map[int][]int{
			266894:   {11551, 14092123},
			20611504: {27126603, 251000},
			21418258: {69880},
		})

		participantsNamesCases := []struct {
			label        string
			id           int
			expectedName string
		}{
			{
				label:        "FactionA1",
				id:           21418258,
				expectedName: "First French Empire",
			},
			{
				label:        "FactionB1",
				id:           20611504,
				expectedName: "Russian Empire",
			},
			{
				label:        "FactionB2",
				id:           266894,
				expectedName: "Austrian Empire",
			},
			{
				label:        "CommanderA1",
				id:           69880,
				expectedName: "Napoleon",
			},
			{
				label:        "CommanderB1",
				id:           27126603,
				expectedName: "Alexander I of Russia",
			},
			{
				label:        "CommanderB2",
				id:           251000,
				expectedName: "Mikhail Kutuzov",
			},
			{
				label:        "CommanderB3",
				id:           11551,
				expectedName: "Francis II, Holy Roman Emperor",
			},
			{
				label:        "CommanderB4",
				id:           14092123,
				expectedName: "Franz von Weyrother",
			},
		}
		for _, pc := range participantsNamesCases {
			kind := domain.FactionKind
			if strings.HasPrefix(strings.ToLower(pc.label), "commander") {
				kind = domain.CommanderKind
			}
			participant := participantsStore.Find(kind, pc.id)
			if participant == nil {
				t.Fatalf("No participant found with id %d for %q", pc.id, pc.label)
			}
			got := participant.Name
			if pc.expectedName != got {
				t.Errorf("Expected %s to have name %q, but instead got %q", pc.label, pc.expectedName, got)
			}
		}
		if !t.Failed() {
			t.Log("Builds the battle properly")
		}
	})

	t.Run("WithOverallCasualtiesAndLossesOnly", func(t *testing.T) {
		t.Parallel()
		battle := assertBattle(t, "https://en.wikipedia.org/wiki/Indian_Rebellion_of_1857")
		assertStruct(t, "casualties", battle.Casualties, domain.SideNumbers{
			A:  "",
			B:  "",
			AB: "6,000 Europeans killed. As many as 800,000 Indians and possibly more, both in the rebellion and in famines and epidemics of disease in its wake, by comparison of 1857 population estimates with Indian Census of 1871.",
		})
		if !t.Failed() {
			t.Log("Stores the information in the appropriate place")
		}
	})

	t.Run("WithSpecificAndOverallCasualtiesAndLosses", func(t *testing.T) {
		t.Parallel()
		battle := assertBattle(t, "https://en.wikipedia.org/wiki/Chilean_Civil_War_of_1891")
		assertStruct(t, "casualties", battle.Casualties, domain.SideNumbers{
			A:  "",
			B:  "1 armoured frigate",
			AB: "5,000",
		})
		if !t.Failed() {
			t.Log("Stores the information in the appropriate places")
		}
	})
}
