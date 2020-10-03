package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/domain/battles"
	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/domain/factions"
	"github.com/sasalatart/batcoms/http/httptest"
	"github.com/sasalatart/batcoms/mocks"
	uuid "github.com/satori/go.uuid"
)

func TestBattlesHandlers(t *testing.T) {
	t.Run("GET /battles/:battleID", func(t *testing.T) {
		t.Parallel()

		t.Run("ValidPersistedUUID", func(t *testing.T) {
			battleMock := mocks.Battle()
			app, _, _, battlesRepoMock := appWithReposMocks()
			battlesRepoMock.On("FindOne", battles.FindOneQuery{
				ID: battleMock.ID,
			}).Return(battleMock, nil)

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

			httptest.AssertFailedFiberGET(t, app, "/battles/"+uuid.String(), http.StatusNotFound, "Battle not found")
			battlesRepoMock.AssertExpectations(t)
		})

		t.Run("InvalidUUID", func(t *testing.T) {
			app, _, _, battlesRepoMock := appWithReposMocks()
			httptest.AssertFailedFiberGET(t, app, "/battles/invalid-uuid", http.StatusBadRequest, "Invalid BattleID")
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
		for _, c := range buildInvalidDatesCases(baseURL) {
			t.Run(c.description, func(t *testing.T) {
				app, _, _, battlesRepoMock := appWithReposMocks()
				httptest.AssertFiberGET(t, app, c.url, http.StatusBadRequest, func(res *http.Response) {
					battlesRepoMock.AssertNotCalled(t, "FindMany")
					httptest.AssertErrorMessage(t, res, c.expectedMessage)
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
			fromFactionURL := baseURL(factionMock.ID.String())

			cases := buildBattlesCases(fromFactionURL, func(q battles.FindManyQuery) battles.FindManyQuery {
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
			for _, c := range buildInvalidDatesCases(fromFactionURL) {
				t.Run(c.description, func(t *testing.T) {
					app, factionsRepoMock, _, battlesRepoMock := appWithReposMocks()
					factionsRepoMock.On("FindOne", factions.FindOneQuery{
						ID: factionMock.ID,
					}).Return(factionMock, nil)

					httptest.AssertFiberGET(t, app, c.url, http.StatusBadRequest, func(res *http.Response) {
						factionsRepoMock.AssertExpectations(t)
						battlesRepoMock.AssertNotCalled(t, "FindMany")
						httptest.AssertErrorMessage(t, res, c.expectedMessage)
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

			httptest.AssertFailedFiberGET(t, app, baseURL(uuid.String()), http.StatusNotFound, "Faction not found")
			factionsRepoMock.AssertExpectations(t)
			battlesRepoMock.AssertNotCalled(t, "FindMany")
		})

		t.Run("InvalidFactionUUID", func(t *testing.T) {
			app, factionsRepoMock, _, battlesRepoMock := appWithReposMocks()
			httptest.AssertFailedFiberGET(t, app, baseURL("invalid-uuid"), http.StatusBadRequest, "Invalid FactionID")
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
			fromCommanderURL := baseURL(commanderMock.ID.String())

			cases := buildBattlesCases(fromCommanderURL, func(q battles.FindManyQuery) battles.FindManyQuery {
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
			for _, c := range buildInvalidDatesCases(fromCommanderURL) {
				t.Run(c.description, func(t *testing.T) {
					app, _, commandersRepoMock, battlesRepoMock := appWithReposMocks()
					commandersRepoMock.On("FindOne", commanders.FindOneQuery{
						ID: commanderMock.ID,
					}).Return(commanderMock, nil)

					httptest.AssertFiberGET(t, app, c.url, http.StatusBadRequest, func(res *http.Response) {
						commandersRepoMock.AssertExpectations(t)
						battlesRepoMock.AssertNotCalled(t, "FindMany")
						httptest.AssertErrorMessage(t, res, c.expectedMessage)
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

			httptest.AssertFailedFiberGET(t, app, baseURL(uuid.String()), http.StatusNotFound, "Commander not found")
			commandersRepoMock.AssertExpectations(t)
			battlesRepoMock.AssertNotCalled(t, "FindMany")
		})

		t.Run("InvalidCommanderUUID", func(t *testing.T) {
			app, _, commandersRepoMock, battlesRepoMock := appWithReposMocks()
			httptest.AssertFailedFiberGET(t, app, baseURL("invalid-uuid"), http.StatusBadRequest, "Invalid CommanderID")
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
			description: "With fromDate filter",
			url:         baseURL + "&fromDate=1769-08-15",
			calledWith:  decorateQuery(battles.FindManyQuery{FromDate: "1769-08-15"}),
		},
		{
			description: "With partial fromDate filter",
			url:         baseURL + "&fromDate=1769-08",
			calledWith:  decorateQuery(battles.FindManyQuery{FromDate: "1769-08-01"}),
		},
		{
			description: "With toDate filter",
			url:         baseURL + "&toDate=1821-05-05",
			calledWith:  decorateQuery(battles.FindManyQuery{ToDate: "1821-05-05"}),
		},
		{
			description: "With partial toDate filter",
			url:         baseURL + "&toDate=1821",
			calledWith:  decorateQuery(battles.FindManyQuery{ToDate: "1821-12-31"}),
		},
		{
			description: "With name, summary, place, result, fromDate and toDate filters",
			url: baseURL +
				"&name=Austerlitz" +
				"&summary=napoleonic" +
				"&place=Moravia" +
				"&result=Treaty+of+Pressburg" +
				"&fromDate=1805-12-02" +
				"&toDate=1805-12-02",
			calledWith: decorateQuery(battles.FindManyQuery{
				Name:     "Austerlitz",
				Summary:  "napoleonic",
				Place:    "Moravia",
				Result:   "Treaty of Pressburg",
				FromDate: "1805-12-02",
				ToDate:   "1805-12-02",
			}),
		},
	}
}

type invalidDatesTableCase struct {
	description     string
	url             string
	expectedMessage string
}

func buildInvalidDatesCases(baseURL string) []invalidDatesTableCase {
	const invalidFromDateMessage = "Invalid fromDate, must be in YYYY-MM-DD format"
	const invalidToDateMessage = "Invalid toDate, must be in YYYY-MM-DD format"
	return []invalidDatesTableCase{
		{
			description:     "Invalid fromDate",
			url:             baseURL + "&fromDate=x",
			expectedMessage: invalidFromDateMessage,
		},
		{
			description:     "Invalid toDate",
			url:             baseURL + "&toDate=x",
			expectedMessage: invalidToDateMessage,
		},
		{
			description:     "Invalid fromDate with valid toDate",
			url:             baseURL + "&fromDate=x&toDate=1821-05-05",
			expectedMessage: invalidFromDateMessage,
		},
		{
			description:     "Valid fromDate with invalid toDate",
			url:             baseURL + "&fromDate=1769-08-15&toDate=x",
			expectedMessage: invalidToDateMessage,
		},
		{
			description:     "Invalid fromDate with invalid toDate",
			url:             baseURL + "&fromDate=x&toDate=y",
			expectedMessage: invalidFromDateMessage,
		},
	}
}
