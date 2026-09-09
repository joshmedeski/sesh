package icon

import (
	"fmt"
	"strings"

	"github.com/joshmedeski/sesh/v2/model"
)

type Icon interface {
	AddIcon(session model.SeshSession) string
	AddIconNoColor(session model.SeshSession) string
	RemoveIcon(name string) string
}

type RealIcon struct {
	config model.Config
	// mux names the running multiplexer backend; live "tmux" sessions take
	// its glyph so herdr workspaces don't show up with a tmux icon.
	mux string
}

func NewIcon(config model.Config, mux string) Icon {
	return &RealIcon{config: config, mux: mux}
}

var (
	zoxideIcon     string = ""
	tmuxIcon       string = ""
	herdrIcon      string = "󰳆"
	configIcon     string = ""
	tmuxinatorIcon string = ""
	tmuxPaneIcon   string = ""
)

// Glyph holds the icon character and ANSI color code for a session source.
type Glyph struct {
	Icon      string
	ColorCode int
}

// Glyphs maps session source names to their icon and color.
var Glyphs = map[string]Glyph{
	"tmux":       {Icon: tmuxIcon, ColorCode: 34},
	"herdr":      {Icon: herdrIcon, ColorCode: 34},
	"config":     {Icon: configIcon, ColorCode: 90},
	"zoxide":     {Icon: zoxideIcon, ColorCode: 36},
	"tmuxinator": {Icon: tmuxinatorIcon, ColorCode: 33},
	"tmux-pane":  {Icon: tmuxPaneIcon, ColorCode: 32},
}

func ansiString(code int, s string) string {
	return fmt.Sprintf("\033[%dm%s\033[39m", code, s)
}

// glyph resolves the glyph for a session. Live sessions are tagged "tmux"
// whichever multiplexer produced them, so those follow the backend name.
func (i *RealIcon) glyph(s model.SeshSession) (Glyph, bool) {
	src := s.Src
	if src == "tmux" && i.mux != "" {
		if g, ok := Glyphs[i.mux]; ok {
			return g, true
		}
	}
	g, ok := Glyphs[src]
	return g, ok
}

func (i *RealIcon) AddIcon(s model.SeshSession) string {
	if g, ok := i.glyph(s); ok {
		return fmt.Sprintf("%s %s", ansiString(g.ColorCode, g.Icon), s.Name)
	}
	return s.Name
}

func (i *RealIcon) AddIconNoColor(s model.SeshSession) string {
	if g, ok := i.glyph(s); ok {
		return fmt.Sprintf("%s %s", g.Icon, s.Name)
	}
	return s.Name
}

func (i *RealIcon) RemoveIcon(name string) string {
	for _, g := range Glyphs {
		if strings.HasPrefix(name, g.Icon+" ") {
			return name[len(g.Icon)+1:]
		}
	}
	return name
}
