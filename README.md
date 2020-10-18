# Battles and Commanders &middot; [![m][bdg-mit]][mit] [![b][bdg-build]][build]

## About

_Battles and Commanders_ is a Wikipedia scraper and API that serves historical battles, commanders
and their factions. You can try the API and read its documentation [here][swaggerhub].

**This project is still work in progress**. Also, consider that scrapers are brittle: they are
subject to the webpage's HTML structure updates, and if Wikipedia decides to alter its HTML, this
scraper may stop working if not properly updated.

## Setup

### Development

For both the scraper and the API, settings may be changed by editing the file in `config/config.yaml`,
but I see no reason to change their current values.

#### Scraper

```sh
$ go run cmd/scraper/main.go
```

After a successful run, the file `data.json` should have been created at the root dir of this
project, containing normalized battles, factions and commanders.

#### API

1. Make sure you have a Postgres instance running.

2. Turn on the API on port 3000:

   ```sh
   $ go run cmd/api/main.go
   ```

3. If you have not yet, then you will need to also run the seeder. The seeder can setup the database
   by either reading from local files or from remote files whose URLs have been supplied:

   ```sh
   $ go run cmd/seeder/main.go -dataURL="https://bit.ly/3dOeqyZ"
   ```

   If you already ran the scraper and the `data.json` file is available in the root dir of the
   project, then the `dataURL` option is unnecessary.

### Running via Docker

```sh
# Start database and app server via docker-compose
$ docker-compose up -d

# Additionally, seeders should be ran
$ docker exec api ./seeder -dataURL="https://bit.ly/3dOeqyZ"

# Stop docker containers:
$ docker-compose stop
```

Now the api should be available on port 3000 of your machine.

## Installing for use with your own Go projects

Some of the functionality used by both the scraper and the API is publicly available for use outside
this project, inside the `pkg` dir. To install, simply run:

```sh
$ go get github.com/sasalatart/batcoms
```

Some usage examples include:

1. Scraping a list of **potential** battles (names and urls only, false positives may be included):

   ```go
   package main

   import (
      "github.com/sasalatart/batcoms/pkg/logger"
      "github.com/sasalatart/batcoms/pkg/scraper/list"
   )

   func main() {
      loggerService := logger.NewDiscard() // Or anything that implements logger.Interface
      potentialBattles := list.Scrape(loggerService)
   }
   ```

2. Scraping specific battles:

   ```go
   package main

   import (
      "github.com/sasalatart/batcoms/pkg/logger"
      "github.com/sasalatart/batcoms/pkg/scraper/battles"
   )

   func main() {
      loggerService := logger.NewDiscard() // Or anything that implements logger.Interface
      scraperService := battles.NewScraper(loggerService)

      austerlitz, err := scraperService.ScrapeOne("https://en.wikipedia.org/wiki/Battle_of_Austerlitz")
      // Handle error and optionally do something with normalized Battle of Austerlitz...
      actium, err := scraperService.ScrapeOne("https://en.wikipedia.org/wiki/Battle_of_Actium")
      // Handle error and optionally do something with normalized Battle of Actium...

      // Each battle contains normalized data (ids of factions and commanders), so we export everything
      data := scraperService.Data()
      // Do something with data.BattlesByID, data.FactionsByID and/or data.CommandersByID...
   }
   ```

3. Cleaning scraped text:

   ```go
   package main

   import "github.com/sasalatart/batcoms/pkg/strclean"

   func main() {
      input := "Soviet victory:[1]\n\nDestruction of the German 6th Army"
      output := strclean.Apply(input) // "Soviet victory: Destruction of the German 6th Army"
   }
   ```

4. Parsing Wikipedia's **Info Box** date text (accuracy improvements are still WIP...):

   ```go
   package main

   import "github.com/sasalatart/batcoms/pkg/dates"

   func main() {
      d1 := "January-March 309 B.C."
      parsed1, err := dates.Parse(d1)
      // []dates.Historic{
      //   dates.Historic{ Year: 309, Month: 1, Day: 0, IsBCE: true },
      //   dates.Historic{ Year: 309, Month: 3, Day: 0, IsBCE: true },
      // }
      // Handle error and do something with parsed1...

      d2 := "July 6, 1950; 69 years ago (1950-07-06)"
      parsed2, err := dates.Parse(d2)
      // []dates.Historic{dates.Historic{ Year: 1950, Month: 7, Day: 6, IsBCE: false }}
      // Handle error and do something with parsed2...
   }
   ```

## Testing

Unit & integration tests have been included. In order to run the integration tests properly, the API
must be running in `test` mode, and Postgres should also be running. Tests can be ran as follows:

```sh
# Shell 1: Turn on the API in test mode
$ go run cmd/api/main.go -test

# Shell 2: Run the actual tests
$ go test ./...
```

## Credits

Special thanks to [Wikipedia][wikipedia] and the content-creators that have provided the historical
data served and scraped by this app, and which are [licensed][cc-share-alike] as Creative Commons
Attribution-ShareAlike 3.0 Unported License.

## License

Copyright (c) 2020, Sebastián Salata Ruiz-Tagle

_Battles and Commanders_ is [MIT licensed](./LICENSE).

[bdg-build]: https://circleci.com/gh/sasalatart/battles-and-commanders.svg?style=svg
[bdg-mit]: https://img.shields.io/badge/License-MIT-blue.svg
[build]: https://circleci.com/gh/battles-and-commanders
[cc-share-alike]: https://creativecommons.org/licenses/by-sa/3.0/
[mit]: https://opensource.org/licenses/MIT
[swaggerhub]: https://app.swaggerhub.com/apis-docs/sasalatart/Battles-and-Commanders/1.0.0
[wikipedia]: https://www.wikipedia.org/
