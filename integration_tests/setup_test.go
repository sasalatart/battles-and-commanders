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
	battleOfAusterlitz, err := battlesDataStore.FindOne(domain.Battle{Name: "Battle of Austerlitz"})
	require.NoError(t, err, "UNEXPECTED FROM SEEDER: Searching for Battle of Austerlitz")
	return battleOfAusterlitz
}

func FirstFrenchEmpire(t *testing.T) domain.Faction {
	firstFrenchEmpire, err := factionsDataStore.FindOne(domain.Faction{Name: "First French Empire"})
	require.NoError(t, err, "UNEXPECTED FROM SEEDER: Searching for First French Empire")
	return firstFrenchEmpire
}

func Napoleon(t *testing.T) domain.Commander {
	napoleon, err := commandersDataStore.FindOne(domain.Commander{Name: "Napoleon"})
	require.NoError(t, err, "UNEXPECTED FROM SEEDER: Searching for Napoleon")
	return napoleon
}
