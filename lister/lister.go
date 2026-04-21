package lister

import (
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/tmuxinator"
	"github.com/joshmedeski/sesh/v2/wezterm"
	"github.com/joshmedeski/sesh/v2/zoxide"
)

type Lister interface {
	List(opts ListOptions) (model.SeshSessions, error)
	FindTmuxSession(name string) (model.SeshSession, bool)
	GetAttachedTmuxSession() (model.SeshSession, bool)
	GetLastTmuxSession() (model.SeshSession, bool)
	FindConfigSession(name string) (model.SeshSession, bool)
	FindConfigWildcard(path string) (model.WildcardConfig, bool)
	FindZoxideSession(name string) (model.SeshSession, bool)
	FindTmuxinatorConfig(name string) (model.SeshSession, bool)
	FindWeztermWorkspace(name string) (model.SeshSession, bool)
	GetActiveWeztermWorkspace() (model.SeshSession, bool)
}

type RealLister struct {
	config     model.Config
	home       home.Home
	tmux       tmux.Tmux
	wezterm    wezterm.Wezterm
	zoxide     zoxide.Zoxide
	tmuxinator tmuxinator.Tmuxinator
}

func NewLister(config model.Config, home home.Home, tmux tmux.Tmux, wezterm wezterm.Wezterm, zoxide zoxide.Zoxide, tmuxinator tmuxinator.Tmuxinator) Lister {
	return &RealLister{config, home, tmux, wezterm, zoxide, tmuxinator}
}
