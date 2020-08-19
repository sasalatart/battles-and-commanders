package postgresql_test

import (
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/mocks"
	"github.com/sasalatart/batcoms/store/postgresql"
	uuid "github.com/satori/go.uuid"
)

func TestBattlesStore(t *testing.T) {
	t.Run("CreateOne", func(t *testing.T) {
		mustSetupCreateOne := func(t *testing.T, mockUUID uuid.UUID, input domain.CreateBattleInput, executes bool) (*gorm.DB, sqlmock.Sqlmock) {
			t.Helper()
			_, gormDB, mock := mustSetupDB(t)
			if !executes {
				return gormDB, mock
			}
			strength, err := json.Marshal(input.Strength)
			if err != nil {
				t.Fatalf("Unable to stringify strength: %s", err)
			}
			casualties, err := json.Marshal(input.Casualties)
			if err != nil {
				t.Fatalf("Unable to stringify casualties: %s", err)
			}
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
			for range input.CommandersByFaction {
				mock.ExpectExec(`^UPDATE "battle_commander_factions"`).
					WithArgs(mockUUID, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}
			mock.ExpectCommit()
			return gormDB, mock
		}
		t.Run("WithValidInput", func(t *testing.T) {
			mockUUID := uuid.NewV4()
			input, err := mocks.CreateBattleInput(domain.CreateBattleInput{})
			db, mock := mustSetupCreateOne(t, mockUUID, input, true)
			defer db.Close()
			store := postgresql.NewBattlesDataStore(db)
			id, err := store.CreateOne(input)
			if err != nil {
				t.Errorf("Unexpected error creating battle: %v", err)
			}
			if id != mockUUID {
				t.Errorf("Expected to return an ID %s, but instead got %s", mockUUID, id)
			}
			assertMeetsExpectations(t, mock)
			if !t.Failed() {
				t.Log("Creates the battle and its corresponding relations in the database")
			}
		})
		t.Run("WithInvalidInput", func(t *testing.T) {
			mockUUID := uuid.NewV4()
			input, err := mocks.CreateBattleInput(domain.CreateBattleInput{URL: "not-a-url"})
			db, mock := mustSetupCreateOne(t, mockUUID, input, false)
			defer db.Close()
			store := postgresql.NewBattlesDataStore(db)
			_, err = store.CreateOne(input)
			if err == nil {
				t.Error("Expected error when creating battle, but got none")
			}
			if _, isValidationError := errors.Cause(err).(validator.ValidationErrors); !isValidationError {
				t.Error("Expected error to be a validation error, but it was not")
			}
			assertMeetsExpectations(t, mock)
			if !t.Failed() {
				t.Log("Fails validation and does not create the battle nor its relations in the database")
			}
		})
	})
}
