package commanders

import (
	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/http/middleware"
	"github.com/sasalatart/batcoms/store"
	uuid "github.com/satori/go.uuid"
)

func handleCommander(s store.CommandersFinder) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		commander, err := s.FindOne(domain.Commander{
			ID: c.Locals("commanderID").(uuid.UUID),
		})
		if err != nil {
			c.Next(err)
		} else if err := c.JSON(commander); err != nil {
			c.Next(err)
		}
	}
}

// RegisterRoutes registers all commanders routes and their handlers in the given *fiber.App
func RegisterRoutes(app *fiber.App, cs store.CommandersFinder) {
	app.Get("/commanders/:commanderID", middleware.IDFrom("commanderID"), handleCommander(cs))
}
