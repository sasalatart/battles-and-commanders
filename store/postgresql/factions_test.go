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
		mustSetupCreateOne := func(t *testing.T, mockUUID uuid.UUID, input domain.CreateFactionInput) (*gorm.DB, sqlmock.Sqlmock) {
			db, mock := mustSetupDB(t)
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
			db, mock := mustSetupCreateOne(t, mockUUID, input)
			defer db.Close()
			fs := postgresql.NewFactionsDataStore(db)

			id, err := fs.CreateOne(input)
			require.NoError(t, err, "Creating faction with valid input")
			assert.Equal(t, mockUUID, id, "Should return the corresponding ID")
			assert.NoError(t, mock.ExpectationsWereMet(), "Not all SQL expectations were met")
		})
		t.Run("WithInvalidInput", func(t *testing.T) {
			input := mocks.CreateFactionInput()
			input.URL = "not-a-url"
			db, mock := mustSetupDB(t)
			defer db.Close()
			fs := postgresql.NewFactionsDataStore(db)

			_, err := fs.CreateOne(input)
			require.Error(t, err, "Creating faction with invalid input")
			_, isValidationError := errors.Cause(err).(validator.ValidationErrors)
			assert.True(t, isValidationError, "Error should be a validator.ValidationErrors")
			assert.NoError(t, mock.ExpectationsWereMet(), "Not all SQL expectations were met")
		})
	})
}
