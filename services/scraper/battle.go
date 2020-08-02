package scraper

import (
	"github.com/gocolly/colly"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/services/summaries"
)

// SBattle scrapes information about the battle found in the URL passed to it
func (s *Scraper) SBattle(url string) (domain.SBattle, error) {
	b := domain.SBattle{URL: url, CommandersByFaction: make(map[int][]int)}
	var err error

	summary, err := summaries.Get(url)
	if err != nil {
		return b, errors.Wrapf(err, "Fetching summary for %s", url)
	}
	b.ID = summary.PageID
	b.Name = summary.Title
	b.Description = summary.Description
	b.Extract = summary.Extract

	err = s.do(url, func(c *colly.Collector) {
		s.subscribeMeta(c, &b)
		s.subscribeParticipants(c, &b)
		s.subscribeStrength(c, &b)
		s.subscribeCasualties(c, &b)

		c.OnScraped(func(_ *colly.Response) {
			err = s.SBattlesStore.Save(b)
		})
	})

	return b, err
}
