package lister

import (
	"sync"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/tmuxinator"
	"github.com/joshmedeski/sesh/v2/zoxide"
)

type Lister interface {
	List(opts ListOptions) (model.SeshSessions, error)
	ListTmuxPanes() (model.SeshSessions, error)
	FindTmuxSession(name string) (model.SeshSession, bool)
	FindTmuxSessionByBase(base string) (model.SeshSession, bool)
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

	// wildcards caches config.WildcardConfigs with their patterns expanded, so
	// resolving a wildcard for every session in a list expands each pattern
	// once. See expandedWildcards.
	wildcardsOnce sync.Once
	wildcards     []expandedWildcard
}

func NewLister(config model.Config, home home.Home, tmux tmux.Tmux, zoxide zoxide.Zoxide, tmuxinator tmuxinator.Tmuxinator) Lister {
	return &RealLister{
		config:     config,
		home:       home,
		tmux:       tmux,
		zoxide:     zoxide,
		tmuxinator: tmuxinator,
	}
}
