package integration_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/sasalatart/batcoms/domain"
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

	t.Run("GET /commanders", func(t *testing.T) {
		t.Parallel()

		var expectedPages uint = 1
		expectedCommanders := []domain.Commander{
			Napoleon(t),
			MikhailKutuzov(t),
			FranzVonWeyrother(t),
			FrancisII(t),
			AlexanderI(t),
		}
		res, err := http.Get(URL("/commanders"))
		require.NoError(t, err, "Requesting commanders")
		defer res.Body.Close()
		httptest.AssertHeaderPages(t, res, expectedPages)
		httptest.AssertJSONCommanders(t, res, expectedCommanders)
	})

	t.Run("GET /factions/:factionID/commanders", func(t *testing.T) {
		t.Parallel()

		route := func(factionID string) string {
			return URL(fmt.Sprintf("/factions/%s/commanders", factionID))
		}

		t.Run("ValidPersistedFactionUUID", func(t *testing.T) {
			factionID := AustrianEmpire(t).ID
			var expectedPages uint = 1
			expectedCommanders := []domain.Commander{FranzVonWeyrother(t), FrancisII(t)}
			res, err := http.Get(route(factionID.String()))
			require.NoError(t, err, "Requesting commanders from Austrian Empire")
			defer res.Body.Close()
			httptest.AssertHeaderPages(t, res, expectedPages)
			httptest.AssertJSONCommanders(t, res, expectedCommanders)
		})

		t.Run("ValidNonPersistedFactionUUID", func(t *testing.T) {
			url := route(uuid.NewV4().String())
			httptest.AssertFailedGET(t, url, http.StatusNotFound, "Not Found")
		})

		t.Run("InvalidFactionUUID", func(t *testing.T) {
			url := route("invalid-uuid")
			httptest.AssertFailedGET(t, url, http.StatusBadRequest, "Bad Request")
		})
	})
}
