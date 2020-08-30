package http_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/domain"
	batcomshttp "github.com/sasalatart/batcoms/http"
	"github.com/sasalatart/batcoms/mocks"
	"github.com/sasalatart/batcoms/store"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFactionsRoutes(t *testing.T) {
	t.Run("GET /factions/:factionID", func(t *testing.T) {
		t.Parallel()
		assertFindOne := func(app *fiber.App, route string, expectedResponse response) {
			t.Helper()
			res := mustGet(t, app, route)
			assert.Equalf(t, expectedResponse.status, res.StatusCode, "HTTP status for %q", route)
			if expectedResponse.errorMessage != "" {
				assertErrorMessage(t, res, expectedResponse.errorMessage)
				return
			}
			factionFromBody := new(domain.Faction)
			err := json.NewDecoder(res.Body).Decode(factionFromBody)
			require.NoError(t, err, "Decoding body into faction struct")
			expectedFaction := expectedResponse.body.(domain.Faction)
			assert.True(t, assert.ObjectsAreEqual(expectedFaction, *factionFromBody), "Comparing body with expected faction")
		}
		t.Run("ValidPersistedUUID", func(t *testing.T) {
			factionMock := mocks.Faction()
			factionsStoreMock := new(mocks.FactionsDataStore)
			factionsStoreMock.On("FindOne", domain.Faction{
				ID: factionMock.ID,
			}).Return(factionMock, nil)
			app := batcomshttp.Setup(factionsStoreMock, new(mocks.CommandersDataStore), new(mocks.BattlesDataStore), true)
			expectedResponse := response{
				status: http.StatusOK,
				body:   factionMock,
			}
			assertFindOne(app, fmt.Sprintf("/factions/%s", factionMock.ID), expectedResponse)
			factionsStoreMock.AssertExpectations(t)
		})
		t.Run("ValidNonPersistedUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			factionsStoreMock := new(mocks.FactionsDataStore)
			factionsStoreMock.On("FindOne", domain.Faction{
				ID: uuid,
			}).Return(domain.Faction{}, store.ErrNotFound)
			app := batcomshttp.Setup(factionsStoreMock, new(mocks.CommandersDataStore), new(mocks.BattlesDataStore), true)
			expectedResponse := response{
				status:       http.StatusNotFound,
				errorMessage: fiber.ErrNotFound.Message,
			}
			assertFindOne(app, fmt.Sprintf("/factions/%s", uuid), expectedResponse)
			factionsStoreMock.AssertExpectations(t)
		})
		t.Run("InvalidUUID", func(t *testing.T) {
			invalidUUID := "invalid-uuid"
			factionsStoreMock := new(mocks.FactionsDataStore)
			app := batcomshttp.Setup(factionsStoreMock, new(mocks.CommandersDataStore), new(mocks.BattlesDataStore), true)
			expectedResponse := response{
				status:       http.StatusBadRequest,
				errorMessage: fiber.ErrBadRequest.Message,
			}
			assertFindOne(app, fmt.Sprintf("/factions/%s", invalidUUID), expectedResponse)
			factionsStoreMock.AssertNotCalled(t, "FindOne")
		})
	})
}
