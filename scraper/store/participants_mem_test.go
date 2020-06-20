package store_test

import (
	"reflect"
	"testing"

	"github.com/sasalatart/batcoms/scraper/domain"
	"github.com/sasalatart/batcoms/scraper/store"
)

func TestParticipantsMem(t *testing.T) {
	store := store.NewParticipantsMem()
	kind := domain.FactionKind

	if err := store.Save(domain.Participant{}); err == nil {
		t.Error("Expected to fail saving a participant with no ID nor URL")
	}
	if err := store.Save(domain.Participant{ID: 3}); err == nil {
		t.Error("Expected to fail saving a participant with no URL")
	}
	if err := store.Save(domain.Participant{URL: "www.3.example.org"}); err == nil {
		t.Error("Expected to fail saving a participant with no ID")
	}

	p1 := domain.Participant{Kind: kind, ID: 1, URL: "www.1.example.org", Flag: "www.flag1.example.org"}
	p2 := domain.Participant{Kind: kind, ID: 2, URL: "www.2.example.org", Flag: "www.flag2.example.org"}

	if err := store.Save(p1); err != nil {
		t.Errorf("Expected no error when saving participant %+v, but instead got: %s", p1, err)
	}
	if err := store.Save(p2); err != nil {
		t.Errorf("Expected no error when saving participant %+v, but instead got: %s", p2, err)
	}

	if inexistentParticipant := store.Find(kind, 100); inexistentParticipant != nil {
		t.Errorf("Expected to get nil when finding an inexistent participant, but instead got %+v", *inexistentParticipant)
	}

	foundP1 := store.Find(kind, p1.ID)
	if !reflect.DeepEqual(*foundP1, p1) {
		t.Errorf("Expected to find participant %+v when searching via its ID, but instead got %+v", p1, *foundP1)
	}
	foundP1 = store.FindByURL(kind, p1.URL)
	if !reflect.DeepEqual(*foundP1, p1) {
		t.Errorf("Expected to find participant %+v when searching via its URL, but instead got %+v", p1, *foundP1)
	}
	foundP1 = store.FindFactionByFlag(p1.Flag)
	if !reflect.DeepEqual(*foundP1, p1) {
		t.Errorf("Expected to find faction %+v when searching via its FlagURL, but instead got %+v", p1, *foundP1)
	}
}
