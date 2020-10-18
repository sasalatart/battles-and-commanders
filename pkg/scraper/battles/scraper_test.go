package battles_test

import (
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain/statistics"
	"github.com/sasalatart/batcoms/domain/wikibattles"
	"github.com/sasalatart/batcoms/pkg/logger"
	"github.com/sasalatart/batcoms/pkg/scraper/battles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWikiBattle(t *testing.T) {
	requireBattle := func(t *testing.T, s *battles.Scraper, url string) wikibattles.Battle {
		t.Helper()
		battle, err := s.ScrapeOne(url)
		require.NoErrorf(t, err, "When scraping %q", url)
		return battle
	}

	t.Run("UsualStructure", func(t *testing.T) {
		t.Parallel()

		scraper := battles.NewScraper(logger.NewDiscard())
		const battleURL = "https://en.wikipedia.org/wiki/Battle_of_Austerlitz"
		battle := requireBattle(t, &scraper, battleURL)
		data := scraper.Data()

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
				expected: wikibattles.SideActors{
					A: []int{21418258},
					B: []int{20611504, 266894},
				},
			},
			{
				attr: "Commanders",
				got:  battle.Commanders,
				expected: wikibattles.SideActors{
					A: []int{69880},
					B: []int{27126603, 251000, 11551, 14092123},
				},
			},
			{
				attr: "Strength",
				got:  battle.Strength,
				expected: statistics.SideNumbers{
					A: "65,000–75,000",
					B: "84,000–95,000",
				},
			},
			{
				attr: "Casualties",
				got:  battle.Casualties,
				expected: statistics.SideNumbers{
					A: "1,305 killed 6,991 wounded 573 captured",
					B: "16,000 killed and wounded 20,000 captured",
				},
			},
			{
				attr: "CommandersByFaction",
				got:  battle.CommandersByFaction,
				expected: wikibattles.CommandersByFaction{
					266894:   {11551, 14092123},
					20611504: {27126603, 251000},
					21418258: {69880},
				},
			},
		}
		for _, c := range attrsTests {
			assert.Equal(t, c.expected, c.got, "Comparing %q", c.attr)
		}

		actorsNamesTests := []struct {
			label    string
			expected string
		}{
			{
				label:    "FactionA1",
				expected: "First French Empire",
			},
			{
				label:    "FactionB1",
				expected: "Russian Empire",
			},
			{
				label:    "FactionB2",
				expected: "Austrian Empire",
			},
			{
				label:    "CommanderA1",
				expected: "Napoleon",
			},
			{
				label:    "CommanderB1",
				expected: "Alexander I of Russia",
			},
			{
				label:    "CommanderB2",
				expected: "Mikhail Kutuzov",
			},
			{
				label:    "CommanderB3",
				expected: "Francis II, Holy Roman Emperor",
			},
			{
				label:    "CommanderB4",
				expected: "Franz von Weyrother",
			},
		}

		factionsNames := []string{}
		for _, f := range data.FactionsByID {
			factionsNames = append(factionsNames, f.Name)
		}
		commandersNames := []string{}
		for _, c := range data.CommandersByID {
			commandersNames = append(commandersNames, c.Name)
		}
		for _, pc := range actorsNamesTests {
			isFaction := strings.HasPrefix(strings.ToLower(pc.label), "faction")
			if isFaction {
				require.Contains(t, factionsNames, pc.expected, "Should contain faction with name %s", pc.expected)
				return
			}
			require.Contains(t, commandersNames, pc.expected, "Should contain commander with name %s", pc.expected)
		}
	})

	t.Run("WithOverallCasualtiesAndLossesOnly", func(t *testing.T) {
		t.Parallel()
		scraper := battles.NewScraper(logger.NewDiscard())
		battle := requireBattle(t, &scraper, "https://en.wikipedia.org/wiki/Indian_Rebellion_of_1857")
		expected := statistics.SideNumbers{
			A:  "",
			B:  "",
			AB: "6,000 British killed. As many as 800,000 Indians and possibly more, both in the rebellion and in famines and epidemics of disease in its wake, by comparison of 1857 population estimates with Indian Census of 1871.",
		}
		assert.Equal(t, expected, battle.Casualties)
	})

	t.Run("WithSpecificAndOverallCasualtiesAndLosses", func(t *testing.T) {
		t.Parallel()
		scraper := battles.NewScraper(logger.NewDiscard())
		battle := requireBattle(t, &scraper, "https://en.wikipedia.org/wiki/Chilean_Civil_War_of_1891")
		expected := statistics.SideNumbers{
			A:  "",
			B:  "1 armoured frigate",
			AB: "5,000",
		}
		assert.Equal(t, expected, battle.Casualties)
	})

	t.Run("WithNoInfoBox", func(t *testing.T) {
		t.Parallel()
		scraper := battles.NewScraper(logger.NewDiscard())
		_, err := scraper.ScrapeOne("https://en.wikipedia.org/wiki/Diplomatic_Revolution")
		require.Errorf(t, err, "Scraping a URL with no info box")
		assert.EqualError(t, errors.Cause(err), battles.ErrNoInfoBox.Error(), "Error should be a battles.ErrNoInfoBox")
	})

	t.Run("WithMoreThanOneInfoBox", func(t *testing.T) {
		t.Parallel()
		scraper := battles.NewScraper(logger.NewDiscard())
		_, err := scraper.ScrapeOne("https://en.wikipedia.org/wiki/Big_Sandy_Expedition")
		require.Errorf(t, err, "Scraping a URL with more than one info box")
		assert.EqualError(t, errors.Cause(err), battles.ErrMoreThanOneInfoBox.Error(), "Error should be a battles.ErrMoreThanOneInfoBox")
	})
}
