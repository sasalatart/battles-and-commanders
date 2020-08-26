package http

import (
	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/domain"
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

func registerBattlesRoutes(app *fiber.App, cs store.BattlesFinder) {
	app.Get("/battles/:battleID", idFrom("battleID"), handleBattle(cs))
}
