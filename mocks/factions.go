package mocks

import (
	"github.com/sasalatart/batcoms/domain/factions"
	"github.com/sasalatart/batcoms/domain/wikiactors"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

var factionUUID = uuid.NewV4()
var factionUUID2 = uuid.NewV4()
var factionUUID3 = uuid.NewV4()

// FactionsRepository mocks datastores used to find & create factions
type FactionsRepository struct {
	mock.Mock
}

// FindOne mocks finding one faction via FactionsRepository
func (fs *FactionsRepository) FindOne(query interface{}, args ...interface{}) (factions.Faction, error) {
	mockArgs := fs.Called(append([]interface{}{query}, args...)...)
	return mockArgs.Get(0).(factions.Faction), mockArgs.Error(1)
}

// CreateOne mocks creating one faction via FactionsRepository
func (fs *FactionsRepository) CreateOne(data factions.CreationInput) (uuid.UUID, error) {
	mockArgs := fs.Called(data)
	return mockArgs.Get(0).(uuid.UUID), mockArgs.Error(1)
}

// Faction returns an instance of factions.Faction that may be used for mocking purposes
func Faction() factions.Faction {
	return factionFromScraped(WikiFaction(), factionUUID)
}

// Faction2 returns an instance of factions.Faction that may be used for mocking purposes
func Faction2() factions.Faction {
	return factionFromScraped(WikiFaction2(), factionUUID2)
}

// Faction3 returns an instance of factions.Faction that may be used for mocking purposes
func Faction3() factions.Faction {
	return factionFromScraped(WikiFaction3(), factionUUID3)
}

// FactionCreationInput returns an instance of factions.CreationInput that may be used for mocking
// inputs to create factions
func FactionCreationInput() factions.CreationInput {
	return createFactionInputFromFaction(Faction())
}

// FactionCreationInput2 returns an instance of factions.CreationInput that may be used for mocking
// inputs to create factions
func FactionCreationInput2() factions.CreationInput {
	return createFactionInputFromFaction(Faction2())
}

// FactionCreationInput3 returns an instance of factions.CreationInput that may be used for mocking
// inputs to create factions
func FactionCreationInput3() factions.CreationInput {
	return createFactionInputFromFaction(Faction3())
}

func factionFromScraped(wf wikiactors.Actor, uuid uuid.UUID) factions.Faction {
	return factions.Faction{
		ID:      uuid,
		WikiID:  wf.ID,
		URL:     wf.URL,
		Name:    wf.Name,
		Summary: wf.Extract,
	}
}

func createFactionInputFromFaction(f factions.Faction) factions.CreationInput {
	return factions.CreationInput{
		WikiID:  f.WikiID,
		URL:     f.URL,
		Name:    f.Name,
		Summary: f.Summary,
	}
}
