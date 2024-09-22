package main

import (
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/joshmedeski/sesh/seshcli"
)

var version = "dev"

func main() {
  slog.Debug("Debug")
  slog.Info("Information")
  slog.Warn("Warning")
  slog.Error("Error")

	app := seshcli.App(version)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
func init() {
	env := os.Getenv("ENV")

  handlerOptions := &slog.HandlerOptions{}

  switch strings.ToLower(env) {
  case "debug":
    handlerOptions.Level = slog.LevelDebug
  case "info":
    handlerOptions.Level = slog.LevelInfo
  case "error":
    handlerOptions.Level = slog.LevelError
  default:
    handlerOptions.Level = slog.LevelWarn
  }

  loggerHandler := slog.NewJSONHandler(os.Stdout, handlerOptions)
  slog.SetDefault(slog.New(loggerHandler))
}
