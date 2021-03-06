package integration_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/http/httptest"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandersEndpoints(t *testing.T) {
	type commandersEndpointCase struct {
		description          string
		url                  string
		expectedCommanders   []commanders.Commander
		expectedErrorCode    int
		expectedErrorMessage string
	}

	assertCommandersEndpointCase := func(t *testing.T, c commandersEndpointCase, expectedPages int) {
		t.Helper()
		res, err := http.Get(c.url)
		require.NoError(t, err, "Requesting commanders")
		defer res.Body.Close()
		if c.expectedErrorMessage == "" {
			assert.Equal(t, http.StatusOK, res.StatusCode)
			httptest.AssertHeaderPages(t, res, expectedPages)
			httptest.AssertJSONCommanders(t, res, c.expectedCommanders)
		} else {
			assert.Equal(t, c.expectedErrorCode, res.StatusCode)
			httptest.AssertErrorMessage(t, res, c.expectedErrorMessage)
		}
	}

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
			httptest.AssertFailedGET(t, url, http.StatusNotFound, "Commander not found")
		})

		t.Run("InvalidUUID", func(t *testing.T) {
			url := route("invalid-uuid")
			httptest.AssertFailedGET(t, url, http.StatusBadRequest, "Invalid CommanderID")
		})
	})

	t.Run("GET /commanders", func(t *testing.T) {
		t.Parallel()

		const expectedPages = 1
		cases := []commandersEndpointCase{
			{
				description: "With no filters",
				url:         URL("/commanders"),
				expectedCommanders: []commanders.Commander{
					ThutmoseIII(t),
					Napoleon(t),
					MikhailKutuzov(t),
					KarlPhilippSebottendorf(t),
					JozsefAlvinczi(t),
					JohannPeterBeaulieu(t),
					FranzVonWeyrother(t),
					FrancisII(t),
					AlexanderI(t),
				},
			},
			{
				description:        "With name filter",
				url:                URL("/commanders?name=napoleon"),
				expectedCommanders: []commanders.Commander{Napoleon(t)},
			},
			{
				description:        "With summary filter",
				url:                URL("/commanders?summary=emperor"),
				expectedCommanders: []commanders.Commander{Napoleon(t), FrancisII(t), AlexanderI(t)},
			},
			{
				description:        "With name and summary filters",
				url:                URL("/commanders?name=alexander&summary=emperor"),
				expectedCommanders: []commanders.Commander{AlexanderI(t)},
			},
		}
		for _, c := range cases {
			t.Run(c.description, func(t *testing.T) {
				assertCommandersEndpointCase(t, c, expectedPages)
			})
		}
	})

	t.Run("GET /factions/:factionID/commanders", func(t *testing.T) {
		t.Parallel()

		route := func(factionID string) string {
			return URL(fmt.Sprintf("/factions/%s/commanders", factionID))
		}

		const expectedPages = 1
		factionID := AustrianEmpire(t).ID
		cases := []commandersEndpointCase{
			{
				description:        "With no filters",
				url:                route(factionID.String()),
				expectedCommanders: []commanders.Commander{FranzVonWeyrother(t), FrancisII(t)},
			},
			{
				description:        "With name filter",
				url:                route(factionID.String()) + "?name=franz",
				expectedCommanders: []commanders.Commander{FranzVonWeyrother(t)},
			},
			{
				description:        "With summary filter",
				url:                route(factionID.String()) + "?summary=emperor",
				expectedCommanders: []commanders.Commander{FrancisII(t)},
			},
			{
				description:        "With name and summary filters",
				url:                route(factionID.String()) + "?name=franz&summary=emperor",
				expectedCommanders: []commanders.Commander{},
			},
			{
				description:          "With valid non-persisted FactionID",
				url:                  route(uuid.NewV4().String()),
				expectedErrorCode:    http.StatusNotFound,
				expectedErrorMessage: "Faction not found",
			},
			{
				description:          "With invalid FactionID",
				url:                  route("invalid-uuid"),
				expectedErrorCode:    http.StatusBadRequest,
				expectedErrorMessage: "Invalid FactionID",
			},
		}
		for _, c := range cases {
			t.Run(c.description, func(t *testing.T) {
				assertCommandersEndpointCase(t, c, expectedPages)
			})
		}
	})
}
