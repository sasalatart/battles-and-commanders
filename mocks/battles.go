package mocks

import (
	"log"

	"github.com/sasalatart/batcoms/domain/battles"
	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/domain/factions"
	"github.com/sasalatart/batcoms/domain/locations"
	"github.com/sasalatart/batcoms/domain/statistics"
	"github.com/sasalatart/batcoms/pkg/dates"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

var battleUUID = uuid.NewV4()

// BattlesRepository mocks datastores used to find & create battles
type BattlesRepository struct {
	mock.Mock
}

// FindOne mocks finding one battle via BattlesRepository
func (bs *BattlesRepository) FindOne(query interface{}, args ...interface{}) (battles.Battle, error) {
	mockArgs := bs.Called(append([]interface{}{query}, args...)...)
	return mockArgs.Get(0).(battles.Battle), mockArgs.Error(1)
}

// CreateOne mocks creating one battle via BattlesRepository
func (bs *BattlesRepository) CreateOne(data battles.CreationInput) (uuid.UUID, error) {
	mockArgs := bs.Called(data)
	return mockArgs.Get(0).(uuid.UUID), mockArgs.Error(1)
}

// Battle returns an instance of battles.Battle that may be used for mocking purposes
func Battle() battles.Battle {
	wb := WikiBattle()
	dates, err := dates.Parse(wb.Date)
	if err != nil {
		log.Fatalf("Cannot parse date %q, and therefore mock.Battle is not valid: %s", wb.Date, err)
	}
	return battles.Battle{
		ID:        battleUUID,
		WikiID:    wb.ID,
		URL:       wb.URL,
		Name:      wb.Name,
		PartOf:    wb.PartOf,
		Summary:   wb.Extract,
		StartDate: dates[0],
		EndDate:   dates[len(dates)-1],
		Location: locations.Location{
			Place:     wb.Location.Place,
			Latitude:  wb.Location.Latitude,
			Longitude: wb.Location.Longitude,
		},
		Result:             wb.Result,
		TerritorialChanges: wb.TerritorialChanges,
		Strength: statistics.SideNumbers{
			A:  wb.Strength.A,
			B:  wb.Strength.B,
			AB: wb.Strength.AB,
		},
		Casualties: statistics.SideNumbers{
			A:  wb.Casualties.A,
			B:  wb.Casualties.B,
			AB: wb.Casualties.AB,
		},
		Factions: battles.FactionsBySide{
			A: []factions.Faction{Faction()},
			B: []factions.Faction{Faction2(), Faction3()},
		},
		Commanders: battles.CommandersBySide{
			A: []commanders.Commander{Commander()},
			B: []commanders.Commander{Commander2(), Commander3(), Commander4(), Commander5()},
		},
		CommandersByFaction: battles.CommandersByFaction{
			Faction().ID:  []uuid.UUID{Commander().ID},
			Faction2().ID: []uuid.UUID{Commander2().ID, Commander3().ID},
			Faction3().ID: []uuid.UUID{Commander4().ID, Commander5().ID},
		},
	}
}

// BattleCreationInput returns an instance of battles.CreationInput that may be used for mocking
// inputs to create battles
func BattleCreationInput() battles.CreationInput {
	b := Battle()
	return battles.CreationInput{
		WikiID:    b.WikiID,
		URL:       b.URL,
		Name:      b.Name,
		PartOf:    b.PartOf,
		Summary:   b.Summary,
		StartDate: b.StartDate,
		EndDate:   b.EndDate,
		Location: locations.Location{
			Place:     b.Location.Place,
			Latitude:  b.Location.Latitude,
			Longitude: b.Location.Longitude,
		},
		Result:             b.Result,
		TerritorialChanges: b.TerritorialChanges,
		Strength: statistics.SideNumbers{
			A:  b.Strength.A,
			B:  b.Strength.B,
			AB: b.Strength.AB,
		},
		Casualties: statistics.SideNumbers{
			A:  b.Casualties.A,
			B:  b.Casualties.B,
			AB: b.Casualties.AB,
		},
		FactionsBySide: battles.IDsBySide{
			A: []uuid.UUID{Faction().ID},
			B: []uuid.UUID{Faction2().ID, Faction3().ID},
		},
		CommandersBySide: battles.IDsBySide{
			A: []uuid.UUID{Commander().ID},
			B: []uuid.UUID{Commander2().ID, Commander3().ID, Commander4().ID, Commander5().ID},
		},
		CommandersByFaction: battles.CommandersByFaction{
			Faction().ID:  []uuid.UUID{Commander().ID},
			Faction2().ID: []uuid.UUID{Commander2().ID, Commander3().ID},
			Faction3().ID: []uuid.UUID{Commander4().ID, Commander5().ID},
		},
	}
}
