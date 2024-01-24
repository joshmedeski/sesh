package zoxide

import (
	"fmt"
	"os/exec"
	"path"
)

func Add(result string) error {
	if !path.IsAbs(result) {
		return fmt.Errorf("can't add relative %q path to zoxide", result)
	}

	cmd := exec.Command("zoxide", "add", result)
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to add %q to zoxide: %w", result, err)
	}

	return nil
}
