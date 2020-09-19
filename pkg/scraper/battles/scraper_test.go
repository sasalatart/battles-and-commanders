package battles_test

import (
	"io/ioutil"
	"testing"

	"github.com/sasalatart/batcoms/db/memory"
	"github.com/sasalatart/batcoms/mocks"
	"github.com/sasalatart/batcoms/pkg/logger"
	"github.com/sasalatart/batcoms/pkg/scraper/battles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExport(t *testing.T) {
	exporterMock := mocks.Exporter{}
	scraper := battles.NewScraper(
		memory.NewWikiActorsRepo(),
		memory.NewWikiBattlesRepo(),
		exporterMock.Export,
		logger.New(ioutil.Discard, ioutil.Discard),
	)
	actorsFileName := "mocked-actors.json"
	battlesFileName := "mocked-battles.json"
	err := scraper.ExportAll(actorsFileName, battlesFileName)

	require.NoError(t, err, "Exporting data from scraper")
	assert.Equal(t, 2, exporterMock.CalledTimes, "Amount of files exported")
	assert.Equal(t, []string{actorsFileName, battlesFileName}, exporterMock.FileNamesUsed, "Files exported order")
}
