package icon

import (
	"fmt"
	"strings"

	"github.com/joshmedeski/sesh/v2/model"
)

type Icon interface {
	AddIcon(session model.SeshSession) string
	RemoveIcon(name string) string
}

type RealIcon struct {
	config model.Config
}

func NewIcon(config model.Config) Icon {
	return &RealIcon{config}
}

var (
	zoxideIcon     string = ""
	tmuxIcon       string = ""
	configIcon     string = ""
	tmuxinatorIcon string = ""
)

func ansiString(code int, s string) string {
	return fmt.Sprintf("\033[%dm%s\033[39m", code, s)
}

func ansiBackground(code int, s string) string {
	return fmt.Sprintf("\033[%dm%s\033[49m", code, s)
}

func getAlertBackground(alertLevel int) int {
	switch alertLevel {
	case 1:
		return 100 // Light gray background (rose-pine surface)
	case 2:
		return 101 // Light red background (rose-pine love)
	case 3:
		return 41  // Red background with emphasis
	default:
		return 0 // No background
	}
}

func (i *RealIcon) AddIcon(s model.SeshSession) string {
	var icon string
	var colorCode int
	switch s.Src {
	case "tmux":
		icon = tmuxIcon
		colorCode = 34 // blue
	case "tmuxinator":
		icon = tmuxinatorIcon
		colorCode = 33 // yellow
	case "zoxide":
		icon = zoxideIcon
		colorCode = 36 // cyan
	case "config":
		icon = configIcon
		colorCode = 90 // gray
	}
	
	var result string
	if icon != "" {
		result = fmt.Sprintf("%s %s", ansiString(colorCode, icon), s.Name)
	} else {
		result = s.Name
	}
	
	if s.Marked {
		result = fmt.Sprintf("📌 %s", result)
		
		if s.AlertLevel > 0 {
			bgCode := getAlertBackground(s.AlertLevel)
			if bgCode > 0 {
				result = ansiBackground(bgCode, result)
			}
		}
	}
	
	return result
}

func (i *RealIcon) RemoveIcon(name string) string {
	if strings.HasPrefix(name, tmuxIcon) || strings.HasPrefix(name, zoxideIcon) || strings.HasPrefix(name, configIcon) || strings.HasPrefix(name, tmuxinatorIcon) {
		return name[4:]
	}
	return name
}
