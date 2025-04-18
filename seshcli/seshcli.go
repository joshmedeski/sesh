package seshcli

import (
	"log/slog"

	"github.com/urfave/cli/v2"

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
	"github.com/joshmedeski/sesh/v2/runtimewrap"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/joshmedeski/sesh/v2/startup"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/tmuxinator"
	"github.com/joshmedeski/sesh/v2/zoxide"
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
