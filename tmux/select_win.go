package tmux

func (t *RealTmux) SelectWindow(targetWindow string) (string, error) {
	return t.shell.Cmd(t.bin, "select-window", "-t", targetWindow)
}
