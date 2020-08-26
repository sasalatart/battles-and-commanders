package mocks

import (
	"github.com/sasalatart/batcoms/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

var baseBattleMock = domain.Battle{
	ID:        uuid.NewV4(),
	WikiID:    118372,
	URL:       "https://en.wikipedia.org/wiki/Battle_of_Austerlitz",
	Name:      "Battle of Austerlitz",
	PartOf:    "Part of the War of the Third Coalition",
	Summary:   "The Battle of Austerlitz, also known as the Battle of the Three Emperors, was one of the most important and decisive engagements of the Napoleonic Wars. In what is widely regarded as the greatest victory achieved by Napoleon, the Grande Armée of France defeated a larger Russian and Austrian army led by Emperor Alexander I and Holy Roman Emperor Francis II. The battle occurred near the town of Austerlitz in the Austrian Empire. Austerlitz brought the War of the Third Coalition to a rapid end, with the Treaty of Pressburg signed by the Austrians later in the month. The battle is often cited as a tactical masterpiece, in the same league as other historic engagements like Cannae or Gaugamela.",
	StartDate: "1805-12-02",
	EndDate:   "1805-12-02",
	Location: domain.Location{
		Place:     "Austerlitz, Moravia, Austria",
		Latitude:  "49°8′N",
		Longitude: "16°46′E",
	},
	Result:             "Decisive French victory. Treaty of Pressburg. Effective end of the Third Coalition",
	TerritorialChanges: "Dissolution of the Holy Roman Empire and creation of the Confederation of the Rhine",
	Strength: domain.SideNumbers{
		A: "65,000–75,000",
		B: "84,000–95,000",
	},
	Casualties: domain.SideNumbers{
		A: "1,305 killed 6,991 wounded 573 captured",
		B: "16,000 killed and wounded 20,000 captured",
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

// BattlesStore mocks the behaviour of store.BattlesFinder interface
type BattlesStore struct {
	mock.Mock
}

// FindOne mocks finding one commander via BattlesStore
func (bs *BattlesStore) FindOne(query interface{}, args ...interface{}) (domain.Battle, error) {
	mockArgs := bs.Called(append([]interface{}{query}, args...)...)
	return mockArgs.Get(0).(domain.Battle), mockArgs.Error(1)
}

// Battle returns an instance of domain.Battle that may be used for mocking purposes
func Battle() domain.Battle {
	mock := baseBattleMock
	return mock
}

// CreateBattleInput returns an instance of domain.CreateBattleInput that may be used for mocking
// inputs to create battles
func CreateBattleInput() domain.CreateBattleInput {
	return domain.CreateBattleInput{
		WikiID:    baseBattleMock.WikiID,
		URL:       baseBattleMock.URL,
		Name:      baseBattleMock.Name,
		PartOf:    baseBattleMock.PartOf,
		Summary:   baseBattleMock.Summary,
		StartDate: baseBattleMock.StartDate,
		EndDate:   baseBattleMock.EndDate,
		Location: domain.Location{
			Place:     baseBattleMock.Location.Place,
			Latitude:  baseBattleMock.Location.Latitude,
			Longitude: baseBattleMock.Location.Longitude,
		},
		Result:             baseBattleMock.Result,
		TerritorialChanges: baseBattleMock.TerritorialChanges,
		Strength: domain.SideNumbers{
			A: baseBattleMock.Strength.A,
			B: baseBattleMock.Strength.B,
		},
		Casualties: domain.SideNumbers{
			A: baseBattleMock.Casualties.A,
			B: baseBattleMock.Casualties.B,
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
