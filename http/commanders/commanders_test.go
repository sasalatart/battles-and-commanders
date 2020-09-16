package commanders_test

import (
	"fmt"
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
	appWithStoreMocks := func() (*fiber.App, *mocks.FactionsDataStore, *mocks.CommandersDataStore) {
		factionsStoreMock := new(mocks.FactionsDataStore)
		commandersStoreMock := new(mocks.CommandersDataStore)
		app := batcomshttp.Setup(factionsStoreMock, commandersStoreMock, new(mocks.BattlesDataStore), true)
		return app, factionsStoreMock, commandersStoreMock
	}

	t.Run("GET /commanders/:commanderID", func(t *testing.T) {
		t.Parallel()

		t.Run("ValidPersistedUUID", func(t *testing.T) {
			commanderMock := mocks.Commander()
			app, _, commandersStoreMock := appWithStoreMocks()
			commandersStoreMock.On("FindOne", domain.Commander{
				ID: commanderMock.ID,
			}).Return(commanderMock, nil)

			httptest.AssertFiberGET(t, app, "/commanders/"+commanderMock.ID.String(), http.StatusOK, func(res *http.Response) {
				httptest.AssertJSONCommander(t, res, commanderMock)
			})
			commandersStoreMock.AssertExpectations(t)
		})

		t.Run("ValidNonPersistedUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			app, _, commandersStoreMock := appWithStoreMocks()
			commandersStoreMock.On("FindOne", domain.Commander{
				ID: uuid,
			}).Return(domain.Commander{}, store.ErrNotFound)

			httptest.AssertFailedFiberGET(t, app, "/commanders/"+uuid.String(), *fiber.ErrNotFound)
			commandersStoreMock.AssertExpectations(t)
		})

		t.Run("InvalidUUID", func(t *testing.T) {
			invalidUUID := "invalid-uuid"
			app, _, commandersStoreMock := appWithStoreMocks()

			httptest.AssertFailedFiberGET(t, app, "/commanders/"+invalidUUID, *fiber.ErrBadRequest)
			commandersStoreMock.AssertNotCalled(t, "FindOne")
		})
	})

	t.Run("GET /commanders", func(t *testing.T) {
		t.Parallel()

		var page uint = 2
		var url = fmt.Sprintf("/commanders?page=%d", page)

		commandersMock := []domain.Commander{mocks.Commander(), mocks.Commander2()}
		expectedPages := 3
		app, _, commandersStoreMock := appWithStoreMocks()
		commandersStoreMock.On("FindMany", store.CommandersQuery{}, page).
			Return(commandersMock, expectedPages, nil)

		httptest.AssertFiberGET(t, app, url, http.StatusOK, func(res *http.Response) {
			httptest.AssertHeaderPages(t, res, uint(expectedPages))
			httptest.AssertJSONCommanders(t, res, commandersMock)
		})
		commandersStoreMock.AssertExpectations(t)
	})

	t.Run("GET /factions/:factionID/commanders", func(t *testing.T) {
		t.Parallel()

		var page uint = 2
		buildURL := func(factionID string) string {
			return fmt.Sprintf("/factions/%s/commanders?page=%d", factionID, page)
		}

		t.Run("ValidPersistedFactionUUID", func(t *testing.T) {
			factionMock := mocks.Faction()
			commandersMock := []domain.Commander{mocks.Commander()}
			expectedPages := 3
			app, factionsStoreMock, commandersStoreMock := appWithStoreMocks()
			factionsStoreMock.On("FindOne", domain.Faction{
				ID: factionMock.ID,
			}).Return(factionMock, nil)
			commandersStoreMock.On("FindMany", store.CommandersQuery{
				FactionID: factionMock.ID,
			}, page).Return(commandersMock, expectedPages, nil)

			httptest.AssertFiberGET(t, app, buildURL(factionMock.ID.String()), http.StatusOK, func(res *http.Response) {
				httptest.AssertHeaderPages(t, res, uint(expectedPages))
				httptest.AssertJSONCommanders(t, res, commandersMock)
			})
			factionsStoreMock.AssertExpectations(t)
			commandersStoreMock.AssertExpectations(t)
		})

		t.Run("ValidNonPersistedFactionUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			app, factionsStoreMock, commandersStoreMock := appWithStoreMocks()
			factionsStoreMock.On("FindOne", domain.Faction{
				ID: uuid,
			}).Return(domain.Faction{}, store.ErrNotFound)

			httptest.AssertFailedFiberGET(t, app, buildURL(uuid.String()), *fiber.ErrNotFound)
			factionsStoreMock.AssertExpectations(t)
			commandersStoreMock.AssertNotCalled(t, "FindMany")
		})

		t.Run("InvalidFactionUUID", func(t *testing.T) {
			invalidUUID := "invalid-uuid"
			app, factionsStoreMock, commandersStoreMock := appWithStoreMocks()

			httptest.AssertFailedFiberGET(t, app, buildURL(invalidUUID), *fiber.ErrBadRequest)
			factionsStoreMock.AssertNotCalled(t, "FindOne")
			commandersStoreMock.AssertNotCalled(t, "FindMany")
		})
	})
}
