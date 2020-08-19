package mocks

import (
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain"
)

// CreateCommanderInput returns an instance of domain.CreateCommanderInput that may be used for
// mocking inputs used to create commanders. The overrides parameter may contain values used to
// override the fallback values used by the default mock
func CreateCommanderInput(overrides domain.CreateCommanderInput) (domain.CreateCommanderInput, error) {
	mock := domain.CreateCommanderInput{
		WikiID:  69880,
		URL:     "https://en.wikipedia.org/wiki/Emperor_Napoleon_I",
		Name:    "Napoleon",
		Summary: `Napoleon Bonaparte, born Napoleone di Buonaparte, byname "Le Corse" or "Le Petit Caporal", was a French statesman and military leader who became notorious as an artillery commander during the French Revolution. He led many successful campaigns during the French Revolutionary Wars and was Emperor of the French as Napoleon I from 1804 until 1814 and again briefly in 1815 during the Hundred Days. Napoleon dominated European and global affairs for more than a decade while leading France against a series of coalitions during the Napoleonic Wars. He won many of these wars and a vast majority of his battles, building a large empire that ruled over much of continental Europe before its final collapse in 1815. He is regarded as one of the greatest military commanders in history, and his wars and campaigns are studied at military schools worldwide. Napoleon's political and cultural legacy has made him one of the most celebrated and controversial leaders in human history.`,
	}
	if err := mergo.Merge(&mock, overrides, mergo.WithOverride); err != nil {
		return mock, errors.Wrap(err, "Creating commander mock")
	}
	return mock, nil
}
