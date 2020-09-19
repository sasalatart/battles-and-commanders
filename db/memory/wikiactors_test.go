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

	f := mocks.WikiFaction()
	require.NoError(t, repo.Save(f), "Saving valid faction")

	c := mocks.WikiCommander()
	c.ID = f.ID
	require.NoError(t, repo.Save(c), "Saving valid commander, even with the same ID of existing faction")

	assert.Nil(t, repo.Find(f.Kind, 100), "Finding inexistent faction")
	assert.Nil(t, repo.Find(c.Kind, 100), "Finding inexistent commander")

	byID := repo.Find(f.Kind, f.ID)
	require.NotNil(t, byID, "Finding existing faction by ID")
	assert.Equal(t, f, *byID, "Finding existing faction by ID")

	byURL := repo.FindByURL(f.Kind, f.URL)
	require.NotNil(t, byURL, "Finding existing faction by URL")
	require.Equal(t, f, *byURL, "Finding existing faction by URL")
}
