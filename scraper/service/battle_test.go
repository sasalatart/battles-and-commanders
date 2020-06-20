package service_test

import (
	"reflect"
	"testing"

	"github.com/sasalatart/batcoms/mocks"
	"github.com/sasalatart/batcoms/scraper/domain"
	"github.com/sasalatart/batcoms/scraper/service"
	"github.com/sasalatart/batcoms/scraper/store"
)

func TestBattle(t *testing.T) {
	const battleURL = "https://en.wikipedia.org/wiki/Battle_of_Austerlitz"

	loggerMock := &mocks.Logger{}
	battlesStore := store.NewBattlesMem()
	participantsStore := store.NewParticipantsMem()
	service := service.NewScraper(battlesStore, participantsStore, loggerMock)

	battle, err := service.Battle(battleURL)
	if err != nil {
		t.Errorf("Unexpected service.Battle error: %s", err)
	}

	stringCases := []struct {
		attrName string
		got      string
		expected string
	}{
		{"URL", battle.URL, battleURL},
		{"Name", battle.Name, "Battle of Austerlitz"},
		{"PartOf", battle.PartOf, "Part of the War of the Third Coalition"},
		{"Description", battle.Description, "Decisive battle of the Napoleonic Wars"},
		{"Extract", battle.Extract, "The Battle of Austerlitz, also known as the Battle of the Three Emperors, was one of the most important and decisive engagements of the Napoleonic Wars. In what is widely regarded as the greatest victory achieved by Napoleon, the Grande Armée of France defeated a larger Russian and Austrian army led by Emperor Alexander I and Holy Roman Emperor Francis II. The battle occurred near the town of Austerlitz in the Austrian Empire. Austerlitz brought the War of the Third Coalition to a rapid end, with the Treaty of Pressburg signed by the Austrians later in the month. The battle is often cited as a tactical masterpiece, in the same league as other historic engagements like Cannae or Gaugamela."},
		{"Date", battle.Date, "2 December 1805"},
		{"Place", battle.Location.Place, "Austerlitz, Moravia, Austrian Empire (now Slavkov u Brna, Czech Republic)"},
		{"Latitude", battle.Location.Latitude, "49°8′N"},
		{"Longitude", battle.Location.Longitude, "16°46′E"},
		{"Result", battle.Result, "Decisive French victory. Treaty of Pressburg. Effective end of the Third Coalition"},
		{"TerritorialChanges", battle.TerritorialChanges, "Dissolution of the Holy Roman Empire and creation of the Confederation of the Rhine"},
	}

	for _, sc := range stringCases {
		if sc.got != sc.expected {
			t.Errorf("Expected battle %s to be %q, but instead got %q", sc.attrName, sc.expected, sc.got)
		}
	}

	expectedStrength := map[string]string{
		"A": "65,000–68,000 (not including III Corps)",
		"B": "84,000–95,000",
	}
	if battle.Strength.A != expectedStrength["A"] || battle.Strength.B != expectedStrength["B"] {
		t.Errorf("Expected battle Strength to be %+v, but instead got %+v", expectedStrength, battle.Strength)
	}

	expectedCasualties := map[string]string{
		"A": "1,305 dead. 6,991 wounded. 8,279 total casualties. 573 captured. 1 standard lost. Total: 9,000",
		"B": "16,000 dead or wounded. 20,000 captured. 186 guns lost. 45 standards lost. Total: 36,000",
	}
	if battle.Casualties.A != expectedCasualties["A"] || battle.Casualties.B != expectedCasualties["B"] {
		t.Errorf("Expected battle Casualties to be %+v, but instead got %+v", expectedCasualties, battle.Casualties)
	}

	expectedFactions := map[string][]int{
		"A": {21418258},
		"B": {20611504, 13277},
	}
	if !reflect.DeepEqual(battle.Factions.A, expectedFactions["A"]) || !reflect.DeepEqual(battle.Factions.B, expectedFactions["B"]) {
		t.Errorf("Expected battle Factions to be %+v, but instead got %+v", expectedFactions, battle.Factions)
	}

	type participantCase struct {
		attrName     string
		expectedName string
		id           int
	}

	assertParticipantIs := func(kind domain.ParticipantKind, c participantCase) {
		t.Helper()

		participant := participantsStore.Find(kind, c.id)
		if participant == nil {
			t.Fatalf("No participant found with id %d for %q", c.id, c.attrName)
		}
		got := participant.Name
		if c.expectedName != got {
			t.Errorf("Expected %s to have name %q, but instead got %q", c.attrName, c.expectedName, got)
		}
	}

	factionsCases := []participantCase{
		{"FactionA1", "First French Empire", expectedFactions["A"][0]},
		{"FactionB1", "Russian Empire", expectedFactions["B"][0]},
		{"FactionB2", "Holy Roman Empire", expectedFactions["B"][1]},
	}
	for _, bc := range factionsCases {
		assertParticipantIs(domain.FactionKind, bc)
	}

	expectedCommanders := map[string][]int{
		"A": {69880},
		"B": {27126603, 251000, 11551, 14092123},
	}
	if !reflect.DeepEqual(battle.Commanders.A, expectedCommanders["A"]) || !reflect.DeepEqual(battle.Commanders.B, expectedCommanders["B"]) {
		t.Errorf("Expected battle Commanders to be %+v, but instead got %+v", expectedCommanders, battle.Commanders)
	}

	commandersCases := []participantCase{
		{"CommanderA1", "Napoleon", expectedCommanders["A"][0]},
		{"CommanderB1", "Alexander I of Russia", expectedCommanders["B"][0]},
		{"CommanderB2", "Mikhail Kutuzov", expectedCommanders["B"][1]},
		{"CommanderB3", "Francis II, Holy Roman Emperor", expectedCommanders["B"][2]},
		{"CommanderB4", "Franz von Weyrother", expectedCommanders["B"][3]},
	}
	for _, cc := range commandersCases {
		assertParticipantIs(domain.CommanderKind, cc)
	}

	expectedCommandersGrouping := map[int][]int{
		13277:    {11551, 14092123},
		20611504: {27126603, 251000},
		21418258: {69880},
	}
	if !reflect.DeepEqual(battle.CommandersByFaction, expectedCommandersGrouping) {
		t.Errorf(
			"Expected commanders to be grouped as %+v, but instead got %+v",
			expectedCommandersGrouping,
			battle.CommandersByFaction,
		)
	}
}
