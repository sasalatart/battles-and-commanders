package battles

import (
	"strings"

	"github.com/gocolly/colly"
	"github.com/sasalatart/batcoms/domain/statistics"
	"github.com/sasalatart/batcoms/pkg/strclean"
)

func subscribeNumbersFor(ctx *battleCtx, title string, sideNumbers *statistics.SideNumbers) {
	twoSides := false
	customID := customID(title)
	subscribeSetInfoBoxID(ctx, title, customID)

	ctx.collector.OnHTML(sideNumbersSelector(sideASelector, customID), ctx.abortable(func(e *colly.HTMLElement) {
		sideNumbers.A = strclean.Apply(e.Text)
	}))

	ctx.collector.OnHTML(sideNumbersSelector(sideBSelector, customID), ctx.abortable(func(e *colly.HTMLElement) {
		sideNumbers.B = strclean.Apply(e.Text)
		twoSides = true
	}))

	ctx.collector.OnHTML(sideNumbersSelector(sideABSelector, customID), ctx.abortable(func(e *colly.HTMLElement) {
		textInTags := e.DOM.ChildrenFiltered("*").Text()
		sideNumbers.AB = strclean.Apply(strings.ReplaceAll(e.Text, textInTags, ""))
	}))

	ctx.collector.OnScraped(func(_ *colly.Response) {
		if ctx.err != nil || twoSides {
			return
		}
		sideNumbers.AB = sideNumbers.A
		sideNumbers.A = ""
		sideNumbers.B = ""
	})
}

func (s *Scraper) subscribeStrength(ctx *battleCtx) {
	subscribeNumbersFor(ctx, "Strength", &ctx.battle.Strength)
}

func (s *Scraper) subscribeCasualties(ctx *battleCtx) {
	subscribeNumbersFor(ctx, "Casualties and losses", &ctx.battle.Casualties)
}
