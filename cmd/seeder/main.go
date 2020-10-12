package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/config"
	"github.com/sasalatart/batcoms/db/postgresql"
	"github.com/sasalatart/batcoms/db/seeder"
	"github.com/sasalatart/batcoms/pkg/logger"
	"github.com/spf13/viper"
)

var actorsURL = flag.String("actorsURL", "", "The URL from which actors file can be downloaded")
var battlesURL = flag.String("battlesURL", "", "The URL from which battles file can be downloaded")

func init() {
	config.Setup()
	flag.Parse()
}

func main() {
	db, sqlDB := postgresql.Connect(nil)
	defer sqlDB.Close()

	loggerService := logger.New(log.Writer(), os.Stderr)

	actorsFileName := fileNameFor(*actorsURL, viper.GetString("SCRAPER_RESULTS.ACTORS"), loggerService)
	if *actorsURL != "" {
		defer os.Remove(actorsFileName)
	}
	battlesFileName := fileNameFor(*battlesURL, viper.GetString("SCRAPER_RESULTS.BATTLES"), loggerService)
	if *battlesURL != "" {
		defer os.Remove(battlesFileName)
	}

	importedData := new(seeder.ImportedData)
	if err := seeder.JSONImport(importedData, actorsFileName, battlesFileName); err != nil {
		log.Fatalf("Failed seeding: %s", err)
	}

	postgresql.Reset(db)
	seeder.Seed(
		importedData,
		postgresql.NewFactionsRepository(db),
		postgresql.NewCommandersRepository(db),
		postgresql.NewBattlesRepository(db),
		loggerService,
	)
}

func fileNameFor(url, defaultName string, loggerService logger.Service) string {
	handleError := func(err error, from string) {
		if err != nil {
			log.Fatalf("No data may be read from %s: %s", from, err)
		}
	}

	if url == "" {
		loggerService.Info(fmt.Sprintf("No URL supplied, falling back to %s...\n", defaultName))
		_, err := os.Stat(defaultName)
		handleError(err, defaultName)
		return defaultName
	}

	tmpFile, err := download(url, loggerService)
	handleError(err, url)
	return tmpFile.Name()
}

func download(url string, logger logger.Service) (*os.File, error) {
	logger.Info(fmt.Sprintf("Downloading from %s...\n", url))
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "Downloading file")
	}
	defer resp.Body.Close()
	out, err := ioutil.TempFile("", "downloaded.*")
	if err != nil {
		return nil, errors.Wrapf(err, "Creating tmp file for file from %s", url)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return out, err
}
