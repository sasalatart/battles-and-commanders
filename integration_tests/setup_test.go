package integration_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/sasalatart/batcoms/config"
	"github.com/sasalatart/batcoms/db/postgresql"
	"github.com/sasalatart/batcoms/db/seeder"
	"github.com/sasalatart/batcoms/domain/battles"
	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/domain/factions"
	"github.com/sasalatart/batcoms/pkg/io/json"
	"github.com/sasalatart/batcoms/pkg/logger"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var db *gorm.DB
var sqlDB *sql.DB
var battlesRepo *postgresql.BattlesRepository
var factionsRepo *postgresql.FactionsRepository
var commandersRepo *postgresql.CommandersRepository

func init() {
	config.Setup()
	db, sqlDB = postgresql.Connect(postgresql.DefaultTestConfig())
	battlesRepo = postgresql.NewBattlesRepository(db)
	factionsRepo = postgresql.NewFactionsRepository(db)
	commandersRepo = postgresql.NewCommandersRepository(db)
}

func TestMain(m *testing.M) {
	dataFileName := viper.GetString("TEST_DATA")
	importedData := new(seeder.ImportedData)
	if err := json.Import(dataFileName, importedData); err != nil {
		log.Fatalf("Error importing data: %s\n", err)
	}

	postgresql.Reset(db)
	seeder.Seed(importedData, factionsRepo, commandersRepo, battlesRepo, logger.NewDiscard())

	code := m.Run()
	sqlDB.Close()
	os.Exit(code)
}

func URL(route string) string {
	return "http://localhost:" + viper.GetString("PORT_TEST") + route
}

func BattleOfAusterlitz(t *testing.T) battles.Battle {
	t.Helper()
	return requireBattle(t, "Battle of Austerlitz")
}

func BattleOfArcole(t *testing.T) battles.Battle {
	t.Helper()
	return requireBattle(t, "Battle of Arcole")
}

func BattleOfLodi(t *testing.T) battles.Battle {
	t.Helper()
	return requireBattle(t, "Battle of Lodi")
}

func BattleOfMegiddo(t *testing.T) battles.Battle {
	t.Helper()
	return requireBattle(t, "Battle of Megiddo (15th century BC)")
}

func FirstFrenchEmpire(t *testing.T) factions.Faction {
	t.Helper()
	return requireFaction(t, "First French Empire")
}

func FrenchFirstRepublic(t *testing.T) factions.Faction {
	t.Helper()
	return requireFaction(t, "French First Republic")
}

func AustrianEmpire(t *testing.T) factions.Faction {
	t.Helper()
	return requireFaction(t, "Austrian Empire")
}

func HabsburgMonarchy(t *testing.T) factions.Faction {
	t.Helper()
	return requireFaction(t, "Habsburg Monarchy")
}

func RussianEmpire(t *testing.T) factions.Faction {
	t.Helper()
	return requireFaction(t, "Russian Empire")
}

func NewKingdomOfEgypt(t *testing.T) factions.Faction {
	t.Helper()
	return requireFaction(t, "New Kingdom of Egypt")
}

func Canaan(t *testing.T) factions.Faction {
	t.Helper()
	return requireFaction(t, "Canaan")
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

func JozsefAlvinczi(t *testing.T) commanders.Commander {
	t.Helper()
	return requireCommander(t, "József Alvinczi")
}

func JohannPeterBeaulieu(t *testing.T) commanders.Commander {
	t.Helper()
	return requireCommander(t, "Johann Peter Beaulieu")
}

func KarlPhilippSebottendorf(t *testing.T) commanders.Commander {
	t.Helper()
	return requireCommander(t, "Karl Philipp Sebottendorf")
}

func AlexanderI(t *testing.T) commanders.Commander {
	t.Helper()
	return requireCommander(t, "Alexander I of Russia")
}

func MikhailKutuzov(t *testing.T) commanders.Commander {
	t.Helper()
	return requireCommander(t, "Mikhail Kutuzov")
}

func ThutmoseIII(t *testing.T) commanders.Commander {
	t.Helper()
	return requireCommander(t, "Thutmose III")
}

func requireNoError(t *testing.T, err error, name string) {
	t.Helper()
	require.NoErrorf(t, err, "UNEXPECTED FROM SEEDER: Searching for %q", name)
}

func requireFaction(t *testing.T, factionName string) factions.Faction {
	t.Helper()
	faction, err := factionsRepo.FindOne(factions.FindOneQuery{Name: factionName})
	requireNoError(t, err, factionName)
	return faction
}

func requireCommander(t *testing.T, commanderName string) commanders.Commander {
	t.Helper()
	commander, err := commandersRepo.FindOne(commanders.FindOneQuery{Name: commanderName})
	requireNoError(t, err, commanderName)
	return commander
}

func requireBattle(t *testing.T, battleName string) battles.Battle {
	t.Helper()
	battle, err := battlesRepo.FindOne(battles.FindOneQuery{Name: battleName})
	requireNoError(t, err, battleName)
	return battle
}
