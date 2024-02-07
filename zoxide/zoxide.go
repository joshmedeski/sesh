package zoxide

import (
	"os/exec"
)

func zoxideCmd(args []string) ([]byte, error) {
	zoxide, err := exec.LookPath("zoxide")
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(zoxide, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return output, nil
}
