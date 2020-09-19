package http

import (
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/sasalatart/batcoms/domain/battles"
	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/domain/factions"
	"github.com/sasalatart/batcoms/http/handlers"
)

// Setup sets up a new fiber server, registers middleware, route handlers, and returns a pointer to it
func Setup(fr factions.Reader, cr commanders.Reader, br battles.Reader, debug bool) *fiber.App {
	app := fiber.New()
	app.Settings.ErrorHandler = handlers.ErrorsHandlerFactory(debug)
	app.Use(middleware.Recover())
	app.Use(middleware.Logger())
	handlers.Register(app, fr, cr, br)
	return app
}
