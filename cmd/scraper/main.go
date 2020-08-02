package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/services/scraper"
	"github.com/sasalatart/batcoms/store/memory"
)

func main() {
	scraperService := scraper.New(
		memory.NewSBattlesStore(),
		memory.NewSParticipantsStore(),
		ioutil.Discard,
	)

	semaphore := make(chan bool, 10)
	list := scraperService.List()
	for i, battle := range list {
		semaphore <- true
		fmt.Printf("\r%d/%d", i, len(list))
		go func(i int, b domain.SBattleItem) {
			if _, err := scraperService.SBattle(b.URL); err != nil {
				log.Printf("Failed scraping %s: %s", b.URL, err)
			}
			<-semaphore
		}(i, battle)
	}

	if err := scraperService.Export("battles.json", "participants.json"); err != nil {
		log.Printf("Error: %s", err)
	}
}
