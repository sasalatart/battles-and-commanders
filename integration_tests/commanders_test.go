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
		cases := []struct {
			description        string
			url                string
			expectedCommanders []domain.Commander
		}{
			{
				description:        "With no filters",
				url:                URL("/commanders"),
				expectedCommanders: []domain.Commander{Napoleon(t), MikhailKutuzov(t), FranzVonWeyrother(t), FrancisII(t), AlexanderI(t)},
			},
			{
				description:        "With name filter",
				url:                URL("/commanders?name=napoleon"),
				expectedCommanders: []domain.Commander{Napoleon(t)},
			},
			{
				description:        "With summary filter",
				url:                URL("/commanders?summary=emperor"),
				expectedCommanders: []domain.Commander{Napoleon(t), FrancisII(t), AlexanderI(t)},
			},
			{
				description:        "With name and summary filters",
				url:                URL("/commanders?name=alexander&summary=emperor"),
				expectedCommanders: []domain.Commander{AlexanderI(t)},
			},
		}
		for _, c := range cases {
			t.Run(c.description, func(t *testing.T) {
				res, err := http.Get(c.url)
				require.NoError(t, err, "Requesting commanders")
				defer res.Body.Close()
				httptest.AssertHeaderPages(t, res, expectedPages)
				httptest.AssertJSONCommanders(t, res, c.expectedCommanders)
			})
		}
	})

	t.Run("GET /factions/:factionID/commanders", func(t *testing.T) {
		t.Parallel()

		route := func(factionID string) string {
			return URL(fmt.Sprintf("/factions/%s/commanders", factionID))
		}

		t.Run("ValidPersistedFactionUUID", func(t *testing.T) {
			var expectedPages uint = 1
			factionID := AustrianEmpire(t).ID
			cases := []struct {
				description        string
				url                string
				expectedCommanders []domain.Commander
			}{
				{
					description:        "With no filters",
					url:                route(factionID.String()),
					expectedCommanders: []domain.Commander{FranzVonWeyrother(t), FrancisII(t)},
				},
				{
					description:        "With name filter",
					url:                route(factionID.String()) + "?name=franz",
					expectedCommanders: []domain.Commander{FranzVonWeyrother(t)},
				},
				{
					description:        "With summary filter",
					url:                route(factionID.String()) + "?summary=emperor",
					expectedCommanders: []domain.Commander{FrancisII(t)},
				},
				{
					description:        "With name and summary filters",
					url:                route(factionID.String()) + "?name=franz&summary=emperor",
					expectedCommanders: []domain.Commander{},
				},
			}
			for _, c := range cases {
				t.Run(c.description, func(t *testing.T) {
					res, err := http.Get(c.url)
					require.NoError(t, err, "Requesting commanders")
					defer res.Body.Close()
					httptest.AssertHeaderPages(t, res, expectedPages)
					httptest.AssertJSONCommanders(t, res, c.expectedCommanders)
				})
			}
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
