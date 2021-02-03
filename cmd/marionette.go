package main

import (
	"log"

	"github.com/dollarkillerx/marionette_go/internal/config"
	"github.com/dollarkillerx/marionette_go/internal/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	config.InitConfig()

	server.RunServers()
}
