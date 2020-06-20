package service_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/sasalatart/batcoms/scraper/service"
	"github.com/sasalatart/batcoms/scraper/store"

	"github.com/sasalatart/batcoms/mocks"
)

func TestList(t *testing.T) {
	loggerMock := &mocks.Logger{}
	service := service.NewScraper(
		store.NewBattlesMem(),
		store.NewParticipantsMem(),
		loggerMock,
	)
	battlesList := service.List()

	t.Run("Logs for each list and with a final count of scraped battles", func(t *testing.T) {
		cc := []struct {
			log               string
			expectedSubstring string
		}{
			{loggerMock.Logs[0], "before_301"},
			{loggerMock.Logs[1], "301-1300"},
			{loggerMock.Logs[2], "1301-1600"},
			{loggerMock.Logs[3], "1601-1800"},
			{loggerMock.Logs[4], "1801-1900"},
			{loggerMock.Logs[5], "1901-2000"},
			{loggerMock.Logs[6], "since_2001"},
		}

		for _, c := range cc {
			if !strings.Contains(c.log, c.expectedSubstring) {
				t.Errorf("Scraper did not log results for %s", c.expectedSubstring)
			}
		}

		gotAmount := loggerMock.Logs[len(loggerMock.Logs)-1]
		expectedAmount := len(battlesList)
		if !strings.Contains(gotAmount, strconv.Itoa(expectedAmount)) {
			t.Errorf("Expected the amounts log to contain the number %d, but instead logged %q", expectedAmount, gotAmount)
		}
	})

	t.Run("Scraps more than 4500 battles", func(t *testing.T) {
		got := len(battlesList)
		expectedMin := 4500
		if got < 4500 {
			t.Errorf("Expected to scrap more than %d battles, but only got %d", expectedMin, got)
		}
	})

	t.Run("Contains battles for each one of the indexed lists", func(t *testing.T) {
		indexedBattlesNames := make(map[string]string)
		for _, battle := range battlesList {
			indexedBattlesNames[battle.Name] = battle.Name
		}

		cc := []struct {
			battleName string
			group      string
		}{
			{"Siege of Lachish", "before 301"},
			{"Battle of Actium", "before 301"},
			{"Battle of Adrianople", "between 301-1300"},
			{"Battle of Hattin", "between 301-1300"},
			{"Battle of Angora", "between 1301-1600"},
			{"Battle of Constantinople", "between 1301-1600"},
			{"Battle of Moscow (1612)", "between 1601-1800"},
			{"Battle of Zenta", "between 1601-1800"},
			{"Battle of Ulm", "between 1801-1900"},
			{"Battle of Austerlitz", "between 1801-1900"},
			{"Battle of Stalingrad", "between 1901-2000"},
			{"Battle of Chuncheon", "between 1901-2000"},
			{"Operation Defensive Shield", "after 2001"},
			{"Battle of Aguelhok", "after 2001"},
		}

		for _, c := range cc {
			if indexedBattlesNames[c.battleName] != c.battleName {
				t.Errorf("%q was not found in the list corresponding to battles %s", c.battleName, c.group)
			}
		}
	})
}
