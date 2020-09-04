package seeder

import (
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/services/io"
	"github.com/sasalatart/batcoms/services/io/json"
)

// JSONImport reads from files containing scraped battles and participants, and stores results
// into the input *io.ImportedData
func JSONImport(battlesFileName, participantsFileName string, importedData *io.ImportedData) error {
	if err := json.Import(battlesFileName, &importedData.SBattlesByID); err != nil {
		return errors.Wrapf(err, "Importing battles from %s", battlesFileName)
	}
	if err := json.Import(participantsFileName, &importedData); err != nil {
		return errors.Wrapf(err, "Importing participants from %s", participantsFileName)
	}
	return nil
}
