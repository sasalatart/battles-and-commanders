package store

import (
	"fmt"
	"sync"

	"github.com/sasalatart/batcoms/scraper/domain"
	"github.com/sasalatart/batcoms/scraper/exports"
)

// ParticipantsMem is an in-memory implementation of ParticipantsStore
type ParticipantsMem struct {
	mutex          sync.Mutex
	byIDByKind     map[domain.ParticipantKind]map[int]*domain.Participant
	byURLByKind    map[domain.ParticipantKind]map[string]int
	factionsByFlag map[string]int
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
		byIDByKind:     byIDByKind,
		byURLByKind:    byURLByKind,
		factionsByFlag: make(map[string]int),
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

// FindFactionByFlag does the same as Find, but instead of searching by ID, it searches by the flag
// URL of the faction
func (r *ParticipantsMem) FindFactionByFlag(flag string) *domain.Participant {
	r.mutex.Lock()

	id, found := r.factionsByFlag[flag]
	r.mutex.Unlock()
	if found != true {
		return nil
	}
	return r.Find(domain.FactionKind, id)
}

// Save stores the given participant. It returns an error if the participant does not have an ID or
// a URL
func (r *ParticipantsMem) Save(p domain.Participant) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if p.ID == 0 {
		return fmt.Errorf("Expected %+v to have an ID, but none was present", p)
	}
	if p.URL == "" {
		return fmt.Errorf("Expected %+v to have an URL, but none was present", p)
	}
	r.byIDByKind[p.Kind][p.ID] = &p
	r.byURLByKind[p.Kind][p.URL] = p.ID
	if p.Kind == domain.FactionKind && p.Flag != "" {
		r.factionsByFlag[p.Flag] = p.ID
	}
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