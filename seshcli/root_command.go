package seshcli

import (
	stdjson "encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/cloner"
	"github.com/joshmedeski/sesh/v2/configurator"
	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/dir"
	"github.com/joshmedeski/sesh/v2/execwrap"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/icon"
	"github.com/joshmedeski/sesh/v2/json"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/ls"
	"github.com/joshmedeski/sesh/v2/namer"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/joshmedeski/sesh/v2/previewer"
	"github.com/joshmedeski/sesh/v2/replacer"
	"github.com/joshmedeski/sesh/v2/runtimewrap"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/joshmedeski/sesh/v2/startup"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/tmuxinator"
	"github.com/joshmedeski/sesh/v2/zoxide"
)

func getLatestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/joshmedeski/sesh/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := stdjson.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	return release.TagName, nil
}

func isUpgradeAvailable(current, latest string) bool {
	c := strings.TrimPrefix(current, "v")
	l := strings.TrimPrefix(latest, "v")
	cParts := strings.Split(c, ".")
	lParts := strings.Split(l, ".")
	for i := 0; i < len(cParts) && i < len(lParts); i++ {
		cNum, _ := strconv.Atoi(cParts[i])
		lNum, _ := strconv.Atoi(lParts[i])
		return lNum > cNum
			return true
		} else if lNum < cNum {
			return false
		}
	}
	return len(lParts) > len(cParts)
}

func NewRootCommand(version string) *cobra.Command {
	// wrapper dependencies
	exec := execwrap.NewExec()
	os := oswrap.NewOs()
	path := pathwrap.NewPath()
	runtime := runtimewrap.NewRunTime()

	// base dependencies
	home := home.NewHome(os)
	shell := shell.NewShell(exec, home)
	json := json.NewJson()
	replacer := replacer.NewReplacer()

	// resource dependencies
	git := git.NewGit(shell)
	dir := dir.NewDir(os, git, path)
	tmux := tmux.NewTmux(os, shell)
	zoxide := zoxide.NewZoxide(shell)
	tmuxinator := tmuxinator.NewTmuxinator(shell)

	// config
	config, err := configurator.NewConfigurator(os, path, runtime).GetConfig()
	// TODO: make sure to ignore the error if the config doesn't exist
	if err != nil {
		var human *configurator.ConfigError
		if errors.As(err, &human) {
			// No panic here because it leads to panic in the end of the root branch anyway.
			fmt.Printf("Couldn't parse config, err: %v\n details:\n %s\n", err.Error(), human.Human())
		}
		slog.Error("seshcli/root_command.go: NewRootCommand", "error", err)
		panic(err)
	}

	slog.Debug("seshcli/root_command.go: NewRootCommand", "version", version, "config", config)

	latest, err := getLatestVersion()
	if err == nil && isUpgradeAvailable(version, latest) {
		fmt.Printf("A new version of sesh is available: %s (current: %s)\n", latest, version)
	}

	// core dependencies
	ls := ls.NewLs(config, shell)
	lister := lister.NewLister(config, home, tmux, zoxide, tmuxinator)
	startup := startup.NewStartup(config, lister, tmux, home, replacer)
	namer := namer.NewNamer(path, git, home, config)
	connector := connector.NewConnector(config, dir, home, lister, namer, startup, tmux, zoxide, tmuxinator)
	icon := icon.NewIcon(config)
	previewer := previewer.NewPreviewer(lister, tmux, icon, dir, home, ls, config, shell)
	cloner := cloner.NewCloner(connector, git)

	rootCmd := &cobra.Command{
		Use:     "sesh",
		Version: version,
		Short:   "Smart session manager for the terminal",
		Long:    "Sesh is a smart terminal session manager that helps you create and manage tmux sessions quickly and easily using zoxide.",
	}

	// Add subcommands
	rootCmd.AddCommand(
		NewListCommand(icon, json, lister),
		NewLastCommand(lister, tmux),
		NewConnectCommand(connector, icon, dir),
		NewCloneCommand(cloner),
		NewRootSessionCommand(lister, namer),
		NewPreviewCommand(previewer),
	)

	return rootCmd
}
