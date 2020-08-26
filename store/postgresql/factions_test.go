package postgresql_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/mocks"
	"github.com/sasalatart/batcoms/store/postgresql"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFactionsStore(t *testing.T) {
	t.Run("CreateOne", func(t *testing.T) {
		mustSetupCreateOne := func(t *testing.T, mockUUID uuid.UUID, input domain.CreateFactionInput, executes bool) (*gorm.DB, sqlmock.Sqlmock) {
			db, mock := mustSetupDB(t)
			if !executes {
				return db, mock
			}
			mock.ExpectBegin()
			mock.ExpectQuery(`^INSERT INTO "factions" (.*)`).
				WithArgs(input.WikiID, input.URL, input.Name, input.Summary).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mockUUID))
			mock.ExpectQuery(`^SELECT "id" FROM "factions"`).
				WithArgs(mockUUID).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mockUUID))
			mock.ExpectCommit()
			return db, mock
		}
		t.Run("WithValidInput", func(t *testing.T) {
			mockUUID := uuid.NewV4()
			input := mocks.CreateFactionInput()
			db, mock := mustSetupCreateOne(t, mockUUID, input, true)
			defer db.Close()
			fs := postgresql.NewFactionsDataStore(db)
			id, err := fs.CreateOne(input)
			require.NoError(t, err, "Creating faction")
			assert.Equal(t, mockUUID, id, "ID of created faction")
			assertMeetsExpectations(t, mock)
		})
		t.Run("WithInvalidInput", func(t *testing.T) {
			mockUUID := uuid.NewV4()
			input := mocks.CreateFactionInput()
			input.URL = "not-a-url"
			db, mock := mustSetupCreateOne(t, mockUUID, input, false)
			defer db.Close()
			fs := postgresql.NewFactionsDataStore(db)
			_, err := fs.CreateOne(input)
			require.Error(t, err, "Creating faction")
			if _, isValidationError := errors.Cause(err).(validator.ValidationErrors); !isValidationError {
				t.Error("Expected error to be a validation error, but it was not")
			}
			assertMeetsExpectations(t, mock)
		})
	})
}
