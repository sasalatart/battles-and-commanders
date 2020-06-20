package service

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sasalatart/batcoms/parser"
	"github.com/sasalatart/batcoms/scraper/domain"
)

func subscribeNumbersFor(c *colly.Collector, trTitle string, storeIn *struct{ A, B string }) {
	trID := "batcoms-" + strings.ReplaceAll(strings.ToLower(trTitle), " ", "")
	setInfoBoxID(c, trTitle, trID)
	c.OnHTML(fmt.Sprintf("#%s > td:first-child", trID), func(e *colly.HTMLElement) {
		storeIn.A = parser.Clean(e.Text)
	})
	c.OnHTML(fmt.Sprintf("#%s > td:nth-child(2)", trID), func(e *colly.HTMLElement) {
		storeIn.B = parser.Clean(e.Text)
	})
}

func (s *Scraper) subscribeStrength(c *colly.Collector, b *domain.Battle) {
	subscribeNumbersFor(c, "Strength", &b.Strength)
}

func (s *Scraper) subscribeCasualties(c *colly.Collector, b *domain.Battle) {
	subscribeNumbersFor(c, "Casualties and losses", &b.Casualties)
}
