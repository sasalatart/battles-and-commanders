package handlers_test

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sasalatart/batcoms/http"
	"github.com/sasalatart/batcoms/mocks"
)

func appWithReposMocks() (*fiber.App, *mocks.FactionsRepository, *mocks.CommandersRepository, *mocks.BattlesRepository) {
	factionsRepoMock := new(mocks.FactionsRepository)
	commandersRepoMock := new(mocks.CommandersRepository)
	battlesRepoMock := new(mocks.BattlesRepository)
	app := http.Setup(factionsRepoMock, commandersRepoMock, battlesRepoMock, true)
	return app, factionsRepoMock, commandersRepoMock, battlesRepoMock
}
