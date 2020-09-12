package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/services/io/json"
	"github.com/sasalatart/batcoms/services/logger"
	"github.com/sasalatart/batcoms/services/scraper/battles"
	"github.com/sasalatart/batcoms/services/scraper/list"
	"github.com/sasalatart/batcoms/store/memory"
)

func main() {
	loggerService := logger.New(ioutil.Discard, os.Stderr)
	battlesScraper := battles.NewScraper(
		memory.NewSBattlesStore(),
		memory.NewSParticipantsStore(),
		json.Export,
		loggerService,
	)

	var failedCount int
	semaphore := make(chan bool, 10)
	list := list.Scrape(loggerService)
	for i, battle := range list {
		semaphore <- true
		fmt.Printf("\r%d/%d (failed: %d)", i, len(list), failedCount)
		go func(i int, b domain.SBattleItem) {
			if _, err := battlesScraper.ScrapeOne(b.URL); err != nil {
				failedCount++
				loggerService.Error(errors.Wrapf(err, "Scraping %s", b.URL))
			}
			<-semaphore
		}(i, battle)
	}

	if err := battlesScraper.ExportAll("battles.json", "participants.json"); err != nil {
		loggerService.Error(err)
	}
}
