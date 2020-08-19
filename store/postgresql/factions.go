package postgresql

import (
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/store/postgresql/schema"
	uuid "github.com/satori/go.uuid"
)

// FactionsDataStore is the repository that abstracts access to the underlying database operations
// used to query and mutate data relating to factions. This implementation relies on GORM and also
// executes validations before interacting with the database
type FactionsDataStore struct {
	db        *gorm.DB
	validator *validator.Validate
}

// NewFactionsDataStore returns a pointer to ready-to-use postgresql.FactionsDataStore
func NewFactionsDataStore(db *gorm.DB) *FactionsDataStore {
	return &FactionsDataStore{db, validator.New()}
}

func serializeFaction(f domain.Faction) *schema.Faction {
	return &schema.Faction{
		WikiID:  f.WikiID,
		URL:     f.URL,
		Name:    f.Name,
		Summary: f.Summary,
	}
}

func deserializeFaction(f *schema.Faction) domain.Faction {
	return domain.Faction{
		ID:      f.ID,
		WikiID:  f.WikiID,
		URL:     f.URL,
		Name:    f.Name,
		Summary: f.Summary,
	}
}

// CreateOne creates a faction in the database. The operation returns the UUID of the new faction
func (s *FactionsDataStore) CreateOne(data domain.CreateFactionInput) (uuid.UUID, error) {
	if err := s.validator.Struct(data); err != nil {
		return uuid.UUID{}, errors.Wrap(err, "Validating faction creation input")
	}
	f := serializeFaction(domain.Faction{
		WikiID:  data.WikiID,
		URL:     data.URL,
		Name:    data.Name,
		Summary: data.Summary,
	})
	if err := s.db.Create(f).Error; err != nil {
		return uuid.UUID{}, errors.Wrap(err, "Creating a faction")
	}
	return f.ID, nil
}
