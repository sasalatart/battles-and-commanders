package factions

import (
	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/http/middleware"
	"github.com/sasalatart/batcoms/store"
	uuid "github.com/satori/go.uuid"
)

func handleFaction(s store.FactionsFinder) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		faction, err := s.FindOne(domain.Faction{
			ID: c.Locals("factionID").(uuid.UUID),
		})
		if err != nil {
			c.Next(err)
		} else if err := c.JSON(faction); err != nil {
			c.Next(err)
		}
	}
}

// RegisterRoutes registers all factions routes and their handlers in the given *fiber.App
func RegisterRoutes(app *fiber.App, fs store.FactionsFinder) {
	app.Get("/factions/:factionID", middleware.IDFrom("factionID"), handleFaction(fs))
}
