package battles_test

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

func TestBattlesRoutes(t *testing.T) {
	t.Run("GET /battles/:battleID", func(t *testing.T) {
		t.Parallel()

		t.Run("ValidPersistedUUID", func(t *testing.T) {
			battleMock := mocks.Battle()
			battlesStoreMock := new(mocks.BattlesDataStore)
			battlesStoreMock.On("FindOne", domain.Battle{
				ID: battleMock.ID,
			}).Return(battleMock, nil)
			app := batcomshttp.Setup(new(mocks.FactionsDataStore), new(mocks.CommandersDataStore), battlesStoreMock, true)

			httptest.AssertFiberGET(t, app, "/battles/"+battleMock.ID.String(), http.StatusOK, func(res *http.Response) {
				httptest.AssertJSONBattle(t, res, battleMock)
			})
			battlesStoreMock.AssertExpectations(t)
		})

		t.Run("ValidNonPersistedUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			battlesStoreMock := new(mocks.BattlesDataStore)
			battlesStoreMock.On("FindOne", domain.Battle{
				ID: uuid,
			}).Return(domain.Battle{}, store.ErrNotFound)
			app := batcomshttp.Setup(new(mocks.FactionsDataStore), new(mocks.CommandersDataStore), battlesStoreMock, true)

			httptest.AssertFailedFiberGET(t, app, "/battles/"+uuid.String(), *fiber.ErrNotFound)
			battlesStoreMock.AssertExpectations(t)
		})

		t.Run("InvalidUUID", func(t *testing.T) {
			invalidUUID := "invalid-uuid"
			battlesStoreMock := new(mocks.BattlesDataStore)
			app := batcomshttp.Setup(new(mocks.FactionsDataStore), new(mocks.CommandersDataStore), battlesStoreMock, true)

			httptest.AssertFailedFiberGET(t, app, "/battles/"+invalidUUID, *fiber.ErrBadRequest)
			battlesStoreMock.AssertNotCalled(t, "FindOne")
		})
	})
}
