package list_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sasalatart/batcoms/mocks"
	"github.com/sasalatart/batcoms/pkg/logger"
	"github.com/sasalatart/batcoms/pkg/scraper/list"
)

func TestList(t *testing.T) {
	infoWriter := new(mocks.Writer)
	logger := logger.New(infoWriter, new(mocks.Writer))

	battlesList := list.Scrape(logger)

	minListLen := 4500
	assert.GreaterOrEqualf(t, len(battlesList), minListLen, "Should obtain more than %d battles", minListLen)

	t.Run("BattlesNames", func(t *testing.T) {
		t.Parallel()

		var gotNames []string
		for _, battle := range battlesList {
			gotNames = append(gotNames, battle.Name)
		}
		for _, expected := range expectedNames {
			assert.Contains(t, gotNames, expected, "Expected %q to be present in the scraped list", expected)
		}
	})

	t.Run("Logs", func(t *testing.T) {
		t.Parallel()

		for _, urlPart := range expectedLists {
			found := false
			for _, log := range infoWriter.Writes {
				if strings.Contains(log, urlPart) {
					found = true
				}
			}
			assert.Truef(t, found, "Scraper did not log results for %q", urlPart)
		}
		expectedLogsAmount := len(battlesList)
		gotAmountText := infoWriter.Writes[len(infoWriter.Writes)-1]
		assert.Containsf(
			t,
			gotAmountText,
			strconv.Itoa(expectedLogsAmount),
			"Expected the amounts log to contain the number %d", expectedLogsAmount,
		)
	})
}

var expectedNames = []string{
	"Action at Blue Mills Landing",
	"Assault on Copenhagen",
	"Attack on Sydney Harbour",
	"Battle at Chedabucto",
	"Battle for Baby 700",
	"Battle in Shakhtarsk Raion",
	"Battle of 1st Bull Run",
	"Battle on Lake Peipus",
	"Battles of Barfleur and La Hogue",
	"Blockade of Germany",
	"Burning of Dungannon",
	"Campaign to Defend Siping",
	"Capture of Amara",
	"Cesena Bloodbath",
	"Combat of the Thirty",
	"Defense of Oguta",
	"Deir ez-Zor Governorate clashes",
	"East Ghouta inter-rebel conflict",
	"Fall of Ashdod",
	"Fourth Crusade",
	"Francisco's Fight",
	"Geary Ambush",
	"Great Stand on the Ugra River",
	"Greater Poland Uprising",
	"Gulf of Sidra incident",
	"Gustafsen Lake Standoff",
	"Invasion of Guantánamo Bay",
	"Jacobite rising of 1689",
	"Kemp's Landing",
	"Liberation of Paris",
	"Long Run Massacre",
	"Lungi Lol confrontation",
	"Manchu conquest of China",
	"Naval bombardment of Japan",
	"Ndop prison break",
	"Nootka Crisis",
	"Operatie Kraai",
	"Operation Achilles",
	"Prayer Book Rebellion",
	"Raid on Chignecto",
	"Recovery of Ré island",
	"Relief of Goes",
	"Revolt of Babylon (626 BC)",
	"Sack of Rome",
	"Second Arab siege of Constantinople",
	"Sherman's March",
	"Sinking of CSS Alabama",
	"Skirmish at Blackwater Creek",
	"Syria missile strikes",
	"The Crossing",
	"United Arab Emirates takeover of Socotra",
	"United States occupation of Veracruz",
	"Vardar Offensive",
	"Wolseley Expedition",
}

var expectedLists = []string{
	"American_Civil_War_battles",
	"American_Revolutionary_War_battles",
	"battles_(alphabetical)",
	"battles_(geographic)",
	"battles_before_301",
	"battles_301-1300",
	"battles_1301-1600",
	"battles_1601-1800",
	"battles_1801-1900",
	"battles_1901-2000",
	"battles_since_2001",
	"Hundred_Years%27_War_battles",
	"military_engagements_of_World_War_I",
	"military_engagements_of_World_War_II",
	"Napoleonic_battles",
}
