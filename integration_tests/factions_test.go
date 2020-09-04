package integration_test

import (
	"net/http"
	"testing"

	"github.com/sasalatart/batcoms/http/httptest"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestFactionsEndpoints(t *testing.T) {
	t.Run("GET /factions/:factionID", func(t *testing.T) {
		t.Parallel()

		route := func(id string) string {
			return URL("/factions/" + id)
		}

		t.Run("ValidPersistedUUID", func(t *testing.T) {
			expectedFaction := FirstFrenchEmpire(t)
			res, err := http.Get(route(expectedFaction.ID.String()))
			require.NoError(t, err, "Requesting First French Empire")
			defer res.Body.Close()
			httptest.AssertJSONFaction(t, res, expectedFaction)
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
