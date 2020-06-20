package service

import (
	"fmt"
	"io"

	"github.com/gocolly/colly"
	"github.com/sasalatart/batcoms/scraper/store"
)

// Scraper is the starting point of any scraping activity done with Wikipedia, and contains all data
// that has been scraped up until a moment
type Scraper struct {
	BattlesStore      store.BattlesStore
	ParticipantsStore store.ParticipantsStore
	logger            io.Writer
}

// NewScraper creates a new scraper instance
func NewScraper(battlesStore store.BattlesStore, participantsStore store.ParticipantsStore, logger io.Writer) Scraper {
	return Scraper{battlesStore, participantsStore, logger}
}

// Export exports the scrapers' relevant information into two main normalized JSON files: one for
// battles, and another one for the participants in each one of those battles.
func (s *Scraper) Export(battlesFileName, participantsFileName string) error {
	if err := s.BattlesStore.Export(battlesFileName); err != nil {
		return fmt.Errorf("Failed exporting the Scraper's results: %s", err.Error())
	}
	if err := s.ParticipantsStore.Export(participantsFileName); err != nil {
		return fmt.Errorf("Failed exporting the Scraper's results: %s", err.Error())
	}
	return nil
}

func (s *Scraper) do(url string, subscribe func(c *colly.Collector)) {
	c := colly.NewCollector()
	subscribe(c)
	c.OnRequest(func(r *colly.Request) {
		message := fmt.Sprintf("Scraping %s\n", r.URL)
		s.logger.Write([]byte(message))
	})
	c.Visit(url)
}
