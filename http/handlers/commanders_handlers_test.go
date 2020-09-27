package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/domain/factions"
	"github.com/sasalatart/batcoms/http/httptest"
	"github.com/sasalatart/batcoms/mocks"
	uuid "github.com/satori/go.uuid"
)

func TestCommandersHandlers(t *testing.T) {
	t.Run("GET /commanders/:commanderID", func(t *testing.T) {
		t.Parallel()

		t.Run("ValidPersistedUUID", func(t *testing.T) {
			commanderMock := mocks.Commander()
			app, _, commandersRepoMock, _ := appWithReposMocks()
			commandersRepoMock.On("FindOne", commanders.FindOneQuery{
				ID: commanderMock.ID,
			}).Return(commanderMock, nil)

			httptest.AssertFiberGET(t, app, "/commanders/"+commanderMock.ID.String(), http.StatusOK, func(res *http.Response) {
				commandersRepoMock.AssertExpectations(t)
				httptest.AssertJSONCommander(t, res, commanderMock)
			})
		})

		t.Run("ValidNonPersistedUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			app, _, commandersRepoMock, _ := appWithReposMocks()
			commandersRepoMock.On("FindOne", commanders.FindOneQuery{
				ID: uuid,
			}).Return(commanders.Commander{}, domain.ErrNotFound)

			httptest.AssertFailedFiberGET(t, app, "/commanders/"+uuid.String(), *fiber.ErrNotFound)
			commandersRepoMock.AssertExpectations(t)
		})

		t.Run("InvalidUUID", func(t *testing.T) {
			app, _, commandersRepoMock, _ := appWithReposMocks()
			httptest.AssertFailedFiberGET(t, app, "/commanders/invalid-uuid", *fiber.ErrBadRequest)
			commandersRepoMock.AssertNotCalled(t, "FindOne")
		})
	})

	t.Run("GET /commanders", func(t *testing.T) {
		t.Parallel()

		const page = 2
		const pagesMock = 3
		baseURL := fmt.Sprintf("/commanders?page=%d", page)
		commandersMock := []commanders.Commander{mocks.Commander(), mocks.Commander2()}

		cases := buildCommandersCases(baseURL, func(q commanders.FindManyQuery) commanders.FindManyQuery {
			return q
		})
		for _, c := range cases {
			t.Run(c.description, func(t *testing.T) {
				app, _, commandersRepoMock, _ := appWithReposMocks()
				commandersRepoMock.On("FindMany", c.calledWith, page).
					Return(commandersMock, pagesMock, nil)
				httptest.AssertFiberGET(t, app, c.url, http.StatusOK, func(res *http.Response) {
					commandersRepoMock.AssertExpectations(t)
					httptest.AssertHeaderPages(t, res, pagesMock)
					httptest.AssertJSONCommanders(t, res, commandersMock)
				})
			})
		}
	})

	t.Run("GET /factions/:factionID/commanders", func(t *testing.T) {
		t.Parallel()

		const page = 2
		baseURL := func(factionID string) string {
			return fmt.Sprintf("/factions/%s/commanders?page=%d", factionID, page)
		}

		t.Run("ValidPersistedFactionUUID", func(t *testing.T) {
			const pagesMock = 3
			factionMock := mocks.Faction()
			commandersMock := []commanders.Commander{mocks.Commander(), mocks.Commander2()}

			cases := buildCommandersCases(baseURL(factionMock.ID.String()), func(q commanders.FindManyQuery) commanders.FindManyQuery {
				q.FactionID = factionMock.ID
				return q
			})
			for _, c := range cases {
				t.Run(c.description, func(t *testing.T) {
					app, factionsRepoMock, commandersRepoMock, _ := appWithReposMocks()
					factionsRepoMock.On("FindOne", factions.FindOneQuery{
						ID: factionMock.ID,
					}).Return(factionMock, nil)
					commandersRepoMock.On("FindMany", c.calledWith, page).
						Return(commandersMock, pagesMock, nil)

					httptest.AssertFiberGET(t, app, c.url, http.StatusOK, func(res *http.Response) {
						commandersRepoMock.AssertExpectations(t)
						httptest.AssertHeaderPages(t, res, pagesMock)
						httptest.AssertJSONCommanders(t, res, commandersMock)
					})
				})
			}
		})

		t.Run("ValidNonPersistedFactionUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			app, factionsRepoMock, commandersRepoMock, _ := appWithReposMocks()
			factionsRepoMock.On("FindOne", factions.FindOneQuery{
				ID: uuid,
			}).Return(factions.Faction{}, domain.ErrNotFound)

			httptest.AssertFailedFiberGET(t, app, baseURL(uuid.String()), *fiber.ErrNotFound)
			factionsRepoMock.AssertExpectations(t)
			commandersRepoMock.AssertNotCalled(t, "FindMany")
		})

		t.Run("InvalidFactionUUID", func(t *testing.T) {
			app, factionsRepoMock, commandersRepoMock, _ := appWithReposMocks()
			httptest.AssertFailedFiberGET(t, app, baseURL("invalid-uuid"), *fiber.ErrBadRequest)
			factionsRepoMock.AssertNotCalled(t, "FindOne")
			commandersRepoMock.AssertNotCalled(t, "FindMany")
		})
	})
}

type commandersTableCase struct {
	description string
	url         string
	calledWith  commanders.FindManyQuery
}

func buildCommandersCases(baseURL string, decorateQuery func(commanders.FindManyQuery) commanders.FindManyQuery) []commandersTableCase {
	return []commandersTableCase{
		{
			description: "With no filters",
			url:         baseURL,
			calledWith:  decorateQuery(commanders.FindManyQuery{}),
		},
		{
			description: "With name filter",
			url:         baseURL + "&name=napoleon",
			calledWith:  decorateQuery(commanders.FindManyQuery{Name: "napoleon"}),
		},
		{
			description: "With summary filter",
			url:         baseURL + "&summary=napoleonic",
			calledWith:  decorateQuery(commanders.FindManyQuery{Summary: "napoleonic"}),
		},
		{
			description: "With name and summary filters",
			url:         baseURL + "&name=napoleon&summary=napoleonic",
			calledWith:  decorateQuery(commanders.FindManyQuery{Name: "napoleon", Summary: "napoleonic"}),
		},
	}
}
