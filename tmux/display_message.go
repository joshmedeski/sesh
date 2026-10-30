package tmux

func (t *RealTmux) DisplayMessage(message string) (string, error) {
	return t.shell.Cmd(t.bin, "display-message", "-d", "0", escapeFormat(message))
}
