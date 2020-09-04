package integration_test

import (
	"net/http"
	"testing"

	"github.com/sasalatart/batcoms/http/httptest"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestCommandersEndpoints(t *testing.T) {
	t.Run("GET /commanders/:commanderID", func(t *testing.T) {
		t.Parallel()

		route := func(id string) string {
			return URL("/commanders/" + id)
		}

		t.Run("ValidPersistedUUID", func(t *testing.T) {
			expectedCommander := Napoleon(t)
			res, err := http.Get(route(expectedCommander.ID.String()))
			require.NoError(t, err, "Requesting Napoleon")
			defer res.Body.Close()
			httptest.AssertJSONCommander(t, res, expectedCommander)
		})

		t.Run("ValidNonPersistedUUID", func(t *testing.T) {
			url := route(uuid.NewV4().String())
			httptest.AssertFailedGET(t, url, http.StatusNotFound, "Not Found")
		})

		t.Run("InvalidUUID", func(t *testing.T) {
			url := route("invalid-uuid")
			httptest.AssertFailedGET(t, url, http.StatusBadRequest, "Bad Request")
		})
	})
}
