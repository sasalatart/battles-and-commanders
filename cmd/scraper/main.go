package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/config"
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
	scraperService := battles.NewScraper(loggerService)

	var failedCount int
	semaphore := make(chan bool, 10)
	list := list.Scrape(loggerService)
	for i, battle := range list {
		semaphore <- true
		fmt.Printf("\r%d/%d (failed: %d)", i, len(list), failedCount)
		go func(i int, b wikibattles.BattleItem) {
			if _, err := scraperService.ScrapeOne(b.URL); err != nil {
				failedCount++
				loggerService.Error(errors.Wrapf(err, "Error scraping %s", b.URL))
			}
			<-semaphore
		}(i, battle)
	}

	data := scraperService.Data()
	fileName := viper.GetString("SCRAPER_DATA")
	if err := json.Export(fileName, data); err != nil {
		log.Fatalln(err)
	}
}
