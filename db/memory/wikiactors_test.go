package memory_test

import (
	"testing"

	"github.com/sasalatart/batcoms/db/memory"
	"github.com/sasalatart/batcoms/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWikiActorsMemRepository(t *testing.T) {
	repo := memory.NewWikiActorsRepo()

	wikiFactionMock := mocks.WikiFaction()
	require.NoError(t, repo.Save(wikiFactionMock), "Saving valid faction")

	wikiCommanderMock := mocks.WikiCommander()
	wikiCommanderMock.ID = wikiFactionMock.ID
	require.NoError(t, repo.Save(wikiCommanderMock), "Saving valid commander")

	assert.Nil(t, repo.Find(wikiFactionMock.Kind, 100), "Finding inexistent faction")
	assert.Nil(t, repo.Find(wikiCommanderMock.Kind, 100), "Finding inexistent commander")

	byID := repo.Find(wikiFactionMock.Kind, wikiFactionMock.ID)
	require.NotNil(t, byID, "Finding existing faction by ID")
	assert.Equal(t, wikiFactionMock, *byID, "Finding existing faction by ID")

	byURL := repo.FindByURL(wikiFactionMock.Kind, wikiFactionMock.URL)
	require.NotNil(t, byURL, "Finding existing faction by URL")
	assert.Equal(t, wikiFactionMock, *byURL, "Finding existing faction by URL")
}
