package postgresql

import (
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/db/postgresql/schema"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/domain/factions"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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

// FindOne finds the first faction in the database that matches the query
func (r *FactionsRepository) FindOne(query factions.FindOneQuery) (factions.Faction, error) {
	f := &schema.Faction{}
	if err := r.db.Where(query).First(f).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return factions.Faction{}, domain.ErrNotFound
	} else if err != nil {
		return factions.Faction{}, errors.Wrap(err, "Executing FactionsRepository.FindOne")
	}
	return deserializeFaction(f), nil
}

// FindMany does a paginated search of all factions matching the given query
func (r *FactionsRepository) FindMany(query factions.FindManyQuery, page int) ([]factions.Faction, int, error) {
	var records int64
	result := &[]schema.Faction{}

	var db = r.db.Model(&schema.Faction{})
	if query.CommanderID != uuid.Nil {
		var fIDs []uuid.UUID
		err := r.db.
			Model(&schema.BattleCommanderFaction{}).
			Where(schema.BattleCommanderFaction{CommanderID: query.CommanderID}).
			Pluck("faction_id", &fIDs).
			Error
		if err != nil {
			return []factions.Faction{}, 0, err
		}
		db = db.Where("id IN ?", fIDs)
	}
	db = ts(db, "name", query.Name)
	db = ts(db, "summary", query.Summary)

	if err := db.Count(&records).Error; err != nil {
		return []factions.Faction{}, 0, err
	}
	pages := int((records / perPage) + 1)

	if err := paginate(db.Order("name DESC"), page, perPage).Find(result).Error; err != nil {
		return []factions.Faction{}, pages, err
	}

	return deserializeFactions(result), pages, nil
}

// CreateOne creates a faction in the database. The operation returns the ID of the new faction
func (r *FactionsRepository) CreateOne(data factions.CreationInput) (uuid.UUID, error) {
	if err := r.validator.Struct(data); err != nil {
		return uuid.Nil, errors.Wrap(err, "Validating faction creation input")
	}
	f := serializeFaction(factions.Faction{
		WikiID:  data.WikiID,
		URL:     data.URL,
		Name:    data.Name,
		Summary: data.Summary,
	})
	if err := r.db.Create(f).Error; err != nil {
		return uuid.Nil, errors.Wrap(err, "Creating a faction")
	}
	return f.ID, nil
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

func deserializeFactions(ff *[]schema.Faction) []factions.Faction {
	results := []factions.Faction{}
	for _, f := range *ff {
		results = append(results, deserializeFaction(&f))
	}
	return results
}
