package store_test

import (
	"reflect"
	"testing"

	"github.com/sasalatart/batcoms/scraper/domain"
	"github.com/sasalatart/batcoms/scraper/store"
)

func TestBattlesMem(t *testing.T) {
	store := store.NewBattlesMem()

	if err := store.Save(domain.Battle{}); err == nil {
		t.Error("Expected to fail saving a battle with no ID nor URL")
	}
	if err := store.Save(domain.Battle{ID: 3}); err == nil {
		t.Error("Expected to fail saving a battle with no URL")
	}
	if err := store.Save(domain.Battle{URL: "www.3.example.org"}); err == nil {
		t.Error("Expected to fail saving a battle with no ID")
	}

	b1 := domain.Battle{ID: 1, URL: "www.1.example.org"}
	b2 := domain.Battle{ID: 2, URL: "www.2.example.org"}

	if err := store.Save(b1); err != nil {
		t.Errorf("Expected no error when saving battle %+v, but instead got: %s", b1, err)
	}
	if err := store.Save(b2); err != nil {
		t.Errorf("Expected no error when saving battle %+v, but instead got: %s", b2, err)
	}

	if inexistentBattle := store.Find(100); inexistentBattle != nil {
		t.Errorf("Expected to get nil when finding an inexistent battle, but instead got %+v", *inexistentBattle)
	}

	foundB1 := store.Find(b1.ID)
	if !reflect.DeepEqual(*foundB1, b1) {
		t.Errorf("Expected to find battle %+v when searching via its ID, but instead got %+v", b1, *foundB1)
	}
}
