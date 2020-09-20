package handlers_test

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/domain/battles"
	batcomshttp "github.com/sasalatart/batcoms/http"
	"github.com/sasalatart/batcoms/http/httptest"
	"github.com/sasalatart/batcoms/mocks"
	uuid "github.com/satori/go.uuid"
)

func TestBattlesHandlers(t *testing.T) {
	t.Run("GET /battles/:battleID", func(t *testing.T) {
		t.Parallel()

		t.Run("ValidPersistedUUID", func(t *testing.T) {
			battleMock := mocks.Battle()
			battlesRepoMock := new(mocks.BattlesRepository)
			battlesRepoMock.On("FindOne", battles.Battle{
				ID: battleMock.ID,
			}).Return(battleMock, nil)
			app := batcomshttp.Setup(new(mocks.FactionsRepository), new(mocks.CommandersRepository), battlesRepoMock, true)

			httptest.AssertFiberGET(t, app, "/battles/"+battleMock.ID.String(), http.StatusOK, func(res *http.Response) {
				battlesRepoMock.AssertExpectations(t)
				httptest.AssertJSONBattle(t, res, battleMock)
			})
		})

		t.Run("ValidNonPersistedUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			battlesRepoMock := new(mocks.BattlesRepository)
			battlesRepoMock.On("FindOne", battles.Battle{
				ID: uuid,
			}).Return(battles.Battle{}, domain.ErrNotFound)
			app := batcomshttp.Setup(new(mocks.FactionsRepository), new(mocks.CommandersRepository), battlesRepoMock, true)

			httptest.AssertFailedFiberGET(t, app, "/battles/"+uuid.String(), *fiber.ErrNotFound)
			battlesRepoMock.AssertExpectations(t)
		})

		t.Run("InvalidUUID", func(t *testing.T) {
			invalidUUID := "invalid-uuid"
			battlesRepoMock := new(mocks.BattlesRepository)
			app := batcomshttp.Setup(new(mocks.FactionsRepository), new(mocks.CommandersRepository), battlesRepoMock, true)

			httptest.AssertFailedFiberGET(t, app, "/battles/"+invalidUUID, *fiber.ErrBadRequest)
			battlesRepoMock.AssertNotCalled(t, "FindOne")
		})
	})
}
