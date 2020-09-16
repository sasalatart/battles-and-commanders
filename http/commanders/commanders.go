package commanders

import (
	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/http/middleware"
	"github.com/sasalatart/batcoms/store"
)

// RegisterRoutes registers all commanders routes and their handlers in the given *fiber.App
func RegisterRoutes(app *fiber.App, ff store.FactionsFinder, cf store.CommandersFinder) {
	app.Get("/commanders/:commanderID",
		middleware.WithCommander(cf),
		middleware.JSONFrom("commander"),
	)

	app.Get("/commanders",
		middleware.WithPage(),
		middleware.WithCommanders(cf),
		middleware.JSONFrom("commanders"),
	)

	app.Get("/factions/:factionID/commanders",
		middleware.WithPage(),
		middleware.WithFaction(ff),
		middleware.WithCommanders(cf),
		middleware.JSONFrom("commanders"),
	)
}
