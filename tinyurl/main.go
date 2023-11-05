package main

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	LocalUrlPrefix  string = "http://localhost:3000/"
	DockerUrlPrefix string = "http://172.17.0.1:3000/"
	RemoteUrlPrefix string = "https://mdtiny.net/"
	AppPort         int    = 3000
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	isDebug := flag.Bool("debug", false, "sets log level to debug")
	isLocal := flag.Bool("local", false, "sets whether is local deployment or not")
	isDocker := flag.Bool("docker", false, "sets whether using Docker deployment")

	flag.Parse()

	// logging
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *isDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Info().Msgf("Docker deployment: %t", *isDocker)
	log.Info().Msgf("Local deployment: %t", *isLocal)

	// DB client and context
	client := GetLocalDbClient()

	// set the URL prefix for the application
	// depending on the deployment environment
	var urlPrefix string
	if *isLocal {
		urlPrefix = LocalUrlPrefix
	} else if *isDocker {
		urlPrefix = DockerUrlPrefix
	} else {
		urlPrefix = RemoteUrlPrefix
	}

	log.Info().Msgf("Application URL: %s", urlPrefix)

	log.Debug().Msg("Starting server")
	app := CreateServer(urlPrefix, client)
	app.Listen(fmt.Sprintf("0.0.0.0:%d", AppPort))
}
