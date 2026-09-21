package tmux

// KillSession kills the tmux session identified by target. The trailing colon
// forces session (not window) target resolution, so a session whose name looks
// like a window index still resolves to itself.
func (t *RealTmux) KillSession(target string) (string, error) {
	return t.shell.Cmd(t.bin, "kill-session", "-t", target+":")
}
