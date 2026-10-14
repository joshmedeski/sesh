package icon

import (
	"errors"
	"fmt"
	"strings"

	"github.com/joshmedeski/sesh/v2/ansi"
	"github.com/joshmedeski/sesh/v2/model"
)

type Icon interface {
	AddIcon(session model.SeshSession) string
	AddIconNoColor(session model.SeshSession) string
	RemoveIcon(name string) string
}

type RealIcon struct {
	glyphs Glyphs
}

func NewIcon(config model.Config) Icon {
	glyphs, _ := SourceGlyphs(config.Icons)
	return &RealIcon{glyphs}
}

// Glyph holds the icon character and ANSI color code for a session source.
type Glyph struct {
	Icon      string
	ColorCode int
}

// Glyphs maps session source names to their icon and color.
type Glyphs map[string]Glyph

var defaultGlyphs = Glyphs{
	"tmux":       {Icon: "", ColorCode: ansi.Blue},
	"config":     {Icon: "", ColorCode: ansi.BrightBlack},
	"zoxide":     {Icon: "", ColorCode: ansi.Cyan},
	"tmuxinator": {Icon: "", ColorCode: ansi.Yellow},
	"tmux-pane":  {Icon: "", ColorCode: ansi.Green},
}

func SourceGlyphs(icons map[string]model.SourceIcon) (Glyphs, error) {
	glyphs := make(Glyphs, len(defaultGlyphs))
	for src, g := range defaultGlyphs {
		glyphs[src] = g
	}

	var errs []error
	for src, override := range icons {
		g, ok := glyphs[src]
		if !ok {
			errs = append(errs, fmt.Errorf("icons: unknown source %q (use tmux, config, zoxide, tmuxinator, or tmux-pane)", src))
			continue
		}
		if override.Color != "" {
			code, ok := ansi.ColorCode(override.Color)
			if !ok {
				errs = append(errs, fmt.Errorf("icons.%s.color: unknown color %q (use black, red, green, yellow, blue, magenta, cyan, white, or a bright-* variant)", src, override.Color))
				continue
			}
			g.ColorCode = code
		}
		if override.Icon != "" {
			g.Icon = override.Icon
		}
		glyphs[src] = g
	}
	return glyphs, errors.Join(errs...)
}

func ansiString(code int, s string) string {
	return fmt.Sprintf("\033[%dm%s\033[39m", code, s)
}

func (i *RealIcon) AddIcon(s model.SeshSession) string {
	if g, ok := i.glyphs[s.Src]; ok {
		return fmt.Sprintf("%s %s", ansiString(g.ColorCode, g.Icon), s.Name)
	}
	return s.Name
}

func (i *RealIcon) AddIconNoColor(s model.SeshSession) string {
	if g, ok := i.glyphs[s.Src]; ok {
		return fmt.Sprintf("%s %s", g.Icon, s.Name)
	}
	return s.Name
}

func (i *RealIcon) RemoveIcon(name string) string {
	for _, g := range i.glyphs {
		if rest, ok := strings.CutPrefix(name, g.Icon+" "); ok {
			return rest
		}
	}
	return name
}
