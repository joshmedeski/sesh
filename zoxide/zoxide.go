package zoxide

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func zoxideCmd(args []string) ([]byte, error) {
	tmux, err := exec.LookPath("zoxide")
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(tmux, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return output, nil
}

func Add(result string) {
	if !path.IsAbs(result) {
		return
	}
	cmd := exec.Command("zoxide", "add", result)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
