package battles

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sasalatart/batcoms/services/parser"
)

var coordsSep = regexp.MustCompile(`^(.*?)\d+°.*`)

func (s *Scraper) subscribeMeta(ctx *battleCtx) {
	ctx.collector.OnHTML(infoBoxSelector, ctx.abortable(func(e *colly.HTMLElement) {
		e.ForEachWithBreak(partOfSelector, func(_ int, c *colly.HTMLElement) bool {
			if !strings.Contains(strings.ToLower(c.Text), "part of") {
				return true
			}
			ctx.battle.PartOf = parser.Clean(c.Text)
			return false
		})
	}))

	ctx.collector.OnHTML(infoBoxSelector, ctx.abortable(func(e *colly.HTMLElement) {
		search := func(childSelector string, toSearch ...string) string {
			var res string
			e.ForEachWithBreak("tr", func(_ int, c *colly.HTMLElement) bool {
				for _, s := range toSearch {
					if s == strings.ToLower(c.ChildText("th")) {
						res = parser.Clean(c.ChildText(childSelector))
						return false
					}
				}
				return true
			})
			return res
		}

		ctx.battle.Date = search("td", "date")
		if ctx.battle.Date == "" {
			ctx.err = ErrNoDate
			return
		}

		ctx.battle.Result = search("td", "result", "status")
		if ctx.battle.Result == "" {
			ctx.err = ErrNoResult
			return
		}

		ctx.battle.Location.Place = strings.Trim(coordsSep.ReplaceAllString(search(".location", "location"), "$1"), " ")
		if ctx.battle.Location.Place == "" {
			ctx.err = ErrNoPlace
			return
		}

		ctx.battle.TerritorialChanges = search("td", "territorialchanges")
	}))

	ctx.collector.OnHTML(coordinatesSelector, ctx.abortable(func(e *colly.HTMLElement) {
		ctx.battle.Location.Latitude = e.ChildText(".latitude")
		ctx.battle.Location.Longitude = e.ChildText(".longitude")
	}))
}
