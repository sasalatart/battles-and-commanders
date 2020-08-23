package main

import (
	"log"

	"github.com/sasalatart/batcoms/config"
	"github.com/sasalatart/batcoms/http"
	"github.com/sasalatart/batcoms/store/postgresql"
	"github.com/spf13/viper"
)

func init() {
	config.Setup()
}

func main() {
	db := postgresql.Connect()
	defer db.Close()

	fs := postgresql.NewFactionsDataStore(db)
	cs := postgresql.NewCommandersDataStore(db)
	server := http.Setup(fs, cs, false)
	log.Fatal(server.Listen(viper.GetInt("PORT")))
}
