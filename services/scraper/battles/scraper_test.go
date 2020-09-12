package battles_test

import (
	"io/ioutil"
	"testing"

	"github.com/sasalatart/batcoms/mocks"
	"github.com/sasalatart/batcoms/services/logger"
	"github.com/sasalatart/batcoms/services/scraper/battles"
	"github.com/sasalatart/batcoms/store/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExport(t *testing.T) {
	exporterMock := mocks.Exporter{}
	scraper := battles.NewScraper(
		memory.NewSBattlesStore(),
		memory.NewSParticipantsStore(),
		exporterMock.Export,
		logger.New(ioutil.Discard, ioutil.Discard),
	)
	battlesFileName := "mocked-battles.json"
	participantsFileName := "mocked-participants.json"
	err := scraper.ExportAll(battlesFileName, participantsFileName)

	require.NoError(t, err, "Exporting data from scraper")
	assert.Equal(t, 2, exporterMock.CalledTimes, "Amount of files exported")
	assert.Equal(t, []string{battlesFileName, participantsFileName}, exporterMock.FileNamesUsed, "Files exported order")
}
