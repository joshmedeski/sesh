package main

import (
	"log"
	"log/slog"
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

func init() {
	env := os.Getenv("ENV")

	handlerOptions := &slog.HandlerOptions{}

	switch env {
	case "debug":
		handlerOptions.Level = slog.LevelDebug
	case "info":
		handlerOptions.Level = slog.LevelInfo
	case "warn":
		handlerOptions.Level = slog.LevelWarn
	default:
		handlerOptions.Level = slog.LevelError
	}

	loggerHandler := slog.NewTextHandler(os.Stdout, handlerOptions)
	logger := slog.New(loggerHandler)

	slog.SetDefault(logger)
}
