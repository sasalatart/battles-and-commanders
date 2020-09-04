package commanders_test

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/domain"
	batcomshttp "github.com/sasalatart/batcoms/http"
	"github.com/sasalatart/batcoms/http/httptest"
	"github.com/sasalatart/batcoms/mocks"
	"github.com/sasalatart/batcoms/store"
	uuid "github.com/satori/go.uuid"
)

func TestCommandersRoutes(t *testing.T) {
	t.Run("GET /commanders/:commanderID", func(t *testing.T) {
		t.Parallel()

		t.Run("ValidPersistedUUID", func(t *testing.T) {
			commanderMock := mocks.Commander()
			commandersStoreMock := new(mocks.CommandersDataStore)
			commandersStoreMock.On("FindOne", domain.Commander{
				ID: commanderMock.ID,
			}).Return(commanderMock, nil)
			app := batcomshttp.Setup(new(mocks.FactionsDataStore), commandersStoreMock, new(mocks.BattlesDataStore), true)

			httptest.AssertFiberGET(t, app, "/commanders/"+commanderMock.ID.String(), http.StatusOK, func(res *http.Response) {
				httptest.AssertJSONCommander(t, res, commanderMock)
			})
			commandersStoreMock.AssertExpectations(t)
		})

		t.Run("ValidNonPersistedUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			commandersStoreMock := new(mocks.CommandersDataStore)
			commandersStoreMock.On("FindOne", domain.Commander{
				ID: uuid,
			}).Return(domain.Commander{}, store.ErrNotFound)
			app := batcomshttp.Setup(new(mocks.FactionsDataStore), commandersStoreMock, new(mocks.BattlesDataStore), true)

			httptest.AssertFailedFiberGET(t, app, "/commanders/"+uuid.String(), *fiber.ErrNotFound)
			commandersStoreMock.AssertExpectations(t)
		})

		t.Run("InvalidUUID", func(t *testing.T) {
			invalidUUID := "invalid-uuid"
			commandersStoreMock := new(mocks.CommandersDataStore)
			app := batcomshttp.Setup(new(mocks.FactionsDataStore), commandersStoreMock, new(mocks.BattlesDataStore), true)

			httptest.AssertFailedFiberGET(t, app, "/commanders/"+invalidUUID, *fiber.ErrBadRequest)
			commandersStoreMock.AssertNotCalled(t, "FindOne")
		})
	})
}
