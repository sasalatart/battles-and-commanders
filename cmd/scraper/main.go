package main

import (
	"fmt"
	"log"

	"github.com/sasalatart/batcoms/scraper/service"
	"github.com/sasalatart/batcoms/scraper/store"
)

func main() {
	scraperService := service.NewScraper(
		store.NewBattlesMem(),
		store.NewParticipantsMem(),
		log.Writer(),
	)

	scraperService.Battle("https://en.wikipedia.org/wiki/Battle_of_Megiddo_(15th_century_BC)")
	scraperService.Battle("https://en.wikipedia.org/wiki/Battle_of_Austerlitz")
	scraperService.Battle("https://en.wikipedia.org/wiki/Battle_of_Stalingrad")

	if err := scraperService.Export("battles.json", "participants.json"); err != nil {
		fmt.Printf("Error: %s", err)
	}
}
