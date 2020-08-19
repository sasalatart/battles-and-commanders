package mocks

import (
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain"
	uuid "github.com/satori/go.uuid"
)

// CreateBattleInput returns an instance of domain.CreateBattleInput that may be used for mocking
// inputs to create battles. The overrides parameter may contain values used to override the
// fallback values used by the default mock
func CreateBattleInput(overrides domain.CreateBattleInput) (domain.CreateBattleInput, error) {
	mock := domain.CreateBattleInput{
		WikiID:    118372,
		URL:       "https://en.wikipedia.org/wiki/Battle_of_Austerlitz",
		Name:      "Battle of Austerlitz",
		PartOf:    "Part of the War of the Third Coalition",
		Summary:   "The Battle of Austerlitz, also known as the Battle of the Three Emperors, was one of the most important and decisive engagements of the Napoleonic Wars. In what is widely regarded as the greatest victory achieved by Napoleon, the Grande Armée of France defeated a larger Russian and Austrian army led by Emperor Alexander I and Holy Roman Emperor Francis II. The battle occurred near the town of Austerlitz in the Austrian Empire. Austerlitz brought the War of the Third Coalition to a rapid end, with the Treaty of Pressburg signed by the Austrians later in the month. The battle is often cited as a tactical masterpiece, in the same league as other historic engagements like Cannae or Gaugamela.",
		StartDate: "1805-12-02",
		EndDate:   "1805-12-02",
		Location: domain.Location{
			Place:     "Austerlitz, Moravia, Austria",
			Latitude:  "49°8′N",
			Longitude: "16°46′E",
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
		FactionsBySide: domain.ParticipantsIDsBySide{
			A: []uuid.UUID{uuid.NewV4()},
			B: []uuid.UUID{uuid.NewV4(), uuid.NewV4()},
		},
		CommandersBySide: domain.ParticipantsIDsBySide{
			A: []uuid.UUID{uuid.NewV4()},
			B: []uuid.UUID{uuid.NewV4(), uuid.NewV4(), uuid.NewV4(), uuid.NewV4()},
		},
	}
	if err := mergo.Merge(&mock, overrides, mergo.WithOverride); err != nil {
		return mock, errors.Wrap(err, "Creating battle mock")
	}
	return mock, nil
}
