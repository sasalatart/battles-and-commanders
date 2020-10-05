package postgresql_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/db/postgresql"
	"github.com/sasalatart/batcoms/domain/battles"
	"github.com/sasalatart/batcoms/mocks"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func TestBattlesRepository(t *testing.T) {
	t.Run("CreateOne", func(t *testing.T) {
		mustSetupCreateOne := func(t *testing.T, mockUUID uuid.UUID, input battles.CreationInput) (*gorm.DB, *sql.DB, sqlmock.Sqlmock) {
			t.Helper()
			db, sqlDB, mock := mustSetupDB(t)

			strength, err := json.Marshal(input.Strength)
			require.NoError(t, err, "Stringifying strength")
			casualties, err := json.Marshal(input.Casualties)
			require.NoError(t, err, "Stringifying casualties")

			mock.ExpectBegin()
			mock.ExpectQuery(`^INSERT INTO "battles"`).
				WithArgs(
					input.WikiID,
					input.URL,
					input.Name,
					input.PartOf,
					input.Summary,
					input.StartDate.String(),
					input.StartDate.ToNum(),
					input.EndDate.String(),
					input.EndDate.ToNum(),
					input.Location.Place,
					input.Location.Latitude,
					input.Location.Longitude,
					input.Result,
					input.TerritorialChanges,
					datatypes.JSON(strength),
					datatypes.JSON(casualties),
				).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mockUUID))

			battleCommandersArgs := []driver.Value{}
			for i := 0; i < len(input.CommandersBySide.A)+len(input.CommandersBySide.B); i++ {
				battleCommandersArgs = append(battleCommandersArgs, mockUUID, sqlmock.AnyArg(), sqlmock.AnyArg())
			}
			mock.ExpectExec(`^INSERT INTO "battle_commanders"`).
				WithArgs(battleCommandersArgs...).
				WillReturnResult(sqlmock.NewResult(1, 1))

			battleFactionsArgs := []driver.Value{}
			for i := 0; i < len(input.FactionsBySide.A)+len(input.FactionsBySide.B); i++ {
				battleFactionsArgs = append(battleFactionsArgs, mockUUID, sqlmock.AnyArg(), sqlmock.AnyArg())
			}
			mock.ExpectExec(`^INSERT INTO "battle_factions"`).
				WithArgs(battleFactionsArgs...).
				WillReturnResult(sqlmock.NewResult(1, 1))

			battleCommanderFactionsArgs := []driver.Value{}
			for _, cIDs := range input.CommandersByFaction {
				for range cIDs {
					battleCommanderFactionsArgs = append(battleCommanderFactionsArgs, mockUUID, sqlmock.AnyArg(), sqlmock.AnyArg())
				}
			}
			mock.ExpectExec(`^INSERT INTO "battle_commander_factions"`).
				WithArgs(battleCommanderFactionsArgs...).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectCommit()
			return db, sqlDB, mock
		}

		t.Run("WithValidInput", func(t *testing.T) {
			mockUUID := uuid.NewV4()
			input := mocks.BattleCreationInput()
			db, sqlDB, mock := mustSetupCreateOne(t, mockUUID, input)
			defer sqlDB.Close()
			repo := postgresql.NewBattlesRepository(db)

			id, err := repo.CreateOne(input)
			require.NoError(t, err, "Creating battle with valid input")
			assert.Equal(t, mockUUID, id, "Should return the corresponding ID")
			assert.NoError(t, mock.ExpectationsWereMet(), "Not all SQL expectations were met")
		})

		t.Run("WithInvalidInput", func(t *testing.T) {
			input := mocks.BattleCreationInput()
			input.URL = "not-a-url"
			db, sqlDB, mock := mustSetupDB(t)
			defer sqlDB.Close()
			repo := postgresql.NewBattlesRepository(db)

			_, err := repo.CreateOne(input)
			require.Error(t, err, "Creating battle with invalid input")
			_, isValidationError := errors.Cause(err).(validator.ValidationErrors)
			assert.True(t, isValidationError, "Error should be a validator.ValidationErrors")
			assert.NoError(t, mock.ExpectationsWereMet(), "Not all SQL expectations were met")
		})
	})
}
