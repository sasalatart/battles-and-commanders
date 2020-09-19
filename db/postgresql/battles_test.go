package postgresql_test

import (
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/db/postgresql"
	"github.com/sasalatart/batcoms/domain/battles"
	"github.com/sasalatart/batcoms/mocks"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBattlesRepository(t *testing.T) {
	t.Run("CreateOne", func(t *testing.T) {
		mustSetupCreateOne := func(t *testing.T, mockUUID uuid.UUID, input battles.CreationInput) (*gorm.DB, sqlmock.Sqlmock) {
			t.Helper()
			db, mock := mustSetupDB(t)

			strength, err := json.Marshal(input.Strength)
			require.NoError(t, err, "Stringifying strength")
			casualties, err := json.Marshal(input.Casualties)
			require.NoError(t, err, "Stringifying casualties")

			mock.ExpectBegin()
			mock.ExpectQuery(`^INSERT INTO "battles" (.*)`).
				WithArgs(
					input.WikiID,
					input.URL,
					input.Name,
					input.PartOf,
					input.Summary,
					input.StartDate,
					input.EndDate,
					input.Location.Place,
					input.Location.Latitude,
					input.Location.Longitude,
					input.Result,
					input.TerritorialChanges,
					postgres.Jsonb{RawMessage: strength},
					postgres.Jsonb{RawMessage: casualties},
				).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mockUUID))
			mock.ExpectQuery(`^SELECT "id" FROM "battles"`).
				WithArgs(mockUUID).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mockUUID))
			for i := 0; i < len(input.CommandersBySide.A)+len(input.CommandersBySide.B); i++ {
				mock.ExpectExec(`^UPDATE "battle_commanders"`).
					WithArgs(sqlmock.AnyArg(), mockUUID, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}
			for i := 0; i < len(input.FactionsBySide.A)+len(input.FactionsBySide.B); i++ {
				mock.ExpectExec(`^UPDATE "battle_factions"`).
					WithArgs(sqlmock.AnyArg(), mockUUID, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}
			for fID, cIDs := range input.CommandersByFaction {
				for _, cID := range cIDs {
					mock.ExpectQuery(`^SELECT \* FROM "battle_commander_factions"`).
						WithArgs(mockUUID, sqlmock.AnyArg(), sqlmock.AnyArg()).
						WillReturnRows(sqlmock.NewRows([]string{"battle_id", "commander_id", "faction_id"}).AddRow(mockUUID, cID, fID))
				}
			}
			mock.ExpectCommit()
			return db, mock
		}

		t.Run("WithValidInput", func(t *testing.T) {
			mockUUID := uuid.NewV4()
			input := mocks.BattleCreationInput()
			db, mock := mustSetupCreateOne(t, mockUUID, input)
			defer db.Close()
			repo := postgresql.NewBattlesRepository(db)

			id, err := repo.CreateOne(input)
			require.NoError(t, err, "Creating battle with valid input")
			assert.Equal(t, mockUUID, id, "Should return the corresponding ID")
			assert.NoError(t, mock.ExpectationsWereMet(), "Not all SQL expectations were met")
		})

		t.Run("WithInvalidInput", func(t *testing.T) {
			input := mocks.BattleCreationInput()
			input.URL = "not-a-url"
			db, mock := mustSetupDB(t)
			defer db.Close()
			repo := postgresql.NewBattlesRepository(db)

			_, err := repo.CreateOne(input)
			require.Error(t, err, "Creating battle with invalid input")
			_, isValidationError := errors.Cause(err).(validator.ValidationErrors)
			assert.True(t, isValidationError, "Error should be a validator.ValidationErrors")
			assert.NoError(t, mock.ExpectationsWereMet(), "Not all SQL expectations were met")
		})
	})
}
