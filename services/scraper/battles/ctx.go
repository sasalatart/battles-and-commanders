package battles

import (
	"github.com/gocolly/colly"
	"github.com/sasalatart/batcoms/domain"
)

type battleCtx struct {
	battle    *domain.SBattle
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
