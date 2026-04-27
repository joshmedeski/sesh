package tmux

func (t *RealTmux) NewWindowInSession(name string, startDir string, targetSession string, shellCommand string) (string, error) {
	args := []string{"new-window", "-n", name, "-c", startDir}
	if targetSession != "" {
		args = append(args, "-t", targetSession)
	}
	if shellCommand != "" {
		args = append(args, shellCommand)
	}
	return t.shell.Cmd(t.bin, args...)
}
