package factions

import (
	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/http/middleware"
	"github.com/sasalatart/batcoms/store"
)

// RegisterRoutes registers all factions routes and their handlers in the given *fiber.App
func RegisterRoutes(app *fiber.App, ff store.FactionsFinder) {
	app.Get("/factions/:factionID",
		middleware.WithFaction(ff),
		middleware.JSONFrom("faction"),
	)
}
