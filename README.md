# Battles and Commanders &middot; [![m][bdg-mit]][mit] [![b][bdg-build]][build]

## About

_Battles and Commanders_ is a Wikipedia scraper and API that serves historical battles, commanders
and their factions. You can try the API and read its documentation [here][swaggerhub]. This project
was built with [Go][go] and [Postgres][postgres].

> **The work here is still in progress**. Also, consider that scrapers are brittle: they are subject
> to the webpage's HTML structure updates, and if Wikipedia decides to alter its HTML, this scraper
> may stop working if not properly updated.

## Development setup

This project requires [Docker][docker] and [docker-compose][docker-compose] to be installed.

In development, Docker is used together with `Make` to let you spin up your local environment
without worrying about complex commands and installing dependencies. You may still run the code
natively with Go without Docker, although you will lose some features such as auto-reload. For a
list and description of all the available `Make` commands, just run `make help`.

Most settings may be changed by editing the file in `config/config.yaml`, although you will probably
not need to change them. You might, however, want to override some, such as the database password.

### Scraper

```sh
# Run the scraper inside a Docker container
$ make scrape
```

The resulting `data.json` file at the root dir of this project will contain normalized battles,
factions and commanders. You may use this file for seeding the API (see next section), or for some
other project.

### API

```sh
# Turn on the API (http://localhost:3000) and a Postgres container (port 14000). The API has
# auto-reload configured
$ make dev_up

# Run the seeder (just needed once). Alternatively, you may run "make dev_seed_local" if you have
# the scraper results file in the root dir of this project
$ make dev_seed_url

# (Optional) remove Docker containers and volumes created by "make dev_up"
$ make dev_destroy
```

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

```sh
# Shell 1: Turn on the API in test mode (http://localhost:8888) and a Postgres container (port 14001)
$ make test_up

# Shell 2: Run the actual tests inside the container
$ make test
```

Just like when running the API in dev mode, you may run the `make test_destroy` command to remove
Docker containers and volumes created for running tests.

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
[docker-compose]: https://docs.docker.com/compose/
[docker]: https://www.docker.com/
[go]: https://golang.org/
[mit]: https://opensource.org/licenses/MIT
[postgres]: https://www.postgresql.org/
[swaggerhub]: https://app.swaggerhub.com/apis-docs/sasalatart/Battles-and-Commanders/1.0.0
[wikipedia]: https://www.wikipedia.org/
