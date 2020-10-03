package integration_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/sasalatart/batcoms/domain/battles"
	"github.com/sasalatart/batcoms/http/httptest"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBattlesEndpoints(t *testing.T) {
	type battlesEndpointCase struct {
		description          string
		url                  string
		expectedBattles      []battles.Battle
		expectedErrorCode    int
		expectedErrorMessage string
	}

	assertBattlesEndpointCase := func(t *testing.T, c battlesEndpointCase, expectedPages int) {
		t.Helper()
		res, err := http.Get(c.url)
		require.NoError(t, err, "Requesting battles")
		defer res.Body.Close()
		if c.expectedErrorMessage == "" {
			assert.Equal(t, http.StatusOK, res.StatusCode)
			httptest.AssertHeaderPages(t, res, expectedPages)
			httptest.AssertJSONBattles(t, res, c.expectedBattles)
		} else {
			assert.Equal(t, c.expectedErrorCode, res.StatusCode)
			httptest.AssertErrorMessage(t, res, c.expectedErrorMessage)
		}
	}

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
			httptest.AssertFailedGET(t, url, http.StatusNotFound, "Battle not found")
		})

		t.Run("InvalidUUID", func(t *testing.T) {
			url := route("invalid-uuid")
			httptest.AssertFailedGET(t, url, http.StatusBadRequest, "Invalid BattleID")
		})
	})

	t.Run("GET /battles", func(t *testing.T) {
		t.Parallel()

		const expectedPages = 1
		baseURL := URL("/battles")
		cases := []battlesEndpointCase{
			{
				description: "With no filters",
				url:         baseURL,
				expectedBattles: []battles.Battle{
					BattleOfMegiddo(t),
					BattleOfLodi(t),
					BattleOfArcole(t),
					BattleOfAusterlitz(t),
				},
			},
			{
				description:     "With name filter",
				url:             baseURL + "?name=Arcole",
				expectedBattles: []battles.Battle{BattleOfArcole(t)},
			},
			{
				description:     "With summary filter",
				url:             baseURL + "?summary=retreat",
				expectedBattles: []battles.Battle{BattleOfLodi(t), BattleOfArcole(t)},
			},
			{
				description:     "With place filter",
				url:             baseURL + "?place=Canaan",
				expectedBattles: []battles.Battle{BattleOfMegiddo(t)},
			},
			{
				description:     "With result filter",
				url:             baseURL + "?result=French+Victory",
				expectedBattles: []battles.Battle{BattleOfLodi(t), BattleOfArcole(t), BattleOfAusterlitz(t)},
			},
			{
				description:     "With fromDate filter",
				url:             baseURL + "?fromDate=1700",
				expectedBattles: []battles.Battle{BattleOfLodi(t), BattleOfArcole(t), BattleOfAusterlitz(t)},
			},
			{
				description:     "With toDate filter",
				url:             baseURL + "?toDate=1700",
				expectedBattles: []battles.Battle{BattleOfMegiddo(t)},
			},
			{
				description:     "With fromDate and toDate filter",
				url:             baseURL + "?fromDate=1805-12-02&toDate=1805-12-02",
				expectedBattles: []battles.Battle{BattleOfAusterlitz(t)},
			},
			{
				description: "With name, summary, place, result, fromDate and toDate filters",
				url: baseURL +
					"?name=Austerlitz" +
					"&summary=napoleonic" +
					"&place=Moravia" +
					"&result=Treaty+of+Pressburg" +
					"&fromDate=1805-12-02" +
					"&toDate=1805-12-02",
				expectedBattles: []battles.Battle{BattleOfAusterlitz(t)},
			},
			{
				description:          "With invalid fromDate",
				url:                  baseURL + "?fromDate=x",
				expectedErrorCode:    http.StatusBadRequest,
				expectedErrorMessage: invalidFromDateMessage,
			},
			{
				description:          "With invalid toDate",
				url:                  baseURL + "?toDate=x",
				expectedErrorCode:    http.StatusBadRequest,
				expectedErrorMessage: invalidToDateMessage,
			},
			{
				description:          "With invalid page",
				url:                  baseURL + "?page=-1",
				expectedErrorCode:    http.StatusBadRequest,
				expectedErrorMessage: "Invalid page",
			},
		}
		for _, c := range cases {
			t.Run(c.description, func(t *testing.T) {
				assertBattlesEndpointCase(t, c, expectedPages)
			})
		}
	})

	t.Run("GET /factions/:factionID/battles", func(t *testing.T) {
		t.Parallel()

		route := func(factionID string) string {
			return URL(fmt.Sprintf("/factions/%s/battles", factionID))
		}

		const expectedPages = 1
		frenchFirstRepublicURL := route(FrenchFirstRepublic(t).ID.String())
		newKingdomOfEgyptURL := route(NewKingdomOfEgypt(t).ID.String())
		cases := []battlesEndpointCase{
			{
				description:     "With no filters",
				url:             frenchFirstRepublicURL,
				expectedBattles: []battles.Battle{BattleOfLodi(t), BattleOfArcole(t)},
			},
			{
				description:     "With name filter",
				url:             frenchFirstRepublicURL + "?name=Lodi",
				expectedBattles: []battles.Battle{BattleOfLodi(t)},
			},
			{
				description:     "With summary filter",
				url:             frenchFirstRepublicURL + "?summary=line+of+retreat",
				expectedBattles: []battles.Battle{BattleOfArcole(t)},
			},
			{
				description:     "With place filter",
				url:             newKingdomOfEgyptURL + "?place=Canaan",
				expectedBattles: []battles.Battle{BattleOfMegiddo(t)},
			},
			{
				description:     "With result filter",
				url:             frenchFirstRepublicURL + "?result=French+Victory",
				expectedBattles: []battles.Battle{BattleOfLodi(t), BattleOfArcole(t)},
			},
			{
				description:     "With fromDate filter",
				url:             frenchFirstRepublicURL + "?fromDate=1796-11",
				expectedBattles: []battles.Battle{BattleOfArcole(t)},
			},
			{
				description:     "With toDate filter",
				url:             frenchFirstRepublicURL + "?toDate=1796",
				expectedBattles: []battles.Battle{BattleOfLodi(t), BattleOfArcole(t)},
			},
			{
				description:     "With fromDate and toDate filter",
				url:             frenchFirstRepublicURL + "?fromDate=1796-05&toDate=1796-06",
				expectedBattles: []battles.Battle{BattleOfLodi(t)},
			},
			{
				description: "With name, summary, place, result, fromDate and toDate filters",
				url: frenchFirstRepublicURL +
					"?name=Lodi" +
					"&summary=time+to+retreat" +
					"&place=present+day+Italy" +
					"&result=French+victory" +
					"&fromDate=1796-05" +
					"&toDate=1796-06",
				expectedBattles: []battles.Battle{BattleOfLodi(t)},
			},
			{
				description:          "With valid, non-persisted FactionID",
				url:                  route(uuid.NewV4().String()),
				expectedErrorCode:    http.StatusNotFound,
				expectedErrorMessage: "Faction not found",
			},
			{
				description:          "With invalid FactionID",
				url:                  route("invalid-id"),
				expectedErrorCode:    http.StatusBadRequest,
				expectedErrorMessage: "Invalid FactionID",
			},
			{
				description:          "With invalid fromDate",
				url:                  frenchFirstRepublicURL + "?fromDate=x",
				expectedErrorCode:    http.StatusBadRequest,
				expectedErrorMessage: invalidFromDateMessage,
			},
			{
				description:          "With invalid toDate",
				url:                  frenchFirstRepublicURL + "?toDate=x",
				expectedErrorCode:    http.StatusBadRequest,
				expectedErrorMessage: invalidToDateMessage,
			},
			{
				description:          "With invalid page",
				url:                  frenchFirstRepublicURL + "?page=-1",
				expectedErrorCode:    http.StatusBadRequest,
				expectedErrorMessage: "Invalid page",
			},
		}
		for _, c := range cases {
			t.Run(c.description, func(t *testing.T) {
				assertBattlesEndpointCase(t, c, expectedPages)
			})
		}
	})

	t.Run("GET /commanders/:commanderID/battles", func(t *testing.T) {
		t.Parallel()

		route := func(commanderID string) string {
			return URL(fmt.Sprintf("/commanders/%s/battles", commanderID))
		}

		const expectedPages = 1
		napoleonURL := route(Napoleon(t).ID.String())
		cases := []battlesEndpointCase{
			{
				description:     "With no filters",
				url:             napoleonURL,
				expectedBattles: []battles.Battle{BattleOfLodi(t), BattleOfArcole(t), BattleOfAusterlitz(t)},
			},
			{
				description:     "With name filter",
				url:             napoleonURL + "?name=Arcole",
				expectedBattles: []battles.Battle{BattleOfArcole(t)},
			},
			{
				description:     "With summary filter",
				url:             napoleonURL + "?summary=line+of+retreat",
				expectedBattles: []battles.Battle{BattleOfArcole(t)},
			},
			{
				description:     "With place filter",
				url:             napoleonURL + "?place=present-day+Italy",
				expectedBattles: []battles.Battle{BattleOfLodi(t)},
			},
			{
				description:     "With result filter",
				url:             napoleonURL + "?result=French+Victory",
				expectedBattles: []battles.Battle{BattleOfLodi(t), BattleOfArcole(t), BattleOfAusterlitz(t)},
			},
			{
				description:     "With fromDate filter",
				url:             napoleonURL + "?fromDate=1796-11",
				expectedBattles: []battles.Battle{BattleOfArcole(t), BattleOfAusterlitz(t)},
			},
			{
				description:     "With toDate filter",
				url:             napoleonURL + "?toDate=1796",
				expectedBattles: []battles.Battle{BattleOfLodi(t), BattleOfArcole(t)},
			},
			{
				description:     "With fromDate and toDate filter",
				url:             napoleonURL + "?fromDate=1796-11&toDate=1805",
				expectedBattles: []battles.Battle{BattleOfArcole(t), BattleOfAusterlitz(t)},
			},
			{
				description: "With name, summary, place, result, fromDate and toDate filters",
				url: napoleonURL +
					"?name=Lodi" +
					"&summary=time+to+retreat" +
					"&place=present+day+Italy" +
					"&result=French+victory" +
					"&fromDate=1796-05" +
					"&toDate=1796-06",
				expectedBattles: []battles.Battle{BattleOfLodi(t)},
			},
			{
				description:          "With valid, non-persisted CommanderID",
				url:                  route(uuid.NewV4().String()),
				expectedErrorCode:    http.StatusNotFound,
				expectedErrorMessage: "Commander not found",
			},
			{
				description:          "With invalid CommanderID",
				url:                  route("invalid-id"),
				expectedErrorCode:    http.StatusBadRequest,
				expectedErrorMessage: "Invalid CommanderID",
			},
			{
				description:          "With invalid fromDate",
				url:                  napoleonURL + "?fromDate=x",
				expectedErrorCode:    http.StatusBadRequest,
				expectedErrorMessage: invalidFromDateMessage,
			},
			{
				description:          "With invalid toDate",
				url:                  napoleonURL + "?toDate=x",
				expectedErrorCode:    http.StatusBadRequest,
				expectedErrorMessage: invalidToDateMessage,
			},
			{
				description:          "With invalid page",
				url:                  napoleonURL + "?page=-1",
				expectedErrorCode:    http.StatusBadRequest,
				expectedErrorMessage: "Invalid page",
			},
		}
		for _, c := range cases {
			t.Run(c.description, func(t *testing.T) {
				assertBattlesEndpointCase(t, c, expectedPages)
			})
		}
	})
}

const invalidFromDateMessage = "Invalid fromDate, must be in YYYY-MM-DD format"
const invalidToDateMessage = "Invalid toDate, must be in YYYY-MM-DD format"
