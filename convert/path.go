package convert

import (
	"fmt"
	"joshmedeski/sesh/dir"
	"os"
)

func PathToPretty(path string) string {
	prettyPath, err := dir.PrettyPath(path)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return prettyPath
}
