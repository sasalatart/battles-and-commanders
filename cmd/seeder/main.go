package main

import (
	"github.com/sasalatart/batcoms/config"
	"github.com/sasalatart/batcoms/repository/postgresql"
)

func main() {
	config.Setup()
	postgresql.Connect()
}
