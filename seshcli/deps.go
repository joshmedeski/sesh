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
	"github.com/joshmedeski/sesh/v2/zoxide"
)

// BaseDeps holds config-free dependencies that can be constructed eagerly.
type BaseDeps struct {
	Exec       execwrap.Exec
	Os         oswrap.Os
	Path       pathwrap.Path
	Runtime    runtimewrap.Runtime
	Home       home.Home
	Shell      shell.Shell
	Json       json.Json
	Replacer   replacer.Replacer
	Git        git.Git
	Dir        dir.Dir
	Tmux       tmux.Tmux
	Zoxide     zoxide.Zoxide
	Tmuxinator tmuxinator.Tmuxinator
}

// Deps holds all dependencies including config-dependent ones.
type Deps struct {
	BaseDeps
	Lister    lister.Lister
	Startup   startup.Startup
	Namer     namer.Namer
	Connector connector.Connector
	Icon      icon.Icon
	Previewer previewer.Previewer
	Cloner    cloner.Cloner
}

// NewBaseDeps constructs all config-free dependencies.
func NewBaseDeps() *BaseDeps {
	exec := execwrap.NewExec()
	os := oswrap.NewOs()
	path := pathwrap.NewPath()
	runtime := runtimewrap.NewRunTime()

	h := home.NewHome(os)
	sh := shell.NewShell(exec, h)
	j := json.NewJson()
	r := replacer.NewReplacer()

	g := git.NewGit(sh)
	d := dir.NewDir(os, g, path)
	t := tmux.NewTmux(os, sh)
	z := zoxide.NewZoxide(sh)
	ti := tmuxinator.NewTmuxinator(sh)

	return &BaseDeps{
		Exec:       exec,
		Os:         os,
		Path:       path,
		Runtime:    runtime,
		Home:       h,
		Shell:      sh,
		Json:       j,
		Replacer:   r,
		Git:        g,
		Dir:        d,
		Tmux:       t,
		Zoxide:     z,
		Tmuxinator: ti,
	}
}

// BuildAll loads config and constructs all config-dependent dependencies.
func (b *BaseDeps) BuildAll(configPath string) (*Deps, error) {
	config, err := configurator.NewConfiguratorWithPath(b.Os, b.Path, b.Runtime, configPath).GetConfig()
	if err != nil {
		return nil, err
	}

	slog.Debug("deps: BuildAll", "config", config)

	l := ls.NewLs(config, b.Shell)
	li := lister.NewLister(config, b.Home, b.Tmux, b.Zoxide, b.Tmuxinator)
	s := startup.NewStartup(config, li, b.Tmux, b.Home, b.Replacer)
	n := namer.NewNamer(b.Path, b.Git, b.Home, config)
	c := connector.NewConnector(config, b.Dir, b.Home, li, n, s, b.Tmux, b.Zoxide, b.Tmuxinator)
	ic := icon.NewIcon(config)
	p := previewer.NewPreviewer(li, b.Tmux, ic, b.Dir, b.Home, l, config, b.Shell)
	cl := cloner.NewCloner(c, b.Git)

	return &Deps{
		BaseDeps:  *b,
		Lister:    li,
		Startup:   s,
		Namer:     n,
		Connector: c,
		Icon:      ic,
		Previewer: p,
		Cloner:    cl,
	}, nil
}

// buildDeps reads the --config flag from cobra and builds all dependencies.
func buildDeps(cmd *cobra.Command, base *BaseDeps) (*Deps, error) {
	configPath, _ := cmd.Root().PersistentFlags().GetString("config")
	deps, err := base.BuildAll(configPath)
	if err != nil {
		var human *configurator.ConfigError
		if errors.As(err, &human) {
			fmt.Printf("Couldn't parse config, err: %v\n details:\n %s\n", err.Error(), human.Human())
		}
		slog.Error("buildDeps", "error", err)
		return nil, err
	}
	return deps, nil
}
