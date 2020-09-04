package main

import (
	"flag"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/sasalatart/batcoms/config"
	"github.com/sasalatart/batcoms/http"
	"github.com/sasalatart/batcoms/store/postgresql"
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
		db = postgresql.Connect(postgresql.TestConfig())
		port = viper.GetInt("PORT_TEST")
	} else {
		db = postgresql.Connect(nil)
		port = viper.GetInt("PORT")
	}
	defer db.Close()

	fs := postgresql.NewFactionsDataStore(db)
	cs := postgresql.NewCommandersDataStore(db)
	bs := postgresql.NewBattlesDataStore(db)
	server := http.Setup(fs, cs, bs, false)
	log.Fatal(server.Listen(port))
}
