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
var faction2UUID = uuid.NewV4()
var faction3UUID = uuid.NewV4()

// FactionsStore mocks the behaviour of store.FactionsFinder interface
type FactionsStore struct {
	mock.Mock
}

// FindOne mocks finding one faction via FactionsStore
func (fs *FactionsStore) FindOne(query interface{}, args ...interface{}) (domain.Faction, error) {
	mockArgs := fs.Called(append([]interface{}{query}, args...)...)
	return mockArgs.Get(0).(domain.Faction), mockArgs.Error(1)
}

// Faction returns an instance of domain.Faction that may be used for mocking purposes
func Faction() domain.Faction {
	mock := baseFactionMock
	return mock
}

// Faction2 returns an instance of domain.Faction that may be used for mocking purposes
func Faction2() domain.Faction {
	return domain.Faction{
		ID:      faction2UUID,
		WikiID:  20611504,
		URL:     "https://en.wikipedia.org/wiki/Imperial_Russia",
		Name:    "Russian Empire",
		Summary: "The Russian Empire was an empire that extended across Eurasia and North America from 1721, following the end of the Great Northern War, until the Republic was proclaimed by the Provisional Government that took power after the February Revolution of 1917. The third-largest empire in history, at its greatest extent stretching over three continents, Europe, Asia, and North America, the Russian Empire was surpassed in size only by the British and Mongol empires. The rise of the Russian Empire coincided with the decline of neighboring rival powers: the Swedish Empire, the Polish–Lithuanian Commonwealth, Persia and the Ottoman Empire. It played a major role in 1812–1814 in defeating Napoleon's ambitions to control Europe and expanded to the west and south.",
	}
}

// Faction3 returns an instance of domain.Faction that may be used for mocking purposes
func Faction3() domain.Faction {
	return domain.Faction{
		ID:      faction3UUID,
		WikiID:  266894,
		URL:     "https://en.wikipedia.org/wiki/Austrian_Empire",
		Name:    "Austrian Empire",
		Summary: "The Austrian Empire was a Central European multinational great power from 1804 to 1867, created by proclamation out of the realms of the Habsburgs. During its existence, it was the third most populous empire after the Russian Empire and the United Kingdom in Europe. Along with Prussia, it was one of the two major powers of the German Confederation. Geographically, it was the third largest empire in Europe after the Russian Empire and the First French Empire. Proclaimed in response to the First French Empire, it partially overlapped with the Holy Roman Empire until the latter's dissolution in 1806.",
	}
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
