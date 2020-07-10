package service_test

import (
	"strings"
	"testing"

	"github.com/sasalatart/batcoms/scraper/domain"
	"github.com/sasalatart/batcoms/scraper/service"
)

func TestPageSummary(t *testing.T) {
	t.Run("Returns an error when given an invalid URL", func(t *testing.T) {
		_, err := service.PageSummary("https://i-do-not-exist.org")
		if err == nil {
			t.Error("Expected to return an error but instead returned nil")
		}
		expectedSubstring := "not a wiki page"
		got := err.Error()
		if !strings.Contains(got, expectedSubstring) {
			t.Errorf("Expected error message to contain %q, but instead was %q", expectedSubstring, got)
		}
	})

	t.Run("Fills a domain.Summary when given a valid Wikipedia page URL", func(t *testing.T) {
		got, err := service.PageSummary("https://en.wikipedia.org/wiki/Battle_of_Stalingrad")
		if err != nil {
			t.Errorf("Expected to not have an error, but instead received %v", err)
		}
		expected := domain.Summary{
			PageID:       4284,
			Type:         "standard",
			Title:        "Battle of Stalingrad",
			DisplayTitle: "Battle of Stalingrad",
			Description:  "Major battle of World War II",
			Extract:      "In the Battle of Stalingrad, Germany and its allies fought the Soviet Union for control of the city of Stalingrad in Southern Russia. Marked by fierce close-quarters combat and direct assaults on civilians in air raids, it is the bloodiest battle in the history of warfare, with an estimated 2 million total casualties. After their defeat at Stalingrad, the German High Command had to withdraw considerable military forces from the Western Front to replace their losses.",
		}
		if got != expected {
			t.Errorf("Expected %v, but got %v instead", expected, got)
		}
	})
}
