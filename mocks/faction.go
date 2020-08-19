package mocks

import (
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain"
)

// CreateFactionInput returns an instance of domain.CreateFactionInput that may be used for mocking
// inputs used to create factions. The overrides parameter may contain values used to override the
// fallback values used by the default mock
func CreateFactionInput(overrides domain.CreateFactionInput) (domain.CreateFactionInput, error) {
	mock := domain.CreateFactionInput{
		WikiID:  21418258,
		URL:     "https://en.wikipedia.org/wiki/French_First_Empire",
		Name:    "First French Empire",
		Summary: "The First French Empire, officially the French Empire or the Napoleonic Empire, was the empire of Napoleon Bonaparte of France and the dominant power in much of continental Europe at the beginning of the 19th century. Although France had already established an overseas colonial empire beginning in the 17th century, the French state had remained a kingdom under the Bourbons and a republic after the French Revolution. Historians refer to Napoleon's regime as the First Empire to distinguish it from the restorationist Second Empire (1852–1870) ruled by his nephew Napoleon III.",
	}
	if err := mergo.Merge(&mock, overrides, mergo.WithOverride); err != nil {
		return mock, errors.Wrap(err, "Creating faction mock")
	}
	return mock, nil
}
