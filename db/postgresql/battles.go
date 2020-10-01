package postgresql

import (
	"encoding/json"

	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/db/postgresql/schema"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/domain/battles"
	"github.com/sasalatart/batcoms/domain/locations"
	"github.com/sasalatart/batcoms/domain/statistics"
	"github.com/sasalatart/batcoms/pkg/dates"
	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// BattlesRepository is the repository that abstracts access to the underlying database operations
// used to query and mutate data relating to battles. This implementation relies on GORM and also
// executes validations before interacting with the database
type BattlesRepository struct {
	db        *gorm.DB
	validator *validator.Validate
}

// NewBattlesRepository returns a pointer to a ready-to-use postgresql.BattlesRepository
func NewBattlesRepository(db *gorm.DB) *BattlesRepository {
	return &BattlesRepository{db, validator.New()}
}

// FindOne finds the first battle in the database that matches the query, together with its related
// factions and commanders
func (r *BattlesRepository) FindOne(query battles.FindOneQuery) (battles.Battle, error) {
	b := new(schema.Battle)
	db := r.db.
		Preload("BattleFactions.Faction").
		Preload("BattleCommanders.Commander").
		Preload("BattleCommanderFactions")
	if err := db.Where(query).First(b).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return battles.Battle{}, domain.ErrNotFound
	} else if err != nil {
		return battles.Battle{}, errors.Wrap(err, "Finding a battle")
	}
	return deserializeBattle(b)
}

// FindMany does a paginated search of all battles matching the given query
func (r *BattlesRepository) FindMany(query battles.FindManyQuery, page int) ([]battles.Battle, int, error) {
	var records int64
	result := &[]schema.Battle{}

	var db = r.db.
		Model(&schema.Battle{}).
		Preload("BattleFactions.Faction").
		Preload("BattleCommanders.Commander").
		Preload("BattleCommanderFactions")

	if query.FactionID != uuid.Nil {
		db = db.Joins("JOIN battle_factions bf ON bf.battle_id = battles.id").
			Where("bf.faction_id = ?", query.FactionID)
	}
	if query.CommanderID != uuid.Nil {
		db = db.Joins("JOIN battle_commanders bc ON bc.battle_id = battles.id").
			Where("bc.commander_id = ?", query.CommanderID)
	}
	db = ts(db, "name", query.Name)
	db = ts(db, "summary", query.Summary)
	db = ts(db, "place", query.Place)
	db = ts(db, "result", query.Result)

	if err := db.Count(&records).Error; err != nil {
		return []battles.Battle{}, 0, err
	}
	pages := int((records / perPage) + 1)

	if err := paginate(db.Order("name DESC"), page, perPage).Find(result).Error; err != nil {
		return []battles.Battle{}, pages, err
	}

	battles, err := deserializeBattles(result)
	return battles, pages, err
}

// CreateOne creates a battle in the database, together with entries in the corresponding tables
// that let us relate the battle with other factions and commanders. The operation returns the ID of
// the new battle
func (r *BattlesRepository) CreateOne(data battles.CreationInput) (uuid.UUID, error) {
	if err := r.validator.Struct(data); err != nil {
		return uuid.Nil, errors.Wrap(err, "Validating battle creation input")
	}
	b, err := serializeBattle(battles.Battle{
		WikiID:             data.WikiID,
		URL:                data.URL,
		Name:               data.Name,
		PartOf:             data.PartOf,
		Summary:            data.Summary,
		StartDate:          data.StartDate,
		EndDate:            data.EndDate,
		Location:           data.Location,
		Result:             data.Result,
		TerritorialChanges: data.TerritorialChanges,
		Strength:           data.Strength,
		Casualties:         data.Casualties,
	})
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "Serializing battles.CreationInput")
	}
	addBattleFactions := func(fIDs []uuid.UUID, side schema.SideKind) {
		for _, fID := range fIDs {
			b.BattleFactions = append(b.BattleFactions, schema.BattleFaction{FactionID: fID, Side: side})
		}
	}
	addBattleFactions(data.FactionsBySide.A, schema.SideA)
	addBattleFactions(data.FactionsBySide.B, schema.SideB)
	addBattleCommanders := func(cIDs []uuid.UUID, side schema.SideKind) {
		for _, cID := range cIDs {
			b.BattleCommanders = append(b.BattleCommanders, schema.BattleCommander{CommanderID: cID, Side: side})
		}
	}
	addBattleCommanders(data.CommandersBySide.A, schema.SideA)
	addBattleCommanders(data.CommandersBySide.B, schema.SideB)
	for fID, cIDS := range data.CommandersByFaction {
		for _, cID := range cIDS {
			b.BattleCommanderFactions = append(b.BattleCommanderFactions, schema.BattleCommanderFaction{
				FactionID:   fID,
				CommanderID: cID,
			})
		}
	}
	if err := r.db.Create(b).Error; err != nil {
		return uuid.Nil, errors.Wrap(err, "Creating the battle")
	}
	return b.ID, nil
}

func serializeBattle(b battles.Battle) (*schema.Battle, error) {
	strength, err := json.Marshal(b.Strength)
	if err != nil {
		return nil, errors.Wrap(err, "Stringifying strength")
	}
	casualties, err := json.Marshal(b.Casualties)
	if err != nil {
		return nil, errors.Wrap(err, "Stringifying casualties")
	}
	startDateNum, err := dates.ToNum(b.StartDate)
	if err != nil {
		return nil, errors.Wrap(err, "Converting StartDate to number")
	}
	endDateNum, err := dates.ToNum(b.EndDate)
	if err != nil {
		return nil, errors.Wrap(err, "Converting EndDate to number")
	}
	res := &schema.Battle{
		WikiID:             b.WikiID,
		URL:                b.URL,
		Name:               b.Name,
		PartOf:             b.PartOf,
		Summary:            b.Summary,
		StartDate:          b.StartDate,
		StartDateNum:       startDateNum,
		EndDate:            b.EndDate,
		EndDateNum:         endDateNum,
		Place:              b.Location.Place,
		Latitude:           b.Location.Latitude,
		Longitude:          b.Location.Longitude,
		Result:             b.Result,
		TerritorialChanges: b.TerritorialChanges,
		Strength:           datatypes.JSON(strength),
		Casualties:         datatypes.JSON(casualties),
	}
	return res, nil
}

func deserializeBattle(b *schema.Battle) (battles.Battle, error) {
	if b == nil {
		return battles.Battle{}, errors.New("Empty battle to deserialize")
	}
	strength := statistics.SideNumbers{}
	if err := fromJSON(b.Strength, &strength); err != nil {
		return battles.Battle{}, errors.Wrapf(err, "Deserializing strength")
	}
	casualties := statistics.SideNumbers{}
	if err := fromJSON(b.Casualties, &casualties); err != nil {
		return battles.Battle{}, errors.Wrapf(err, "Deserializing casualties")
	}
	commanders := battles.CommandersBySide{}
	for _, bc := range b.BattleCommanders {
		commander := deserializeCommander(&bc.Commander)
		if bc.Side == schema.SideA {
			commanders.A = append(commanders.A, commander)
		} else {
			commanders.B = append(commanders.B, commander)
		}
	}
	factions := battles.FactionsBySide{}
	for _, bf := range b.BattleFactions {
		faction := deserializeFaction(&bf.Faction)
		if bf.Side == schema.SideA {
			factions.A = append(factions.A, faction)
		} else {
			factions.B = append(factions.B, faction)
		}
	}
	commandersByFaction := make(battles.CommandersByFaction)
	for _, bcf := range b.BattleCommanderFactions {
		commandersByFaction[bcf.FactionID] = append(commandersByFaction[bcf.FactionID], bcf.CommanderID)
	}
	res := battles.Battle{
		ID:        b.ID,
		WikiID:    b.WikiID,
		URL:       b.URL,
		Name:      b.Name,
		PartOf:    b.PartOf,
		Summary:   b.Summary,
		StartDate: b.StartDate,
		EndDate:   b.EndDate,
		Location: locations.Location{
			Place:     b.Place,
			Latitude:  b.Latitude,
			Longitude: b.Longitude,
		},
		Result:              b.Result,
		TerritorialChanges:  b.TerritorialChanges,
		Strength:            strength,
		Casualties:          casualties,
		Factions:            factions,
		Commanders:          commanders,
		CommandersByFaction: commandersByFaction,
	}
	return res, nil
}

func deserializeBattles(bb *[]schema.Battle) ([]battles.Battle, error) {
	results := []battles.Battle{}
	for _, b := range *bb {
		battle, err := deserializeBattle(&b)
		if err != nil {
			return results, err
		}
		results = append(results, battle)
	}
	return results, nil
}
