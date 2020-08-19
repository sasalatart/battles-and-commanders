package memory

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/services/io"
)

// SParticipantsStore is an in-memory implementation of store.SParticipants
type SParticipantsStore struct {
	mutex       sync.RWMutex
	validator   *validator.Validate
	byIDByKind  map[domain.ParticipantKind]map[int]*domain.SParticipant
	byURLByKind map[domain.ParticipantKind]map[string]int
}

// NewSParticipantsStore returns a pointer to an empty, ready-to-use memory.SParticipantsStore
func NewSParticipantsStore() *SParticipantsStore {
	byIDByKind := make(map[domain.ParticipantKind]map[int]*domain.SParticipant)
	byURLByKind := make(map[domain.ParticipantKind]map[string]int)
	for _, kind := range []domain.ParticipantKind{domain.FactionKind, domain.CommanderKind} {
		byIDByKind[kind] = make(map[int]*domain.SParticipant)
		byURLByKind[kind] = make(map[string]int)
	}
	return &SParticipantsStore{
		validator:   validator.New(),
		byIDByKind:  byIDByKind,
		byURLByKind: byURLByKind,
	}
}

// Find searches for the pointer to a participant by its kind and ID. If none is found, nil is
// returned
func (s *SParticipantsStore) Find(kind domain.ParticipantKind, id int) *domain.SParticipant {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	p, found := s.byIDByKind[kind][id]
	if found != true {
		return nil
	}
	return p
}

// FindByURL does the same as Find, but instead of searching by ID, it searches by the URL of a
// participant
func (s *SParticipantsStore) FindByURL(kind domain.ParticipantKind, url string) *domain.SParticipant {
	s.mutex.RLock()

	id, found := s.byURLByKind[kind][url]
	s.mutex.RUnlock()
	if found != true {
		return nil
	}
	return s.Find(kind, id)
}

// Save stores the given participant. It returns an error if the participant does not have an ID or
// a URL
func (s *SParticipantsStore) Save(p domain.SParticipant) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.validator.Struct(p); err != nil {
		return errors.Wrapf(err, "Saving participant with URL %s", p.URL)
	}
	s.byIDByKind[p.Kind][p.ID] = &p
	s.byURLByKind[p.Kind][p.URL] = p.ID
	return nil
}

// Export saves data stored to the specified file using its input io.ExporterFunc
func (s *SParticipantsStore) Export(fileName string, exporterFunc io.ExporterFunc) error {
	return exporterFunc(fileName, struct {
		FactionsByID   map[int]*domain.SParticipant
		CommandersByID map[int]*domain.SParticipant
	}{
		FactionsByID:   s.byIDByKind[domain.FactionKind],
		CommandersByID: s.byIDByKind[domain.CommanderKind],
	})
}
