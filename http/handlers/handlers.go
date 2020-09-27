package handlers

import (
	"errors"

	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/domain/battles"
	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/domain/factions"
	"github.com/sasalatart/batcoms/http/middleware"
)

// Register registers all factions, commanders and battles routes together with their handlers in
// the given *fiber.App
func Register(app *fiber.App, fr factions.Reader, cr commanders.Reader, br battles.Reader) {
	app.Get("/factions/:factionID",
		middleware.WithFaction(fr),
		middleware.JSONFrom("faction"),
	)

	app.Get("/factions",
		middleware.WithPage(),
		middleware.WithFactions(fr),
		middleware.JSONFrom("factions"),
	)

	app.Get("/commanders/:commanderID/factions",
		middleware.WithPage(),
		middleware.WithCommander(cr),
		middleware.WithFactions(fr),
		middleware.JSONFrom("factions"),
	)

	app.Get("/commanders/:commanderID",
		middleware.WithCommander(cr),
		middleware.JSONFrom("commander"),
	)

	app.Get("/commanders",
		middleware.WithPage(),
		middleware.WithCommanders(cr),
		middleware.JSONFrom("commanders"),
	)

	app.Get("/factions/:factionID/commanders",
		middleware.WithPage(),
		middleware.WithFaction(fr),
		middleware.WithCommanders(cr),
		middleware.JSONFrom("commanders"),
	)

	app.Get("/battles/:battleID",
		middleware.WithBattle(br),
		middleware.JSONFrom("battle"),
	)

	app.Get("/battles",
		middleware.WithPage(),
		middleware.WithBattles(br),
		middleware.JSONFrom("battles"),
	)

	app.Get("/factions/:factionID/battles",
		middleware.WithPage(),
		middleware.WithFaction(fr),
		middleware.WithBattles(br),
		middleware.JSONFrom("battles"),
	)

	app.Get("/commanders/:commanderID/battles",
		middleware.WithPage(),
		middleware.WithCommander(cr),
		middleware.WithBattles(br),
		middleware.JSONFrom("battles"),
	)
}

// ErrorsHandlerFactory creates a *fiber.App ErrorHandler, used to fine-tune responses when a
// request is not successful. This can be in debug or non-debug mode
func ErrorsHandlerFactory(debug bool) func(ctx *fiber.Ctx, err error) {
	return func(ctx *fiber.Ctx, err error) {
		code := fiber.StatusInternalServerError
		message := "Internal server error"
		if debug {
			message = err.Error()
		}
		if errors.Is(err, domain.ErrNotFound) {
			err = fiber.ErrNotFound
		}
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
			message = e.Message
		}
		ctx.Status(code).SendString(message)
	}
}
