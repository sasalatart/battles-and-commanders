package battles

import (
	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/http/middleware"
	"github.com/sasalatart/batcoms/store"
)

// RegisterRoutes registers all battles routes and their handlers in the given *fiber.App
func RegisterRoutes(app *fiber.App, bf store.BattlesFinder) {
	app.Get("/battles/:battleID", middleware.WithBattle(bf), middleware.JSONFrom("battle"))
}
