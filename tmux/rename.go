package tmux

// RenameSession renames the tmux session identified by target to newName.
// newName may contain spaces (it is passed as a single argv element, so no
// shell quoting is required).
func (t *RealTmux) RenameSession(target string, newName string) (string, error) {
	return t.shell.Cmd(t.bin, "rename-session", "-t", target, newName)
}
