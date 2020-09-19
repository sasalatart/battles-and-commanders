package memory_test

import (
	"testing"

	"github.com/sasalatart/batcoms/db/memory"
	"github.com/sasalatart/batcoms/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWikiBattlesMemRepository(t *testing.T) {
	repo := memory.NewWikiBattlesRepo()
	b := mocks.WikiBattle()

	require.NoError(t, repo.Save(b), "Saving valid battle")
	assert.Nil(t, repo.Find(999), "Finding inexistent battle")
	found := repo.Find(b.ID)
	require.NotNil(t, found, "Finding an existing battle")
	require.Equal(t, b, *found, "Finding an existing battle")
}
