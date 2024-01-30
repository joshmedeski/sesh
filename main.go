package main

import (
	"log"
	"os"

	db "github.com/joshmedeski/sesh/database"
	"github.com/joshmedeski/sesh/seshcli"
)

var version = "dev"

func main() {
	sqlPath := os.ExpandEnv("$HOME/.local/share/sesh/sesh.db")
	storage := db.NewSqliteDatabase(sqlPath)

	app := seshcli.App(version, storage)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
