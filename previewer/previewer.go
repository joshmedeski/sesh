package previewer

import (
	"github.com/joshmedeski/sesh/dir"
	"github.com/joshmedeski/sesh/home"
	"github.com/joshmedeski/sesh/icon"
	"github.com/joshmedeski/sesh/lister"
	"github.com/joshmedeski/sesh/ls"
	"github.com/joshmedeski/sesh/tmux"
)

type Previewer interface {
	// Previews a session or directory
	Preview(name string) (string, error)
}

type RealPreviewer struct {
	icon       icon.Icon
	strategies []PreviewStrategy
}

func NewPreviewer(
	lister lister.Lister,
	tmux tmux.Tmux,
	icon icon.Icon,
	dir dir.Dir,
	home home.Home,
	ls ls.Ls,
) Previewer {
	strategies := make([]PreviewStrategy, 0, 3)
	strategies = append(strategies,
		NewTmuxStrategy(lister, tmux),
		NewConfigStrategy(lister, ls),
		NewDirectoryStrategy(home, dir, ls),
	)

	return &RealPreviewer{
		icon:       icon,
		strategies: strategies,
	}
}

func (p *RealPreviewer) Preview(name string) (string, error) {
	trimmedName := p.icon.RemoveIcon(name)

	for _, strategy := range p.strategies {
		output, err := strategy.Execute(trimmedName)
		if err != nil {
			return "", err
		}
		if output != "" {
			return output, nil
		}
	}
	return "", nil
}
