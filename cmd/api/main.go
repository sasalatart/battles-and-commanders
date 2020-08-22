package main

import (
	"log"

	"github.com/sasalatart/batcoms/config"
	"github.com/sasalatart/batcoms/http"
	"github.com/spf13/viper"
)

func init() {
	config.Setup()
}

func main() {
	server := http.Setup()
	log.Fatal(server.Listen(viper.GetInt("PORT")))
}
