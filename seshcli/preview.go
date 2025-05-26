package seshcli

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/joshmedeski/sesh/v2/previewer"
	cli "github.com/urfave/cli/v2"
)

func Preview(p previewer.Previewer) *cli.Command {
	return &cli.Command{
		Name:                   "preview",
		Aliases:                []string{"p"},
		Usage:                  "Preview a session or directory",
		UseShortOptionHandling: true,
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() != 1 {
				return errors.New("session name or directory is required")
			}

			name := cCtx.Args().First()

			// Clean the name to handle fzf input with icons and ANSI codes
			cleanedName := cleanNameForPreview(name)
			
			output, err := p.Preview(cleanedName)
			if err != nil {
				return cli.Exit(err, 1)
			}

			fmt.Print(output)

			return nil
		},
	}
}

// cleanNameForPreview removes ANSI codes, icons, and preserves window info for tmux strategy
func cleanNameForPreview(name string) string {
	// Remove ANSI color codes
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m|\[[0-9;]*m`)
	cleaned := ansiRegex.ReplaceAllString(name, "")
	
	// Remove emoji icons (📌, etc.) and extra whitespace
	cleaned = strings.TrimSpace(cleaned)
	cleaned = regexp.MustCompile(`^[📌🟢⚙️📁⚡[:space:]]*`).ReplaceAllString(cleaned, "")
	cleaned = strings.TrimSpace(cleaned)
	
	// Handle marked window format: "sesh config:nvim(4)" -> "sesh config:4" for tmux targeting
	if strings.Contains(cleaned, ":") && strings.Contains(cleaned, "(") {
		// Extract session and window number: "sesh config:nvim(4)" -> "sesh config:4"
		re := regexp.MustCompile(`^(.+):(.+)\((\d+)\)$`)
		if matches := re.FindStringSubmatch(cleaned); len(matches) == 4 {
			sessionName := strings.TrimSpace(matches[1])
			windowNumber := matches[3]
			cleaned = sessionName + ":" + windowNumber
		}
	}
	
	return cleaned
}
