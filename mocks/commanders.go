package mocks

import (
	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/domain/wikiactors"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

var commanderUUID = uuid.NewV4()
var commander2UUID = uuid.NewV4()
var commander3UUID = uuid.NewV4()
var commander4UUID = uuid.NewV4()
var commander5UUID = uuid.NewV4()

// CommandersRepository mocks repositories used to read and write commanders
type CommandersRepository struct {
	mock.Mock
}

// FindOne mocks finding one commander via CommandersRepository
func (r *CommandersRepository) FindOne(query commanders.FindOneQuery) (commanders.Commander, error) {
	mockArgs := r.Called(query)
	return mockArgs.Get(0).(commanders.Commander), mockArgs.Error(1)
}

// FindMany mocks finding many commanders via CommandersRepository
func (r *CommandersRepository) FindMany(query commanders.FindManyQuery, page int) ([]commanders.Commander, int, error) {
	mockArgs := r.Called(query, page)
	return mockArgs.Get(0).([]commanders.Commander), mockArgs.Int(1), mockArgs.Error(2)
}

// CreateOne mocks creating one commander via CommandersRepository
func (r *CommandersRepository) CreateOne(data commanders.CreationInput) (uuid.UUID, error) {
	mockArgs := r.Called(data)
	return mockArgs.Get(0).(uuid.UUID), mockArgs.Error(1)
}

// Commander returns an instance of commanders.Commander that may be used for mocking purposes
func Commander() commanders.Commander {
	return commanderFromScraped(WikiCommander(), commanderUUID)
}

// Commander2 returns an instance of commanders.Commander that may be used for mocking purposes
func Commander2() commanders.Commander {
	return commanderFromScraped(WikiCommander2(), commander2UUID)
}

// Commander3 returns an instance of commanders.Commander that may be used for mocking purposes
func Commander3() commanders.Commander {
	return commanderFromScraped(WikiCommander3(), commander3UUID)
}

// Commander4 returns an instance of commanders.Commander that may be used for mocking purposes
func Commander4() commanders.Commander {
	return commanderFromScraped(WikiCommander4(), commander4UUID)
}

// Commander5 returns an instance of commanders.Commander that may be used for mocking purposes
func Commander5() commanders.Commander {
	return commanderFromScraped(WikiCommander5(), commander5UUID)
}

// CommanderCreationInput returns an instance of commanders.CreationInput that may be used for
// mocking inputs to create commanders
func CommanderCreationInput() commanders.CreationInput {
	return createCommanderInputFromCommander(Commander())
}

// CommanderCreationInput2 returns an instance of commanders.CreationInput that may be used for
// mocking inputs to create commanders
func CommanderCreationInput2() commanders.CreationInput {
	return createCommanderInputFromCommander(Commander2())
}

// CommanderCreationInput3 returns an instance of commanders.CreationInput that may be used for
// mocking inputs to create commanders
func CommanderCreationInput3() commanders.CreationInput {
	return createCommanderInputFromCommander(Commander3())
}

// CommanderCreationInput4 returns an instance of commanders.CreationInput that may be used for
// mocking inputs to create commanders
func CommanderCreationInput4() commanders.CreationInput {
	return createCommanderInputFromCommander(Commander4())
}

// CommanderCreationInput5 returns an instance of commanders.CreationInput that may be used for
// mocking inputs to create commanders
func CommanderCreationInput5() commanders.CreationInput {
	return createCommanderInputFromCommander(Commander5())
}

func commanderFromScraped(wc wikiactors.Actor, uuid uuid.UUID) commanders.Commander {
	return commanders.Commander{
		ID:      uuid,
		WikiID:  wc.ID,
		URL:     wc.URL,
		Name:    wc.Name,
		Summary: wc.Extract,
	}
}

func createCommanderInputFromCommander(c commanders.Commander) commanders.CreationInput {
	return commanders.CreationInput{
		WikiID:  c.WikiID,
		URL:     c.URL,
		Name:    c.Name,
		Summary: c.Summary,
	}
}
