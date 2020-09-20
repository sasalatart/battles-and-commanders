package main

import (
	"flag"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/sasalatart/batcoms/config"
	"github.com/sasalatart/batcoms/db/postgresql"
	"github.com/sasalatart/batcoms/http"
	"github.com/spf13/viper"
)

var testModeFlag = flag.Bool("test", false, "Whether the API should run in test mode or not")

func init() {
	config.Setup()
	flag.Parse()
}

func main() {
	var db *gorm.DB
	var port int
	if *testModeFlag {
		db = postgresql.Connect(postgresql.DefaultTestConfig())
		port = viper.GetInt("PORT_TEST")
	} else {
		db = postgresql.Connect(nil)
		port = viper.GetInt("PORT")
	}
	defer db.Close()

	server := http.Setup(
		postgresql.NewFactionsRepository(db),
		postgresql.NewCommandersRepository(db),
		postgresql.NewBattlesRepository(db),
		false,
	)
	log.Fatal(server.Listen(port))
}
