package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/services/parser"
)

func subscribeNumbersFor(c *colly.Collector, trTitle string, storeIn *domain.SideNumbers) {
	trID := "batcoms-" + strings.ReplaceAll(strings.ToLower(trTitle), " ", "")
	subscribeInfoBoxID(c, trTitle, trID)
	twoSides := false
	c.OnHTML(fmt.Sprintf("#%s > td:first-child", trID), func(e *colly.HTMLElement) {
		storeIn.A = parser.Clean(e.Text)
	})
	c.OnHTML(fmt.Sprintf("#%s > td:nth-child(2)", trID), func(e *colly.HTMLElement) {
		storeIn.B = parser.Clean(e.Text)
		twoSides = true
	})
	c.OnHTML(fmt.Sprintf("#%s + tr>td[colspan='2']", trID), func(e *colly.HTMLElement) {
		textInTags := e.DOM.ChildrenFiltered("*").Text()
		storeIn.AB = parser.Clean(strings.ReplaceAll(e.Text, textInTags, ""))
	})
	c.OnScraped(func(_ *colly.Response) {
		if twoSides {
			return
		}
		storeIn.AB = storeIn.A
		storeIn.A = ""
		storeIn.B = ""
	})
}

func (s *Scraper) subscribeStrength(c *colly.Collector, b *domain.SBattle) {
	subscribeNumbersFor(c, "Strength", &b.Strength)
}

func (s *Scraper) subscribeCasualties(c *colly.Collector, b *domain.SBattle) {
	subscribeNumbersFor(c, "Casualties and losses", &b.Casualties)
}
