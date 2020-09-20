package battles

import (
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/db/memory"
	"github.com/sasalatart/batcoms/pkg/io"
	"github.com/sasalatart/batcoms/pkg/logger"
)

// Scraper is the struct encapsulating all the necessary behaviour to scrape battles, one by one, as
// well as exporting the results into files
type Scraper struct {
	wikiActorsRepo  *memory.WikiActorsRepo
	wikiBattlesRepo *memory.WikiBattlesRepo
	exporterFunc    io.ExporterFunc
	logger          logger.Service
}

// NewScraper creates a new instance of battles.Scraper
func NewScraper(
	wikiActorsRepo *memory.WikiActorsRepo,
	wikiBattlesRepo *memory.WikiBattlesRepo,
	exporterFunc io.ExporterFunc,
	l logger.Service,
) Scraper {
	return Scraper{
		wikiActorsRepo:  wikiActorsRepo,
		wikiBattlesRepo: wikiBattlesRepo,
		exporterFunc:    exporterFunc,
		logger:          l,
	}
}

// ExportAll exports the scrapers' relevant information into two main normalized JSON files: one for
// all battles, and another one for all of the factions and commanders in each one of those battles
func (s *Scraper) ExportAll(actorsFileName, battlesFileName string) error {
	if err := s.wikiActorsRepo.Export(actorsFileName, s.exporterFunc); err != nil {
		return errors.Wrap(err, "Exporting actors results")
	}
	if err := s.wikiBattlesRepo.Export(battlesFileName, s.exporterFunc); err != nil {
		return errors.Wrap(err, "Exporting battles results")
	}
	return nil
}
