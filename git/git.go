package git

import (
	"os/exec"
	"regexp"
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

func WorktreePath(path string) string {
	gitWorktreePathCmd := exec.Command("git", "-C", path, "rev-parse", "--git-common-dir")
	gitWorktreePathByteOutput, err := gitWorktreePathCmd.CombinedOutput()
	if err != nil {
		return ""
	}
	gitWorktreePath := strings.TrimSpace(string(gitWorktreePathByteOutput))
	match, _ := regexp.MatchString(`^(\.\./)*\.git$`, gitWorktreePath)
	if match {
		return ""
	}
	suffixes := []string{"/.git", "/.bare"}
	for _, suffix := range suffixes {
		gitWorktreePath = strings.TrimSuffix(gitWorktreePath, suffix)
	}
	return gitWorktreePath
}
