package seshcli

import (
	"errors"
	"fmt"
	"log/slog"

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
	"github.com/joshmedeski/sesh/v2/tui"
	"github.com/joshmedeski/sesh/v2/zoxide"
)

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

	// core dependencies
	ls := ls.NewLs(config, shell)
	lister := lister.NewLister(config, home, tmux, zoxide, tmuxinator)
	startup := startup.NewStartup(config, lister, tmux, home, replacer)
	namer := namer.NewNamer(path, git, home)
	connector := connector.NewConnector(config, dir, home, lister, namer, startup, tmux, zoxide, tmuxinator)
	icon := icon.NewIcon(config)
	previewer := previewer.NewPreviewer(lister, tmux, icon, dir, home, ls, config, shell)
	cloner := cloner.NewCloner(connector, git)

	// tui
	tui := tui.NewTui(lister)

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
		NewTuiCommand(tui),
	)

	return rootCmd
}
