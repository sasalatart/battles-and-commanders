package postgresql

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/db/postgresql/schema"
	"github.com/spf13/viper"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
func Connect(c *ConnectionConfig) (*gorm.DB, *sql.DB) {
	if c == nil {
		c = defaultConfig()
	}

	handleError := func(err error) {
		if err != nil {
			panic(errors.Wrap(err, "Unable to connect to database"))
		}
	}
	dsn := postgres.Open(fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Name,
		c.Pass,
	))
	db, err := gorm.Open(dsn, &gorm.Config{})
	handleError(err)
	sqlDB, err := db.DB()
	handleError(err)
	return db, sqlDB
}

// Reset drops all existing tables, and automigrates them again
func Reset(db *gorm.DB) {
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;`)
	schemas := []interface{}{
		&schema.BattleCommanderFaction{},
		&schema.BattleFaction{},
		&schema.BattleCommander{},
		&schema.Faction{},
		&schema.Commander{},
		&schema.Battle{},
	}
	db.Migrator().DropTable(schemas...)
	db.AutoMigrate(schemas...)

	db.Exec(`CREATE INDEX ts_factions_name_idx ON factions USING GIST(to_tsvector('english', name));`)
	db.Exec(`CREATE INDEX ts_factions_summary_idx ON factions USING GIST(to_tsvector('english', summary));`)

	db.Exec(`CREATE INDEX ts_commanders_name_idx ON commanders USING GIST(to_tsvector('english', name));`)
	db.Exec(`CREATE INDEX ts_commanders_summary_idx ON commanders USING GIST(to_tsvector('english', summary));`)

	db.Exec(`CREATE INDEX ts_battles_name_idx ON battles USING GIST(to_tsvector('english', name));`)
	db.Exec(`CREATE INDEX ts_battles_summary_idx ON battles USING GIST(to_tsvector('english', summary));`)
	db.Exec(`CREATE INDEX ts_battles_place_idx ON battles USING GIST(to_tsvector('english', place));`)
	db.Exec(`CREATE INDEX ts_battles_result_idx ON battles USING GIST(to_tsvector('english', result));`)
}

func fromJSON(data datatypes.JSON, storeTo interface{}) error {
	parsed, err := data.MarshalJSON()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(parsed, storeTo); err != nil {
		return err
	}
	return nil
}

func paginate(db *gorm.DB, page, perPage int) *gorm.DB {
	return db.Offset((page - 1) * perPage).Limit(perPage)
}

func ts(db *gorm.DB, attribute, value string) *gorm.DB {
	if value == "" {
		return db
	}
	return db.Where(fmt.Sprintf("to_tsvector('english', %s) @@ phraseto_tsquery(?)", attribute), value)
}

const perPage = 50
