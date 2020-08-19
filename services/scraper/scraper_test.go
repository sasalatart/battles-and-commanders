package scraper_test

import (
	"testing"

	"github.com/sasalatart/batcoms/mocks"
	"github.com/sasalatart/batcoms/services/scraper"
	"github.com/sasalatart/batcoms/store/memory"
)

func TestExport(t *testing.T) {
	exporterMock := mocks.Exporter{}
	s := scraper.New(
		memory.NewSBattlesStore(),
		memory.NewSParticipantsStore(),
		exporterMock.Export,
		&mocks.Logger{},
	)
	battlesFileName := "mocked-battles.json"
	participantsFileName := "mocked-participants.json"
	if err := s.Export(battlesFileName, participantsFileName); err != nil {
		t.Fatalf("Expected scraper.Export to not throw any errors, but got %s", err)
	}
	if exporterMock.CalledTimes != 2 {
		t.Errorf("Expected to export two files, but instead exported %d", exporterMock.CalledTimes)
	}
	if exporterMock.FileNamesUsed[0] != battlesFileName || exporterMock.FileNamesUsed[1] != participantsFileName {
		t.Errorf(
			"Expected to export first to %s and then to %s, but instead exported to %v",
			battlesFileName,
			participantsFileName,
			exporterMock.FileNamesUsed,
		)
	}
	if !t.Failed() {
		t.Log("Exports scraped data to the specified battles and participants files")
	}
}
