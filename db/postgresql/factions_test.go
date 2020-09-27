package postgresql_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/db/postgresql"
	"github.com/sasalatart/batcoms/domain/factions"
	"github.com/sasalatart/batcoms/mocks"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestFactionsRepository(t *testing.T) {
	t.Run("CreateOne", func(t *testing.T) {
		mustSetupCreateOne := func(t *testing.T, mockUUID uuid.UUID, input factions.CreationInput) (*gorm.DB, *sql.DB, sqlmock.Sqlmock) {
			db, sqlDB, mock := mustSetupDB(t)
			mock.ExpectBegin()
			mock.ExpectQuery(`^INSERT INTO "factions" (.*)`).
				WithArgs(input.WikiID, input.URL, input.Name, input.Summary).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mockUUID))
			mock.ExpectCommit()
			return db, sqlDB, mock
		}
		t.Run("WithValidInput", func(t *testing.T) {
			mockUUID := uuid.NewV4()
			input := mocks.FactionCreationInput()
			db, sqlDB, mock := mustSetupCreateOne(t, mockUUID, input)
			defer sqlDB.Close()
			fs := postgresql.NewFactionsRepository(db)

			id, err := fs.CreateOne(input)
			require.NoError(t, err, "Creating faction with valid input")
			assert.Equal(t, mockUUID, id, "Should return the corresponding ID")
			assert.NoError(t, mock.ExpectationsWereMet(), "Not all SQL expectations were met")
		})
		t.Run("WithInvalidInput", func(t *testing.T) {
			input := mocks.FactionCreationInput()
			input.URL = "not-a-url"
			db, sqlDB, mock := mustSetupDB(t)
			defer sqlDB.Close()
			fs := postgresql.NewFactionsRepository(db)

			_, err := fs.CreateOne(input)
			require.Error(t, err, "Creating faction with invalid input")
			_, isValidationError := errors.Cause(err).(validator.ValidationErrors)
			assert.True(t, isValidationError, "Error should be a validator.ValidationErrors")
			assert.NoError(t, mock.ExpectationsWereMet(), "Not all SQL expectations were met")
		})
	})
}
