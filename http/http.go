package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sasalatart/batcoms/domain/battles"
	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/domain/factions"
	"github.com/sasalatart/batcoms/http/handlers"
)

// Setup sets up a new fiber server, registers middleware, route handlers, and returns a pointer to it
func Setup(fr factions.Reader, cr commanders.Reader, br battles.Reader, debug bool) *fiber.App {
	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())
	handlers.Register(app, fr, cr, br)
	return app
}
