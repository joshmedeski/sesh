package seshcli

import (
	"log/slog"

	"github.com/urfave/cli/v2"

	"github.com/joshmedeski/sesh/cloner"
	"github.com/joshmedeski/sesh/configurator"
	"github.com/joshmedeski/sesh/connector"
	"github.com/joshmedeski/sesh/dir"
	"github.com/joshmedeski/sesh/execwrap"
	"github.com/joshmedeski/sesh/git"
	"github.com/joshmedeski/sesh/home"
	"github.com/joshmedeski/sesh/icon"
	"github.com/joshmedeski/sesh/json"
	"github.com/joshmedeski/sesh/lister"
	"github.com/joshmedeski/sesh/ls"
	"github.com/joshmedeski/sesh/namer"
	"github.com/joshmedeski/sesh/oswrap"
	"github.com/joshmedeski/sesh/pathwrap"
	"github.com/joshmedeski/sesh/previewer"
	"github.com/joshmedeski/sesh/runtimewrap"
	"github.com/joshmedeski/sesh/shell"
	"github.com/joshmedeski/sesh/startup"
	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/tmuxinator"
	"github.com/joshmedeski/sesh/zoxide"
)

func App(version string) cli.App {
	// wrapper dependencies
	exec := execwrap.NewExec()
	os := oswrap.NewOs()
	path := pathwrap.NewPath()
	runtime := runtimewrap.NewRunTime()

	// base dependencies
	home := home.NewHome(os)
	shell := shell.NewShell(exec, home)
	json := json.NewJson()

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
		slog.Error("seshcli/seshcli.go: App", "error", err)
		panic(err)
	}

	slog.Debug("seshcli/seshcli.go: App", "version", version, "config", config)

	// core dependencies
	ls := ls.NewLs(config, shell)
	lister := lister.NewLister(config, home, tmux, zoxide, tmuxinator)
	startup := startup.NewStartup(config, lister, tmux)
	namer := namer.NewNamer(path, git, home)
	connector := connector.NewConnector(config, dir, home, lister, namer, startup, tmux, zoxide, tmuxinator)
	icon := icon.NewIcon(config)
	previewer := previewer.NewPreviewer(lister, tmux, icon, dir, home, ls, config, shell)
	cloner := cloner.NewCloner(connector, git)

	return cli.App{
		Name:    "sesh",
		Version: version,
		Usage:   "Smart session manager for the terminal",
		Commands: []*cli.Command{
			List(icon, json, lister),
			Last(lister, tmux),
			Connect(connector, icon, dir),
			Clone(cloner),
			Root(lister, namer),
			Preview(previewer),
		},
	}
}
