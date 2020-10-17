package seeder

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain/battles"
	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/domain/factions"
	"github.com/sasalatart/batcoms/pkg/dates"
	"github.com/sasalatart/batcoms/pkg/logger"
	uuid "github.com/satori/go.uuid"
)

// idsMap maps WikiIDs to their corresponding UUIDs
type idsMap map[int]uuid.UUID

type seeder struct {
	importedData     *ImportedData
	factionsWriter   factions.Writer
	commandersWriter commanders.Writer
	battlesWriter    battles.Writer
	logger           logger.Service
}

// Seed fills factions, commanders and battles data stores with the available ImportedData
func Seed(
	importedData *ImportedData,
	factionsWriter factions.Writer,
	commandersWriter commanders.Writer,
	battlesWriter battles.Writer,
	logger logger.Service,
) {
	service := seeder{importedData, factionsWriter, commandersWriter, battlesWriter, logger}
	service.battles(service.factions(), service.commanders())
}

func (s *seeder) factions() idsMap {
	fIDsByWikiID := make(idsMap)
	current := 0
	total := len(s.importedData.WikiFactionsByID)
	for _, wf := range s.importedData.WikiFactionsByID {
		fmt.Printf("\rSeeding factions (%d/%d)", current, total)
		current++
		input := factions.CreationInput{
			WikiID:  wf.ID,
			URL:     wf.URL,
			Name:    wf.Name,
			Summary: wf.Extract,
		}
		if fID, err := s.factionsWriter.CreateOne(input); err != nil {
			s.logger.Error(errors.Wrapf(err, "Error creating faction with URL %s", input.URL))
		} else {
			fIDsByWikiID[wf.ID] = fID
		}
	}
	s.logger.Info("\nFinished seeding factions\n")
	return fIDsByWikiID
}

func (s *seeder) commanders() idsMap {
	cIDsByWikiID := make(idsMap)
	current := 0
	total := len(s.importedData.WikiCommandersByID)
	for _, wc := range s.importedData.WikiCommandersByID {
		fmt.Printf("\rSeeding commanders (%d/%d)", current, total)
		current++
		input := commanders.CreationInput{
			WikiID:  wc.ID,
			URL:     wc.URL,
			Name:    wc.Name,
			Summary: wc.Extract,
		}
		if cID, err := s.commandersWriter.CreateOne(input); err != nil {
			s.logger.Error(errors.Wrapf(err, "Error creating commander with URL %s", input.URL))
		} else {
			cIDsByWikiID[wc.ID] = cID
		}
	}
	s.logger.Info("\nFinished seeding commanders\n")
	return cIDsByWikiID
}

func (s *seeder) battles(fIDsByWikiID, cIDsByWikiID idsMap) {
	current := 0
	total := len(s.importedData.WikiBattlesByID)
	for _, wb := range s.importedData.WikiBattlesByID {
		fmt.Printf("\rSeeding battles (%d/%d)", current, total)
		current++
		dates, err := dates.Parse(wb.Date)
		if err != nil {
			s.logger.Error(errors.Wrapf(err, "Error parsing date %q", wb.Date))
			continue
		}
		input := battles.CreationInput{
			WikiID:              wb.ID,
			URL:                 wb.URL,
			Name:                wb.Name,
			PartOf:              wb.PartOf,
			Summary:             wb.Extract,
			StartDate:           dates[0],
			EndDate:             dates[len(dates)-1],
			Location:            wb.Location,
			Result:              wb.Result,
			TerritorialChanges:  wb.TerritorialChanges,
			Strength:            wb.Strength,
			Casualties:          wb.Casualties,
			CommandersByFaction: make(battles.CommandersByFaction),
		}
		input.FactionsBySide.A = s.translateWikiIDs(wb.Factions.A, fIDsByWikiID)
		input.FactionsBySide.B = s.translateWikiIDs(wb.Factions.B, fIDsByWikiID)
		input.CommandersBySide.A = s.translateWikiIDs(wb.Commanders.A, cIDsByWikiID)
		input.CommandersBySide.B = s.translateWikiIDs(wb.Commanders.B, cIDsByWikiID)
		for sfWikiID, scWikiIDs := range wb.CommandersByFaction {
			fID := fIDsByWikiID[sfWikiID]
			input.CommandersByFaction[fID] = s.translateWikiIDs(scWikiIDs, cIDsByWikiID)
		}
		if _, err := s.battlesWriter.CreateOne(input); err != nil {
			s.logger.Error(errors.Wrapf(err, "Error creating battle with URL %s", input.URL))
		}
	}
	s.logger.Info("\nFinished seeding battles\n")
}

func (s *seeder) translateWikiIDs(from []int, idsMapper idsMap) []uuid.UUID {
	result := []uuid.UUID{}
	for _, wikiID := range from {
		if id, ok := idsMapper[wikiID]; ok {
			result = append(result, id)
		} else {
			s.logger.Error(fmt.Errorf("ID not found for WikiID %d", wikiID))
		}
	}
	return result
}
