package mocks

import (
	"github.com/sasalatart/batcoms/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

var baseFactionMock = domain.Faction{
	ID:      uuid.NewV4(),
	WikiID:  21418258,
	URL:     "https://en.wikipedia.org/wiki/French_First_Empire",
	Name:    "First French Empire",
	Summary: "The First French Empire, officially the French Empire or the Napoleonic Empire, was the empire of Napoleon Bonaparte of France and the dominant power in much of continental Europe at the beginning of the 19th century. Although France had already established an overseas colonial empire beginning in the 17th century, the French state had remained a kingdom under the Bourbons and a republic after the French Revolution. Historians refer to Napoleon's regime as the First Empire to distinguish it from the restorationist Second Empire (1852–1870) ruled by his nephew Napoleon III.",
}

// FactionsStore mocks the behaviour of store.Factions interface
type FactionsStore struct {
	mock.Mock
}

// FindOne mocks finding one faction via FactionsStore
func (m *FactionsStore) FindOne(query interface{}, args ...interface{}) (domain.Faction, error) {
	mockArgs := m.Called(append([]interface{}{query}, args...)...)
	return mockArgs.Get(0).(domain.Faction), mockArgs.Error(1)
}

// Faction returns an instance of domain.Faction that may be used for mocking purposes
func Faction() domain.Faction {
	mock := baseFactionMock
	return mock
}

// CreateFactionInput returns an instance of domain.CreateFactionInput that may be used for mocking
// inputs to create factions
func CreateFactionInput() domain.CreateFactionInput {
	return domain.CreateFactionInput{
		WikiID:  baseFactionMock.WikiID,
		URL:     baseFactionMock.URL,
		Name:    baseFactionMock.Name,
		Summary: baseFactionMock.Summary,
	}
}
