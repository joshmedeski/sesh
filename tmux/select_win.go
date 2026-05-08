package tmux

func (t *RealTmux) SelectWindow(targetWindow string) (string, error) {
	return t.shell.Cmd("tmux", "select-window", "-t", targetWindow)
}
