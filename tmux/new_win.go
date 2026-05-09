package tmux

func (t *RealTmux) NewWindowInSession(name string, startDir string, targetSession string) (string, error) {
	args := []string{"new-window", "-n", name, "-c", startDir}
	if targetSession != "" {
		args = append(args, "-t", targetSession)
	}
	return t.shell.Cmd("tmux", args...)
}
