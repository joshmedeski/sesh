package main

import (
	"joshmedeski/sesh/seshcli"
	"log"
	"os"
)

func main() {
	app := seshcli.App()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
