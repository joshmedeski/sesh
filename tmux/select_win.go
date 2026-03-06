package tmux

func (t *RealTmux) SelectWindow(name string) (string, error) {
	return t.shell.Cmd("tmux", "select-window", "-t", name)
}
