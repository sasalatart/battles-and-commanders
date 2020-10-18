package battles

import (
	"github.com/sasalatart/batcoms/db/memory"
	"github.com/sasalatart/batcoms/domain/wikiactors"
	"github.com/sasalatart/batcoms/domain/wikibattles"
	"github.com/sasalatart/batcoms/pkg/logger"
)

// Scraper is the struct encapsulating all the necessary behaviour to scrape battles, one by one, as
// well as exporting the results into files
type Scraper struct {
	wikiActorsRepo  *memory.WikiActorsRepo
	wikiBattlesRepo *memory.WikiBattlesRepo
	logger          logger.Interface
}

// ExportedData is the struct used to retrieve all factions, commanders and battles that have been
// scraped after successive runs of scraper.ScrapeOne. All of these have been normalized by their
// Wikipedia IDs
type ExportedData struct {
	FactionsByID   map[int]*wikiactors.Actor
	CommandersByID map[int]*wikiactors.Actor
	BattlesByID    map[int]*wikibattles.Battle
}

// NewScraper creates a new instance of battles.Scraper
func NewScraper(l logger.Interface) Scraper {
	return Scraper{
		wikiActorsRepo:  memory.NewWikiActorsRepo(),
		wikiBattlesRepo: memory.NewWikiBattlesRepo(),
		logger:          l,
	}
}

// Data builds a battles.ExportedData struct with all of the scraped data obtained from successive
// runs of scraper.ScrapeOne
func (s *Scraper) Data() ExportedData {
	factionsByID, commandersByID := s.wikiActorsRepo.Data()
	return ExportedData{
		FactionsByID:   factionsByID,
		CommandersByID: commandersByID,
		BattlesByID:    s.wikiBattlesRepo.Data(),
	}
}
