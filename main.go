package main

import (
	"log"
	"os"

	"github.com/joshmedeski/sesh/seshcli"
)

var version = "dev"

func main() {
	app := seshcli.App(version)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
