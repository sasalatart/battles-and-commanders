package service

import (
	"strings"

	"github.com/gocolly/colly"
)

const infoBoxSelector = ".infobox.vevent > tbody"

func setInfoBoxID(c *colly.Collector, title string, id string) {
	c.OnHTML(infoBoxSelector, func(e *colly.HTMLElement) {
		e.ForEachWithBreak("th", func(_ int, c *colly.HTMLElement) bool {
			if strings.ToLower(c.Text) != strings.ToLower(title) {
				return true
			}

			c.DOM.Parent().Next().SetAttr("id", id)
			return false
		})
	})
}
