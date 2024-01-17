package git

import (
	"os/exec"
	"regexp"
	"strings"
)

func gitCmd(args []string) ([]byte, error) {
	tmux, err := exec.LookPath("git")
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
