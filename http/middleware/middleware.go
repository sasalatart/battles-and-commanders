package middleware

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/domain/battles"
	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/domain/factions"
	"github.com/sasalatart/batcoms/pkg/dates"
	uuid "github.com/satori/go.uuid"
)

// WithPage middleware parses the optional "page" query parameter, validates it, and then stores its
// value into ctx.Locals under the key "page"
func WithPage() func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		page, err := strconv.Atoi(ctx.Query("page", "1"))
		if err != nil || page <= 0 {
			ctx.Next(newErrBadRequest("Invalid page"))
			return
		}
		ctx.Locals("page", page)
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
			ctx.Next(newErrBadRequest("Invalid FactionID"))
			return
		}
		faction, err := r.FindOne(factions.FindOneQuery{ID: id})
		if err != nil {
			ctx.Next(handleFindOneError(err, "Faction"))
			return
		}
		ctx.Locals("faction", faction)
		ctx.Next()
	}
}

// WithFactions middleware finds factions according to the optional :commanderID URL parameter and
// the optional "page" query parameter (falling back to 1), and sets them into ctx.Locals under the
// key "factions". When present, it will also use the "name" and "summary" query parameters to
// refine this search
func WithFactions(r factions.Reader) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		query := factions.FindManyQuery{
			Name:        ctx.Query("name"),
			Summary:     ctx.Query("summary"),
			CommanderID: commanderIDFromLocals(ctx),
		}
		factions, pages, err := r.FindMany(query, pageFromLocals(ctx))
		if err != nil {
			ctx.Next(err)
			return
		}
		ctx.Set("x-pages", fmt.Sprint(pages))
		ctx.Locals("factions", factions)
		ctx.Next()
	}
}

// WithCommander middleware sets the commander corresponding to the :commanderID URL parameter into
// ctx.Locals under the key "commander"
func WithCommander(r commanders.Reader) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		id, err := uuid.FromString(ctx.Params("commanderID"))
		if err != nil {
			ctx.Next(newErrBadRequest("Invalid CommanderID"))
			return
		}
		commander, err := r.FindOne(commanders.FindOneQuery{ID: id})
		if err != nil {
			ctx.Next(handleFindOneError(err, "Commander"))
			return
		}
		ctx.Locals("commander", commander)
		ctx.Next()
	}
}

// WithCommanders middleware finds commanders according to the optional :factionID URL parameter and
// the optional "page" query parameter (falling back to 1), and sets them into ctx.Locals under the
// key "commanders". When present, it will also use the "name" and "summary" query parameters to
// refine this search
func WithCommanders(r commanders.Reader) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		query := commanders.FindManyQuery{
			Name:      ctx.Query("name"),
			Summary:   ctx.Query("summary"),
			FactionID: factionIDFromLocals(ctx),
		}
		commanders, pages, err := r.FindMany(query, pageFromLocals(ctx))
		if err != nil {
			ctx.Next(err)
			return
		}
		ctx.Set("x-pages", fmt.Sprint(pages))
		ctx.Locals("commanders", commanders)
		ctx.Next()
	}
}

// WithBattles middleware finds battles according to the optional :factionID or :commanderID URL
// parameters and the optional "page" query parameter (falling back to 1), and sets them into
// ctx.Locals under the key "battles". When present, it will also use the "name", "summary", "place"
// and "result" query parameters to refine this search
func WithBattles(r battles.Reader) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		var fromDate dates.Historic
		if ctx.Query("fromDate") != "" {
			date, err := dates.New(ctx.Query("fromDate"))
			if err != nil {
				ctx.Next(newErrBadRequest("Invalid fromDate, must be in YYYY-MM-DD format"))
				return
			}
			fromDate = date.ToBeginning()
		}
		var toDate dates.Historic
		if ctx.Query("toDate") != "" {
			date, err := dates.New(ctx.Query("toDate"))
			if err != nil {
				ctx.Next(newErrBadRequest("Invalid toDate, must be in YYYY-MM-DD format"))
				return
			}
			toDate = date.ToEnd()
		}
		query := battles.FindManyQuery{
			Name:        ctx.Query("name"),
			Summary:     ctx.Query("summary"),
			Place:       ctx.Query("place"),
			Result:      ctx.Query("result"),
			FromDate:    fromDate,
			ToDate:      toDate,
			FactionID:   factionIDFromLocals(ctx),
			CommanderID: commanderIDFromLocals(ctx),
		}
		battles, pages, err := r.FindMany(query, pageFromLocals(ctx))
		if err != nil {
			ctx.Next(err)
			return
		}
		ctx.Set("x-pages", fmt.Sprint(pages))
		ctx.Locals("battles", battles)
		ctx.Next()
	}
}

// WithBattle middleware sets the battle corresponding to the :battleID URL parameter into
// ctx.Locals under the key "battle"
func WithBattle(r battles.Reader) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		id, err := uuid.FromString(ctx.Params("battleID"))
		if err != nil {
			ctx.Next(newErrBadRequest("Invalid BattleID"))
			return
		}
		battle, err := r.FindOne(battles.FindOneQuery{ID: id})
		if err != nil {
			ctx.Next(handleFindOneError(err, "Battle"))
			return
		}
		ctx.Locals("battle", battle)
		ctx.Next()
	}
}

func pageFromLocals(ctx *fiber.Ctx) int {
	page := 1
	if queryPage, hasPage := ctx.Locals("page").(int); hasPage {
		page = queryPage
	}
	return page
}

func factionIDFromLocals(ctx *fiber.Ctx) uuid.UUID {
	if faction, hasFaction := ctx.Locals("faction").(factions.Faction); hasFaction {
		return faction.ID
	}
	return uuid.Nil
}

func commanderIDFromLocals(ctx *fiber.Ctx) uuid.UUID {
	if commander, hasCommander := ctx.Locals("commander").(commanders.Commander); hasCommander {
		return commander.ID
	}
	return uuid.Nil
}

func handleFindOneError(err error, resourceName string) error {
	if err != domain.ErrNotFound {
		return err
	}
	return newErrNotFound(resourceName)
}
