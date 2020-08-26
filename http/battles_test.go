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

func TestBattlesRoutes(t *testing.T) {
	t.Run("GET /battles/:battleID", func(t *testing.T) {
		t.Parallel()
		assertFindOne := func(app *fiber.App, route string, expectedResponse response) {
			t.Helper()
			res := mustGet(t, app, route)
			assert.Equalf(t, expectedResponse.status, res.StatusCode, "HTTP status for %q", route)
			if expectedResponse.errorMessage != "" {
				assertErrorMessage(t, res, expectedResponse.errorMessage)
				return
			}
			battleFromBody := new(domain.Battle)
			err := json.NewDecoder(res.Body).Decode(battleFromBody)
			require.NoError(t, err, "Decoding body into battle struct")
			expectedBattle := expectedResponse.body.(domain.Battle)
			assert.True(t, assert.ObjectsAreEqual(expectedBattle, *battleFromBody), "Comparing body with expected battle")
		}
		t.Run("ValidPersistedUUID", func(t *testing.T) {
			battleMock := mocks.Battle()
			battlesStoreMock := &mocks.BattlesStore{}
			battlesStoreMock.On("FindOne", domain.Battle{
				ID: battleMock.ID,
			}).Return(battleMock, nil)
			app := batcomshttp.Setup(new(mocks.FactionsStore), new(mocks.CommandersStore), battlesStoreMock, true)
			expectedResponse := response{
				status: http.StatusOK,
				body:   battleMock,
			}
			assertFindOne(app, fmt.Sprintf("/battles/%s", battleMock.ID), expectedResponse)
			battlesStoreMock.AssertExpectations(t)
		})
		t.Run("ValidNonPersistedUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			battlesStoreMock := &mocks.BattlesStore{}
			battlesStoreMock.On("FindOne", domain.Battle{
				ID: uuid,
			}).Return(domain.Battle{}, store.ErrNotFound)
			app := batcomshttp.Setup(new(mocks.FactionsStore), new(mocks.CommandersStore), battlesStoreMock, true)
			expectedResponse := response{
				status:       http.StatusNotFound,
				errorMessage: fiber.ErrNotFound.Message,
			}
			assertFindOne(app, fmt.Sprintf("/battles/%s", uuid), expectedResponse)
			battlesStoreMock.AssertExpectations(t)
		})
		t.Run("InvalidUUID", func(t *testing.T) {
			invalidUUID := "invalid-uuid"
			battlesStoreMock := &mocks.BattlesStore{}
			app := batcomshttp.Setup(new(mocks.FactionsStore), new(mocks.CommandersStore), battlesStoreMock, true)
			expectedResponse := response{
				status:       http.StatusBadRequest,
				errorMessage: fiber.ErrBadRequest.Message,
			}
			assertFindOne(app, fmt.Sprintf("/battles/%s", invalidUUID), expectedResponse)
			battlesStoreMock.AssertNotCalled(t, "FindOne")
		})
	})
}
