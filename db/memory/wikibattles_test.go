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
	wikiBattleMock := mocks.WikiBattle()

	require.NoError(t, repo.Save(wikiBattleMock), "Saving valid battle")
	assert.Nil(t, repo.Find(999), "Finding inexistent battle")
	found := repo.Find(wikiBattleMock.ID)
	require.NotNil(t, found, "Finding an existing battle")
	assert.Equal(t, wikiBattleMock, *found, "Finding an existing battle")
}
