package wezterm

import "strconv"

func itoa(n int) string {
	return strconv.Itoa(n)
}

func (w *RealWezterm) SpawnWorkspace(name string, cwd string) (string, error) {
	return w.shell.Cmd("wezterm", "cli", "spawn", "--workspace", name, "--cwd", cwd, "--new-window")
}
