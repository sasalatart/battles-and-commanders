package integration_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/sasalatart/batcoms/config"
	"github.com/sasalatart/batcoms/db/postgresql"
	"github.com/sasalatart/batcoms/db/seeder"
	"github.com/sasalatart/batcoms/domain/battles"
	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/domain/factions"
	"github.com/sasalatart/batcoms/pkg/logger"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

var db *gorm.DB
var battlesRepo *postgresql.BattlesRepository
var factionsRepo *postgresql.FactionsRepository
var commandersRepo *postgresql.CommandersRepository

func init() {
	config.Setup()
	db = postgresql.Connect(postgresql.DefaultTestConfig())
	battlesRepo = postgresql.NewBattlesRepository(db)
	factionsRepo = postgresql.NewFactionsRepository(db)
	commandersRepo = postgresql.NewCommandersRepository(db)
}

func TestMain(m *testing.M) {
	actorsFileName := viper.GetString("TEST_SEEDERS.ACTORS")
	battlesFileName := viper.GetString("TEST_SEEDERS.BATTLES")
	importedData := new(seeder.ImportedData)
	if err := seeder.JSONImport(importedData, actorsFileName, battlesFileName); err != nil {
		log.Fatalf("Failed seeding: %s", err)
	}

	postgresql.Reset(db)
	seeder.Seed(importedData, factionsRepo, commandersRepo, battlesRepo, logger.New(ioutil.Discard, ioutil.Discard))

	code := m.Run()
	db.Close()
	os.Exit(code)
}

func URL(route string) string {
	return "http://localhost:" + viper.GetString("PORT_TEST") + route
}

func BattleOfAusterlitz(t *testing.T) battles.Battle {
	t.Helper()
	return requireBattle(t, "Battle of Austerlitz")
}

func FirstFrenchEmpire(t *testing.T) factions.Faction {
	t.Helper()
	return requireFaction(t, "First French Empire")
}

func AustrianEmpire(t *testing.T) factions.Faction {
	t.Helper()
	return requireFaction(t, "Austrian Empire")
}

func Napoleon(t *testing.T) commanders.Commander {
	t.Helper()
	return requireCommander(t, "Napoleon")
}

func FrancisII(t *testing.T) commanders.Commander {
	t.Helper()
	return requireCommander(t, "Francis II, Holy Roman Emperor")
}

func FranzVonWeyrother(t *testing.T) commanders.Commander {
	t.Helper()
	return requireCommander(t, "Franz von Weyrother")
}

func AlexanderI(t *testing.T) commanders.Commander {
	t.Helper()
	return requireCommander(t, "Alexander I of Russia")
}

func MikhailKutuzov(t *testing.T) commanders.Commander {
	t.Helper()
	return requireCommander(t, "Mikhail Kutuzov")
}

func requireNoError(t *testing.T, err error, name string) {
	t.Helper()
	require.NoErrorf(t, err, "UNEXPECTED FROM SEEDER: Searching for %q", name)
}

func requireFaction(t *testing.T, factionName string) factions.Faction {
	t.Helper()
	faction, err := factionsRepo.FindOne(factions.Faction{Name: factionName})
	requireNoError(t, err, factionName)
	return faction
}

func requireCommander(t *testing.T, commanderName string) commanders.Commander {
	t.Helper()
	commander, err := commandersRepo.FindOne(commanders.Commander{Name: commanderName})
	requireNoError(t, err, commanderName)
	return commander
}

func requireBattle(t *testing.T, battleName string) battles.Battle {
	t.Helper()
	battle, err := battlesRepo.FindOne(battles.Battle{Name: battleName})
	requireNoError(t, err, battleName)
	return battle
}
