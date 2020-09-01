package main

import (
	"log"

	"github.com/sasalatart/batcoms/config"
	"github.com/sasalatart/batcoms/services/seeder"
	"github.com/sasalatart/batcoms/store/postgresql"
	"github.com/spf13/viper"
)

func init() {
	config.Setup()
}

func main() {
	db := postgresql.Connect()
	defer db.Close()

	postgresql.Reset(db)
	battlesFileName := viper.GetString("SCRAPER_RESULTS.BATTLES")
	participantsFileName := viper.GetString("SCRAPER_RESULTS.PARTICIPANTS")
	importedData, err := seeder.JSONImport(battlesFileName, participantsFileName)
	if err != nil {
		log.Fatalf("Failed seeding: %s", err)
	}
	seederService := seeder.Service{
		ImportedData:      importedData,
		FactionsCreator:   postgresql.NewFactionsDataStore(db),
		CommandersCreator: postgresql.NewCommandersDataStore(db),
		BattlesCreator:    postgresql.NewBattlesDataStore(db),
		Logger:            log.Writer(),
	}
	seederService.Seed()
}
