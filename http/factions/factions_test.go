package factions_test

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

func TestFactionsRoutes(t *testing.T) {
	t.Run("GET /factions/:factionID", func(t *testing.T) {
		t.Parallel()

		t.Run("ValidPersistedUUID", func(t *testing.T) {
			factionMock := mocks.Faction()
			factionsStoreMock := new(mocks.FactionsDataStore)
			factionsStoreMock.On("FindOne", domain.Faction{
				ID: factionMock.ID,
			}).Return(factionMock, nil)
			app := batcomshttp.Setup(factionsStoreMock, new(mocks.CommandersDataStore), new(mocks.BattlesDataStore), true)

			httptest.AssertFiberGET(t, app, "/factions/"+factionMock.ID.String(), http.StatusOK, func(res *http.Response) {
				httptest.AssertJSONFaction(t, res, factionMock)
			})
			factionsStoreMock.AssertExpectations(t)
		})

		t.Run("ValidNonPersistedUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			factionsStoreMock := new(mocks.FactionsDataStore)
			factionsStoreMock.On("FindOne", domain.Faction{
				ID: uuid,
			}).Return(domain.Faction{}, store.ErrNotFound)
			app := batcomshttp.Setup(factionsStoreMock, new(mocks.CommandersDataStore), new(mocks.BattlesDataStore), true)

			httptest.AssertFailedFiberGET(t, app, "/factions/"+uuid.String(), *fiber.ErrNotFound)
			factionsStoreMock.AssertExpectations(t)
		})

		t.Run("InvalidUUID", func(t *testing.T) {
			invalidUUID := "invalid-uuid"
			factionsStoreMock := new(mocks.FactionsDataStore)
			app := batcomshttp.Setup(factionsStoreMock, new(mocks.CommandersDataStore), new(mocks.BattlesDataStore), true)

			httptest.AssertFailedFiberGET(t, app, "/factions/"+invalidUUID, *fiber.ErrBadRequest)
			factionsStoreMock.AssertNotCalled(t, "FindOne")
		})
	})
}
