package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/domain/factions"
	batcomshttp "github.com/sasalatart/batcoms/http"
	"github.com/sasalatart/batcoms/http/httptest"
	"github.com/sasalatart/batcoms/mocks"
	uuid "github.com/satori/go.uuid"
)

func TestCommandersHandlers(t *testing.T) {
	appWithReposMocks := func() (*fiber.App, *mocks.FactionsRepository, *mocks.CommandersRepository) {
		factionsRepoMock := new(mocks.FactionsRepository)
		commandersRepoMock := new(mocks.CommandersRepository)
		app := batcomshttp.Setup(factionsRepoMock, commandersRepoMock, new(mocks.BattlesRepository), true)
		return app, factionsRepoMock, commandersRepoMock
	}

	t.Run("GET /commanders/:commanderID", func(t *testing.T) {
		t.Parallel()

		t.Run("ValidPersistedUUID", func(t *testing.T) {
			commanderMock := mocks.Commander()
			app, _, commandersRepoMock := appWithReposMocks()
			commandersRepoMock.On("FindOne", commanders.Commander{
				ID: commanderMock.ID,
			}).Return(commanderMock, nil)

			httptest.AssertFiberGET(t, app, "/commanders/"+commanderMock.ID.String(), http.StatusOK, func(res *http.Response) {
				commandersRepoMock.AssertExpectations(t)
				httptest.AssertJSONCommander(t, res, commanderMock)
			})
		})

		t.Run("ValidNonPersistedUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			app, _, commandersRepoMock := appWithReposMocks()
			commandersRepoMock.On("FindOne", commanders.Commander{
				ID: uuid,
			}).Return(commanders.Commander{}, domain.ErrNotFound)

			httptest.AssertFailedFiberGET(t, app, "/commanders/"+uuid.String(), *fiber.ErrNotFound)
			commandersRepoMock.AssertExpectations(t)
		})

		t.Run("InvalidUUID", func(t *testing.T) {
			invalidUUID := "invalid-uuid"
			app, _, commandersRepoMock := appWithReposMocks()

			httptest.AssertFailedFiberGET(t, app, "/commanders/"+invalidUUID, *fiber.ErrBadRequest)
			commandersRepoMock.AssertNotCalled(t, "FindOne")
		})
	})

	t.Run("GET /commanders", func(t *testing.T) {
		t.Parallel()

		var page uint = 2
		pagesMock := 3
		commandersMock := []commanders.Commander{mocks.Commander(), mocks.Commander2()}

		cases := []struct {
			description string
			url         string
			calledWith  commanders.Query
		}{
			{
				description: "With no filters",
				url:         fmt.Sprintf("/commanders?page=%d", page),
				calledWith:  commanders.Query{},
			},
			{
				description: "With name filter",
				url:         fmt.Sprintf("/commanders?page=%d&name=napoleon", page),
				calledWith:  commanders.Query{Name: "napoleon"},
			},
			{
				description: "With summary filter",
				url:         fmt.Sprintf("/commanders?page=%d&summary=napoleonic", page),
				calledWith:  commanders.Query{Summary: "napoleonic"},
			},
			{
				description: "With name and summary filters",
				url:         fmt.Sprintf("/commanders?page=%d&name=napoleon&summary=napoleonic", page),
				calledWith:  commanders.Query{Name: "napoleon", Summary: "napoleonic"},
			},
		}
		for _, c := range cases {
			t.Run(c.description, func(t *testing.T) {
				app, _, commandersRepoMock := appWithReposMocks()
				commandersRepoMock.On("FindMany", c.calledWith, page).
					Return(commandersMock, pagesMock, nil)
				httptest.AssertFiberGET(t, app, c.url, http.StatusOK, func(res *http.Response) {
					commandersRepoMock.AssertExpectations(t)
					httptest.AssertHeaderPages(t, res, uint(pagesMock))
					httptest.AssertJSONCommanders(t, res, commandersMock)
				})
			})
		}
	})

	t.Run("GET /factions/:factionID/commanders", func(t *testing.T) {
		t.Parallel()

		var page uint = 2
		buildURL := func(factionID string) string {
			return fmt.Sprintf("/factions/%s/commanders?page=%d", factionID, page)
		}

		t.Run("ValidPersistedFactionUUID", func(t *testing.T) {
			pagesMock := 3
			factionMock := mocks.Faction()
			commandersMock := []commanders.Commander{mocks.Commander(), mocks.Commander2()}

			cases := []struct {
				description string
				url         string
				calledWith  commanders.Query
			}{
				{
					description: "With no filters",
					url:         buildURL(factionMock.ID.String()),
					calledWith:  commanders.Query{FactionID: factionMock.ID},
				},
				{
					description: "With name filter",
					url:         buildURL(factionMock.ID.String()) + "&name=napoleon",
					calledWith:  commanders.Query{FactionID: factionMock.ID, Name: "napoleon"},
				},
				{
					description: "With summary filter",
					url:         buildURL(factionMock.ID.String()) + "&summary=napoleonic",
					calledWith:  commanders.Query{FactionID: factionMock.ID, Summary: "napoleonic"},
				},
				{
					description: "With name and summary filters",
					url:         buildURL(factionMock.ID.String()) + "&name=napoleon&summary=napoleonic",
					calledWith:  commanders.Query{FactionID: factionMock.ID, Name: "napoleon", Summary: "napoleonic"},
				},
			}
			for _, c := range cases {
				t.Run(c.description, func(t *testing.T) {
					app, factionsRepoMock, commandersRepoMock := appWithReposMocks()
					factionsRepoMock.On("FindOne", factions.Faction{
						ID: factionMock.ID,
					}).Return(factionMock, nil)
					commandersRepoMock.On("FindMany", c.calledWith, page).
						Return(commandersMock, pagesMock, nil)

					httptest.AssertFiberGET(t, app, c.url, http.StatusOK, func(res *http.Response) {
						commandersRepoMock.AssertExpectations(t)
						httptest.AssertHeaderPages(t, res, uint(pagesMock))
						httptest.AssertJSONCommanders(t, res, commandersMock)
					})
				})
			}
		})

		t.Run("ValidNonPersistedFactionUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			app, factionsRepoMock, commandersRepoMock := appWithReposMocks()
			factionsRepoMock.On("FindOne", factions.Faction{
				ID: uuid,
			}).Return(factions.Faction{}, domain.ErrNotFound)

			httptest.AssertFailedFiberGET(t, app, buildURL(uuid.String()), *fiber.ErrNotFound)
			factionsRepoMock.AssertExpectations(t)
			commandersRepoMock.AssertNotCalled(t, "FindMany")
		})

		t.Run("InvalidFactionUUID", func(t *testing.T) {
			invalidUUID := "invalid-uuid"
			app, factionsRepoMock, commandersRepoMock := appWithReposMocks()

			httptest.AssertFailedFiberGET(t, app, buildURL(invalidUUID), *fiber.ErrBadRequest)
			factionsRepoMock.AssertNotCalled(t, "FindOne")
			commandersRepoMock.AssertNotCalled(t, "FindMany")
		})
	})
}
