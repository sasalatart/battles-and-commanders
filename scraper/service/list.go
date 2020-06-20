package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/sasalatart/batcoms/scraper/domain"

	"github.com/gocolly/colly"
)

type listStrategy struct {
	urlSuffix string
	selector  string
}

// List scrapes and retrieves the full list of Wikipedia's battles (name & URL only) when grouped by
// centuries
func (s *Scraper) List() []domain.BattleItem {
	battles := []domain.BattleItem{}

	var strategies = [...]listStrategy{
		{"before_301", "td:nth-last-child(2) > a"},
		{"301-1300", "td:nth-last-child(3) > a"},
		{"1301-1600", "td:nth-last-child(3) > a"},
		{"1601-1800", "h2+ul a:first-child"},
		{"1801-1900", "h2+ul a:first-child"},
		{"1901-2000", "h2+ul a:first-child"},
		{"since_2001", "td:first-child > a"},
	}
	for _, strategy := range strategies {
		func(ls listStrategy, battles *[]domain.BattleItem) {
			url := fmt.Sprintf("https://en.wikipedia.org/wiki/List_of_battles_%s", ls.urlSuffix)
			subscribe := func(c *colly.Collector) {
				c.OnHTML(fmt.Sprintf("#content %s", ls.selector), func(e *colly.HTMLElement) {
					if e.Attr("class") == "new" {
						return
					}

					href := e.Attr("href")
					if strings.Contains(href, "://") && !strings.Contains(href, "wikipedia.org") {
						return
					}

					url := "https://en.wikipedia.org" + href
					*battles = append(*battles, domain.BattleItem{Name: e.Text, URL: url})
				})
			}
			if err := s.do(url, subscribe); err != nil {
				log.Printf("Error scraping list in %s: %s", url, err)
			}
		}(strategy, &battles)
	}

	feedbackMessage := fmt.Sprintf("There are %d battles that can be scraped\n", len(battles))
	s.logger.Write([]byte(feedbackMessage))
	return battles
}
