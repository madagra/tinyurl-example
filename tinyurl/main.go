package main

import (
	"fmt"
)

var AppPort int = 3000

var LocalUrlPrefix = "http://localhost:3000/"

var ProdUrlPrefix = "https://mdtiny.net/"

func main() {
	app := CreateServer(LocalUrlPrefix)
	app.Listen(fmt.Sprintf(":%d", AppPort))
}
