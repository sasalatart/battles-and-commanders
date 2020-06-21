package service

import (
	"strings"

	"github.com/gocolly/colly"
	"github.com/sasalatart/batcoms/parser"
	"github.com/sasalatart/batcoms/scraper/domain"
)

func (s *Scraper) subscribeMeta(c *colly.Collector, b *domain.Battle) {
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
		searchForTRAndSetIn := func(toSearch, childSelector string, toSet *string) {
			e.ForEachWithBreak("tr", func(_ int, c *colly.HTMLElement) bool {
				if strings.ToLower(c.ChildText("th")) != toSearch {
					return true
				}

				*toSet = parser.Clean(c.ChildText(childSelector))
				return false
			})
		}

		searchForTRAndSetIn("date", "td", &b.Date)
		searchForTRAndSetIn("result", "td", &b.Result)
		searchForTRAndSetIn("territorialchanges", "td", &b.TerritorialChanges)
		searchForTRAndSetIn("location", ".location", &b.Location.Place)

		var redundantCoordinates string
		searchForTRAndSetIn("location", ".location #coordinates", &redundantCoordinates)
		redundantCoordinates = strings.ReplaceAll(redundantCoordinates, "Coordinates: ", "")
		b.Location.Place = strings.ReplaceAll(b.Location.Place, redundantCoordinates, "")
		b.Location.Place = strings.ReplaceAll(b.Location.Place, "Coordinates: ", "")
	})

	c.OnHTML("#coordinates", func(e *colly.HTMLElement) {
		b.Location.Latitude = e.ChildText(".latitude")
		b.Location.Longitude = e.ChildText(".longitude")
	})
}
