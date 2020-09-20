package postgresql_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/db/postgresql"
	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/mocks"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandersRepository(t *testing.T) {
	t.Run("CreateOne", func(t *testing.T) {
		mustSetupCreateOne := func(t *testing.T, mockUUID uuid.UUID, input commanders.CreationInput) (*gorm.DB, sqlmock.Sqlmock) {
			db, mock := mustSetupDB(t)
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
			input := mocks.CommanderCreationInput()
			db, mock := mustSetupCreateOne(t, mockUUID, input)
			defer db.Close()
			repo := postgresql.NewCommandersRepository(db)

			id, err := repo.CreateOne(input)
			require.NoError(t, err, "Creating commander with valid input")
			assert.Equal(t, mockUUID, id, "Should return the corresponding ID")
			assert.NoError(t, mock.ExpectationsWereMet(), "Not all SQL expectations were met")
		})
		t.Run("WithInvalidInput", func(t *testing.T) {
			input := mocks.CommanderCreationInput()
			input.URL = "not-a-url"
			db, mock := mustSetupDB(t)
			defer db.Close()
			repo := postgresql.NewCommandersRepository(db)

			_, err := repo.CreateOne(input)
			require.Error(t, err, "Creating commander with invalid input")
			_, isValidationError := errors.Cause(err).(validator.ValidationErrors)
			assert.True(t, isValidationError, "Error should be a validator.ValidationErrors")
			assert.NoError(t, mock.ExpectationsWereMet(), "Not all SQL expectations were met")
		})
	})
}
