package memory_test

import (
	"testing"

	"github.com/sasalatart/batcoms/mocks"
	"github.com/sasalatart/batcoms/store/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSParticipantsMemStore(t *testing.T) {
	store := memory.NewSParticipantsStore()

	f := mocks.SFaction()
	require.NoError(t, store.Save(f), "Saving valid faction")

	c := mocks.SCommander()
	c.ID = f.ID
	require.NoError(t, store.Save(c), "Saving valid commander, even with the same ID of existing faction")

	assert.Nil(t, store.Find(f.Kind, 100), "Finding inexistent faction")
	assert.Nil(t, store.Find(c.Kind, 100), "Finding inexistent commander")

	byID := store.Find(f.Kind, f.ID)
	require.NotNil(t, byID, "Finding existing faction by ID")
	assert.Equal(t, f, *byID, "Finding existing faction by ID")

	byURL := store.FindByURL(f.Kind, f.URL)
	require.NotNil(t, byURL, "Finding existing faction by URL")
	require.Equal(t, f, *byURL, "Finding existing faction by URL")
}
