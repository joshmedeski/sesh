package seshcli

import (
	"fmt"
	"strings"

	"github.com/joshmedeski/sesh/v2/icon"
	"github.com/joshmedeski/sesh/v2/json"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/marker"
	"github.com/joshmedeski/sesh/v2/model"
	cli "github.com/urfave/cli/v2"
)

func List(icon icon.Icon, json json.Json, list lister.Lister, marker marker.Marker) *cli.Command {
	return &cli.Command{
		Name:                   "list",
		Aliases:                []string{"l"},
		Usage:                  "List sessions",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "show configured sessions",
			},
			&cli.BoolFlag{
				Name:    "json",
				Aliases: []string{"j"},
				Usage:   "output as json",
			},
			&cli.BoolFlag{
				Name:    "tmux",
				Aliases: []string{"t"},
				Usage:   "show tmux sessions",
			},
			&cli.BoolFlag{
				Name:    "zoxide",
				Aliases: []string{"z"},
				Usage:   "show zoxide results",
			},
			&cli.BoolFlag{
				Name:    "hide-attached",
				Aliases: []string{"H"},
				Usage:   "don't show currently attached sessions",
			},
			&cli.BoolFlag{
				Name:    "icons",
				Aliases: []string{"i"},
				Usage:   "show icons",
			},
			&cli.BoolFlag{
				Name:    "tmuxinator",
				Aliases: []string{"T"},
				Usage:   "show tmuxinator configs",
			},
			&cli.BoolFlag{
				Name:    "hide-duplicates",
				Aliases: []string{"d"},
				Usage:   "hide duplicate entries",
			},
			&cli.BoolFlag{
				Name:    "marked",
				Aliases: []string{"m"},
				Usage:   "show only marked sessions",
			},
		},
		Action: func(cCtx *cli.Context) error {
			sessions, err := list.List(lister.ListOptions{
				Config:         cCtx.Bool("config"),
				HideAttached:   cCtx.Bool("hide-attached"),
				Icons:          cCtx.Bool("icons"),
				Json:           cCtx.Bool("json"),
				Tmux:           cCtx.Bool("tmux"),
				Zoxide:         cCtx.Bool("zoxide"),
				Tmuxinator:     cCtx.Bool("tmuxinator"),
				HideDuplicates: cCtx.Bool("hide-duplicates"),
				Marked:         cCtx.Bool("marked"),
			}, marker)
			if err != nil {
				return fmt.Errorf("couldn't list sessions: %q", err)
			}

			if cCtx.Bool("json") {
				var sessionsArray []model.SeshSession
				for _, i := range sessions.OrderedIndex {
					sessionsArray = append(sessionsArray, sessions.Directory[i])
				}
				fmt.Println(json.EncodeSessions(sessionsArray))
				return nil
			}

			for _, i := range sessions.OrderedIndex {
				name := sessions.Directory[i].Name
				if cCtx.Bool("icons") {
					name = icon.AddIcon(sessions.Directory[i])
				}
				// Clean display name for better readability (spaces -> underscores)
				displayName := cleanDisplayName(name)
				fmt.Println(displayName)
			}

			return nil
		},
	}
}
// cleanDisplayName cleans session names for display purposes only
// This prevents visual confusion with spaces while preserving original names for connection
func cleanDisplayName(name string) string {
	// CONSERVATIVE: Only clean spaces in session names, and only if they look like tmux sessions
	// Don't touch config/zoxide paths or anything with slashes
	if strings.Contains(name, " ") && !strings.Contains(name, "/") && !strings.Contains(name, "~") {
		// Only replace spaces with underscores, leave dots and colons alone
		return strings.ReplaceAll(name, " ", "_")
	}
	return name
}
