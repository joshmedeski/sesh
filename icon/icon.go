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

// Glyph holds the icon character and ANSI color code for a session source.
type Glyph struct {
	Icon      string
	ColorCode int
}

// Glyphs maps session source names to their icon and color.
var Glyphs = map[string]Glyph{
	"tmux":       {Icon: tmuxIcon, ColorCode: 34},
	"config":     {Icon: configIcon, ColorCode: 90},
	"zoxide":     {Icon: zoxideIcon, ColorCode: 36},
	"tmuxinator": {Icon: tmuxinatorIcon, ColorCode: 33},
}


func ansiString(code int, s string) string {
	return fmt.Sprintf("\033[%dm%s\033[39m", code, s)
}

func (i *RealIcon) AddIcon(s model.SeshSession) string {
	if g, ok := Glyphs[s.Src]; ok {
		return fmt.Sprintf("%s %s", ansiString(g.ColorCode, g.Icon), s.Name)
	}
	return s.Name
}

func (i *RealIcon) RemoveIcon(name string) string {
	if strings.HasPrefix(name, tmuxIcon) || strings.HasPrefix(name, zoxideIcon) || strings.HasPrefix(name, configIcon) || strings.HasPrefix(name, tmuxinatorIcon) {
		return name[4:]
	}
	return name
}
