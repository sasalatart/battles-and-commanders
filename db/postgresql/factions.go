package postgresql

import (
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/db/postgresql/schema"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/domain/factions"
	uuid "github.com/satori/go.uuid"
)

// FactionsRepository is the repository that abstracts access to the underlying database operations
// used to query and mutate data relating to factions. This implementation relies on GORM and also
// executes validations before interacting with the database
type FactionsRepository struct {
	db        *gorm.DB
	validator *validator.Validate
}

// NewFactionsRepository returns a pointer to a ready-to-use postgresql.FactionsRepository
func NewFactionsRepository(db *gorm.DB) *FactionsRepository {
	return &FactionsRepository{db, validator.New()}
}

func serializeFaction(f factions.Faction) *schema.Faction {
	return &schema.Faction{
		WikiID:  f.WikiID,
		URL:     f.URL,
		Name:    f.Name,
		Summary: f.Summary,
	}
}

func deserializeFaction(f *schema.Faction) factions.Faction {
	return factions.Faction{
		ID:      f.ID,
		WikiID:  f.WikiID,
		URL:     f.URL,
		Name:    f.Name,
		Summary: f.Summary,
	}
}

// FindOne finds the first faction in the database that matches the query
func (r *FactionsRepository) FindOne(query interface{}, args ...interface{}) (factions.Faction, error) {
	f := &schema.Faction{}
	if err := r.db.Where(query, args...).Find(f).Error; gorm.IsRecordNotFoundError(err) {
		return factions.Faction{}, domain.ErrNotFound
	} else if err != nil {
		return factions.Faction{}, errors.Wrap(err, "Executing FactionsRepository.FindOne")
	}
	return deserializeFaction(f), nil
}

// CreateOne creates a faction in the database. The operation returns the ID of the new faction
func (r *FactionsRepository) CreateOne(data factions.CreationInput) (uuid.UUID, error) {
	if err := r.validator.Struct(data); err != nil {
		return uuid.UUID{}, errors.Wrap(err, "Validating faction creation input")
	}
	f := serializeFaction(factions.Faction{
		WikiID:  data.WikiID,
		URL:     data.URL,
		Name:    data.Name,
		Summary: data.Summary,
	})
	if err := r.db.Create(f).Error; err != nil {
		return uuid.UUID{}, errors.Wrap(err, "Creating a faction")
	}
	return f.ID, nil
}
