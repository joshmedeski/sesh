package seshcli

import (
	"github.com/joshmedeski/sesh/configurator"
	"github.com/joshmedeski/sesh/connector"
	"github.com/joshmedeski/sesh/dir"
	"github.com/joshmedeski/sesh/execwrap"
	"github.com/joshmedeski/sesh/git"
	"github.com/joshmedeski/sesh/home"
	"github.com/joshmedeski/sesh/lister"
	"github.com/joshmedeski/sesh/namer"
	"github.com/joshmedeski/sesh/oswrap"
	"github.com/joshmedeski/sesh/pathwrap"
	"github.com/joshmedeski/sesh/runtimewrap"
	"github.com/joshmedeski/sesh/shell"
	"github.com/joshmedeski/sesh/tmux"
	"github.com/joshmedeski/sesh/zoxide"
	"github.com/urfave/cli/v2"
)

func App(version string) cli.App {
	// wrapper dependencies
	exec := execwrap.NewExec()
	os := oswrap.NewOs()
	path := pathwrap.NewPath()
	runtime := runtimewrap.NewRunTime()

	// base dependencies
	dir := dir.NewDir(os, path)
	shell := shell.NewShell(exec)
	home := home.NewHome(os)

	// resource dependencies
	git := git.NewGit(shell)
	tmux := tmux.NewTmux(os, shell)
	zoxide := zoxide.NewZoxide(shell)

	// config
	config, err := configurator.NewConfigurator(os, path, runtime).GetConfig()
	// TODO: make sure to ignore the error if the config doesn't exist
	if err != nil {
		panic(err)
	}

	// core dependencies
	lister := lister.NewLister(config, home, tmux, zoxide)
	namer := namer.NewNamer(path, git)
	connector := connector.NewConnector(config, dir, home, lister, namer, tmux, zoxide)

	return cli.App{
		Name:    "sesh",
		Version: version,
		Usage:   "Smart session manager for the terminal",
		Commands: []*cli.Command{
			List(lister),
			Connect(connector),
			Clone(),
		},
	}
}
