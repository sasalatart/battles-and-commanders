package list

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain/wikibattles"
	"github.com/sasalatart/batcoms/pkg/logger"
	"github.com/sasalatart/batcoms/pkg/scraper/names"
	"github.com/sasalatart/batcoms/pkg/scraper/urls"
	"github.com/sasalatart/batcoms/pkg/strclean"

	"github.com/gocolly/colly"
)

type hrefsCache map[string]struct{}

// Scrape scrapes and retrieves the full list of Wikipedia's battles when grouped by different
// criteria, in the form of wikibattles.BattleItem
func Scrape(l logger.Service) []wikibattles.BattleItem {
	hrefs := make(hrefsCache)
	var items []wikibattles.BattleItem
	for _, urlPart := range battlesLists {
		listURL := "https://en.wikipedia.org/wiki" + urlPart
		if err := do(listURL, &items, hrefs, l); err != nil {
			l.Error(errors.Wrapf(err, "Scraping list in %s", listURL))
		}
	}
	l.Info(fmt.Sprintf("There are %d items that can be scraped", len(items)))
	return items
}

func do(url string, battlesItems *[]wikibattles.BattleItem, hrefs hrefsCache, l logger.Service) error {
	c := colly.NewCollector()

	c.OnHTML(listItemsSelector, func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if _, cached := hrefs[href]; cached || e.Attr("class") == "new" || urls.ShouldSkip(href) || urls.NotSpecific(href) {
			return
		}
		name := strclean.Apply(e.Text)
		if !names.IsBattle(name) {
			return
		}
		hrefs[href] = struct{}{}
		*battlesItems = append(*battlesItems, wikibattles.BattleItem{
			URL:  "https://en.wikipedia.org" + href,
			Name: name,
		})
	})

	c.OnRequest(func(r *colly.Request) {
		l.Info(fmt.Sprintf("Scraping %s", r.URL))
	})

	return c.Visit(url)
}

const listItemsSelector = "#content a[href]"

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
