package tmuxp

import "strings"

func Load(sessionName string, workspaceFile string) error {
	trimmedSessionName := strings.Trim(sessionName, "\"'")
	_, err := tmuxpCmd([]string{"load", "-d", "-s", trimmedSessionName, workspaceFile})
	if err != nil {
		return err
	}
	return nil
}
