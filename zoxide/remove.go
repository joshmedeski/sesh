package zoxide

// Remove drops a path from the frecency backend. Nothing else in sesh writes
// to the backend destructively, so the caller is expected to have confirmed
// the removal with the user first.
func (z *RealZoxide) Remove(path string) error {
	parts, err := z.shell.PrepareCmd(z.removeCommand, map[string]string{"{}": path})
	if err != nil {
		return err
	}
	if _, err := z.shell.Cmd(parts[0], parts[1:]...); err != nil {
		return err
	}
	return nil
}
