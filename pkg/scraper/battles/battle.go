package battles

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain/summaries"
	"github.com/sasalatart/batcoms/domain/wikibattles"
)

// ScrapeOne scrapes information about the battle found in the URL passed to it
func (s *Scraper) ScrapeOne(url string) (wikibattles.Battle, error) {
	battle := wikibattles.Battle{URL: url, CommandersByFaction: make(wikibattles.CommandersByFaction)}
	if err := s.assignSummary(&battle); err != nil {
		return battle, errors.Wrap(err, "Assigning summary")
	}

	ctx := &battleCtx{&battle, colly.NewCollector(), nil}
	ctx.collector.OnRequest(func(r *colly.Request) {
		s.logger.Info(fmt.Sprintf("Scraping %s", r.URL))
	})
	s.assertHasOneInfoBox(ctx)
	s.subscribeMeta(ctx)
	s.subscribeActors(ctx)
	s.subscribeStrength(ctx)
	s.subscribeCasualties(ctx)

	if err := ctx.collector.Visit(url); err != nil {
		return battle, errors.Wrap(err, "Doing the request to scrape")
	}
	if ctx.err != nil {
		return battle, ctx.err
	}
	if err := s.wikiBattlesRepo.Save(battle); err != nil {
		return battle, errors.Wrapf(err, "Saving battle")
	}
	return battle, nil
}

func (s *Scraper) assignSummary(b *wikibattles.Battle) error {
	summary, err := summaries.Fetch(b.URL)
	if err != nil {
		return errors.Wrap(err, "Fetching summary")
	}
	if summary.Extract == "" {
		return ErrNoSummaryExtract
	}
	b.ID = summary.PageID
	b.Name = summary.Title
	b.Description = summary.Description
	b.Extract = summary.Extract
	return nil
}
