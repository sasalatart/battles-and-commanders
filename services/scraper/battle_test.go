package scraper_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/mocks"
	"github.com/sasalatart/batcoms/services/scraper"
	"github.com/sasalatart/batcoms/store/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSBattle(t *testing.T) {
	sbStore := memory.NewSBattlesStore()
	spStore := memory.NewSParticipantsStore()
	exporterMock := mocks.Exporter{}
	scraperService := scraper.New(sbStore, spStore, exporterMock.Export, ioutil.Discard)

	requireBattle := func(t *testing.T, url string) domain.SBattle {
		t.Helper()
		battle, err := scraperService.SBattle(url)
		require.NoErrorf(t, err, "When scraping %q", url)
		return battle
	}

	t.Run("UsualStructure", func(t *testing.T) {
		t.Parallel()
		const battleURL = "https://en.wikipedia.org/wiki/Battle_of_Austerlitz"
		battle := requireBattle(t, battleURL)

		attrsTests := []struct {
			attr     string
			got      interface{}
			expected interface{}
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
			{
				attr: "Factions",
				got:  battle.Factions,
				expected: domain.ScrapedSideParticipants{
					A: []int{21418258},
					B: []int{20611504, 266894},
				},
			},
			{
				attr: "Commanders",
				got:  battle.Commanders,
				expected: domain.ScrapedSideParticipants{
					A: []int{69880},
					B: []int{27126603, 251000, 11551, 14092123},
				},
			},
			{
				attr: "Strength",
				got:  battle.Strength,
				expected: domain.SideNumbers{
					A: "65,000–75,000",
					B: "84,000–95,000",
				},
			},
			{
				attr: "Casualties",
				got:  battle.Casualties,
				expected: domain.SideNumbers{
					A: "1,305 killed 6,991 wounded 573 captured",
					B: "16,000 killed and wounded 20,000 captured",
				},
			},
			{
				attr: "CommandersByFaction",
				got:  battle.CommandersByFaction,
				expected: domain.ScrapedCommandersByFaction{
					266894:   {11551, 14092123},
					20611504: {27126603, 251000},
					21418258: {69880},
				},
			},
		}
		for _, sc := range attrsTests {
			assert.Equal(t, sc.expected, sc.got, "Comparing %q", sc.attr)
		}

		participantsNamesTests := []struct {
			label    string
			id       int
			expected string
		}{
			{
				label:    "FactionA1",
				id:       21418258,
				expected: "First French Empire",
			},
			{
				label:    "FactionB1",
				id:       20611504,
				expected: "Russian Empire",
			},
			{
				label:    "FactionB2",
				id:       266894,
				expected: "Austrian Empire",
			},
			{
				label:    "CommanderA1",
				id:       69880,
				expected: "Napoleon",
			},
			{
				label:    "CommanderB1",
				id:       27126603,
				expected: "Alexander I of Russia",
			},
			{
				label:    "CommanderB2",
				id:       251000,
				expected: "Mikhail Kutuzov",
			},
			{
				label:    "CommanderB3",
				id:       11551,
				expected: "Francis II, Holy Roman Emperor",
			},
			{
				label:    "CommanderB4",
				id:       14092123,
				expected: "Franz von Weyrother",
			},
		}
		for _, pc := range participantsNamesTests {
			kind := domain.FactionKind
			if strings.HasPrefix(strings.ToLower(pc.label), "commander") {
				kind = domain.CommanderKind
			}
			got := spStore.Find(kind, pc.id)
			require.NotNilf(t, got, "Searching for participant with id %d for %q", pc.id, pc.label)
			assert.Equal(t, pc.expected, got.Name)
		}
	})

	t.Run("WithOverallCasualtiesAndLossesOnly", func(t *testing.T) {
		t.Parallel()
		battle := requireBattle(t, "https://en.wikipedia.org/wiki/Indian_Rebellion_of_1857")
		expected := domain.SideNumbers{
			A:  "",
			B:  "",
			AB: "6,000 British killed. As many as 800,000 Indians and possibly more, both in the rebellion and in famines and epidemics of disease in its wake, by comparison of 1857 population estimates with Indian Census of 1871.",
		}
		assert.Equal(t, expected, battle.Casualties)
	})

	t.Run("WithSpecificAndOverallCasualtiesAndLosses", func(t *testing.T) {
		t.Parallel()
		battle := requireBattle(t, "https://en.wikipedia.org/wiki/Chilean_Civil_War_of_1891")
		expected := domain.SideNumbers{
			A:  "",
			B:  "1 armoured frigate",
			AB: "5,000",
		}
		assert.Equal(t, expected, battle.Casualties)
	})
}
