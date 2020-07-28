package service_test

import (
	"strings"
	"testing"

	"github.com/sasalatart/batcoms/scraper/domain"
	"github.com/sasalatart/batcoms/scraper/service"
)

func TestPageSummary(t *testing.T) {
	t.Run("InvalidURL", func(t *testing.T) {
		t.Parallel()
		_, err := service.PageSummary("https://i-do-not-exist.org")
		if err == nil {
			t.Error("Expected to return an error but instead returned nil")
		}
		expectedSubstring := "not a wiki page"
		got := err.Error()
		if !strings.Contains(got, expectedSubstring) {
			t.Errorf("Expected error message to contain %q, but instead was %q", expectedSubstring, got)
		}
		if !t.Failed() {
			t.Log("Returns an error when given an invalid URL")
		}
	})

	t.Run("ValidURL", func(t *testing.T) {
		t.Parallel()
		got, err := service.PageSummary("https://en.wikipedia.org/wiki/Battle_of_Austerlitz")
		if err != nil {
			t.Errorf("Expected to not have an error, but instead received %v", err)
		}
		expected := domain.Summary{
			PageID:       118372,
			Type:         "standard",
			Title:        "Battle of Austerlitz",
			DisplayTitle: "Battle of Austerlitz",
			Description:  "Battle of the Napoleonic Wars",
			Extract:      "The Battle of Austerlitz, also known as the Battle of the Three Emperors, was one of the most important and decisive engagements of the Napoleonic Wars. In what is widely regarded as the greatest victory achieved by Napoleon, the Grande Armée of France defeated a larger Russian and Austrian army led by Emperor Alexander I and Holy Roman Emperor Francis II. The battle occurred near the town of Austerlitz in the Austrian Empire. Austerlitz brought the War of the Third Coalition to a rapid end, with the Treaty of Pressburg signed by the Austrians later in the month. The battle is often cited as a tactical masterpiece, in the same league as other historic engagements like Cannae or Gaugamela.",
		}
		if got != expected {
			t.Errorf("Expected %v, but got %v instead", expected, got)
		}
		if !t.Failed() {
			t.Log("Fills a domain.Summary when given a valid Wikipedia page URL")
		}
	})
}
