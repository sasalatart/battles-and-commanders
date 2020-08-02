package scraper

import (
	"fmt"
	"log"

	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/services/parser"
	"github.com/sasalatart/batcoms/services/scraper/names"
	"github.com/sasalatart/batcoms/services/scraper/urls"

	"github.com/gocolly/colly"
)

var battlesLists = []string{
	"/Battles_of_the_Seven_Years%27_War",
	"/List_of_American_Civil_War_battles",
	"/List_of_American_Revolutionary_War_battles",
	"/List_of_battles_(alphabetical)",
	"/List_of_battles_(geographic)",
	"/List_of_battles_301-1300",
	"/List_of_battles_1301-1600",
	"/List_of_battles_1601-1800",
	"/List_of_battles_1801-1900",
	"/List_of_battles_1901-2000",
	"/List_of_battles_before_301",
	"/List_of_battles_since_2001",
	"/List_of_Hundred_Years%27_War_battles",
	"/List_of_military_engagements_of_World_War_I",
	"/List_of_military_engagements_of_World_War_II",
	"/List_of_Napoleonic_battles",
}

// List scrapes and retrieves the full list of Wikipedia's battles when grouped by different
// criteria, in the form of domain.SBattleItem
func (s *Scraper) List() []domain.SBattleItem {
	urlsByName := make(map[string]string)

	for _, urlPart := range battlesLists {
		listURL := "https://en.wikipedia.org/wiki" + urlPart
		subscribe := func(c *colly.Collector) {
			c.OnHTML("#content a[href]", func(e *colly.HTMLElement) {
				name := parser.Clean(e.Text)
				if _, cached := urlsByName[name]; cached {
					return
				}
				if name == "" || e.Attr("class") == "new" {
					return
				}
				href := e.Attr("href")
				if urls.ShouldSkip(href) || urls.NotSpecific(href) || !names.IsBattle(name) {
					return
				}
				urlsByName[name] = "https://en.wikipedia.org" + href
			})
		}
		if err := s.do(listURL, subscribe); err != nil {
			log.Printf("Error scraping list in %s: %s", listURL, err)
		}
	}

	items := []domain.SBattleItem{}
	for name, url := range urlsByName {
		items = append(items, domain.SBattleItem{URL: url, Name: name})
	}

	feedbackMessage := fmt.Sprintf("There are %d items that can be scraped\n", len(items))
	s.logger.Write([]byte(feedbackMessage))
	return items
}
