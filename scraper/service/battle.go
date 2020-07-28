package service

import (
	"github.com/gocolly/colly"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/scraper/domain"
)

// Battle scrapes information about the battle found in the URL passed to it
func (s *Scraper) Battle(url string) (domain.Battle, error) {
	b := domain.Battle{URL: url, CommandersByFaction: make(map[int][]int)}
	var err error

	summary, err := PageSummary(url)
	if err != nil {
		return b, errors.Wrapf(err, "Fetching summary for %s", url)
	}
	b.ID = summary.PageID
	b.Description = summary.Description
	b.Extract = summary.Extract

	err = s.do(url, func(c *colly.Collector) {
		s.subscribeMeta(c, &b)
		s.subscribeParticipants(c, &b)
		s.subscribeStrength(c, &b)
		s.subscribeCasualties(c, &b)

		c.OnScraped(func(_ *colly.Response) {
			err = s.BattlesStore.Save(b)
		})
	})

	return b, err
}
