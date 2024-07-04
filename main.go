package main

import (
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/joshmedeski/sesh/seshcli"
)

var version = "dev"


func main() {
	app := seshcli.App(version)
    slog.Error("Testing")
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func init() {
	env := os.Getenv("ENV")

    configDir, err := os.UserConfigDir()
    if err != nil {
        log.Fatal(err)
    }

    f, err :=  os.OpenFile(configDir+"/sesh/sesh.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }

    w := io.MultiWriter(os.Stdout, f)

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

	loggerHandler := slog.NewTextHandler(w, handlerOptions)
	logger := slog.New(loggerHandler)

	slog.SetDefault(logger)
}
