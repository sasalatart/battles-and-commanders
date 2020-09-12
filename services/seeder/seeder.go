package seeder

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain"
	batcomsio "github.com/sasalatart/batcoms/services/io"
	"github.com/sasalatart/batcoms/services/logger"
	"github.com/sasalatart/batcoms/services/parser"
	"github.com/sasalatart/batcoms/store"
	uuid "github.com/satori/go.uuid"
)

// Service is the struct through which the seeder service gets initialized
type Service struct {
	ImportedData      *batcomsio.ImportedData
	FactionsCreator   store.FactionsCreator
	CommandersCreator store.CommandersCreator
	BattlesCreator    store.BattlesCreator
	Logger            logger.Service
}

// Seed fills factions, commanders and battles data stores with the service's available ImportedData
func (s *Service) Seed() {
	s.battles(s.commanders(), s.factions())
}

func (s *Service) factions() domain.IDsMap {
	fIDsByWikiID := make(domain.IDsMap)
	current := 0
	total := len(s.ImportedData.SFactionsByID)
	for _, sf := range s.ImportedData.SFactionsByID {
		fmt.Printf("\rSeeding factions (%d/%d)", current, total)
		current++
		input := domain.CreateFactionInput{
			WikiID:  sf.ID,
			URL:     sf.URL,
			Name:    sf.Name,
			Summary: sf.Extract,
		}
		if fID, err := s.FactionsCreator.CreateOne(input); err != nil {
			s.Logger.Error(errors.Wrapf(err, "Creating faction with URL %s", input.URL))
		} else {
			fIDsByWikiID[sf.ID] = fID
		}
	}
	s.Logger.Info("\nFinished seeding factions")
	return fIDsByWikiID
}

func (s *Service) commanders() domain.IDsMap {
	cIDsByWikiID := make(domain.IDsMap)
	current := 0
	total := len(s.ImportedData.SCommandersByID)
	for _, sc := range s.ImportedData.SCommandersByID {
		fmt.Printf("\rSeeding commanders (%d/%d)", current, total)
		current++
		input := domain.CreateCommanderInput{
			WikiID:  sc.ID,
			URL:     sc.URL,
			Name:    sc.Name,
			Summary: sc.Extract,
		}
		if cID, err := s.CommandersCreator.CreateOne(input); err != nil {
			s.Logger.Error(errors.Wrapf(err, "Creating commander with URL %s", input.URL))
		} else {
			cIDsByWikiID[sc.ID] = cID
		}
	}
	s.Logger.Info("\nFinished seeding commanders")
	return cIDsByWikiID
}

func (s *Service) battles(cIDsByWikiID, fIDsByWikiID domain.IDsMap) {
	current := 0
	total := len(s.ImportedData.SBattlesByID)
	for _, sb := range s.ImportedData.SBattlesByID {
		fmt.Printf("\rSeeding battles (%d/%d)", current, total)
		current++
		dates, err := parser.Date(sb.Date)
		if err != nil {
			s.Logger.Error(errors.Wrapf(err, "Parsing date %q", sb.Date))
			continue
		}
		input := domain.CreateBattleInput{
			WikiID:              sb.ID,
			URL:                 sb.URL,
			Name:                sb.Name,
			PartOf:              sb.PartOf,
			Summary:             sb.Extract,
			StartDate:           dates[0],
			EndDate:             dates[len(dates)-1],
			Location:            sb.Location,
			Result:              sb.Result,
			TerritorialChanges:  sb.TerritorialChanges,
			Strength:            sb.Strength,
			Casualties:          sb.Casualties,
			CommandersByFaction: make(domain.CommandersByFaction),
		}
		input.FactionsBySide.A = s.translateWikiIDs(sb.Factions.A, fIDsByWikiID)
		input.FactionsBySide.B = s.translateWikiIDs(sb.Factions.B, fIDsByWikiID)
		input.CommandersBySide.A = s.translateWikiIDs(sb.Commanders.A, cIDsByWikiID)
		input.CommandersBySide.B = s.translateWikiIDs(sb.Commanders.B, cIDsByWikiID)
		for sfWikiID, scWikiIDs := range sb.CommandersByFaction {
			fID := fIDsByWikiID[sfWikiID]
			input.CommandersByFaction[fID] = s.translateWikiIDs(scWikiIDs, cIDsByWikiID)
		}
		if _, err := s.BattlesCreator.CreateOne(input); err != nil {
			s.Logger.Error(errors.Wrapf(err, "Creating battle with URL %s", input.URL))
		}
	}
	s.Logger.Info("\nFinished seeding battles")
}

func (s *Service) translateWikiIDs(from []int, idsMapper domain.IDsMap) []uuid.UUID {
	result := []uuid.UUID{}
	for _, wikiID := range from {
		if id, ok := idsMapper[wikiID]; ok {
			result = append(result, id)
		} else {
			s.Logger.Error(fmt.Errorf("ID not found for WikiID %d", wikiID))
		}
	}
	return result
}
