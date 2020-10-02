package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/domain/factions"
	"github.com/sasalatart/batcoms/http/httptest"
	"github.com/sasalatart/batcoms/mocks"
	uuid "github.com/satori/go.uuid"
)

func TestFactionsHandlers(t *testing.T) {
	t.Run("GET /factions/:factionID", func(t *testing.T) {
		t.Parallel()

		t.Run("ValidPersistedUUID", func(t *testing.T) {
			factionMock := mocks.Faction()
			app, factionsRepoMock, _, _ := appWithReposMocks()
			factionsRepoMock.On("FindOne", factions.FindOneQuery{
				ID: factionMock.ID,
			}).Return(factionMock, nil)

			httptest.AssertFiberGET(t, app, "/factions/"+factionMock.ID.String(), http.StatusOK, func(res *http.Response) {
				factionsRepoMock.AssertExpectations(t)
				httptest.AssertJSONFaction(t, res, factionMock)
			})
		})

		t.Run("ValidNonPersistedUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			app, factionsRepoMock, _, _ := appWithReposMocks()
			factionsRepoMock.On("FindOne", factions.FindOneQuery{
				ID: uuid,
			}).Return(factions.Faction{}, domain.ErrNotFound)

			httptest.AssertFailedFiberGET(t, app, "/factions/"+uuid.String(), http.StatusNotFound, "Faction not found")
			factionsRepoMock.AssertExpectations(t)
		})

		t.Run("InvalidUUID", func(t *testing.T) {
			app, factionsRepoMock, _, _ := appWithReposMocks()
			httptest.AssertFailedFiberGET(t, app, "/factions/invalid-uuid", http.StatusBadRequest, "Invalid FactionID")
			factionsRepoMock.AssertNotCalled(t, "FindOne")
		})
	})

	t.Run("GET /factions", func(t *testing.T) {
		t.Parallel()

		const page = 2
		const pagesMock = 3
		baseURL := fmt.Sprintf("/factions?page=%d", page)
		factionsMock := []factions.Faction{mocks.Faction(), mocks.Faction2()}

		cases := buildFactionsCases(baseURL, func(q factions.FindManyQuery) factions.FindManyQuery {
			return q
		})
		for _, c := range cases {
			t.Run(c.description, func(t *testing.T) {
				app, factionsRepoMock, _, _ := appWithReposMocks()
				factionsRepoMock.On("FindMany", c.calledWith, page).
					Return(factionsMock, pagesMock, nil)
				httptest.AssertFiberGET(t, app, c.url, http.StatusOK, func(res *http.Response) {
					factionsRepoMock.AssertExpectations(t)
					httptest.AssertHeaderPages(t, res, pagesMock)
					httptest.AssertJSONFactions(t, res, factionsMock)
				})
			})
		}
	})

	t.Run("GET /commanders/:commanderID/factions", func(t *testing.T) {
		t.Parallel()

		const page = 2
		buildURL := func(commanderID string) string {
			return fmt.Sprintf("/commanders/%s/factions?page=%d", commanderID, page)
		}

		t.Run("ValidPersistedCommanderUUID", func(t *testing.T) {
			const pagesMock = 3
			commanderMock := mocks.Commander()
			factionsMock := []factions.Faction{mocks.Faction(), mocks.Faction2()}

			cases := buildFactionsCases(buildURL(commanderMock.ID.String()), func(q factions.FindManyQuery) factions.FindManyQuery {
				q.CommanderID = commanderMock.ID
				return q
			})
			for _, c := range cases {
				t.Run(c.description, func(t *testing.T) {
					app, factionsRepoMock, commandersRepoMock, _ := appWithReposMocks()
					commandersRepoMock.On("FindOne", commanders.FindOneQuery{
						ID: commanderMock.ID,
					}).Return(commanderMock, nil)
					factionsRepoMock.On("FindMany", c.calledWith, page).
						Return(factionsMock, pagesMock, nil)

					httptest.AssertFiberGET(t, app, c.url, http.StatusOK, func(res *http.Response) {
						factionsRepoMock.AssertExpectations(t)
						httptest.AssertHeaderPages(t, res, pagesMock)
						httptest.AssertJSONFactions(t, res, factionsMock)
					})
				})
			}
		})

		t.Run("ValidNonPersistedCommanderUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			app, factionsRepoMock, commandersRepoMock, _ := appWithReposMocks()
			commandersRepoMock.On("FindOne", commanders.FindOneQuery{
				ID: uuid,
			}).Return(commanders.Commander{}, domain.ErrNotFound)

			httptest.AssertFailedFiberGET(t, app, buildURL(uuid.String()), http.StatusNotFound, "Commander not found")
			commandersRepoMock.AssertExpectations(t)
			factionsRepoMock.AssertNotCalled(t, "FindMany")
		})

		t.Run("InvalidCommanderUUID", func(t *testing.T) {
			app, factionsRepoMock, commandersRepoMock, _ := appWithReposMocks()
			httptest.AssertFailedFiberGET(t, app, buildURL("invalid-uuid"), http.StatusBadRequest, "Invalid CommanderID")
			commandersRepoMock.AssertNotCalled(t, "FindMany")
			factionsRepoMock.AssertNotCalled(t, "FindOne")
		})
	})
}

type factionsTableCase struct {
	description string
	url         string
	calledWith  factions.FindManyQuery
}

func buildFactionsCases(baseURL string, decorateQuery func(factions.FindManyQuery) factions.FindManyQuery) []factionsTableCase {
	return []factionsTableCase{
		{
			description: "With no filters",
			url:         baseURL,
			calledWith:  decorateQuery(factions.FindManyQuery{}),
		},
		{
			description: "With name filter",
			url:         baseURL + "&name=First+French+Empire",
			calledWith:  decorateQuery(factions.FindManyQuery{Name: "First French Empire"}),
		},
		{
			description: "With summary filter",
			url:         baseURL + "&summary=continental+Europe",
			calledWith:  decorateQuery(factions.FindManyQuery{Summary: "continental Europe"}),
		},
		{
			description: "With name and summary filters",
			url:         baseURL + "&name=First+French+Empire&summary=continental+Europe",
			calledWith:  decorateQuery(factions.FindManyQuery{Name: "First French Empire", Summary: "continental Europe"}),
		},
	}
}
