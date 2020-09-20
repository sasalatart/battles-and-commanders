package memory

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain/wikibattles"
	"github.com/sasalatart/batcoms/pkg/io"
)

// WikiBattlesRepo is an in-memory implementation of wikibattles.Repository
type WikiBattlesRepo struct {
	mutex     sync.RWMutex
	validator *validator.Validate
	byID      map[int]*wikibattles.Battle
}

// NewWikiBattlesRepo returns a pointer to an empty, ready-to-use memory.WikiBattlesRepo
func NewWikiBattlesRepo() *WikiBattlesRepo {
	return &WikiBattlesRepo{validator: validator.New(), byID: make(map[int]*wikibattles.Battle)}
}

// Find searches for the pointer to a battle by its ID. If none is found, nil is returned
func (r *WikiBattlesRepo) Find(id int) *wikibattles.Battle {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	battle, found := r.byID[id]
	if found != true {
		return nil
	}
	return battle
}

// Save stores the given battle. It returns an error if validations on the struct fail
func (r *WikiBattlesRepo) Save(b wikibattles.Battle) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if err := r.validator.Struct(b); err != nil {
		return errors.Wrapf(err, "Saving battle with URL %s", b.URL)
	}
	r.byID[b.ID] = &b
	return nil
}

// Export saves data stored to the specified file using its input io.ExporterFunc
func (r *WikiBattlesRepo) Export(fileName string, exporterFunc io.ExporterFunc) error {
	return exporterFunc(fileName, r.byID)
}
