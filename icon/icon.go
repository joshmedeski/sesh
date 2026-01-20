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
	projectsIcon   string = ""
)

var defaultProjectIcons = map[string]string{
	".git":           "",
	"package.json":   "",
	"Cargo.toml":     "",
	"go.mod":         "",
	"pyproject.toml": "",
	"composer.json":  "",
	"Gemfile":        "",
	"mix.exs":        "",
	"pom.xml":        "",
	"build.gradle":   "",
	"Makefile":       "",
}

func ansiString(code int, s string) string {
	return fmt.Sprintf("\033[%dm%s\033[39m", code, s)
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
	case "projects":
		icon = i.getProjectIcon(s.ProjectType)
		colorCode = 35 // magenta
	}
	if icon != "" {
		return fmt.Sprintf("%s %s", ansiString(colorCode, icon), s.Name)
	}
	return s.Name
}

func (i *RealIcon) getProjectIcon(projectType string) string {
	if i.config.ProjectIcons != nil {
		if icon, ok := i.config.ProjectIcons[projectType]; ok {
			return icon
		}
	}
	if icon, ok := defaultProjectIcons[projectType]; ok {
		return icon
	}
	return projectsIcon
}

func (i *RealIcon) RemoveIcon(name string) string {
	// Format: \033[XXm<icon>\033[39m <name>
	if strings.HasPrefix(name, "\033[") {
		endIdx := strings.Index(name, "m")
		if endIdx != -1 {
			remaining := name[endIdx+1:]
			resetIdx := strings.Index(remaining, "\033[39m ")
			if resetIdx != -1 {
				return remaining[resetIdx+7:]
			}
		}
	}
	if strings.HasPrefix(name, tmuxIcon) || strings.HasPrefix(name, zoxideIcon) || strings.HasPrefix(name, configIcon) || strings.HasPrefix(name, tmuxinatorIcon) {
		return name[4:]
	}
	return name
}
