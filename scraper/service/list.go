package service

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/sasalatart/batcoms/parser"
	"github.com/sasalatart/batcoms/scraper/domain"
	"github.com/sasalatart/batcoms/scraper/urls"

	"github.com/gocolly/colly"
)

type listStrategy struct {
	urlSuffix string
	selector  string
}

var keywords = strings.Join([]string{
	"action", "ambush", "assault", "attack",
	"battle", "blockade", "bloodbath", "bombardment", "bombing", "burning",
	"campaign", "capture", "clash", "clashes", "combat", "conflict", "confrontation", "conquest", "crisis", "crossing", "crusade",
	"defense",
	"engagement", "expedition",
	"fall", "fight",
	"incident", "insurgency", "intervention", "invasion",
	"landing", "liberation",
	"massacre", "march", "mutiny",
	"occupation", "offensive", "operatie", "operation",
	"prison break",
	"raid", "rebellion", "recovery", "relief", "revolt", "rising",
	"sack", "siege", "sinking", "skirmish", "stand", "standoff", "strike",
	"takeover",
	"uprising",
	"war",
}, "|")

func buildRegex(format string) *regexp.Regexp {
	return regexp.MustCompile("(?i)" + fmt.Sprintf(format, keywords))
}

var inside = buildRegex("(%s)s?\\s(at|for|in|of|on|to)")
var suffix = buildRegex("\\s(%s)s?$")
var prefix = buildRegex("^(%s)\\s")

// List scrapes and retrieves the full list of Wikipedia's battles (name & URL only) when grouped by
// centuries
func (s *Scraper) List() []domain.BattleItem {
	urlsByName := make(map[string]string)

	for _, listURL := range urls.BattlesLists() {
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
				if urls.ShouldSkip(href) || urls.NotSpecific(href) {
					return
				}

				bT := []byte(name)
				if !inside.Match(bT) && !suffix.Match(bT) && !prefix.Match(bT) {
					return
				}

				urlsByName[name] = "https://en.wikipedia.org" + href
			})
		}

		if err := s.do(listURL, subscribe); err != nil {
			log.Printf("Error scraping list in %s: %s", listURL, err)
		}
	}

	battles := []domain.BattleItem{}
	for name, url := range urlsByName {
		battles = append(battles, domain.BattleItem{URL: url, Name: name})
	}

	feedbackMessage := fmt.Sprintf("There are %d battles that can be scraped\n", len(battles))
	s.logger.Write([]byte(feedbackMessage))
	return battles
}
