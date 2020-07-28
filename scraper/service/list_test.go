package service_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/sasalatart/batcoms/scraper/service"
	"github.com/sasalatart/batcoms/scraper/store"

	"github.com/sasalatart/batcoms/mocks"
)

func TestList(t *testing.T) {
	loggerMock := &mocks.Logger{}
	service := service.NewScraper(
		store.NewBattlesMem(),
		store.NewParticipantsMem(),
		loggerMock,
	)
	battlesList := service.List()

	gotBattles := len(battlesList)
	expectedMinBattles := 4500
	if gotBattles < expectedMinBattles {
		t.Errorf("Expected to scrap more than %d battles, but only got %d", expectedMinBattles, gotBattles)
	} else {
		t.Logf("Scraps more than %d battles", expectedMinBattles)
	}

	indexedBattlesNames := make(map[string]string)
	for _, battle := range battlesList {
		indexedBattlesNames[battle.Name] = battle.Name
	}
	nn := []string{
		"Action at Blue Mills Landing",
		"Ambush of Geary",
		"Assault on Copenhagen",
		"Attack at Country Harbour",
		"Battle at Chedabucto",
		"Battle for Baby 700",
		"Battle in Shakhtarsk Raion",
		"Battle of 1st Bull Run",
		"Battle on Lake Peipus",
		"Battles at Chignecto",
		"Battles of Barfleur and La Hogue",
		"Blockade of Germany",
		"Cesena Bloodbath",
		"Naval bombardment of Japan",
		"Burning of Dungannon",
		"Campaign to Defend Siping",
		"Capture of Amara",
		"Deir ez-Zor Governorate clashes",
		"Combat of the Thirty",
		"East Ghouta inter-rebel conflict",
		"Lungi Lol confrontation",
		"Manchu conquest of China",
		"Nootka Crisis",
		"The Crossing",
		"Fourth Crusade",
		"Defense of Oguta",
		"Nejd Expedition",
		"Fall of Ashdod",
		"Francisco's Fight",
		"Gulf of Sidra incident",
		"Invasion of Guantánamo Bay",
		"Kemp's Landing",
		"Liberation of Paris",
		"Long Run Massacre",
		"Sherman's March to the Sea",
		"United States occupation of Veracruz",
		"Vardar Offensive",
		"Operatie Kraai",
		"Operation Achilles",
		"Ndop prison break",
		"Niger raid",
		"Prayer Book Rebellion",
		"Recovery of Ré island",
		"Relief of Goes",
		"Revolt of Babylon (626 BC)",
		"Jacobite rising of 1689",
		"Sack of Rome",
		"Second Arab siege of Constantinople",
		"Sinking of CSS Alabama",
		"Skirmish at Blackwater Creek",
		"Great Stand on the Ugra River",
		"Gustafsen Lake Standoff",
		"Syria missile strikes",
		"United Arab Emirates takeover of Socotra",
		"Greater Poland Uprising",
	}
	allBattlesFound := true
	for _, n := range nn {
		if indexedBattlesNames[n] != n {
			t.Errorf("%q was not found in the list of scraped battles", n)
			allBattlesFound = false
		}
	}
	if allBattlesFound {
		t.Log("Contains battles found from different patterns")
	}

	cc := []string{
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
	for _, c := range cc {
		found := false
		for _, log := range loggerMock.Logs {
			if strings.Contains(log, c) {
				found = true
			}
		}
		if !found {
			t.Errorf("Scraper did not log results for %s", c)
		}
	}
	gotLogsAmount := loggerMock.Logs[len(loggerMock.Logs)-1]
	expectedLogsAmount := len(battlesList)
	if !strings.Contains(gotLogsAmount, strconv.Itoa(expectedLogsAmount)) {
		t.Errorf(
			"Expected the amounts log to contain the number %d, but instead logged %q",
			expectedLogsAmount,
			gotLogsAmount,
		)
	} else {
		t.Log("Logs for each list and with a final count of scraped battles")
	}
}
