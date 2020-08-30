package seeder

import (
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/services/io"
	"github.com/sasalatart/batcoms/services/io/json"
)

// JSONImport reads from files containing scraped battles and participants, and returns an
// io.ImportedData containing them
func JSONImport(battlesFileName, participantsFileName string) (*io.ImportedData, error) {
	importedData := new(io.ImportedData)
	if err := json.Import(battlesFileName, &importedData.SBattlesByID); err != nil {
		return nil, errors.Wrapf(err, "Importing battles from %s", battlesFileName)
	}
	if err := json.Import(participantsFileName, &importedData); err != nil {
		return nil, errors.Wrapf(err, "Importing participants from %s", participantsFileName)
	}
	return importedData, nil
}
