package mocks

import (
	"github.com/sasalatart/batcoms/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

var commanderUUID = uuid.NewV4()
var commander2UUID = uuid.NewV4()
var commander3UUID = uuid.NewV4()
var commander4UUID = uuid.NewV4()
var commander5UUID = uuid.NewV4()

// CommandersDataStore mocks datastores used to find & create commanders
type CommandersDataStore struct {
	mock.Mock
}

// FindOne mocks finding one commander via CommandersDataStore
func (cs *CommandersDataStore) FindOne(query interface{}, args ...interface{}) (domain.Commander, error) {
	mockArgs := cs.Called(append([]interface{}{query}, args...)...)
	return mockArgs.Get(0).(domain.Commander), mockArgs.Error(1)
}

// CreateOne mocks creating one commander via CommandersDataStore
func (cs *CommandersDataStore) CreateOne(data domain.CreateCommanderInput) (uuid.UUID, error) {
	mockArgs := cs.Called(data)
	return mockArgs.Get(0).(uuid.UUID), mockArgs.Error(1)
}

// Commander returns an instance of domain.Commander that may be used for mocking purposes
func Commander() domain.Commander {
	return commanderFromScraped(SCommander(), commanderUUID)
}

// Commander2 returns an instance of domain.Commander that may be used for mocking purposes
func Commander2() domain.Commander {
	return commanderFromScraped(SCommander2(), commander2UUID)
}

// Commander3 returns an instance of domain.Commander that may be used for mocking purposes
func Commander3() domain.Commander {
	return commanderFromScraped(SCommander3(), commander3UUID)
}

// Commander4 returns an instance of domain.Commander that may be used for mocking purposes
func Commander4() domain.Commander {
	return commanderFromScraped(SCommander4(), commander4UUID)
}

// Commander5 returns an instance of domain.Commander that may be used for mocking purposes
func Commander5() domain.Commander {
	return commanderFromScraped(SCommander5(), commander5UUID)
}

// CreateCommanderInput returns an instance of domain.CreateCommanderInput that may be used for
// mocking inputs to create commanders
func CreateCommanderInput() domain.CreateCommanderInput {
	return createCommanderInputFromCommander(Commander())
}

// CreateCommanderInput2 returns an instance of domain.CreateCommanderInput that may be used for
// mocking inputs to create commanders
func CreateCommanderInput2() domain.CreateCommanderInput {
	return createCommanderInputFromCommander(Commander2())
}

// CreateCommanderInput3 returns an instance of domain.CreateCommanderInput that may be used for
// mocking inputs to create commanders
func CreateCommanderInput3() domain.CreateCommanderInput {
	return createCommanderInputFromCommander(Commander3())
}

// CreateCommanderInput4 returns an instance of domain.CreateCommanderInput that may be used for
// mocking inputs to create commanders
func CreateCommanderInput4() domain.CreateCommanderInput {
	return createCommanderInputFromCommander(Commander4())
}

// CreateCommanderInput5 returns an instance of domain.CreateCommanderInput that may be used for
// mocking inputs to create commanders
func CreateCommanderInput5() domain.CreateCommanderInput {
	return createCommanderInputFromCommander(Commander5())
}

func commanderFromScraped(sc domain.SParticipant, uuid uuid.UUID) domain.Commander {
	return domain.Commander{
		ID:      uuid,
		WikiID:  sc.ID,
		URL:     sc.URL,
		Name:    sc.Name,
		Summary: sc.Extract,
	}
}

func createCommanderInputFromCommander(c domain.Commander) domain.CreateCommanderInput {
	return domain.CreateCommanderInput{
		WikiID:  c.WikiID,
		URL:     c.URL,
		Name:    c.Name,
		Summary: c.Summary,
	}
}
