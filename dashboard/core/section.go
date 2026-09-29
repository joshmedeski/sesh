// section.go
package core

import (
	tea "charm.land/bubbletea/v2"
	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/joshmedeski/sesh/v2/tmux"
)

// Section is the contract every dashboard pane implements: the two permanent
// lists (sessions, configured) and the optional widget sections.
type Section interface {
	Name() string
	Init() tea.Cmd
	Update(msg tea.Msg) (Section, tea.Cmd)
	// ViewBorderless returns the pane title (drawn on the shared frame border)
	// and the borderless content (already clipped/padded to the given height).
	ViewBorderless(width, height int, focused bool) (title string, content string)
	Chosen() string
	TotalItems() int
	Width() float64
}

// Sorter is implemented by sections whose list can be re-sorted by the `s`
// key. SortLabel returns the current sort mode's name.
type Sorter interface {
	SortLabel() string
}

// Filterer is implemented by sections that support type-to-filter. While
// Filtering() is true the Model routes every key (except ctrl+c) to the pane
// so quit/switch/nav keys don't fire mid-typing.
type Filterer interface {
	Filtering() bool
	FilterQuery() string
}

// Clicker is implemented by list sections that move their selection to a
// specific row in response to a mouse click. row is 0-based and
// view-relative (the line within the pane's content area).
type Clicker interface {
	ClickAt(row int)
}

// SectionDeps carries the external dependencies every section needs.
type SectionDeps struct {
	Tmux      tmux.Tmux
	Lister    lister.Lister
	Git       git.Git
	Connector connector.Connector
	Shell     shell.Shell
	Runner    CommandRunner
	HomeDir   string
}
