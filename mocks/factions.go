package mocks

import (
	"github.com/sasalatart/batcoms/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

var factionUUID = uuid.NewV4()
var factionUUID2 = uuid.NewV4()
var factionUUID3 = uuid.NewV4()

// FactionsDataStore mocks datastores used to find & create factions
type FactionsDataStore struct {
	mock.Mock
}

// FindOne mocks finding one faction via FactionsDataStore
func (fs *FactionsDataStore) FindOne(query interface{}, args ...interface{}) (domain.Faction, error) {
	mockArgs := fs.Called(append([]interface{}{query}, args...)...)
	return mockArgs.Get(0).(domain.Faction), mockArgs.Error(1)
}

// CreateOne mocks creating one faction via FactionsDataStore
func (fs *FactionsDataStore) CreateOne(data domain.CreateFactionInput) (uuid.UUID, error) {
	mockArgs := fs.Called(data)
	return mockArgs.Get(0).(uuid.UUID), mockArgs.Error(1)
}

// Faction returns an instance of domain.Faction that may be used for mocking purposes
func Faction() domain.Faction {
	return factionFromScraped(SFaction(), factionUUID)
}

// Faction2 returns an instance of domain.Faction that may be used for mocking purposes
func Faction2() domain.Faction {
	return factionFromScraped(SFaction2(), factionUUID2)
}

// Faction3 returns an instance of domain.Faction that may be used for mocking purposes
func Faction3() domain.Faction {
	return factionFromScraped(SFaction3(), factionUUID3)
}

// CreateFactionInput returns an instance of domain.CreateFactionInput that may be used for mocking
// inputs to create factions
func CreateFactionInput() domain.CreateFactionInput {
	return createFactionInputFromFaction(Faction())
}

// CreateFactionInput2 returns an instance of domain.CreateFactionInput that may be used for mocking
// inputs to create factions
func CreateFactionInput2() domain.CreateFactionInput {
	return createFactionInputFromFaction(Faction2())
}

// CreateFactionInput3 returns an instance of domain.CreateFactionInput that may be used for mocking
// inputs to create factions
func CreateFactionInput3() domain.CreateFactionInput {
	return createFactionInputFromFaction(Faction3())
}

func factionFromScraped(sf domain.SParticipant, uuid uuid.UUID) domain.Faction {
	return domain.Faction{
		ID:      uuid,
		WikiID:  sf.ID,
		URL:     sf.URL,
		Name:    sf.Name,
		Summary: sf.Extract,
	}
}

func createFactionInputFromFaction(f domain.Faction) domain.CreateFactionInput {
	return domain.CreateFactionInput{
		WikiID:  f.WikiID,
		URL:     f.URL,
		Name:    f.Name,
		Summary: f.Summary,
	}
}
