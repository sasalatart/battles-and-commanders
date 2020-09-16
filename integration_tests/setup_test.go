package integration_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/sasalatart/batcoms/config"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/services/io"
	"github.com/sasalatart/batcoms/services/logger"
	"github.com/sasalatart/batcoms/services/seeder"
	"github.com/sasalatart/batcoms/store/postgresql"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

var db *gorm.DB
var battlesDataStore *postgresql.BattlesDataStore
var factionsDataStore *postgresql.FactionsDataStore
var commandersDataStore *postgresql.CommandersDataStore

func init() {
	config.Setup()
	db = postgresql.Connect(postgresql.TestConfig())
	battlesDataStore = postgresql.NewBattlesDataStore(db)
	factionsDataStore = postgresql.NewFactionsDataStore(db)
	commandersDataStore = postgresql.NewCommandersDataStore(db)
}

func TestMain(m *testing.M) {
	battlesFileName := viper.GetString("TEST_SEEDERS.BATTLES")
	participantsFileName := viper.GetString("TEST_SEEDERS.PARTICIPANTS")
	importedData := new(io.ImportedData)
	if err := seeder.JSONImport(battlesFileName, participantsFileName, importedData); err != nil {
		log.Fatalf("Failed seeding: %s", err)
	}
	seederService := seeder.Service{
		ImportedData:      importedData,
		FactionsCreator:   factionsDataStore,
		CommandersCreator: commandersDataStore,
		BattlesCreator:    battlesDataStore,
		Logger:            logger.New(ioutil.Discard, ioutil.Discard),
	}
	postgresql.Reset(db)
	seederService.Seed()

	code := m.Run()
	db.Close()
	os.Exit(code)
}

func URL(route string) string {
	return "http://localhost:" + viper.GetString("PORT_TEST") + route
}

func BattleOfAusterlitz(t *testing.T) domain.Battle {
	t.Helper()
	return requireBattle(t, "Battle of Austerlitz")
}

func FirstFrenchEmpire(t *testing.T) domain.Faction {
	t.Helper()
	return requireFaction(t, "First French Empire")
}

func AustrianEmpire(t *testing.T) domain.Faction {
	t.Helper()
	return requireFaction(t, "Austrian Empire")
}

func Napoleon(t *testing.T) domain.Commander {
	t.Helper()
	return requireCommander(t, "Napoleon")
}

func FrancisII(t *testing.T) domain.Commander {
	t.Helper()
	return requireCommander(t, "Francis II, Holy Roman Emperor")
}

func FranzVonWeyrother(t *testing.T) domain.Commander {
	t.Helper()
	return requireCommander(t, "Franz von Weyrother")
}

func AlexanderI(t *testing.T) domain.Commander {
	t.Helper()
	return requireCommander(t, "Alexander I of Russia")
}

func MikhailKutuzov(t *testing.T) domain.Commander {
	t.Helper()
	return requireCommander(t, "Mikhail Kutuzov")
}

func requireNoError(t *testing.T, err error, name string) {
	t.Helper()
	require.NoErrorf(t, err, "UNEXPECTED FROM SEEDER: Searching for %q", name)
}

func requireFaction(t *testing.T, factionName string) domain.Faction {
	t.Helper()
	faction, err := factionsDataStore.FindOne(domain.Faction{Name: factionName})
	requireNoError(t, err, factionName)
	return faction
}

func requireCommander(t *testing.T, commanderName string) domain.Commander {
	t.Helper()
	commander, err := commandersDataStore.FindOne(domain.Commander{Name: commanderName})
	requireNoError(t, err, commanderName)
	return commander
}

func requireBattle(t *testing.T, battleName string) domain.Battle {
	t.Helper()
	battle, err := battlesDataStore.FindOne(domain.Battle{Name: battleName})
	requireNoError(t, err, battleName)
	return battle
}
