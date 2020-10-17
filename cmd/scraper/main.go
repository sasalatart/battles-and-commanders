package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/config"
	"github.com/sasalatart/batcoms/db/memory"
	"github.com/sasalatart/batcoms/domain/wikibattles"
	"github.com/sasalatart/batcoms/pkg/io/json"
	"github.com/sasalatart/batcoms/pkg/logger"
	"github.com/sasalatart/batcoms/pkg/scraper/battles"
	"github.com/sasalatart/batcoms/pkg/scraper/list"
	"github.com/spf13/viper"
)

func init() {
	config.Setup()
}

func main() {
	loggerService := logger.New(ioutil.Discard, os.Stderr)
	battlesScraper := battles.NewScraper(
		memory.NewWikiActorsRepo(),
		memory.NewWikiBattlesRepo(),
		json.Export,
		loggerService,
	)

	var failedCount int
	semaphore := make(chan bool, 10)
	list := list.Scrape(loggerService)
	for i, battle := range list {
		semaphore <- true
		fmt.Printf("\r%d/%d (failed: %d)", i, len(list), failedCount)
		go func(i int, b wikibattles.BattleItem) {
			if _, err := battlesScraper.ScrapeOne(b.URL); err != nil {
				failedCount++
				loggerService.Error(errors.Wrapf(err, "Error scraping %s", b.URL))
			}
			<-semaphore
		}(i, battle)
	}

	actorsFileName := viper.GetString("SCRAPER_RESULTS.ACTORS")
	battlesFileName := viper.GetString("SCRAPER_RESULTS.BATTLES")
	if err := battlesScraper.ExportAll(actorsFileName, battlesFileName); err != nil {
		loggerService.Error(err)
	}
}
