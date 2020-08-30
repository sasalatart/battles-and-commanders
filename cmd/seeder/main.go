package main

import (
	"log"

	"github.com/sasalatart/batcoms/config"
	"github.com/sasalatart/batcoms/services/seeder"
	"github.com/sasalatart/batcoms/store/postgresql"
	"github.com/sasalatart/batcoms/store/postgresql/schema"
	"github.com/spf13/viper"
)

func init() {
	config.Setup()
}

func main() {
	db := postgresql.Connect()
	defer db.Close()

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;`)
	schemas := []interface{}{&schema.BattleCommanderFaction{}, &schema.BattleFaction{}, &schema.BattleCommander{}, &schema.Faction{}, &schema.Commander{}, &schema.Battle{}}
	for _, s := range schemas {
		db.DropTableIfExists(s)
		db.AutoMigrate(s)
	}
	db.Model(&schema.BattleFaction{}).AddForeignKey("battle_id", "battles(id)", "CASCADE", "CASCADE")
	db.Model(&schema.BattleFaction{}).AddForeignKey("faction_id", "factions(id)", "CASCADE", "CASCADE")
	db.Model(&schema.BattleCommander{}).AddForeignKey("battle_id", "battles(id)", "CASCADE", "CASCADE")
	db.Model(&schema.BattleCommander{}).AddForeignKey("commander_id", "commanders(id)", "CASCADE", "CASCADE")
	db.Model(&schema.BattleCommanderFaction{}).AddForeignKey("battle_id", "battles(id)", "CASCADE", "CASCADE")
	db.Model(&schema.BattleCommanderFaction{}).AddForeignKey("commander_id", "commanders(id)", "CASCADE", "CASCADE")
	db.Model(&schema.BattleCommanderFaction{}).AddForeignKey("faction_id", "factions(id)", "CASCADE", "CASCADE")

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
