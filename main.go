package main

import (
	"context"
	"fmt"
	"io" // This is used to enable the multiwritter and be able to write to the log file and the console at the same time
	"log/slog"
	"os"
	"path"    // This is used to create the path where the log files will be stored
	"strings" // This is required to conpare the evironment variables
	"time"    // This is used to get the current date and create the log file

	"github.com/charmbracelet/fang"
	"github.com/joshmedeski/sesh/v2/seshcli"
)

var version = "dev"

func main() {
	slog.Debug("Debug")
	slog.Info("Information")
	slog.Warn("Warning")
	slog.Error("Error")

	// Extract --config/-C flag from os.Args before cobra builds the command tree,
	// since the DI graph is constructed eagerly in NewRootCommand.
	// The flag is stripped from os.Args so cobra/fang don't see it.
	configPath, remaining := extractConfigFlag(os.Args[1:])
	os.Args = append(os.Args[:1], remaining...)

	cmd := seshcli.NewRootCommand(version, configPath)
	if err := fang.Execute(context.TODO(), cmd, fang.WithColorSchemeFunc(fang.AnsiColorScheme), fang.WithoutVersion()); err != nil {
		slog.Error("main file: ", "error", err)
		os.Exit(1)
	}
}

func extractConfigFlag(args []string) (configPath string, remaining []string) {
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--config" || args[i] == "-C":
			if i+1 < len(args) {
				configPath = args[i+1]
				i++ // skip value
			}
		case strings.HasPrefix(args[i], "--config="):
			configPath = args[i][len("--config="):]
		default:
			remaining = append(remaining, args[i])
		}
	}
	return
}

func init() {
	var f *os.File
	var err error
	fileOnly := false

	// TempDir returns the default directory to use for temporary files.
	//
	// On Unix systems, it returns $TMPDIR if non-empty, else /tmp.
	// On Windows, it uses GetTempPath, returning the first non-empty
	// value from %TMP%, %TEMP%, %USERPROFILE%, or the Windows directory.
	// On Plan 9, it returns /tmp.
	// It does not guarantee the user can write to the directory;
	userTempDir := os.TempDir()
	if f, err = createLoggerFile(userTempDir); err != nil {
		if !strings.Contains(err.Error(), "permission denied") {
			slog.Error("Unable to create logger file", "error", err)
			os.Exit(1)
		}

		// If we can't write to the temp dir, try the user home dir
		userTempDir, err = os.UserHomeDir()
		if err != nil {
			slog.Error("Unable to get user home directory", "error", err)
			os.Exit(1)
		}

		if f, err = createLoggerFile(userTempDir); err != nil {
			slog.Error("Unable to create logger file in user home directory", "error", err)
			os.Exit(1)
		}
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

func createLoggerFile(userTempDir string) (*os.File, error) {
	now := time.Now()
	date := fmt.Sprintf("%s.log", now.Format("2006-01-02"))

	if err := os.MkdirAll(path.Join(userTempDir, ".seshtmp"), 0o755); err != nil {
		return nil, err
	}

	fileFullPath := path.Join(userTempDir, ".seshtmp", date)
	file, err := os.OpenFile(fileFullPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		return nil, err
	}

	return file, nil
}
