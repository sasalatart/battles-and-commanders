package integration_test

import (
	"net/http"
	"testing"

	"github.com/sasalatart/batcoms/http/httptest"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestBattlesEndpoints(t *testing.T) {
	t.Run("GET /battles/:battleID", func(t *testing.T) {
		t.Parallel()

		route := func(id string) string {
			return URL("/battles/" + id)
		}

		t.Run("ValidPersistedUUID", func(t *testing.T) {
			expectedBattle := BattleOfAusterlitz(t)
			res, err := http.Get(route(expectedBattle.ID.String()))
			require.NoError(t, err, "Requesting Battle of Austerlitz")
			defer res.Body.Close()
			httptest.AssertJSONBattle(t, res, expectedBattle)
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
