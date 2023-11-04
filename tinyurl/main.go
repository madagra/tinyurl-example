package main

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var AppPort int = 3000

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	isDebug := flag.Bool("debug", false, "sets log level to debug")
	isLocal := flag.Bool("local", true, "sets whether is local deployment or not")

	// logging
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *isDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// DB client and context
	client := GetLocalDbClient()

	var urlPrefix string
	if *isLocal {
		urlPrefix = LocalUrlPrefix
	} else {
		urlPrefix = RemoteUrlPrefix
	}

	log.Debug().Msg("Starting server")
	app := CreateServer(urlPrefix, client)
	app.Listen(fmt.Sprintf("0.0.0.0:%d", AppPort))
}
