package memory_test

import (
	"testing"

	"github.com/sasalatart/batcoms/mocks"
	"github.com/sasalatart/batcoms/store/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSBattlesMemStore(t *testing.T) {
	store := memory.NewSBattlesStore()
	b := mocks.SBattle()

	require.NoError(t, store.Save(b), "Saving valid battle")
	assert.Nil(t, store.Find(999), "Finding inexistent battle")
	found := store.Find(b.ID)
	require.NotNil(t, found, "Finding an existing battle")
	require.Equal(t, b, *found, "Finding an existing battle")
}
