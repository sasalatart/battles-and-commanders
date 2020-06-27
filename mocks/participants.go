package mocks

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/sasalatart/batcoms/scraper/domain"
)

// Faction returns a faction instance of domain.Participant that may be used for testing purposes.
// Some values may be overriden by the input faction.
func Faction(overrides domain.Participant) (domain.Participant, error) {
	mock := domain.Participant{
		Kind:        domain.FactionKind,
		ID:          21418258,
		URL:         "https://en.wikipedia.org/wiki/French_First_Empire",
		Flag:        "//upload.wikimedia.org/wikipedia/en/thumb/c/c3/Flag_of_France.svg/23px-Flag_of_France.svg.png",
		Name:        "First French Empire",
		Description: "Empire of Napoleon I of France between 1804–1815",
		Extract:     "The First French Empire, officially the French Empire or the Napoleonic Empire, was the empire of Napoleon Bonaparte of France and the dominant power in much of continental Europe at the beginning of the 19th century. Although France had already established an overseas colonial empire beginning in the 17th century, the French state had remained a kingdom under the Bourbons and a republic after the French Revolution. Historians refer to Napoleon's regime as the First Empire to distinguish it from the restorationist Second Empire (1852–1870) ruled by his nephew Napoleon III.",
	}

	if err := mergo.Merge(&mock, overrides, mergo.WithOverride); err != nil {
		return mock, fmt.Errorf("Failed creating faction mock: %s", err)
	}
	return mock, nil
}

// Commander returns a commander instance of domain.Participant that may be used for testing
// purposes. Some values may be overriden by the input commander.
func Commander(overrides domain.Participant) (domain.Participant, error) {
	mock := domain.Participant{
		Kind:        domain.CommanderKind,
		ID:          69880,
		URL:         "https://en.wikipedia.org/wiki/Emperor_Napoleon_I",
		Flag:        "",
		Name:        "Napoleon",
		Description: "19th century French military leader, strategist, and politician",
		Extract:     "Napoleon Bonaparte, born Napoleone di Buonaparte, was a French statesman and military leader who became famous as an artillery commander during the French Revolution. He led many successful campaigns during the French Revolutionary Wars and was Emperor of the French as Napoleon I from 1804 until 1814 and again briefly in 1815 during the Hundred Days. Napoleon dominated European and global affairs for more than a decade while leading France against a series of coalitions during the Napoleonic Wars. He won many of these wars and a vast majority of his battles, building a large empire that ruled over much of continental Europe before its final collapse in 1815. He is considered one of the greatest commanders in history, and his wars and campaigns are studied at military schools worldwide. Napoleon's political and cultural legacy has made him one of the most celebrated and controversial leaders in human history.",
	}

	if err := mergo.Merge(&mock, overrides, mergo.WithOverride); err != nil {
		return mock, fmt.Errorf("Failed creating commander mock: %s", err)
	}
	return mock, nil
}
