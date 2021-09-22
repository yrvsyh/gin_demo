package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yrvsyh/gin_demo/router"
)

func main() {
	log.Logger = log.Level(zerolog.InfoLevel)
	panic(router.InitRouter().Run("0.0.0.0:8080"))
}
