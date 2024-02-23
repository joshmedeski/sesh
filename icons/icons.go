package icons

import (
	"fmt"

	"github.com/joshmedeski/sesh/session"
)

// TODO: add to config to allow for custom icons
var (
	ZoxideIcon string = ""
	TmuxIcon   string = ""
)

func ansiString(code int, s string) string {
	return fmt.Sprintf("\033[%dm%s\033[39m", code, s)
}

func PrintWithIcon(s session.Session) string {
	icon := ZoxideIcon
	colorCode := 36 // cyan
	if s.Src == "tmux" {
		icon = TmuxIcon
		colorCode = 34 // blue
	}
	return fmt.Sprintf("%s %s", ansiString(colorCode, icon), s.Name)
}
