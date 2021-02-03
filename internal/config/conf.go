package config

import (
	"log"

	"github.com/dollarkillerx/env"
)

var CONF *config

type config struct {
	ListenAddr string
}

func InitConfig() {
	var conf config
	if err := env.FillBase(&conf); err != nil {
		log.Fatalln(err)
	}

	CONF = &conf
}
