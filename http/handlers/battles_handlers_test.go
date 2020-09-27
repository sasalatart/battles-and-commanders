package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/domain/battles"
	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/domain/factions"
	batcomshttp "github.com/sasalatart/batcoms/http"
	"github.com/sasalatart/batcoms/http/httptest"
	"github.com/sasalatart/batcoms/mocks"
	uuid "github.com/satori/go.uuid"
)

func TestBattlesHandlers(t *testing.T) {
	t.Run("GET /battles/:battleID", func(t *testing.T) {
		t.Parallel()

		t.Run("ValidPersistedUUID", func(t *testing.T) {
			battleMock := mocks.Battle()
			battlesRepoMock := new(mocks.BattlesRepository)
			battlesRepoMock.On("FindOne", battles.FindOneQuery{
				ID: battleMock.ID,
			}).Return(battleMock, nil)
			app := batcomshttp.Setup(new(mocks.FactionsRepository), new(mocks.CommandersRepository), battlesRepoMock, true)

			httptest.AssertFiberGET(t, app, "/battles/"+battleMock.ID.String(), http.StatusOK, func(res *http.Response) {
				battlesRepoMock.AssertExpectations(t)
				httptest.AssertJSONBattle(t, res, battleMock)
			})
		})

		t.Run("ValidNonPersistedUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			app, _, _, battlesRepoMock := appWithReposMocks()
			battlesRepoMock.On("FindOne", battles.FindOneQuery{
				ID: uuid,
			}).Return(battles.Battle{}, domain.ErrNotFound)

			httptest.AssertFailedFiberGET(t, app, "/battles/"+uuid.String(), *fiber.ErrNotFound)
			battlesRepoMock.AssertExpectations(t)
		})

		t.Run("InvalidUUID", func(t *testing.T) {
			app, _, _, battlesRepoMock := appWithReposMocks()
			httptest.AssertFailedFiberGET(t, app, "/battles/invalid-uuid", *fiber.ErrBadRequest)
			battlesRepoMock.AssertNotCalled(t, "FindOne")
		})
	})

	t.Run("GET /battles", func(t *testing.T) {
		t.Parallel()

		const page = 2
		const pagesMock = 3
		baseURL := fmt.Sprintf("/battles?page=%d", page)
		battlesMock := []battles.Battle{mocks.Battle()}

		cases := buildBattlesCases(baseURL, func(q battles.FindManyQuery) battles.FindManyQuery {
			return q
		})
		for _, c := range cases {
			t.Run(c.description, func(t *testing.T) {
				app, _, _, battlesRepoMock := appWithReposMocks()
				battlesRepoMock.On("FindMany", c.calledWith, page).
					Return(battlesMock, pagesMock, nil)
				httptest.AssertFiberGET(t, app, c.url, http.StatusOK, func(res *http.Response) {
					battlesRepoMock.AssertExpectations(t)
					httptest.AssertHeaderPages(t, res, pagesMock)
					httptest.AssertJSONBattles(t, res, battlesMock)
				})
			})
		}
	})

	t.Run("GET /factions/:factionID/battles", func(t *testing.T) {
		t.Parallel()

		const page = 2
		baseURL := func(factionID string) string {
			return fmt.Sprintf("/factions/%s/battles?page=%d", factionID, page)
		}

		t.Run("ValidPersistedFactionUUID", func(t *testing.T) {
			const pagesMock = 3
			factionMock := mocks.Faction()
			battlesMock := []battles.Battle{mocks.Battle()}

			cases := buildBattlesCases(baseURL(factionMock.ID.String()), func(q battles.FindManyQuery) battles.FindManyQuery {
				q.FactionID = factionMock.ID
				return q
			})
			for _, c := range cases {
				t.Run(c.description, func(t *testing.T) {
					app, factionsRepoMock, _, battlesRepoMock := appWithReposMocks()
					factionsRepoMock.On("FindOne", factions.FindOneQuery{
						ID: factionMock.ID,
					}).Return(factionMock, nil)
					battlesRepoMock.On("FindMany", c.calledWith, page).
						Return(battlesMock, pagesMock, nil)

					httptest.AssertFiberGET(t, app, c.url, http.StatusOK, func(res *http.Response) {
						battlesRepoMock.AssertExpectations(t)
						httptest.AssertHeaderPages(t, res, pagesMock)
						httptest.AssertJSONBattles(t, res, battlesMock)
					})
				})
			}
		})

		t.Run("ValidNonPersistedFactionUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			app, factionsRepoMock, _, battlesRepoMock := appWithReposMocks()
			factionsRepoMock.On("FindOne", factions.FindOneQuery{
				ID: uuid,
			}).Return(factions.Faction{}, domain.ErrNotFound)

			httptest.AssertFailedFiberGET(t, app, baseURL(uuid.String()), *fiber.ErrNotFound)
			factionsRepoMock.AssertExpectations(t)
			battlesRepoMock.AssertNotCalled(t, "FindMany")
		})

		t.Run("InvalidFactionUUID", func(t *testing.T) {
			app, factionsRepoMock, _, battlesRepoMock := appWithReposMocks()
			httptest.AssertFailedFiberGET(t, app, baseURL("invalid-uuid"), *fiber.ErrBadRequest)
			factionsRepoMock.AssertNotCalled(t, "FindOne")
			battlesRepoMock.AssertNotCalled(t, "FindMany")
		})
	})

	t.Run("GET /commanders/:commanderID/battles", func(t *testing.T) {
		t.Parallel()

		const page = 2
		baseURL := func(commanderID string) string {
			return fmt.Sprintf("/commanders/%s/battles?page=%d", commanderID, page)
		}

		t.Run("ValidPersistedCommanderUUID", func(t *testing.T) {
			const pagesMock = 3
			commanderMock := mocks.Commander()
			battlesMock := []battles.Battle{mocks.Battle()}

			cases := buildBattlesCases(baseURL(commanderMock.ID.String()), func(q battles.FindManyQuery) battles.FindManyQuery {
				q.CommanderID = commanderMock.ID
				return q
			})
			for _, c := range cases {
				t.Run(c.description, func(t *testing.T) {
					app, _, commandersRepoMock, battlesRepoMock := appWithReposMocks()
					commandersRepoMock.On("FindOne", commanders.FindOneQuery{
						ID: commanderMock.ID,
					}).Return(commanderMock, nil)
					battlesRepoMock.On("FindMany", c.calledWith, page).
						Return(battlesMock, pagesMock, nil)

					httptest.AssertFiberGET(t, app, c.url, http.StatusOK, func(res *http.Response) {
						battlesRepoMock.AssertExpectations(t)
						httptest.AssertHeaderPages(t, res, pagesMock)
						httptest.AssertJSONBattles(t, res, battlesMock)
					})
				})
			}
		})

		t.Run("ValidNonPersistedCommanderUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			app, _, commandersRepoMock, battlesRepoMock := appWithReposMocks()
			commandersRepoMock.On("FindOne", commanders.FindOneQuery{
				ID: uuid,
			}).Return(commanders.Commander{}, domain.ErrNotFound)

			httptest.AssertFailedFiberGET(t, app, baseURL(uuid.String()), *fiber.ErrNotFound)
			commandersRepoMock.AssertExpectations(t)
			battlesRepoMock.AssertNotCalled(t, "FindMany")
		})

		t.Run("InvalidCommanderUUID", func(t *testing.T) {
			app, _, commandersRepoMock, battlesRepoMock := appWithReposMocks()
			httptest.AssertFailedFiberGET(t, app, baseURL("invalid-uuid"), *fiber.ErrBadRequest)
			commandersRepoMock.AssertNotCalled(t, "FindOne")
			battlesRepoMock.AssertNotCalled(t, "FindMany")
		})
	})
}

type battlesTableCase struct {
	description string
	url         string
	calledWith  battles.FindManyQuery
}

func buildBattlesCases(baseURL string, decorateQuery func(battles.FindManyQuery) battles.FindManyQuery) []battlesTableCase {
	return []battlesTableCase{
		{
			description: "With no filters",
			url:         baseURL,
			calledWith:  decorateQuery(battles.FindManyQuery{}),
		},
		{
			description: "With name filter",
			url:         baseURL + "&name=Austerlitz",
			calledWith:  decorateQuery(battles.FindManyQuery{Name: "Austerlitz"}),
		},
		{
			description: "With summary filter",
			url:         baseURL + "&summary=napoleonic",
			calledWith:  decorateQuery(battles.FindManyQuery{Summary: "napoleonic"}),
		},
		{
			description: "With place filter",
			url:         baseURL + "&place=Moravia",
			calledWith:  decorateQuery(battles.FindManyQuery{Place: "Moravia"}),
		},
		{
			description: "With result filter",
			url:         baseURL + "&result=Treaty+of+Pressburg",
			calledWith:  decorateQuery(battles.FindManyQuery{Result: "Treaty of Pressburg"}),
		},
		{
			description: "With name, summary, place and result filters",
			url:         baseURL + "&name=Austerlitz&summary=napoleonic&place=Moravia&result=Treaty+of+Pressburg",
			calledWith: decorateQuery(battles.FindManyQuery{
				Name:    "Austerlitz",
				Summary: "napoleonic",
				Place:   "Moravia",
				Result:  "Treaty of Pressburg",
			}),
		},
	}
}
