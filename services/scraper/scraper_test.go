package scraper_test

import (
	"testing"

	"github.com/sasalatart/batcoms/mocks"
	"github.com/sasalatart/batcoms/services/scraper"
	"github.com/sasalatart/batcoms/store/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExport(t *testing.T) {
	exporterMock := mocks.Exporter{}
	s := scraper.New(
		memory.NewSBattlesStore(),
		memory.NewSParticipantsStore(),
		exporterMock.Export,
		new(mocks.Logger),
	)
	battlesFileName := "mocked-battles.json"
	participantsFileName := "mocked-participants.json"
	err := s.Export(battlesFileName, participantsFileName)

	require.NoError(t, err, "Exporting data from scraper")
	assert.Equal(t, 2, exporterMock.CalledTimes, "Amount of files exported")
	assert.Equal(t, []string{battlesFileName, participantsFileName}, exporterMock.FileNamesUsed, "Files exported order")
}
