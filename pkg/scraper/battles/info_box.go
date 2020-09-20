package battles

import (
	"strings"

	"github.com/gocolly/colly"
)

func (s *Scraper) assertHasOneInfoBox(ctx *battleCtx) {
	ctx.collector.OnHTML(contentSelector, ctx.abortable(func(e *colly.HTMLElement) {
		infoBoxAmount := 0
		e.ForEach(infoBoxSelector, func(infoBoxIndex int, c *colly.HTMLElement) {
			infoBoxAmount++
			if infoBoxAmount > 1 {
				ctx.err = ErrMoreThanOneInfoBox
				return
			}
		})
		if infoBoxAmount == 0 {
			ctx.err = ErrNoInfoBox
		}
	}))
}

func subscribeSetInfoBoxID(ctx *battleCtx, title string, id string) {
	ctx.collector.OnHTML(infoBoxSelector, ctx.abortable(func(e *colly.HTMLElement) {
		e.ForEachWithBreak("th", func(_ int, c *colly.HTMLElement) bool {
			if strings.ToLower(c.Text) != strings.ToLower(title) {
				return true
			}
			c.DOM.Parent().Next().SetAttr("id", id)
			return false
		})
	}))
}
