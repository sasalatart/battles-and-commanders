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
var commander2UUID = uuid.NewV4()
var commander3UUID = uuid.NewV4()
var commander4UUID = uuid.NewV4()
var commander5UUID = uuid.NewV4()

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

// Commander2 returns an instance of domain.Commander that may be used for mocking purposes
func Commander2() domain.Commander {
	return domain.Commander{
		ID:      commander2UUID,
		WikiID:  27126603,
		URL:     "https://en.wikipedia.org/wiki/Alexander_I_of_Russia",
		Name:    "Alexander I of Russia",
		Summary: "Alexander I was the Emperor of Russia (Tsar) between 1801 and 1825. He was the eldest son of Paul I and Sophie Dorothea of Württemberg. Alexander was the first king of Congress Poland, reigning from 1815 to 1825, as well as the first Russian Grand Duke of Finland, reigning from 1809 to 1825.",
	}
}

// Commander3 returns an instance of domain.Commander that may be used for mocking purposes
func Commander3() domain.Commander {
	return domain.Commander{
		ID:      commander3UUID,
		WikiID:  251000,
		URL:     "https://en.wikipedia.org/wiki/Mikhail_Illarionovich_Kutuzov",
		Name:    "Mikhail Kutuzov",
		Summary: "Prince Mikhail Illarionovich Golenishchev-Kutuzov was a Field Marshal of the Russian Empire. He served as one of the finest military officers and diplomats of Russia under the reign of three Romanov Tsars: Catherine II, Paul I and Alexander I. His military career was closely associated with the rising period of Russia from the end of the 18th century to the beginning of the 19th century. Kutuzov is considered to have been one of the best Russian generals.",
	}
}

// Commander4 returns an instance of domain.Commander that may be used for mocking purposes
func Commander4() domain.Commander {
	return domain.Commander{
		ID:      commander4UUID,
		WikiID:  11551,
		URL:     "https://en.wikipedia.org/wiki/Francis_II,_Holy_Roman_Emperor",
		Name:    "Francis II, Holy Roman Emperor",
		Summary: "Francis II was the last Holy Roman Emperor, ruling from 1792 until 6 August 1806, when he dissolved the Holy Roman Empire after the decisive defeat at the hands of the First French Empire led by Napoleon at the Battle of Austerlitz. In 1804, he had founded the Austrian Empire and became Francis I, the first Emperor of Austria, ruling from 1804 to 1835, so later he was named the first Doppelkaiser in history.. For the two years between 1804 and 1806, Francis used the title and style by the Grace of God elected Roman Emperor, ever Augustus, hereditary Emperor of Austria and he was called the Emperor of both the Holy Roman Empire and Austria. He was also Apostolic King of Hungary, Croatia and Bohemia as Francis I. He also served as the first president of the German Confederation following its establishment in 1815.",
	}
}

// Commander5 returns an instance of domain.Commander that may be used for mocking purposes
func Commander5() domain.Commander {
	return domain.Commander{
		ID:      commander5UUID,
		WikiID:  14092123,
		URL:     "https://en.wikipedia.org/wiki/Franz_von_Weyrother",
		Name:    "Franz von Weyrother",
		Summary: "Franz von Weyrother was an Austrian staff officer and general who fought during the French Revolutionary Wars and the Napoleonic Wars. He drew up the plans for the disastrous defeats at the Battle of Rivoli, Battle of Hohenlinden and the Battle of Austerlitz, in which the Austrian army was defeated by Napoleon Bonaparte twice and Jean Moreau once.",
	}
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
