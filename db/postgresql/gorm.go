package postgresql

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres drivers
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/db/postgresql/schema"
	"github.com/spf13/viper"
)

// ConnectionConfig is a struct that contains database connection configuration options
type ConnectionConfig struct {
	Host string
	Port string
	Name string
	User string
	Pass string
}

func defaultConfig() *ConnectionConfig {
	return &ConnectionConfig{
		Host: viper.GetString("PSQL_HOST"),
		Port: viper.GetString("PSQL_PORT"),
		Name: viper.GetString("PSQL_NAME"),
		User: viper.GetString("PSQL_USER"),
		Pass: viper.GetString("PSQL_PASS"),
	}
}

// DefaultTestConfig exposes a ConnectionConfig set up to work with the default test environment
func DefaultTestConfig() *ConnectionConfig {
	c := defaultConfig()
	c.Name = viper.GetString("PSQL_NAME_TEST")
	return c
}

// Connect establishes a database connection to the PostgreSQL instance
func Connect(c *ConnectionConfig) *gorm.DB {
	if c == nil {
		c = defaultConfig()
	}
	db, err := gorm.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Name,
		c.Pass,
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

	db.Exec(`CREATE INDEX ts_factions_name_idx ON factions USING GIST(to_tsvector('english', name));`)
	db.Exec(`CREATE INDEX ts_factions_summary_idx ON factions USING GIST(to_tsvector('english', summary));`)
	db.Exec(`CREATE INDEX ts_commanders_name_idx ON commanders USING GIST(to_tsvector('english', name));`)
	db.Exec(`CREATE INDEX ts_commanders_summary_idx ON commanders USING GIST(to_tsvector('english', summary));`)
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

const perPage = 50