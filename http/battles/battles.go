package battles

import (
	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/http/middleware"
	"github.com/sasalatart/batcoms/store"
	uuid "github.com/satori/go.uuid"
)

func handleBattle(s store.BattlesFinder) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		battle, err := s.FindOne(domain.Battle{
			ID: c.Locals("battleID").(uuid.UUID),
		})
		if err != nil {
			c.Next(err)
		} else if err := c.JSON(battle); err != nil {
			c.Next(err)
		}
	}
}

// RegisterRoutes registers all battles routes and their handlers in the given *fiber.App
func RegisterRoutes(app *fiber.App, cs store.BattlesFinder) {
	app.Get("/battles/:battleID", middleware.IDFrom("battleID"), handleBattle(cs))
}
