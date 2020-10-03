package integration_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/sasalatart/batcoms/domain/factions"
	"github.com/sasalatart/batcoms/http/httptest"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFactionsEndpoints(t *testing.T) {
	type factionsEndpointCase struct {
		description          string
		url                  string
		expectedFactions     []factions.Faction
		expectedErrorCode    int
		expectedErrorMessage string
	}

	assertFactionsEndpointCase := func(t *testing.T, c factionsEndpointCase, expectedPages int) {
		t.Helper()
		res, err := http.Get(c.url)
		require.NoError(t, err, "Requesting factions")
		defer res.Body.Close()
		if c.expectedErrorMessage == "" {
			assert.Equal(t, http.StatusOK, res.StatusCode)
			httptest.AssertHeaderPages(t, res, expectedPages)
			httptest.AssertJSONFactions(t, res, c.expectedFactions)
		} else {
			assert.Equal(t, c.expectedErrorCode, res.StatusCode)
			httptest.AssertErrorMessage(t, res, c.expectedErrorMessage)
		}
	}

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
			httptest.AssertFailedGET(t, url, http.StatusNotFound, "Faction not found")
		})

		t.Run("InvalidUUID", func(t *testing.T) {
			url := route("invalid-uuid")
			httptest.AssertFailedGET(t, url, http.StatusBadRequest, "Invalid FactionID")
		})
	})

	t.Run("GET /factions", func(t *testing.T) {
		t.Parallel()

		const expectedPages = 1
		cases := []factionsEndpointCase{
			{
				description: "With no filters",
				url:         URL("/factions"),
				expectedFactions: []factions.Faction{
					RussianEmpire(t),
					NewKingdomOfEgypt(t),
					HabsburgMonarchy(t),
					FrenchFirstRepublic(t),
					FirstFrenchEmpire(t),
					Canaan(t),
					AustrianEmpire(t),
				},
			},
			{
				description:      "With name filter",
				url:              URL("/factions?name=First+French+Empire"),
				expectedFactions: []factions.Faction{FirstFrenchEmpire(t)},
			},
			{
				description:      "With summary filter",
				url:              URL("/factions?summary=ruled+by+Napoleon"),
				expectedFactions: []factions.Faction{FirstFrenchEmpire(t)},
			},
			{
				description:      "With name and summary filters",
				url:              URL("/factions?name=russian&summary=eurasia"),
				expectedFactions: []factions.Faction{RussianEmpire(t)},
			},
		}
		for _, c := range cases {
			t.Run(c.description, func(t *testing.T) {
				assertFactionsEndpointCase(t, c, expectedPages)
			})
		}
	})

	t.Run("GET /commanders/:commanderID/factions", func(t *testing.T) {
		t.Parallel()

		route := func(commanderID string) string {
			return URL(fmt.Sprintf("/commanders/%s/factions", commanderID))
		}

		const expectedPages = 1
		commanderID := Napoleon(t).ID
		cases := []struct {
			description          string
			url                  string
			expectedFactions     []factions.Faction
			expectedErrorCode    int
			expectedErrorMessage string
		}{
			{
				description:      "With no filters",
				url:              route(commanderID.String()),
				expectedFactions: []factions.Faction{FrenchFirstRepublic(t), FirstFrenchEmpire(t)},
			},
			{
				description:      "With name filter",
				url:              route(commanderID.String()) + "?name=empire",
				expectedFactions: []factions.Faction{FirstFrenchEmpire(t)},
			},
			{
				description:      "With summary filter",
				url:              route(commanderID.String()) + "?summary=first+republic",
				expectedFactions: []factions.Faction{FrenchFirstRepublic(t)},
			},
			{
				description:      "With name and summary filters",
				url:              route(commanderID.String()) + "?name=empire&summary=ruled+by+Napoleon+Bonaparte",
				expectedFactions: []factions.Faction{FirstFrenchEmpire(t)},
			},
			{
				description:          "With valid non-persisted CommanderID",
				url:                  route(uuid.NewV4().String()),
				expectedErrorCode:    http.StatusNotFound,
				expectedErrorMessage: "Commander not found",
			},
			{
				description:          "With invalid CommanderID",
				url:                  route("invalid-uuid"),
				expectedErrorCode:    http.StatusBadRequest,
				expectedErrorMessage: "Invalid CommanderID",
			},
		}
		for _, c := range cases {
			t.Run(c.description, func(t *testing.T) {
				assertFactionsEndpointCase(t, c, expectedPages)
			})
		}
	})
}
