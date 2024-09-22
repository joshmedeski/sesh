package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path"
	"strings"
	"time"

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
	var f *os.File
	var err error
	fileOnly := false

	if f, err = createLoggerFile(); err != nil {
		log.Fatalf("Unable to create logger file: %v", err)
	}
	if f == nil {
		log.Println("init: no file")
	}

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
		fileOnly = true
	}

	var loggerHandler *slog.JSONHandler
	if !fileOnly {
		multiWriter := io.MultiWriter(os.Stdout, f)
		loggerHandler = slog.NewJSONHandler(multiWriter, handlerOptions)
	} else {
		loggerHandler = slog.NewJSONHandler(f, handlerOptions)
	}
	slog.SetDefault(slog.New(loggerHandler))
}

func createLoggerFile() (*os.File, error) {
	now := time.Now()
	date := fmt.Sprintf("%s.log", now.Format("2006-01-02"))
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(path.Join(userHomeDir, ".config", "sesh"), 0755); err != nil {
		return nil, err
	}

	fileFullPath := path.Join(userHomeDir, ".config", "sesh", date)
	file, err := os.OpenFile(fileFullPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	if file == nil {
		log.Println("createLoggerFle: nil file")
	}

	return file, nil
}
