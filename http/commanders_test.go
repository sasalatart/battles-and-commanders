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

func TestCommandersRoutes(t *testing.T) {
	t.Run("FindOne", func(t *testing.T) {
		t.Parallel()
		assertFindOne := func(app *fiber.App, route string, expectedResponse response) {
			t.Helper()
			res := mustGet(t, app, route)
			assert.Equalf(t, expectedResponse.status, res.StatusCode, "HTTP status for %q", route)
			if expectedResponse.errorMessage != "" {
				assertErrorMessage(t, res, expectedResponse.errorMessage)
				return
			}
			commanderFromBody := &domain.Commander{}
			err := json.NewDecoder(res.Body).Decode(commanderFromBody)
			require.NoError(t, err, "Decoding body into commander struct")
			expectedCommander := expectedResponse.body.(domain.Commander)
			assert.True(t, assert.ObjectsAreEqual(expectedCommander, *commanderFromBody), "Comparing body with expected commander")
		}
		t.Run("ValidPersistedUUID", func(t *testing.T) {
			commanderMock := mocks.Commander()
			commandersStoreMock := &mocks.CommandersStore{}
			commandersStoreMock.On("FindOne", domain.Commander{
				ID: commanderMock.ID,
			}).Return(commanderMock, nil)
			app := batcomshttp.Setup(new(mocks.FactionsStore), commandersStoreMock, true)
			expectedResponse := response{
				status: http.StatusOK,
				body:   commanderMock,
			}
			assertFindOne(app, fmt.Sprintf("/commanders/%s", commanderMock.ID), expectedResponse)
			commandersStoreMock.AssertExpectations(t)
		})
		t.Run("ValidNonPersistedUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			commandersStoreMock := &mocks.CommandersStore{}
			commandersStoreMock.On("FindOne", domain.Commander{
				ID: uuid,
			}).Return(domain.Commander{}, store.ErrNotFound)
			app := batcomshttp.Setup(new(mocks.FactionsStore), commandersStoreMock, true)
			expectedResponse := response{
				status:       http.StatusNotFound,
				errorMessage: fiber.ErrNotFound.Message,
			}
			assertFindOne(app, fmt.Sprintf("/commanders/%s", uuid), expectedResponse)
			commandersStoreMock.AssertExpectations(t)
		})
		t.Run("InvalidUUID", func(t *testing.T) {
			invalidUUID := "invalid-uuid"
			commandersStoreMock := &mocks.CommandersStore{}
			app := batcomshttp.Setup(new(mocks.FactionsStore), commandersStoreMock, true)
			expectedResponse := response{
				status:       http.StatusBadRequest,
				errorMessage: fiber.ErrBadRequest.Message,
			}
			assertFindOne(app, fmt.Sprintf("/commanders/%s", invalidUUID), expectedResponse)
			commandersStoreMock.AssertNotCalled(t, "FindOne")
		})
	})
}
