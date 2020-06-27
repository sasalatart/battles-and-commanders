package store

import (
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/sasalatart/batcoms/scraper/domain"
	"github.com/sasalatart/batcoms/scraper/exports"
)

// ParticipantsMem is an in-memory implementation of ParticipantsStore
type ParticipantsMem struct {
	mutex       sync.Mutex
	validator   *validator.Validate
	byIDByKind  map[domain.ParticipantKind]map[int]*domain.Participant
	byURLByKind map[domain.ParticipantKind]map[string]int
}

// NewParticipantsMem returns a pointer to an empty, ready-to-use ParticipantsMem store
func NewParticipantsMem() *ParticipantsMem {
	byIDByKind := make(map[domain.ParticipantKind]map[int]*domain.Participant)
	byURLByKind := make(map[domain.ParticipantKind]map[string]int)
	for _, kind := range []domain.ParticipantKind{domain.FactionKind, domain.CommanderKind} {
		byIDByKind[kind] = make(map[int]*domain.Participant)
		byURLByKind[kind] = make(map[string]int)
	}
	return &ParticipantsMem{
		validator:   validator.New(),
		byIDByKind:  byIDByKind,
		byURLByKind: byURLByKind,
	}
}

// Find searches for the pointer to a participant by its kind and ID. If none is found, nil is
// returned
func (r *ParticipantsMem) Find(kind domain.ParticipantKind, id int) *domain.Participant {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	p, found := r.byIDByKind[kind][id]
	if found != true {
		return nil
	}
	return p
}

// FindByURL does the same as Find, but instead of searching by ID, it searches by the URL of a
// participant
func (r *ParticipantsMem) FindByURL(kind domain.ParticipantKind, url string) *domain.Participant {
	r.mutex.Lock()

	id, found := r.byURLByKind[kind][url]
	r.mutex.Unlock()
	if found != true {
		return nil
	}
	return r.Find(kind, id)
}

// Save stores the given participant. It returns an error if the participant does not have an ID or
// a URL
func (r *ParticipantsMem) Save(p domain.Participant) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if err := r.validator.Struct(p); err != nil {
		return fmt.Errorf("Failed saving participant with URL %s: %s", p.URL, err)
	}
	r.byIDByKind[p.Kind][p.ID] = &p
	r.byURLByKind[p.Kind][p.URL] = p.ID
	return nil
}

// Export saves data stored to the specified file, in JSON format
func (r *ParticipantsMem) Export(fileName string) error {
	return exports.JSON(fileName, struct {
		FactionsByID   map[int]*domain.Participant
		CommandersByID map[int]*domain.Participant
	}{
		FactionsByID:   r.byIDByKind[domain.FactionKind],
		CommandersByID: r.byIDByKind[domain.CommanderKind],
	})
}
