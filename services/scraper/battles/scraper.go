package battles

import (
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/services/io"
	"github.com/sasalatart/batcoms/services/logger"
	"github.com/sasalatart/batcoms/store/memory"
)

// Scraper is the struct encapsulating all the necessary behaviour to scrape battles, one by one, as
// well as exporting the results into files
type Scraper struct {
	sBattlesStore      *memory.SBattlesStore
	sParticipantsStore *memory.SParticipantsStore
	exporterFunc       io.ExporterFunc
	logger             logger.Service
}

// NewScraper creates a new instance of battles.Scraper
func NewScraper(
	sBattlesStore *memory.SBattlesStore,
	sParticipantsStore *memory.SParticipantsStore,
	exporterFunc io.ExporterFunc,
	l logger.Service,
) Scraper {
	return Scraper{
		sBattlesStore:      sBattlesStore,
		sParticipantsStore: sParticipantsStore,
		exporterFunc:       exporterFunc,
		logger:             l,
	}
}

// ExportAll exports the scrapers' relevant information into two main normalized JSON files: one for
// all battles, and another one for all of the participants in each one of those battles
func (s *Scraper) ExportAll(battlesFileName, participantsFileName string) error {
	if err := s.sBattlesStore.Export(battlesFileName, s.exporterFunc); err != nil {
		return errors.Wrap(err, "Exporting battles results")
	}
	if err := s.sParticipantsStore.Export(participantsFileName, s.exporterFunc); err != nil {
		return errors.Wrap(err, "Exporting participants results")
	}
	return nil
}
