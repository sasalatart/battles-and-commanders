package postgresql

import (
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/store/postgresql/schema"
	uuid "github.com/satori/go.uuid"
)

// CommandersDataStore is the repository that abstracts access to the underlying database operations
// used to query and mutate data relating to commanders. This implementation relies on GORM and also
// executes validations before interacting with the database
type CommandersDataStore struct {
	db        *gorm.DB
	validator *validator.Validate
}

// NewCommandersDataStore returns a pointer to ready-to-use postgresql.CommandersDataStore
func NewCommandersDataStore(db *gorm.DB) *CommandersDataStore {
	return &CommandersDataStore{db, validator.New()}
}

func serializeCommander(c domain.Commander) *schema.Commander {
	return &schema.Commander{
		WikiID:  c.WikiID,
		URL:     c.URL,
		Name:    c.Name,
		Summary: c.Summary,
	}
}

func deserializeCommander(c *schema.Commander) domain.Commander {
	return domain.Commander{
		ID:      c.ID,
		WikiID:  c.WikiID,
		URL:     c.URL,
		Name:    c.Name,
		Summary: c.Summary,
	}
}

// CreateOne creates a commander in the database. The operation returns the UUID of the new commander
func (s *CommandersDataStore) CreateOne(data domain.CreateCommanderInput) (uuid.UUID, error) {
	if err := s.validator.Struct(data); err != nil {
		return uuid.UUID{}, errors.Wrapf(err, "Validating commander creation input with URL %s", data.URL)
	}
	c := serializeCommander(domain.Commander{
		WikiID:  data.WikiID,
		URL:     data.URL,
		Name:    data.Name,
		Summary: data.Summary,
	})
	if err := s.db.Create(c).Error; err != nil {
		return uuid.UUID{}, errors.Wrap(err, "Creating a commander")
	}
	return c.ID, nil
}
