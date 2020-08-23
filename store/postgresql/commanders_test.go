package postgresql_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/mocks"
	"github.com/sasalatart/batcoms/store"
	"github.com/sasalatart/batcoms/store/postgresql"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandersStore(t *testing.T) {
	t.Run("FindOne", func(t *testing.T) {
		t.Run("WithPersistedUUID", func(t *testing.T) {
			db, mock := mustSetupDB(t)
			defer db.Close()
			commanderMock := mocks.Commander()
			mock.ExpectQuery(`^SELECT \* FROM "commanders"`).
				WithArgs(commanderMock.ID).
				WillReturnRows(sqlmock.NewRows(
					[]string{"id", "wiki_id", "url", "name", "summary"},
				).AddRow(commanderMock.ID, commanderMock.WikiID, commanderMock.URL, commanderMock.Name, commanderMock.Summary))
			fs := postgresql.NewCommandersDataStore(db)
			foundCommander, err := fs.FindOne(domain.Commander{
				ID: commanderMock.ID,
			})
			require.NoError(t, err, "Finding commander")
			assert.True(t, assert.ObjectsAreEqual(commanderMock, foundCommander), "Comparing found commander with queried one")
			assertMeetsExpectations(t, mock)
		})
		t.Run("WithNonPersistedUUID", func(t *testing.T) {
			db, mock := mustSetupDB(t)
			defer db.Close()
			uuid := uuid.NewV4()
			mock.ExpectQuery(`^SELECT \* FROM "commanders"`).
				WithArgs(uuid).
				WillReturnError(gorm.ErrRecordNotFound)
			fs := postgresql.NewCommandersDataStore(db)
			_, err := fs.FindOne(domain.Commander{
				ID: uuid,
			})
			require.Error(t, err, "Finding commander")
			assert.IsType(t, store.ErrNotFound, err, "Comparing store error")
			assertMeetsExpectations(t, mock)
		})
	})

	t.Run("CreateOne", func(t *testing.T) {
		mustSetupCreateOne := func(t *testing.T, mockUUID uuid.UUID, input domain.CreateCommanderInput, executes bool) (*gorm.DB, sqlmock.Sqlmock) {
			db, mock := mustSetupDB(t)
			if !executes {
				return db, mock
			}
			mock.ExpectBegin()
			mock.ExpectQuery(`^INSERT INTO "commanders" (.*)`).
				WithArgs(input.WikiID, input.URL, input.Name, input.Summary).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mockUUID))
			mock.ExpectQuery(`^SELECT "id" FROM "commanders"`).
				WithArgs(mockUUID).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mockUUID))
			mock.ExpectCommit()
			return db, mock
		}
		t.Run("WithValidInput", func(t *testing.T) {
			mockUUID := uuid.NewV4()
			input := mocks.CreateCommanderInput()
			db, mock := mustSetupCreateOne(t, mockUUID, input, true)
			defer db.Close()
			store := postgresql.NewCommandersDataStore(db)
			id, err := store.CreateOne(input)
			if err != nil {
				t.Errorf("Unexpected error creating commander: %v", err)
			}
			if id != mockUUID {
				t.Errorf("Expected to return an ID %s, but instead got %s", mockUUID, id)
			}
			assertMeetsExpectations(t, mock)
			if !t.Failed() {
				t.Log("Creates the commander in the database")
			}
		})
		t.Run("WithInvalidInput", func(t *testing.T) {
			mockUUID := uuid.NewV4()
			input := mocks.CreateCommanderInput()
			input.URL = "not-a-url"
			db, mock := mustSetupCreateOne(t, mockUUID, input, false)
			defer db.Close()
			store := postgresql.NewCommandersDataStore(db)
			_, err := store.CreateOne(input)
			if err == nil {
				t.Error("Expected error when creating commander, but got none")
			}
			if _, isValidationError := errors.Cause(err).(validator.ValidationErrors); !isValidationError {
				t.Error("Expected error to be a validation error, but it was not")
			}
			assertMeetsExpectations(t, mock)
			if !t.Failed() {
				t.Log("Fails validation and does not create the commander in the database")
			}
		})
	})
}
