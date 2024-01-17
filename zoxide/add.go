package zoxide

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

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
