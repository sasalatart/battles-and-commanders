package memory_test

import (
	"reflect"
	"testing"

	"github.com/sasalatart/batcoms/mocks"
	"github.com/sasalatart/batcoms/store/memory"
)

func TestSBattlesMemStore(t *testing.T) {
	store := memory.NewSBattlesStore()

	b := mocks.SBattle()
	if err := store.Save(b); err != nil {
		t.Fatalf("Expected no error when saving battle %+v, but instead got: %s", b, err)
	} else {
		t.Log("Saves battles")
	}

	if inexistent := store.Find(999); inexistent != nil {
		t.Errorf("Expected to get nil when finding an inexistent battle, but instead got %+v", *inexistent)
	} else {
		t.Log("Returns errors when trying to find battles that do not exist")
	}

	found := store.Find(b.ID)
	if !reflect.DeepEqual(*found, b) {
		t.Errorf("Expected to find battle %+v when searching via its ID, but instead got %+v", b, *found)
	} else {
		t.Log("Finds battles given their IDs")
	}
}
