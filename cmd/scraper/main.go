package main

import (
	"log"

	"github.com/sasalatart/batcoms/scraper"
)

func main() {
	scraper.ScrapeList(log.Writer())
}
