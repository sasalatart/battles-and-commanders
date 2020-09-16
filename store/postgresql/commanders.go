package postgresql

import (
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/store"
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

// NewCommandersDataStore returns a pointer to a ready-to-use postgresql.CommandersDataStore
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

func deserializeCommanders(cc *[]schema.Commander) []domain.Commander {
	results := []domain.Commander{}
	for _, c := range *cc {
		results = append(results, deserializeCommander(&c))
	}
	return results
}

// FindOne finds the first commander in the database that matches the query
func (s *CommandersDataStore) FindOne(query domain.Commander) (domain.Commander, error) {
	c := &schema.Commander{}
	if err := s.db.Where(query).Find(c).Error; gorm.IsRecordNotFoundError(err) {
		return domain.Commander{}, store.ErrNotFound
	} else if err != nil {
		return domain.Commander{}, errors.Wrap(err, "Executing CommandersDataStore.FindOne")
	}
	return deserializeCommander(c), nil
}

// FindMany does a paginated search of all commanders matching the given query
func (s *CommandersDataStore) FindMany(query store.CommandersQuery, page uint) ([]domain.Commander, uint, error) {
	var records uint
	commanders := &[]schema.Commander{}
	var db = s.db.Model(&schema.Commander{}).Order("name DESC")
	if query.FactionID != uuid.Nil {
		db = db.Joins("JOIN battle_commander_factions bcf ON bcf.commander_id = commanders.id").
			Where("bcf.faction_id = ?", query.FactionID)
	}

	if err := db.Count(&records).Error; err != nil {
		return []domain.Commander{}, records, err
	}

	db = db.Offset((page - 1) * perPage).Limit(perPage)
	if err := db.Find(commanders).Error; err != nil {
		return []domain.Commander{}, records, err
	}

	return deserializeCommanders(commanders), (records / perPage) + 1, nil
}

// CreateOne creates a commander in the database. The operation returns the ID of the new commander
func (s *CommandersDataStore) CreateOne(data domain.CreateCommanderInput) (uuid.UUID, error) {
	if err := s.validator.Struct(data); err != nil {
		return uuid.UUID{}, errors.Wrap(err, "Validating commander creation input")
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
