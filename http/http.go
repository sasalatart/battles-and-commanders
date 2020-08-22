package http

import (
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
)

func registerRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) {
		c.Send("Hello world!")
	})
}

// Setup sets up a new fiber server, registers the route handlers, and returns a pointer to it.
func Setup() *fiber.App {
	app := fiber.New()
	app.Use(middleware.Logger())
	registerRoutes(app)
	return app
}
