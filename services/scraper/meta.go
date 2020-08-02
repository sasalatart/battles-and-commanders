package scraper

import (
	"strings"

	"github.com/gocolly/colly"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/services/parser"
)

func (s *Scraper) subscribeMeta(c *colly.Collector, b *domain.SBattle) {
	c.OnHTML(infoBoxSelector, func(e *colly.HTMLElement) {
		b.Name = parser.Clean(e.ChildText("th.summary"))
		e.ForEachWithBreak("tr:nth-child(2) > td", func(_ int, c *colly.HTMLElement) bool {
			if !strings.Contains(strings.ToLower(c.Text), "part of") {
				return true
			}

			b.PartOf = parser.Clean(c.Text)
			return false
		})
	})

	c.OnHTML(infoBoxSelector+" tbody", func(e *colly.HTMLElement) {
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

		b.Date = search("td", "date")
		b.Result = search("td", "result", "status")
		b.TerritorialChanges = search("td", "territorialchanges")
		b.Location.Place = search(".location", "location")

		redundantCoordinates := search(".location #coordinates", "location")
		redundantCoordinates = strings.ReplaceAll(redundantCoordinates, "Coordinates: ", "")
		b.Location.Place = strings.ReplaceAll(b.Location.Place, redundantCoordinates, "")
		b.Location.Place = strings.ReplaceAll(b.Location.Place, "Coordinates: ", "")
	})

	c.OnHTML("#coordinates", func(e *colly.HTMLElement) {
		b.Location.Latitude = e.ChildText(".latitude")
		b.Location.Longitude = e.ChildText(".longitude")
	})
}
