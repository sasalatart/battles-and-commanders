package summaries_test

import (
	"testing"

	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/services/summaries"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSummary(t *testing.T) {
	t.Run("InvalidURL", func(t *testing.T) {
		t.Parallel()
		_, err := summaries.Get("https://i-do-not-exist.org")
		require.Error(t, err, "Fetching summary for invalid URL")
		assert.Contains(t, err.Error(), "not a wiki page", "Comparing error message contents")
	})

	t.Run("ValidURL", func(t *testing.T) {
		t.Parallel()
		got, err := summaries.Get("https://en.wikipedia.org/wiki/Battle_of_Austerlitz")
		require.NoError(t, err, "Fetching summary for valid URL")
		expected := domain.Summary{
			PageID:       118372,
			Type:         "standard",
			Title:        "Battle of Austerlitz",
			DisplayTitle: "Battle of Austerlitz",
			Description:  "Battle of the Napoleonic Wars",
			Extract:      "The Battle of Austerlitz, also known as the Battle of the Three Emperors, was one of the most important and decisive engagements of the Napoleonic Wars. In what is widely regarded as the greatest victory achieved by Napoleon, the Grande Armée of France defeated a larger Russian and Austrian army led by Emperor Alexander I and Holy Roman Emperor Francis II. The battle occurred near the town of Austerlitz in the Austrian Empire. Austerlitz brought the War of the Third Coalition to a rapid end, with the Treaty of Pressburg signed by the Austrians later in the month. The battle is often cited as a tactical masterpiece, in the same league as other historic engagements like Cannae or Gaugamela.",
		}
		assert.Equal(t, expected, got, "Comparing obtained summary with expected one")
	})
}
