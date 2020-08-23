package http_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type response struct {
	status       int
	errorMessage string
	faction      domain.Faction
}

func TestFactionsRoutes(t *testing.T) {
	t.Run("FindOne", func(t *testing.T) {
		t.Parallel()
		assertFindOne := func(app *fiber.App, route string, expectedResponse response) {
			t.Helper()
			req, err := http.NewRequest("GET", route, nil)
			require.NoError(t, err, route)
			res, err := app.Test(req, -1)
			require.NoError(t, err, route)
			assert.Equalf(t, expectedResponse.status, res.StatusCode, "HTTP status for %q", route)
			if expectedResponse.errorMessage != "" {
				body, err := ioutil.ReadAll(res.Body)
				require.NoError(t, err, "Reading from body")
				assert.Equal(t, expectedResponse.errorMessage, string(body), "Comparing body with expected error message")
				return
			}
			factionFromBody := &domain.Faction{}
			err = json.NewDecoder(res.Body).Decode(factionFromBody)
			require.NoError(t, err, "Decoding body into faction struct")
			assert.True(t, assert.ObjectsAreEqual(expectedResponse.faction, *factionFromBody), "Comparing body with expected faction")
		}
		t.Run("ValidPersistedUUID", func(t *testing.T) {
			factionMock := mocks.Faction()
			factionsStoreMock := &mocks.FactionsStore{}
			factionsStoreMock.On("FindOne", domain.Faction{
				ID: factionMock.ID,
			}).Return(factionMock, nil)
			app := batcomshttp.Setup(factionsStoreMock, true)
			expectedResponse := response{
				status:  http.StatusOK,
				faction: factionMock,
			}
			assertFindOne(app, fmt.Sprintf("/factions/%s", factionMock.ID), expectedResponse)
			factionsStoreMock.AssertExpectations(t)
		})
		t.Run("ValidNonPersistedUUID", func(t *testing.T) {
			uuidMock := uuid.NewV4()
			factionsStoreMock := &mocks.FactionsStore{}
			factionsStoreMock.On("FindOne", domain.Faction{
				ID: uuidMock,
			}).Return(domain.Faction{}, store.ErrNotFound)
			app := batcomshttp.Setup(factionsStoreMock, true)
			expectedResponse := response{
				status:       http.StatusNotFound,
				errorMessage: fiber.ErrNotFound.Message,
			}
			assertFindOne(app, fmt.Sprintf("/factions/%s", uuidMock), expectedResponse)
			factionsStoreMock.AssertExpectations(t)
		})
		t.Run("InvalidUUID", func(t *testing.T) {
			invalidUUID := "invalid-uuid"
			factionsStoreMock := &mocks.FactionsStore{}
			app := batcomshttp.Setup(factionsStoreMock, true)
			expectedResponse := response{
				status:       http.StatusBadRequest,
				errorMessage: fiber.ErrBadRequest.Message,
			}
			assertFindOne(app, fmt.Sprintf("/factions/%s", invalidUUID), expectedResponse)
			factionsStoreMock.AssertNotCalled(t, "FindOne")
		})
	})
}
