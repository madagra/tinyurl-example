package main

import (
	"fmt"
	"flag"
	"github.com/rs/zerolog"
)

var AppPort int = 3000

var LocalUrlPrefix = "http://localhost:3000/"

var ProdUrlPrefix = "https://mdtiny.net/"

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	debug := flag.Bool("debug", false, "sets log level to debug")

    zerolog.SetGlobalLevel(zerolog.InfoLevel)
    if *debug {
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
    }	
	
	app := CreateServer(LocalUrlPrefix)
	app.Listen(fmt.Sprintf(":%d", AppPort))
}
