package middleware

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/domain/battles"
	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/domain/factions"
	uuid "github.com/satori/go.uuid"
)

// WithPage middleware parses the optional "page" query parameter, validates it, and then stores its
// value into ctx.Locals under the key "page"
func WithPage() func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		page, err := strconv.Atoi(ctx.Query("page", "1"))
		if err != nil || page <= 0 {
			ctx.Next(fiber.ErrBadRequest)
			return
		}
		ctx.Locals("page", uint(page))
		ctx.Next()
	}
}

// JSONFrom middleware renders a JSON response from the provided ctx.Locals key
func JSONFrom(key string) func(ctx *fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		if err := ctx.JSON(ctx.Locals(key)); err != nil {
			ctx.Next(err)
		}
	}
}

// WithFaction middleware sets the faction corresponding to the :factionID URL parameter into
// ctx.Locals under the key "faction"
func WithFaction(r factions.Reader) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		id, err := uuid.FromString(ctx.Params("factionID"))
		if err != nil {
			ctx.Next(fiber.ErrBadRequest)
			return
		}
		faction, err := r.FindOne(factions.Faction{ID: id})
		if err != nil {
			ctx.Next(err)
			return
		}
		ctx.Locals("faction", faction)
		ctx.Next()
	}
}

// WithCommander middleware sets the commander corresponding to the :commanderID URL parameter into
// ctx.Locals under the key "commander"
func WithCommander(r commanders.Reader) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		id, err := uuid.FromString(ctx.Params("commanderID"))
		if err != nil {
			ctx.Next(fiber.ErrBadRequest)
			return
		}
		commander, err := r.FindOne(commanders.Commander{ID: id})
		if err != nil {
			ctx.Next(err)
			return
		}
		ctx.Locals("commander", commander)
		ctx.Next()
	}
}

// WithFactions middleware finds factions according to the optional :commanderID URL parameter and
// the optional "page" query parameter (falling back to 1), and sets them into ctx.Locals under the
// key "factions". When present, it will also use the "name" and "summary" query parameters to
// refine this search
func WithFactions(r factions.Reader) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		var page uint = 1
		if queryPage, hasPage := ctx.Locals("page").(uint); hasPage {
			page = queryPage
		}
		query := factions.Query{Name: ctx.Query("name"), Summary: ctx.Query("summary")}
		if commander, hasCommander := ctx.Locals("commander").(commanders.Commander); hasCommander {
			query.CommanderID = commander.ID
		}
		factions, pages, err := r.FindMany(query, page)
		if err != nil {
			ctx.Next(err)
			return
		}
		ctx.Set("x-pages", fmt.Sprint(pages))
		ctx.Locals("factions", factions)
		ctx.Next()
	}
}

// WithCommanders middleware finds commanders according to the optional :factionID URL parameter and
// the optional "page" query parameter (falling back to 1), and sets them into ctx.Locals under the
// key "commanders". When present, it will also use the "name" and "summary" query parameters to
// refine this search
func WithCommanders(r commanders.Reader) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		var page uint = 1
		if queryPage, hasPage := ctx.Locals("page").(uint); hasPage {
			page = queryPage
		}
		query := commanders.Query{Name: ctx.Query("name"), Summary: ctx.Query("summary")}
		if faction, hasFaction := ctx.Locals("faction").(factions.Faction); hasFaction {
			query.FactionID = faction.ID
		}
		commanders, pages, err := r.FindMany(query, page)
		if err != nil {
			ctx.Next(err)
			return
		}
		ctx.Set("x-pages", fmt.Sprint(pages))
		ctx.Locals("commanders", commanders)
		ctx.Next()
	}
}

// WithBattle middleware sets the battle corresponding to the :battleID URL parameter into
// ctx.Locals under the key "battle"
func WithBattle(r battles.Reader) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		id, err := uuid.FromString(ctx.Params("battleID"))
		if err != nil {
			ctx.Next(fiber.ErrBadRequest)
			return
		}
		battle, err := r.FindOne(battles.Battle{ID: id})
		if err != nil {
			ctx.Next(err)
			return
		}
		ctx.Locals("battle", battle)
		ctx.Next()
	}
}
