package tmux

func (t *RealTmux) SelectLayout(targetWindow string, layout string) (string, error) {
	return t.shell.Cmd(t.bin, "select-layout", "-t", sessionTarget(targetWindow), layout)
}
