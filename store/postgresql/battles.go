package postgresql

import (
	"encoding/json"

	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/store/postgresql/schema"
	uuid "github.com/satori/go.uuid"
)

// BattlesDataStore is the repository that abstracts access to the underlying database operations
// used to query and mutate data relating to battles. This implementation relies on GORM and also
// executes validations before interacting with the database
type BattlesDataStore struct {
	db        *gorm.DB
	validator *validator.Validate
}

// NewBattlesDataStore returns a pointer to a ready-to-use postgresql.BattlesDataStore
func NewBattlesDataStore(db *gorm.DB) *BattlesDataStore {
	return &BattlesDataStore{db, validator.New()}
}

func serializeBattle(b domain.Battle) (*schema.Battle, error) {
	strength, err := json.Marshal(b.Strength)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to stringify strength for %s", b.Name)
	}
	casualties, err := json.Marshal(b.Casualties)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to stringify casualties for %s", b.Name)
	}
	res := &schema.Battle{
		WikiID:             b.WikiID,
		URL:                b.URL,
		Name:               b.Name,
		PartOf:             b.PartOf,
		Summary:            b.Summary,
		StartDate:          b.StartDate,
		EndDate:            b.EndDate,
		Place:              b.Location.Place,
		Latitude:           b.Location.Latitude,
		Longitude:          b.Location.Longitude,
		Result:             b.Result,
		TerritorialChanges: b.TerritorialChanges,
		Strength:           postgres.Jsonb{RawMessage: strength},
		Casualties:         postgres.Jsonb{RawMessage: casualties},
	}
	return res, nil
}

func deserializeBattle(b *schema.Battle) (domain.Battle, error) {
	if b == nil {
		return domain.Battle{}, errors.New("Empty battle to deserialize")
	}
	strength := domain.SideNumbers{}
	if err := fromJSONB(b.Strength, &strength); err != nil {
		return domain.Battle{}, errors.Wrapf(err, "Deserializing strength")
	}
	casualties := domain.SideNumbers{}
	if err := fromJSONB(b.Casualties, &casualties); err != nil {
		return domain.Battle{}, errors.Wrapf(err, "Deserializing casualties")
	}
	commanders := domain.CommandersBySide{}
	for _, bc := range b.BattleCommanders {
		commander := deserializeCommander(bc.Commander)
		if bc.Side == schema.SideA {
			commanders.A = append(commanders.A, commander)
		} else {
			commanders.B = append(commanders.B, commander)
		}
	}
	factions := domain.FactionsBySide{}
	for _, bf := range b.BattleFactions {
		faction := deserializeFaction(bf.Faction)
		if bf.Side == schema.SideA {
			factions.A = append(factions.A, faction)
		} else {
			factions.B = append(factions.B, faction)
		}
	}
	commandersByFaction := make(domain.CommandersByFaction)
	for _, bcf := range b.BattleCommanderFactions {
		commandersByFaction[bcf.FactionID] = append(commandersByFaction[bcf.FactionID], bcf.CommanderID)
	}
	res := domain.Battle{
		ID:        b.ID,
		WikiID:    b.WikiID,
		URL:       b.URL,
		Name:      b.Name,
		PartOf:    b.PartOf,
		Summary:   b.Summary,
		StartDate: b.StartDate,
		EndDate:   b.EndDate,
		Location: domain.Location{
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

// CreateOne creates a battle in the database, together with entries in the corresponding tables
// that let us relate the battle with other participants. The operation returns the ID of the new
// battle
func (s *BattlesDataStore) CreateOne(data domain.CreateBattleInput) (uuid.UUID, error) {
	if err := s.validator.Struct(data); err != nil {
		return uuid.UUID{}, errors.Wrap(err, "Validating battle creation input")
	}
	b, err := serializeBattle(domain.Battle{
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
		return uuid.UUID{}, errors.Wrap(err, "Serializing domain.CreateBattleInput")
	}
	addBattleCommanders := func(cIDs []uuid.UUID, side schema.SideKind) {
		for _, cID := range cIDs {
			b.BattleCommanders = append(b.BattleCommanders, schema.BattleCommander{CommanderID: cID, Side: side})
		}
	}
	addBattleCommanders(data.CommandersBySide.A, schema.SideA)
	addBattleCommanders(data.CommandersBySide.B, schema.SideB)
	addBattleFactions := func(fIDs []uuid.UUID, side schema.SideKind) {
		for _, fID := range fIDs {
			b.BattleFactions = append(b.BattleFactions, schema.BattleFaction{FactionID: fID, Side: side})
		}
	}
	addBattleFactions(data.FactionsBySide.A, schema.SideA)
	addBattleFactions(data.FactionsBySide.B, schema.SideB)
	for fID, cIDS := range data.CommandersByFaction {
		for _, cID := range cIDS {
			b.BattleCommanderFactions = append(b.BattleCommanderFactions, schema.BattleCommanderFaction{
				FactionID:   fID,
				CommanderID: cID,
			})
		}
	}
	if err := s.db.Create(b).Error; err != nil {
		return uuid.UUID{}, errors.Wrap(err, "Creating the battle")
	}
	return b.ID, nil
}
