package tmuxp

import (
	"bytes"
	"fmt"
	"os/exec"
)

func tmuxpCmd(args []string) ([]byte, error) {
	tmuxp, err := exec.LookPath("tmuxp")
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(tmuxp, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return nil, err
	}
	return output, nil
}
