package lister

import (
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/tmuxinator"
	"github.com/joshmedeski/sesh/v2/zoxide"
)

type Lister interface {
	List(opts ListOptions) (model.SeshSessions, error)
	ListTmuxPanes() (model.SeshSessions, error)
	FindTmuxSession(name string) (model.SeshSession, bool)
	GetAttachedTmuxSession() (model.SeshSession, bool)
	GetLastTmuxSession() (model.SeshSession, bool)
	FindConfigSession(name string) (model.SeshSession, bool)
	FindConfigWildcard(path string) (model.WildcardConfig, bool)
	FindZoxideSession(name string) (model.SeshSession, bool)
	FindTmuxinatorConfig(name string) (model.SeshSession, bool)
}

type RealLister struct {
	os         oswrap.Os
	config     model.Config
	home       home.Home
	tmux       tmux.Tmux
	zoxide     zoxide.Zoxide
	tmuxinator tmuxinator.Tmuxinator
}

func NewLister(os oswrap.Os, config model.Config, home home.Home, tmux tmux.Tmux, zoxide zoxide.Zoxide, tmuxinator tmuxinator.Tmuxinator) Lister {
	return &RealLister{os, config, home, tmux, zoxide, tmuxinator}
}
