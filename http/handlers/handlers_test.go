package handlers_test

import (
	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/http"
	"github.com/sasalatart/batcoms/mocks"
)

func appWithReposMocks() (*fiber.App, *mocks.FactionsRepository, *mocks.CommandersRepository) {
	factionsRepoMock := new(mocks.FactionsRepository)
	commandersRepoMock := new(mocks.CommandersRepository)
	app := http.Setup(factionsRepoMock, commandersRepoMock, new(mocks.BattlesRepository), true)
	return app, factionsRepoMock, commandersRepoMock
}
