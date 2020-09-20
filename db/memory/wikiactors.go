package memory

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain/wikiactors"
	"github.com/sasalatart/batcoms/pkg/io"
)

// WikiActorsRepo is an in-memory implementation of wikiactors.Repository
type WikiActorsRepo struct {
	mutex       sync.RWMutex
	validator   *validator.Validate
	byIDByKind  map[wikiactors.Kind]map[int]*wikiactors.Actor
	byURLByKind map[wikiactors.Kind]map[string]int
}

// NewWikiActorsRepo returns a pointer to an empty, ready-to-use memory.WikiActorsRepo
func NewWikiActorsRepo() *WikiActorsRepo {
	byIDByKind := make(map[wikiactors.Kind]map[int]*wikiactors.Actor)
	byURLByKind := make(map[wikiactors.Kind]map[string]int)
	for _, kind := range []wikiactors.Kind{wikiactors.FactionKind, wikiactors.CommanderKind} {
		byIDByKind[kind] = make(map[int]*wikiactors.Actor)
		byURLByKind[kind] = make(map[string]int)
	}
	return &WikiActorsRepo{
		validator:   validator.New(),
		byIDByKind:  byIDByKind,
		byURLByKind: byURLByKind,
	}
}

// Find searches for the pointer to an actor by its kind and ID. If none is found, nil is returned
func (r *WikiActorsRepo) Find(kind wikiactors.Kind, id int) *wikiactors.Actor {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	p, found := r.byIDByKind[kind][id]
	if found != true {
		return nil
	}
	return p
}

// FindByURL does the same as Find, but instead of searching by ID, it searches by the URL of an actor
func (r *WikiActorsRepo) FindByURL(kind wikiactors.Kind, url string) *wikiactors.Actor {
	r.mutex.RLock()

	id, found := r.byURLByKind[kind][url]
	r.mutex.RUnlock()
	if found != true {
		return nil
	}
	return r.Find(kind, id)
}

// Save stores the given actor. It returns an error if validations on the struct fail
func (r *WikiActorsRepo) Save(p wikiactors.Actor) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if err := r.validator.Struct(p); err != nil {
		return errors.Wrap(err, "Saving actor")
	}
	r.byIDByKind[p.Kind][p.ID] = &p
	r.byURLByKind[p.Kind][p.URL] = p.ID
	return nil
}

// Export saves data stored to the specified file using its input io.ExporterFunc
func (r *WikiActorsRepo) Export(fileName string, exporterFunc io.ExporterFunc) error {
	return exporterFunc(fileName, struct {
		FactionsByID   map[int]*wikiactors.Actor
		CommandersByID map[int]*wikiactors.Actor
	}{
		FactionsByID:   r.byIDByKind[wikiactors.FactionKind],
		CommandersByID: r.byIDByKind[wikiactors.CommanderKind],
	})
}
