package postgresql

import (
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/db/postgresql/schema"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/domain/commanders"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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

// FindOne finds the first commander in the database that matches the query
func (r *CommandersRepository) FindOne(query commanders.FindOneQuery) (commanders.Commander, error) {
	c := &schema.Commander{}
	if err := r.db.Where(query).First(c).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return commanders.Commander{}, domain.ErrNotFound
	} else if err != nil {
		return commanders.Commander{}, errors.Wrap(err, "Executing CommandersRepository.FindOne")
	}
	return deserializeCommander(c), nil
}

// FindMany does a paginated search of all commanders matching the given query
func (r *CommandersRepository) FindMany(query commanders.FindManyQuery, page int) ([]commanders.Commander, int, error) {
	var records int64
	result := &[]schema.Commander{}

	var db = r.db.Model(&schema.Commander{})
	if query.FactionID != uuid.Nil {
		var cIDs []uuid.UUID
		err := r.db.
			Model(&schema.BattleCommanderFaction{}).
			Where(schema.BattleCommanderFaction{FactionID: query.FactionID}).
			Pluck("commander_id", &cIDs).
			Error
		if err != nil {
			return []commanders.Commander{}, 0, err
		}
		db = db.Where("id IN ?", cIDs)
	}
	db = ts(db, "name", query.Name)
	db = ts(db, "summary", query.Summary)

	if err := db.Count(&records).Error; err != nil {
		return []commanders.Commander{}, 0, err
	}
	pages := int((records / perPage) + 1)

	if err := paginate(db.Order("name DESC"), page, perPage).Find(result).Error; err != nil {
		return []commanders.Commander{}, pages, err
	}

	return deserializeCommanders(result), pages, nil
}

// CreateOne creates a commander in the database. The operation returns the ID of the new commander
func (r *CommandersRepository) CreateOne(data commanders.CreationInput) (uuid.UUID, error) {
	if err := r.validator.Struct(data); err != nil {
		return uuid.Nil, errors.Wrap(err, "Validating commander creation input")
	}
	c := serializeCommander(commanders.Commander{
		WikiID:  data.WikiID,
		URL:     data.URL,
		Name:    data.Name,
		Summary: data.Summary,
	})
	if err := r.db.Create(c).Error; err != nil {
		return uuid.Nil, errors.Wrap(err, "Creating a commander")
	}
	return c.ID, nil
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
