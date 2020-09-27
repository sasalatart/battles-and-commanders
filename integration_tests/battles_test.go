package integration_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/sasalatart/batcoms/domain/battles"
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

	t.Run("GET /battles", func(t *testing.T) {
		t.Parallel()

		const expectedPages = 1
		baseURL := URL("/battles")
		cases := []struct {
			description     string
			url             string
			expectedBattles []battles.Battle
		}{
			{
				description: "With no filters",
				url:         baseURL,
				expectedBattles: []battles.Battle{
					BattleOfMegiddo(t),
					BattleOfLodi(t),
					BattleOfAusterlitz(t),
					BattleOfArcole(t),
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
				expectedBattles: []battles.Battle{BattleOfLodi(t), BattleOfAusterlitz(t), BattleOfArcole(t)},
			},
			{
				description:     "With name, summary, place and result filters",
				url:             baseURL + "?name=Austerlitz&summary=napoleonic&place=Moravia&result=Treaty+of+Pressburg",
				expectedBattles: []battles.Battle{BattleOfAusterlitz(t)},
			},
		}
		for _, c := range cases {
			t.Run(c.description, func(t *testing.T) {
				res, err := http.Get(c.url)
				require.NoError(t, err, "Requesting battles")
				defer res.Body.Close()
				httptest.AssertHeaderPages(t, res, expectedPages)
				httptest.AssertJSONBattles(t, res, c.expectedBattles)
			})
		}
	})

	t.Run("GET /factions/:factionID/battles", func(t *testing.T) {
		t.Parallel()

		route := func(factionID string) string {
			return URL(fmt.Sprintf("/factions/%s/battles", factionID))
		}

		t.Run("ValidPersistedFactionUUID", func(t *testing.T) {
			const expectedPages = 1
			frenchFirstRepublicURL := route(FrenchFirstRepublic(t).ID.String())
			newKingdomOfEgyptURL := route(NewKingdomOfEgypt(t).ID.String())
			cases := []struct {
				description     string
				url             string
				expectedBattles []battles.Battle
			}{
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
					description:     "With name, summary, place and result filters",
					url:             frenchFirstRepublicURL + "?name=Lodi&summary=time+to+retreat&place=present+day+Italy&result=French+victory",
					expectedBattles: []battles.Battle{BattleOfLodi(t)},
				},
			}
			for _, c := range cases {
				t.Run(c.description, func(t *testing.T) {
					res, err := http.Get(c.url)
					require.NoError(t, err, "Requesting battles")
					defer res.Body.Close()
					httptest.AssertHeaderPages(t, res, expectedPages)
					httptest.AssertJSONBattles(t, res, c.expectedBattles)
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

	t.Run("GET /commanders/:commanderID/battles", func(t *testing.T) {
		t.Parallel()

		route := func(commanderID string) string {
			return URL(fmt.Sprintf("/commanders/%s/battles", commanderID))
		}

		t.Run("ValidPersistedCommanderUUID", func(t *testing.T) {
			const expectedPages = 1
			napoleonURL := route(Napoleon(t).ID.String())
			cases := []struct {
				description     string
				url             string
				expectedBattles []battles.Battle
			}{
				{
					description:     "With no filters",
					url:             napoleonURL,
					expectedBattles: []battles.Battle{BattleOfLodi(t), BattleOfAusterlitz(t), BattleOfArcole(t)},
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
					expectedBattles: []battles.Battle{BattleOfLodi(t), BattleOfAusterlitz(t), BattleOfArcole(t)},
				},
				{
					description:     "With name, summary, place and result filters",
					url:             napoleonURL + "?name=Lodi&summary=time+to+retreat&place=present+day+Italy&result=French+victory",
					expectedBattles: []battles.Battle{BattleOfLodi(t)},
				},
			}
			for _, c := range cases {
				t.Run(c.description, func(t *testing.T) {
					res, err := http.Get(c.url)
					require.NoError(t, err, "Requesting battles")
					defer res.Body.Close()
					httptest.AssertHeaderPages(t, res, expectedPages)
					httptest.AssertJSONBattles(t, res, c.expectedBattles)
				})
			}
		})

		t.Run("ValidNonPersistedCommanderUUID", func(t *testing.T) {
			url := route(uuid.NewV4().String())
			httptest.AssertFailedGET(t, url, http.StatusNotFound, "Not Found")
		})

		t.Run("InvalidCommanderUUID", func(t *testing.T) {
			url := route("invalid-uuid")
			httptest.AssertFailedGET(t, url, http.StatusBadRequest, "Bad Request")
		})
	})
}
