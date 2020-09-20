package seeder

import (
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain/wikiactors"
	"github.com/sasalatart/batcoms/domain/wikibattles"
	"github.com/sasalatart/batcoms/pkg/io/json"
)

// ImportedData contains scraped battles and actors that have been read from a previously exported
// file. These have been indexed by their Wikipedia IDs
type ImportedData struct {
	WikiBattlesByID    map[string]wikibattles.Battle
	WikiFactionsByID   map[string]wikiactors.Actor `json:"FactionsByID"`
	WikiCommandersByID map[string]wikiactors.Actor `json:"CommandersByID"`
}

// JSONImport reads from files containing scraped battles and actors, and stores results into the
// given *seeder.ImportedData
func JSONImport(importedData *ImportedData, actorsFileName, battlesFileName string) error {
	if err := json.Import(actorsFileName, &importedData); err != nil {
		return errors.Wrapf(err, "Importing actors from %s", actorsFileName)
	}
	if err := json.Import(battlesFileName, &importedData.WikiBattlesByID); err != nil {
		return errors.Wrapf(err, "Importing battles from %s", battlesFileName)
	}
	return nil
}
