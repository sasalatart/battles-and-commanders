package postgresql

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres drivers
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/store/postgresql/schema"
	"github.com/spf13/viper"
)

// Connect establishes a database connection to the PostgreSQL instance
func Connect() *gorm.DB {
	db, err := gorm.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		viper.GetString("PSQL_HOST"),
		viper.GetString("PSQL_PORT"),
		viper.GetString("PSQL_USER"),
		viper.GetString("PSQL_NAME"),
		viper.GetString("PSQL_PASS"),
	))
	if err != nil {
		panic(errors.Wrap(err, "Unable to connect to database"))
	}
	return db
}

// Reset drops all existing tables, and automigrates them again
func Reset(db *gorm.DB) {
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
}

func fromJSONB(data postgres.Jsonb, storeTo interface{}) error {
	parsed, err := data.MarshalJSON()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(parsed, storeTo); err != nil {
		return err
	}
	return nil
}
