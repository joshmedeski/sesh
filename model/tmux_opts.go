package model

// TmuxWindowOpts are the arguments for `tmux new-window`.
type TmuxWindowOpts struct {
	Name          string
	StartDir      string
	TargetSession string
	// Command becomes the new pane's process, so the window closes when the
	// command exits. Empty starts the default shell.
	Command    string
	Background bool // -d: create the window without making it the active one
}
