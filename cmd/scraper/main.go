package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/sasalatart/batcoms/scraper/domain"
	"github.com/sasalatart/batcoms/scraper/service"
	"github.com/sasalatart/batcoms/scraper/store"
)

func main() {
	scraperService := service.NewScraper(
		store.NewBattlesMem(),
		store.NewParticipantsMem(),
		ioutil.Discard,
	)

	semaphore := make(chan bool, 10)
	list := scraperService.List()
	for i, battle := range list {
		semaphore <- true
		go func(i int, b domain.BattleItem) {
			_, err := scraperService.Battle(b.URL)
			if err != nil {
				log.Printf("Failed scraping %s: %s", b.URL, err)
			}
			fmt.Printf("\r%d/%d", i, len(list))
			<-semaphore
		}(i, battle)
	}

	if err := scraperService.Export("battles.json", "participants.json"); err != nil {
		log.Printf("Error: %s", err)
	}
}
