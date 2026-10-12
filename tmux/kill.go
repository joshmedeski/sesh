package tmux

func (t *RealTmux) KillSession(target string) (string, error) {
	return t.shell.Cmd(t.bin, "kill-session", "-t", sessionTarget(target))
}
