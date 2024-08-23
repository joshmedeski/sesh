package zoxide

func (z *RealZoxide) Add(path string) error {
	_, err := z.shell.Cmd("zoxide", "add", path)
	if err != nil {
		return err
	}
	return nil
}
