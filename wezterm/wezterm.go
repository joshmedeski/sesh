package wezterm

import (
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/shell"
)

type Wezterm interface {
	ListAllPanes() ([]model.WeztermPane, error)
	ListWorkspaces() ([]*model.WeztermWorkspace, error)
	SpawnWorkspace(name string, cwd string) (string, error)
	SendText(paneID int, text string) (string, error)
	GetText(paneID int) (string, error)
	IsInside() bool
	ActivatePane(paneID int) (string, error)
}

type RealWezterm struct {
	os    oswrap.Os
	shell shell.Shell
}

func NewWezterm(os oswrap.Os, shell shell.Shell) Wezterm {
	return &RealWezterm{os, shell}
}

func (w *RealWezterm) IsInside() bool {
	return len(w.os.Getenv("WEZTERM_PANE")) > 0
}

func (w *RealWezterm) SendText(paneID int, text string) (string, error) {
	return w.shell.Cmd("wezterm", "cli", "send-text", "--pane-id", itoa(paneID), "--", text)
}

func (w *RealWezterm) GetText(paneID int) (string, error) {
	return w.shell.Cmd("wezterm", "cli", "get-text", "--pane-id", itoa(paneID))
}

func (w *RealWezterm) ActivatePane(paneID int) (string, error) {
	return w.shell.Cmd("wezterm", "cli", "activate-pane", "--pane-id", itoa(paneID))
}
