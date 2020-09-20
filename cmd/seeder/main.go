package main

import (
	"log"
	"os"

	"github.com/sasalatart/batcoms/config"
	"github.com/sasalatart/batcoms/db/postgresql"
	"github.com/sasalatart/batcoms/db/seeder"
	"github.com/sasalatart/batcoms/pkg/logger"
	"github.com/spf13/viper"
)

func init() {
	config.Setup()
}

func main() {
	db := postgresql.Connect(nil)
	defer db.Close()

	actorsFileName := viper.GetString("SCRAPER_RESULTS.ACTORS")
	battlesFileName := viper.GetString("SCRAPER_RESULTS.BATTLES")
	importedData := new(seeder.ImportedData)
	if err := seeder.JSONImport(importedData, actorsFileName, battlesFileName); err != nil {
		log.Fatalf("Failed seeding: %s", err)
	}

	postgresql.Reset(db)
	seeder.Seed(
		importedData,
		postgresql.NewFactionsRepository(db),
		postgresql.NewCommandersRepository(db),
		postgresql.NewBattlesRepository(db),
		logger.New(log.Writer(), os.Stderr),
	)
}
