package builder

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/model"
)

func determineCmd(b *RealBuilder, path string) []string {
	if _, err := b.os.Stat(path + "/Makefile"); err == nil {
		return []string{"make", "build"}
	}
	if _, err := b.os.Stat(path + "/package-lock.json"); err == nil {
		return []string{"npm", "run", "build"}
	}
	if _, err := b.os.Stat(path + "/pnpm-lock.yaml"); err == nil {
		return []string{"pnpm", "run", "build"}
	}
	if _, err := b.os.Stat(path + "/yarn.lock"); err == nil {
		return []string{"yarn", "build"}
	}
	return []string{}
}

func (b *RealBuilder) Build(session model.SeshSession) (string, error) {
	cmd := determineCmd(b, session.Path)
	if len(cmd) == 0 {
		return "", fmt.Errorf("no build command found for path: %s", session.Path)
	}
	fmt.Println("Building with command:", fmt.Sprintf("%v", cmd))
	cmdOutput, err := b.shell.Cmd(cmd[0], cmd[1:]...)
	if err != nil {
		return "", err
	}
	return cmdOutput, nil
}
