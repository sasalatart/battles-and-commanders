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

func TestFactionsHandlers(t *testing.T) {
	t.Run("GET /factions/:factionID", func(t *testing.T) {
		t.Parallel()

		t.Run("ValidPersistedUUID", func(t *testing.T) {
			factionMock := mocks.Faction()
			app, factionsRepoMock, _ := appWithReposMocks()
			factionsRepoMock.On("FindOne", factions.Faction{
				ID: factionMock.ID,
			}).Return(factionMock, nil)

			httptest.AssertFiberGET(t, app, "/factions/"+factionMock.ID.String(), http.StatusOK, func(res *http.Response) {
				factionsRepoMock.AssertExpectations(t)
				httptest.AssertJSONFaction(t, res, factionMock)
			})
		})

		t.Run("ValidNonPersistedUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			app, factionsRepoMock, _ := appWithReposMocks()
			factionsRepoMock.On("FindOne", factions.Faction{
				ID: uuid,
			}).Return(factions.Faction{}, domain.ErrNotFound)

			httptest.AssertFailedFiberGET(t, app, "/factions/"+uuid.String(), *fiber.ErrNotFound)
			factionsRepoMock.AssertExpectations(t)
		})

		t.Run("InvalidUUID", func(t *testing.T) {
			invalidUUID := "invalid-uuid"
			app, factionsRepoMock, _ := appWithReposMocks()

			httptest.AssertFailedFiberGET(t, app, "/factions/"+invalidUUID, *fiber.ErrBadRequest)
			factionsRepoMock.AssertNotCalled(t, "FindOne")
		})
	})

	t.Run("GET /factions", func(t *testing.T) {
		t.Parallel()

		const page uint = 2
		const pagesMock = 3
		factionsMock := []factions.Faction{mocks.Faction(), mocks.Faction2()}

		cases := []struct {
			description string
			url         string
			calledWith  factions.Query
		}{
			{
				description: "With no filters",
				url:         fmt.Sprintf("/factions?page=%d", page),
				calledWith:  factions.Query{},
			},
			{
				description: "With name filter",
				url:         fmt.Sprintf("/factions?page=%d&name=First+French+Empire", page),
				calledWith:  factions.Query{Name: "First French Empire"},
			},
			{
				description: "With summary filter",
				url:         fmt.Sprintf("/factions?page=%d&summary=continental+Europe", page),
				calledWith:  factions.Query{Summary: "continental Europe"},
			},
			{
				description: "With name and summary filters",
				url:         fmt.Sprintf("/factions?page=%d&name=First+French+Empire&summary=continental+Europe", page),
				calledWith:  factions.Query{Name: "First French Empire", Summary: "continental Europe"},
			},
		}
		for _, c := range cases {
			t.Run(c.description, func(t *testing.T) {
				app, factionsRepoMock, _ := appWithReposMocks()
				factionsRepoMock.On("FindMany", c.calledWith, page).
					Return(factionsMock, pagesMock, nil)
				httptest.AssertFiberGET(t, app, c.url, http.StatusOK, func(res *http.Response) {
					factionsRepoMock.AssertExpectations(t)
					httptest.AssertHeaderPages(t, res, uint(pagesMock))
					httptest.AssertJSONFactions(t, res, factionsMock)
				})
			})
		}
	})

	t.Run("GET /commanders/:commanderID/factions", func(t *testing.T) {
		t.Parallel()

		const page uint = 2
		buildURL := func(commanderID string) string {
			return fmt.Sprintf("/commanders/%s/factions?page=%d", commanderID, page)
		}

		t.Run("ValidPersistedCommanderUUID", func(t *testing.T) {
			const pagesMock = 3
			commanderMock := mocks.Commander()
			factionsMock := []factions.Faction{mocks.Faction(), mocks.Faction2()}

			cases := []struct {
				description string
				url         string
				calledWith  factions.Query
			}{
				{
					description: "With no filters",
					url:         buildURL(commanderMock.ID.String()),
					calledWith:  factions.Query{CommanderID: commanderMock.ID},
				},
				{
					description: "With name filter",
					url:         buildURL(commanderMock.ID.String()) + "&name=First+French+Empire",
					calledWith:  factions.Query{CommanderID: commanderMock.ID, Name: "First French Empire"},
				},
				{
					description: "With summary filter",
					url:         buildURL(commanderMock.ID.String()) + "&summary=continental+Europe",
					calledWith:  factions.Query{CommanderID: commanderMock.ID, Summary: "continental Europe"},
				},
				{
					description: "With name and summary filters",
					url:         buildURL(commanderMock.ID.String()) + "&name=First+French+Empire&summary=continental+Europe",
					calledWith:  factions.Query{CommanderID: commanderMock.ID, Name: "First French Empire", Summary: "continental Europe"},
				},
			}
			for _, c := range cases {
				t.Run(c.description, func(t *testing.T) {
					app, factionsRepoMock, commandersRepoMock := appWithReposMocks()
					commandersRepoMock.On("FindOne", commanders.Commander{
						ID: commanderMock.ID,
					}).Return(commanderMock, nil)
					factionsRepoMock.On("FindMany", c.calledWith, page).
						Return(factionsMock, pagesMock, nil)

					httptest.AssertFiberGET(t, app, c.url, http.StatusOK, func(res *http.Response) {
						factionsRepoMock.AssertExpectations(t)
						httptest.AssertHeaderPages(t, res, uint(pagesMock))
						httptest.AssertJSONFactions(t, res, factionsMock)
					})
				})
			}
		})

		t.Run("ValidNonPersistedCommanderUUID", func(t *testing.T) {
			uuid := uuid.NewV4()
			app, factionsRepoMock, commandersRepoMock := appWithReposMocks()
			commandersRepoMock.On("FindOne", commanders.Commander{
				ID: uuid,
			}).Return(commanders.Commander{}, domain.ErrNotFound)

			httptest.AssertFailedFiberGET(t, app, buildURL(uuid.String()), *fiber.ErrNotFound)
			commandersRepoMock.AssertExpectations(t)
			factionsRepoMock.AssertNotCalled(t, "FindMany")
		})

		t.Run("InvalidCommanderUUID", func(t *testing.T) {
			invalidUUID := "invalid-uuid"
			app, factionsRepoMock, commandersRepoMock := appWithReposMocks()

			httptest.AssertFailedFiberGET(t, app, buildURL(invalidUUID), *fiber.ErrBadRequest)
			commandersRepoMock.AssertNotCalled(t, "FindMany")
			factionsRepoMock.AssertNotCalled(t, "FindOne")
		})
	})
}
