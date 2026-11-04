package tmux

func (t *RealTmux) DisplayMessage(message string) (string, error) {
	return t.shell.Cmd(t.bin, displayMessageArgs(message)...)
}

func displayMessageArgs(message string) []string {
	return []string{"display-message", "-d", "0", escapeFormat(message)}
}
