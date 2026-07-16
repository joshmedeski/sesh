package dashboard

import (
	"runtime"
	"strings"

	"github.com/joshmedeski/sesh/v2/execwrap"
)

func runShellCommand(cmd string) ([]byte, error) {
	if runtime.GOOS == "windows" {
		useCmd := execwrap.NewExec().Command("cmd", "/c", cmd)
		return useCmd.CombinedOutput()
	}
	useCmd := execwrap.NewExec().Command("sh", "-c", cmd)
	return useCmd.CombinedOutput()
}

func runCommand(name string, args ...string) (string, error) {
	out, err := execwrap.NewExec().Command(name, args...).Output()
	if err != nil {
		errString := strings.TrimSpace(err.Error())
		if strings.Contains(errString, "no server running on") {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSuffix(string(out), "\n"), nil
}
