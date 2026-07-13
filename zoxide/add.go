package zoxide

func (z *RealZoxide) Add(path string) error {
	parts, err := z.shell.PrepareCmd(z.addCommand, map[string]string{"{}": path})
	if err != nil {
		return err
	}
	if _, err := z.shell.Cmd(parts[0], parts[1:]...); err != nil {
		return err
	}
	return nil
}
