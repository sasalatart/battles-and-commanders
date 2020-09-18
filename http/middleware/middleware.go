package middleware

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/store"
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
func WithFaction(s store.FactionsFinder) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		id, err := uuid.FromString(ctx.Params("factionID"))
		if err != nil {
			ctx.Next(fiber.ErrBadRequest)
			return
		}
		faction, err := s.FindOne(domain.Faction{ID: id})
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
func WithCommander(s store.CommandersFinder) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		id, err := uuid.FromString(ctx.Params("commanderID"))
		if err != nil {
			ctx.Next(fiber.ErrBadRequest)
			return
		}
		commander, err := s.FindOne(domain.Commander{ID: id})
		if err != nil {
			ctx.Next(err)
			return
		}
		ctx.Locals("commander", commander)
		ctx.Next()
	}
}

// WithCommanders middleware finds commanders according to the optional :factionID URL parameter and
// the optional "page" query parameter (falling back to 1), and sets them into ctx.Locals under the
// key "commanders"
func WithCommanders(s store.CommandersFinder) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		var page uint = 1
		if queryPage, hasPage := ctx.Locals("page").(uint); hasPage {
			page = queryPage
		}
		query := store.CommandersQuery{Name: ctx.Query("name"), Summary: ctx.Query("summary")}
		if faction, hasFaction := ctx.Locals("faction").(domain.Faction); hasFaction {
			query.FactionID = faction.ID
		}
		commanders, pages, err := s.FindMany(query, page)
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
func WithBattle(s store.BattlesFinder) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		id, err := uuid.FromString(ctx.Params("battleID"))
		if err != nil {
			ctx.Next(fiber.ErrBadRequest)
			return
		}
		battle, err := s.FindOne(domain.Battle{ID: id})
		if err != nil {
			ctx.Next(err)
			return
		}
		ctx.Locals("battle", battle)
		ctx.Next()
	}
}
