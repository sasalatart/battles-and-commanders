package handlers_test

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/domain/factions"
	batcomshttp "github.com/sasalatart/batcoms/http"
	"github.com/sasalatart/batcoms/http/httptest"
	"github.com/sasalatart/batcoms/mocks"
	uuid "github.com/satori/go.uuid"
)

func TestFactionsHandlers(t *testing.T) {
	t.Run("GET /factions/:factionID", func(t *testing.T) {
		t.Parallel()

		t.Run("ValidPersistedUUID", func(t *testing.T) {
			factionMock := mocks.Faction()
			factionsRepoMock := new(mocks.FactionsRepository)
			factionsRepoMock.On("FindOne", factions.Faction{
				ID: factionMock.ID,
			}).Return(factionMock, nil)
			app := batcomshttp.Setup(factionsRepoMock, new(mocks.CommandersRepository), new(mocks.BattlesRepository), true)

			httptest.AssertFiberGET(t, app, "/factions/"+factionMock.ID.String(), http.StatusOK, func(res *http.Response) {
				factionsRepoMock.AssertExpectations(t)
				httptest.AssertJSONFaction(t, res, factionMock)
			})
		})

		t.Run("ValidNonPersistedUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			factionsRepoMock := new(mocks.FactionsRepository)
			factionsRepoMock.On("FindOne", factions.Faction{
				ID: uuid,
			}).Return(factions.Faction{}, domain.ErrNotFound)
			app := batcomshttp.Setup(factionsRepoMock, new(mocks.CommandersRepository), new(mocks.BattlesRepository), true)

			httptest.AssertFailedFiberGET(t, app, "/factions/"+uuid.String(), *fiber.ErrNotFound)
			factionsRepoMock.AssertExpectations(t)
		})

		t.Run("InvalidUUID", func(t *testing.T) {
			invalidUUID := "invalid-uuid"
			factionsRepoMock := new(mocks.FactionsRepository)
			app := batcomshttp.Setup(factionsRepoMock, new(mocks.CommandersRepository), new(mocks.BattlesRepository), true)

			httptest.AssertFailedFiberGET(t, app, "/factions/"+invalidUUID, *fiber.ErrBadRequest)
			factionsRepoMock.AssertNotCalled(t, "FindOne")
		})
	})
}
