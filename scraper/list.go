package scraper

import (
	"fmt"
	"io"

	"github.com/gocolly/colly"
	batcoms "github.com/sasalatart/batcoms"
)

type scrapeParams struct {
	urlSuffix string
	selector  string
}

func scrapeList(sp scrapeParams, battles *[]batcoms.ScrapedBattle, logger io.Writer) {
	c := colly.NewCollector()

	c.OnHTML(fmt.Sprintf("#content %s", sp.selector), func(e *colly.HTMLElement) {
		if e.Attr("class") == "new" {
			return
		}

		*battles = append(*battles, batcoms.ScrapedBattle{Name: e.Text, URL: e.Attr("href")})
	})

	c.OnRequest(func(r *colly.Request) {
		message := fmt.Sprintf("Scraping %s\n", r.URL)
		logger.Write([]byte(message))
	})

	URL := fmt.Sprintf("https://en.wikipedia.org/wiki/List_of_battles_%s", sp.urlSuffix)
	c.Visit(URL)
}

// ScrapeList returns a slice of batcoms.ScrapedBattle produced after running the scraper through
// each one of Wikipedia's indexed lists of battles
func ScrapeList(logger io.Writer) []batcoms.ScrapedBattle {
	battles := []batcoms.ScrapedBattle{}

	var pp = [...]scrapeParams{
		{"before_301", "td:nth-last-child(2) > a"},
		{"301-1300", "td:nth-last-child(3) > a"},
		{"1301-1600", "td:nth-last-child(3) > a"},
		{"1601-1800", "h2+ul a:first-child"},
		{"1801-1900", "h2+ul a:first-child"},
		{"1901-2000", "h2+ul a:first-child"},
		{"since_2001", "td:first-child > a"},
	}

	for _, p := range pp {
		scrapeList(p, &battles, logger)
	}

	feedbackMessage := fmt.Sprintf("There are %d battles that can be scraped\n", len(battles))
	logger.Write([]byte(feedbackMessage))

	return battles
}
