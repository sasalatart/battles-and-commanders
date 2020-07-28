package store_test

import (
	"reflect"
	"testing"

	"github.com/sasalatart/batcoms/mocks"
	"github.com/sasalatart/batcoms/scraper/domain"
	"github.com/sasalatart/batcoms/scraper/store"
)

func TestParticipantsMem(t *testing.T) {
	store := store.NewParticipantsMem()

	f, err := mocks.Faction(domain.Participant{})
	if err = store.Save(f); err != nil {
		t.Fatalf("Expected no error when saving faction %+v, but instead got: %s", f, err)
	} else {
		t.Log("Saves factions")
	}

	c, err := mocks.Commander(domain.Participant{ID: f.ID})
	if err = store.Save(c); err != nil {
		t.Fatalf("Expected no error when saving commander %+v, but instead got: %s", c, err)
	} else {
		t.Log("Saves commanders")
	}

	if inexistentParticipant := store.Find(f.Kind, 100); inexistentParticipant != nil {
		t.Errorf("Expected to get nil when finding an inexistent participant, but instead got %+v", *inexistentParticipant)
	} else {
		t.Log("Returns errors when trying to find participants that do not exist")
	}

	found := store.Find(f.Kind, f.ID)
	if !reflect.DeepEqual(*found, f) {
		t.Errorf("Expected to find participant %+v when searching via its ID, but instead got %+v", f, *found)
	} else {
		t.Log("Finds participants given their IDs")
	}
	found = store.FindByURL(f.Kind, f.URL)
	if !reflect.DeepEqual(*found, f) {
		t.Errorf("Expected to find participant %+v when searching via its URL, but instead got %+v", f, *found)
	} else {
		t.Log("Finds participants given their URLs")
	}
}
