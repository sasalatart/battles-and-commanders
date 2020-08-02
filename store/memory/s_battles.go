package memory

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/store/exports"
)

// SBattlesStore is an in-memory implementation of store.SBattles
type SBattlesStore struct {
	mutex     sync.RWMutex
	validator *validator.Validate
	byID      map[int]*domain.SBattle
}

// NewSBattlesStore returns a pointer to an empty, ready-to-use memory.SBattlesStore
func NewSBattlesStore() *SBattlesStore {
	return &SBattlesStore{validator: validator.New(), byID: make(map[int]*domain.SBattle)}
}

// Find searches for the pointer to a battle by its ID. If none is found, nil is returned
func (s *SBattlesStore) Find(id int) *domain.SBattle {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	battle, found := s.byID[id]
	if found != true {
		return nil
	}
	return battle
}

// Save stores the given battle. It returns an error if the battle does not have an ID or a URL
func (s *SBattlesStore) Save(b domain.SBattle) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.validator.Struct(b); err != nil {
		return errors.Wrapf(err, "Saving battle with URL %s", b.URL)
	}
	s.byID[b.ID] = &b
	return nil
}

// Export saves data stored to the specified file, in JSON format
func (s *SBattlesStore) Export(fileName string) error {
	return exports.JSON(fileName, s.byID)
}
