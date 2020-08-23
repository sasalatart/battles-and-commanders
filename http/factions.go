package http

import (
	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/store"
	uuid "github.com/satori/go.uuid"
)

func handleFaction(s store.FactionsFinder) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		faction, err := s.FindOne(domain.Faction{
			ID: c.Locals("factionId").(uuid.UUID),
		})
		if err != nil {
			c.Next(err)
		} else if err := c.JSON(faction); err != nil {
			c.Next(err)
		}
	}
}

func registerFactionsRoutes(app *fiber.App, fs store.FactionsFinder) {
	app.Get("/factions/:factionId", idFrom("factionId"), handleFaction(fs))
}
