package seshcli

import (
	"github.com/joshmedeski/sesh/config"
	"github.com/joshmedeski/sesh/execwrap"
	"github.com/joshmedeski/sesh/home"
	"github.com/joshmedeski/sesh/lister"
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
	shell := shell.NewShell(exec)
	home := home.NewHome(os)

	// core dependencies
	tmux := tmux.NewTmux(shell)
	zoxide := zoxide.NewZoxide(shell)
	config := config.NewConfig(os, path, runtime)
	lister := lister.NewLister(config, home, tmux, zoxide)

	return cli.App{
		Name:    "sesh",
		Version: version,
		Usage:   "Smart session manager for the terminal",
		Commands: []*cli.Command{
			List(lister),
			Connect(),
			Clone(),
		},
	}
}
