package mocks

import (
	"github.com/sasalatart/batcoms/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

var baseCommanderMock = domain.Commander{
	ID:      uuid.NewV4(),
	WikiID:  69880,
	URL:     "https://en.wikipedia.org/wiki/Emperor_Napoleon_I",
	Name:    "Napoleon",
	Summary: `Napoleon Bonaparte, born Napoleone di Buonaparte, byname "Le Corse" or "Le Petit Caporal", was a French statesman and military leader who became notorious as an artillery commander during the French Revolution. He led many successful campaigns during the French Revolutionary Wars and was Emperor of the French as Napoleon I from 1804 until 1814 and again briefly in 1815 during the Hundred Days. Napoleon dominated European and global affairs for more than a decade while leading France against a series of coalitions during the Napoleonic Wars. He won many of these wars and a vast majority of his battles, building a large empire that ruled over much of continental Europe before its final collapse in 1815. He is regarded as one of the greatest military commanders in history, and his wars and campaigns are studied at military schools worldwide. Napoleon's political and cultural legacy has made him one of the most celebrated and controversial leaders in human history.`,
}

// CommandersStore mocks the behaviour of store.CommandersFinder interface
type CommandersStore struct {
	mock.Mock
}

// FindOne mocks finding one commander via CommandersStore
func (cs *CommandersStore) FindOne(query interface{}, args ...interface{}) (domain.Commander, error) {
	mockArgs := cs.Called(append([]interface{}{query}, args...)...)
	return mockArgs.Get(0).(domain.Commander), mockArgs.Error(1)
}

// Commander returns an instance of domain.Commander that may be used for mocking purposes
func Commander() domain.Commander {
	mock := baseCommanderMock
	return mock
}

// CreateCommanderInput returns an instance of domain.CreateCommanderInput that may be used for
// mocking inputs to create commanders
func CreateCommanderInput() domain.CreateCommanderInput {
	return domain.CreateCommanderInput{
		WikiID:  baseCommanderMock.WikiID,
		URL:     baseCommanderMock.URL,
		Name:    baseCommanderMock.Name,
		Summary: baseCommanderMock.Summary,
	}
}
