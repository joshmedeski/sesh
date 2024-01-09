package git

import (
	"os/exec"
	"strings"
)

func RootPath(path string) string {
	gitRootPathCmd := exec.Command("git", "-C", path, "rev-parse", "--show-toplevel")
	gitRootPathByteOutput, err := gitRootPathCmd.CombinedOutput()
	if err != nil {
		return ""
	}
	gitRootPath := strings.TrimSpace(string(gitRootPathByteOutput))
	return gitRootPath
}
