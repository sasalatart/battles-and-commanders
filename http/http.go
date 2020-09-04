package http

import (
	"errors"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/sasalatart/batcoms/http/battles"
	"github.com/sasalatart/batcoms/http/commanders"
	"github.com/sasalatart/batcoms/http/factions"
	"github.com/sasalatart/batcoms/store"
)

// Setup sets up a new fiber server, registers the route handlers, and returns a pointer to it.
func Setup(fs store.FactionsFinder, cs store.CommandersFinder, bs store.BattlesFinder, debug bool) *fiber.App {
	app := fiber.New()
	app.Settings.ErrorHandler = errorsHandlerFactory(debug)
	app.Use(middleware.Recover())
	app.Use(middleware.Logger())
	factions.RegisterRoutes(app, fs)
	commanders.RegisterRoutes(app, cs)
	battles.RegisterRoutes(app, bs)
	return app
}

func errorsHandlerFactory(debug bool) func(ctx *fiber.Ctx, err error) {
	return func(ctx *fiber.Ctx, err error) {
		code := fiber.StatusInternalServerError
		message := "Internal server error"
		if debug {
			message = err.Error()
		}
		if errors.Is(err, store.ErrNotFound) {
			err = fiber.ErrNotFound
		}
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
			message = e.Message
		}
		ctx.Status(code).SendString(message)
	}
}
