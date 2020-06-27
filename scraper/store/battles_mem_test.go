package store_test

import (
	"reflect"
	"testing"

	"github.com/sasalatart/batcoms/mocks"
	"github.com/sasalatart/batcoms/scraper/domain"
	"github.com/sasalatart/batcoms/scraper/store"
)

func TestBattlesMem(t *testing.T) {
	store := store.NewBattlesMem()

	b, err := mocks.Battle(domain.Battle{})
	if err = store.Save(b); err != nil {
		t.Fatalf("Expected no error when saving battle %+v, but instead got: %s", b, err)
	}

	if inexistent := store.Find(999); inexistent != nil {
		t.Errorf("Expected to get nil when finding an inexistent battle, but instead got %+v", *inexistent)
	}

	found := store.Find(b.ID)
	if !reflect.DeepEqual(*found, b) {
		t.Errorf("Expected to find battle %+v when searching via its ID, but instead got %+v", b, *found)
	}
}
