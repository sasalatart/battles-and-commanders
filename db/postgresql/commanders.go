package postgresql

import (
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/db/postgresql/schema"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/domain/commanders"
	uuid "github.com/satori/go.uuid"
)

// CommandersRepository is the repository that abstracts access to the underlying database operations
// used to query and mutate data relating to commanders. This implementation relies on GORM and also
// executes validations before interacting with the database
type CommandersRepository struct {
	db        *gorm.DB
	validator *validator.Validate
}

// NewCommandersRepository returns a pointer to a ready-to-use postgresql.CommandersRepository
func NewCommandersRepository(db *gorm.DB) *CommandersRepository {
	return &CommandersRepository{db, validator.New()}
}

func serializeCommander(c commanders.Commander) *schema.Commander {
	return &schema.Commander{
		WikiID:  c.WikiID,
		URL:     c.URL,
		Name:    c.Name,
		Summary: c.Summary,
	}
}

func deserializeCommander(c *schema.Commander) commanders.Commander {
	return commanders.Commander{
		ID:      c.ID,
		WikiID:  c.WikiID,
		URL:     c.URL,
		Name:    c.Name,
		Summary: c.Summary,
	}
}

func deserializeCommanders(cc *[]schema.Commander) []commanders.Commander {
	results := []commanders.Commander{}
	for _, c := range *cc {
		results = append(results, deserializeCommander(&c))
	}
	return results
}

// FindOne finds the first commander in the database that matches the query
func (s *CommandersRepository) FindOne(query commanders.Commander) (commanders.Commander, error) {
	c := &schema.Commander{}
	if err := s.db.Where(query).Find(c).Error; gorm.IsRecordNotFoundError(err) {
		return commanders.Commander{}, domain.ErrNotFound
	} else if err != nil {
		return commanders.Commander{}, errors.Wrap(err, "Executing CommandersRepository.FindOne")
	}
	return deserializeCommander(c), nil
}

// FindMany does a paginated search of all commanders matching the given query
func (s *CommandersRepository) FindMany(query commanders.Query, page uint) ([]commanders.Commander, uint, error) {
	var records uint
	result := &[]schema.Commander{}

	var db = s.db.Model(&schema.Commander{}).Order("name DESC")
	if query.FactionID != uuid.Nil {
		db = db.Joins("JOIN battle_commander_factions bcf ON bcf.commander_id = commanders.id").
			Where("bcf.faction_id = ?", query.FactionID)
	}
	if query.Name != "" {
		db = db.Where("to_tsvector('english', name) @@ plainto_tsquery(?)", query.Name)
	}
	if query.Summary != "" {
		db = db.Where("to_tsvector('english', summary) @@ phraseto_tsquery(?)", query.Summary)
	}

	if err := db.Count(&records).Error; err != nil {
		return []commanders.Commander{}, records, err
	}

	db = db.Offset((page - 1) * perPage).Limit(perPage)
	if err := db.Find(result).Error; err != nil {
		return []commanders.Commander{}, records, err
	}

	return deserializeCommanders(result), (records / perPage) + 1, nil
}

// CreateOne creates a commander in the database. The operation returns the ID of the new commander
func (s *CommandersRepository) CreateOne(data commanders.CreationInput) (uuid.UUID, error) {
	if err := s.validator.Struct(data); err != nil {
		return uuid.UUID{}, errors.Wrap(err, "Validating commander creation input")
	}
	c := serializeCommander(commanders.Commander{
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
