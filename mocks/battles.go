package mocks

import (
	"log"

	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/services/parser"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

var battleUUID = uuid.NewV4()

// BattlesDataStore mocks datastores used to find & create battles
type BattlesDataStore struct {
	mock.Mock
}

// FindOne mocks finding one battle via BattlesDataStore
func (bs *BattlesDataStore) FindOne(query interface{}, args ...interface{}) (domain.Battle, error) {
	mockArgs := bs.Called(append([]interface{}{query}, args...)...)
	return mockArgs.Get(0).(domain.Battle), mockArgs.Error(1)
}

// CreateOne mocks creating one battle via BattlesDataStore
func (bs *BattlesDataStore) CreateOne(data domain.CreateBattleInput) (uuid.UUID, error) {
	mockArgs := bs.Called(data)
	return mockArgs.Get(0).(uuid.UUID), mockArgs.Error(1)
}

// Battle returns an instance of domain.Battle that may be used for mocking purposes
func Battle() domain.Battle {
	sb := SBattle()
	dates, err := parser.Date(sb.Date)
	if err != nil {
		log.Fatalf("Cannot parse date %q, and therefore mock.Battle is not valid: %s", sb.Date, err)
	}
	return domain.Battle{
		ID:        battleUUID,
		WikiID:    sb.ID,
		URL:       sb.URL,
		Name:      sb.Name,
		PartOf:    sb.PartOf,
		Summary:   sb.Extract,
		StartDate: dates[0],
		EndDate:   dates[len(dates)-1],
		Location: domain.Location{
			Place:     sb.Location.Place,
			Latitude:  sb.Location.Latitude,
			Longitude: sb.Location.Longitude,
		},
		Result:             sb.Result,
		TerritorialChanges: sb.TerritorialChanges,
		Strength: domain.SideNumbers{
			A:  sb.Strength.A,
			B:  sb.Strength.B,
			AB: sb.Strength.AB,
		},
		Casualties: domain.SideNumbers{
			A:  sb.Casualties.A,
			B:  sb.Casualties.B,
			AB: sb.Casualties.AB,
		},
		Factions: domain.FactionsBySide{
			A: []domain.Faction{Faction()},
			B: []domain.Faction{Faction2(), Faction3()},
		},
		Commanders: domain.CommandersBySide{
			A: []domain.Commander{Commander()},
			B: []domain.Commander{Commander2(), Commander3(), Commander4(), Commander5()},
		},
		CommandersByFaction: domain.CommandersByFaction{
			Faction().ID:  []uuid.UUID{Commander().ID},
			Faction2().ID: []uuid.UUID{Commander2().ID, Commander3().ID},
			Faction3().ID: []uuid.UUID{Commander4().ID, Commander5().ID},
		},
	}
}

// CreateBattleInput returns an instance of domain.CreateBattleInput that may be used for mocking
// inputs to create battles
func CreateBattleInput() domain.CreateBattleInput {
	b := Battle()
	return domain.CreateBattleInput{
		WikiID:    b.WikiID,
		URL:       b.URL,
		Name:      b.Name,
		PartOf:    b.PartOf,
		Summary:   b.Summary,
		StartDate: b.StartDate,
		EndDate:   b.EndDate,
		Location: domain.Location{
			Place:     b.Location.Place,
			Latitude:  b.Location.Latitude,
			Longitude: b.Location.Longitude,
		},
		Result:             b.Result,
		TerritorialChanges: b.TerritorialChanges,
		Strength: domain.SideNumbers{
			A:  b.Strength.A,
			B:  b.Strength.B,
			AB: b.Strength.AB,
		},
		Casualties: domain.SideNumbers{
			A:  b.Casualties.A,
			B:  b.Casualties.B,
			AB: b.Casualties.AB,
		},
		FactionsBySide: domain.ParticipantsIDsBySide{
			A: []uuid.UUID{Faction().ID},
			B: []uuid.UUID{Faction2().ID, Faction3().ID},
		},
		CommandersBySide: domain.ParticipantsIDsBySide{
			A: []uuid.UUID{Commander().ID},
			B: []uuid.UUID{Commander2().ID, Commander3().ID, Commander4().ID, Commander5().ID},
		},
		CommandersByFaction: domain.CommandersByFaction{
			Faction().ID:  []uuid.UUID{Commander().ID},
			Faction2().ID: []uuid.UUID{Commander2().ID, Commander3().ID},
			Faction3().ID: []uuid.UUID{Commander4().ID, Commander5().ID},
		},
	}
}
