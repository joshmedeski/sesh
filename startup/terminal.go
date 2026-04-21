package startup

import (
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/wezterm"
)

// Terminal abstracts the multiplexer operations needed by startup logic.
type Terminal interface {
	CreateWindow(cwd string, name string) (string, error)
	SendCommand(target string, command string) (string, error)
	FocusFirstWindow() (string, error)
}

// TmuxTerminal adapts tmux.Tmux to the Terminal interface.
type TmuxTerminal struct {
	tmux tmux.Tmux
}

func NewTmuxTerminal(tmux tmux.Tmux) Terminal {
	return &TmuxTerminal{tmux: tmux}
}

func (t *TmuxTerminal) CreateWindow(cwd string, name string) (string, error) {
	return t.tmux.NewWindow(cwd, name)
}

func (t *TmuxTerminal) SendCommand(target string, command string) (string, error) {
	return t.tmux.SendKeys(target, command)
}

func (t *TmuxTerminal) FocusFirstWindow() (string, error) {
	return t.tmux.NextWindow()
}

// WeztermTerminal adapts wezterm.Wezterm to the Terminal interface.
type WeztermTerminal struct {
	wezterm       wezterm.Wezterm
	workspaceName string
	lastPaneID    int
}

func NewWeztermTerminal(wezterm wezterm.Wezterm, workspaceName string) Terminal {
	return &WeztermTerminal{wezterm: wezterm, workspaceName: workspaceName}
}

func (w *WeztermTerminal) CreateWindow(cwd string, name string) (string, error) {
	// Spawn a new pane in the workspace. WezTerm's spawn creates a new tab.
	output, err := w.wezterm.SpawnWorkspace(w.workspaceName, cwd)
	if err != nil {
		return output, err
	}
	// Track the last pane we created for SendCommand.
	panes, err := w.wezterm.ListAllPanes()
	if err == nil {
		for _, p := range panes {
			if p.Workspace == w.workspaceName && p.PaneID > w.lastPaneID {
				w.lastPaneID = p.PaneID
			}
		}
	}
	return output, nil
}

func (w *WeztermTerminal) SendCommand(target string, command string) (string, error) {
	if w.lastPaneID > 0 {
		return w.wezterm.SendText(w.lastPaneID, command)
	}
	// Fallback: find the first pane in this workspace.
	panes, err := w.wezterm.ListAllPanes()
	if err != nil {
		return "", err
	}
	for _, p := range panes {
		if p.Workspace == w.workspaceName {
			return w.wezterm.SendText(p.PaneID, command)
		}
	}
	return "", nil
}

func (w *WeztermTerminal) FocusFirstWindow() (string, error) {
	panes, err := w.wezterm.ListAllPanes()
	if err != nil {
		return "", err
	}
	for _, p := range panes {
		if p.Workspace == w.workspaceName {
			return w.wezterm.ActivatePane(p.PaneID)
		}
	}
	return "", nil
}
