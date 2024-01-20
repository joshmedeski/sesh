package git

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type CloneOptions struct {
	Dir    *string
	CmdDir *string
	Repo   string
}

type ClonedRepo struct {
	Name string
	Path string
}

func Clone(o CloneOptions) (ClonedRepo, error) {
	cmdArgs := []string{"clone", o.Repo}
	if o.Dir != nil && strings.TrimSpace(*o.Dir) != "" {
		cmdArgs = append(cmdArgs, *o.Dir)
	}
	cmd := exec.Command("git", cmdArgs...)
	if o.CmdDir != nil && strings.TrimSpace(*o.CmdDir) != "" {
		cmd.Dir = *o.CmdDir
	}

  cmd.Stdin = os.Stdin
  cmd.Stderr = os.Stderr
  cmd.Stdout = os.Stdout
  cmd.Env = os.Environ()

	err := cmd.Run()
	if err != nil {
		return ClonedRepo{}, err
	}
	name := findRepo(o.Repo)
	if o.Dir != nil && strings.TrimSpace(*o.Dir) != "" {
		name = *o.Dir
	}
	path := cmd.Dir + "/" + name
	return ClonedRepo{
		Name: name,
		Path: path,
	}, nil
}

func findRepo(repo string) string {
	repo = strings.TrimSuffix(repo, ".git")
	re := regexp.MustCompile(`([^\/]*)$`)
	match := re.FindString(repo)
	return match
}
