package convert

import (
	"fmt"
	"os"

	"github.com/joshmedeski/sesh/dir"
)

func PathToPretty(path string) string {
	prettyPath, err := dir.PrettyPath(path)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return prettyPath
}
