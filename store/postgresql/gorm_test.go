package postgresql_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func mustSetupDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	handleStubDBError(t, err)
	gormDB, err := gorm.Open("postgres", db)
	gormDB.LogMode(true)
	handleStubDBError(t, err)
	return gormDB, mock
}

func handleStubDBError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Unexpected error when opening stub database: %s", err)
	}
}

func assertMeetsExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not all SQL expectations were met: %s", err)
	}
}
