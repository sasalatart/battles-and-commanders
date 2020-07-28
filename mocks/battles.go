package mocks

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/sasalatart/batcoms/scraper/domain"
)

// Battle returns an instance of domain.Battle that may be used for testing purposes. Some values
// may be overriden by the input battle.
func Battle(overrides domain.Battle) (domain.Battle, error) {
	mock := domain.Battle{
		ID:          118372,
		URL:         "https://en.wikipedia.org/wiki/Battle_of_Austerlitz",
		Name:        "Battle of Austerlitz",
		PartOf:      "Part of the War of the Third Coalition",
		Description: "Battle of the Napoleonic Wars",
		Extract:     "The Battle of Austerlitz, also known as the Battle of the Three Emperors, was one of the most important and decisive engagements of the Napoleonic Wars. In what is widely regarded as the greatest victory achieved by Napoleon, the Grande Armée of France defeated a larger Russian and Austrian army led by Emperor Alexander I and Holy Roman Emperor Francis II. The battle occurred near the town of Austerlitz in the Austrian Empire. Austerlitz brought the War of the Third Coalition to a rapid end, with the Treaty of Pressburg signed by the Austrians later in the month. The battle is often cited as a tactical masterpiece, in the same league as other historic engagements like Cannae or Gaugamela.",
		Date:        "2 December 1805",
		Location: domain.Location{
			Place:     "Austerlitz, Moravia, Austria",
			Latitude:  "49°8'0\"N",
			Longitude: "16°46'0\"E",
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
		Factions: domain.SideParticipants{
			A: []int{21418258},
			B: []int{20611504, 266894},
		},
		Commanders: domain.SideParticipants{
			A: []int{69880},
			B: []int{27126603, 251000, 11551, 14092123},
		},
		CommandersByFaction: map[int][]int{
			21418258: {69880},
			20611504: {27126603, 251000},
			266894:   {11551, 14092123},
		},
	}

	if err := mergo.Merge(&mock, overrides, mergo.WithOverride); err != nil {
		return mock, fmt.Errorf("Failed creating battle mock: %s", err)
	}
	return mock, nil
}
