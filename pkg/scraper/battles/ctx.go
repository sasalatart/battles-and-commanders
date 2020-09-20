package battles

import (
	"github.com/gocolly/colly"
	"github.com/sasalatart/batcoms/domain/wikibattles"
)

type battleCtx struct {
	battle    *wikibattles.Battle
	collector *colly.Collector
	err       error
}

func (ctx *battleCtx) abortable(f colly.HTMLCallback) colly.HTMLCallback {
	return func(e *colly.HTMLElement) {
		if ctx.err == nil {
			f(e)
		}
	}
}
