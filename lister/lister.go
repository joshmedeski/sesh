package lister

import (
	"github.com/Wingsdh/cc-sesh/v2/home"
	"github.com/Wingsdh/cc-sesh/v2/model"
	"github.com/Wingsdh/cc-sesh/v2/tmux"
	"github.com/Wingsdh/cc-sesh/v2/tmuxinator"
	"github.com/Wingsdh/cc-sesh/v2/zoxide"
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
	config     model.Config
	home       home.Home
	tmux       tmux.Tmux
	zoxide     zoxide.Zoxide
	tmuxinator tmuxinator.Tmuxinator
}

func NewLister(config model.Config, home home.Home, tmux tmux.Tmux, zoxide zoxide.Zoxide, tmuxinator tmuxinator.Tmuxinator) Lister {
	return &RealLister{config, home, tmux, zoxide, tmuxinator}
}
