package store

import (
	"fmt"
	"sync"

	"github.com/sasalatart/batcoms/scraper/domain"
	"github.com/sasalatart/batcoms/scraper/exports"
)

// BattlesMem is an in-memory implementation of BattlesStore
type BattlesMem struct {
	mutex sync.Mutex
	byID  map[int]*domain.Battle
}

// NewBattlesMem returns a pointer to an empty, ready-to-use BattlesMem store
func NewBattlesMem() *BattlesMem {
	return &BattlesMem{byID: make(map[int]*domain.Battle)}
}

// Find searches for the pointer to a battle by its ID. If none is found, nil is returned
func (r *BattlesMem) Find(id int) *domain.Battle {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	battle, found := r.byID[id]
	if found != true {
		return nil
	}
	return battle
}

// Save stores the given battle. It returns an error if the battle does not have an ID or a URL
func (r *BattlesMem) Save(b domain.Battle) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if b.ID == 0 {
		return fmt.Errorf("Expected battle %+v to have an ID, but got none", b)
	}
	if b.URL == "" {
		return fmt.Errorf("Expected battle %+v to have an URL, but got none", b)
	}
	r.byID[b.ID] = &b
	return nil
}

// Export saves data stored to the specified file, in JSON format
func (r *BattlesMem) Export(fileName string) error {
	return exports.JSON(fileName, r.byID)
}
