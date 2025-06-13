package namer

import (
	"errors"
	"strings"
)

func parseBareFromWorkTreeList(worktreeList string) (string, error) {
	fields := strings.Fields(worktreeList)
	if len(fields) == 0 {
		return "", errors.New("error parsing git worktree fields")
	}

	bareRoot := fields[0]
	if strings.HasSuffix(bareRoot, "/.bare") {
		bareRoot = strings.TrimSuffix(bareRoot, "/.bare")
	}

	return bareRoot, nil
}
