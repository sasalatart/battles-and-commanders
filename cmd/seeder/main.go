package main

import (
	"github.com/sasalatart/batcoms/config"
	"github.com/sasalatart/batcoms/store/postgresql"
)

func main() {
	config.Setup()
	postgresql.Connect()
}
