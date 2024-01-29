package zoxide

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func Add(result string) error {
	p, err := filepath.Abs(result)
	if err != nil {
		return fmt.Errorf("can't add %q path to zoxide", result)
	}

	cmd := exec.Command("zoxide", "add", p)
	_, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to add %q to zoxide: %w", p, err)
	}

	return nil
}
