package projects

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
)

type Projects interface {
	List() ([]model.SeshSession, error)
}

type RealProjects struct {
	config model.Config
	home   home.Home
}

var markerSpecificity = map[string]int{
	".git":           1,
	"Makefile":       1,
	"package.json":   10,
	"Cargo.toml":     10,
	"go.mod":         10,
	"pyproject.toml": 10,
	"composer.json":  10,
	"Gemfile":        10,
	"mix.exs":        10,
	"pom.xml":        10,
	"build.gradle":   10,
}

func NewProjects(config model.Config, home home.Home) Projects {
	if len(config.ProjectMarkers) == 0 {
		config.ProjectMarkers = []string{
			".git",
			"package.json",
			"Cargo.toml",
			"go.mod",
			"pyproject.toml",
			"composer.json",
			"Gemfile",
			"mix.exs",
			"pom.xml",
			"build.gradle",
			"Makefile",
		}
	}
	if config.MaxDepth <= 0 {
		config.MaxDepth = 3
	}
	return &RealProjects{config, home}
}

func (p *RealProjects) List() ([]model.SeshSession, error) {
	var sessions []model.SeshSession
	seen := make(map[string]bool)

	for _, root := range p.config.ProjectRoots {
		expandedRoot, err := p.home.ExpandHome(root)
		if err != nil {
			continue
		}

		err = filepath.WalkDir(expandedRoot, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			if !d.IsDir() {
				return nil
			}

			rel, err := filepath.Rel(expandedRoot, path)
			if err != nil {
				return nil
			}

			if rel != "." {
				depth := len(strings.Split(rel, string(filepath.Separator)))
				if depth > p.config.MaxDepth {
					return filepath.SkipDir
				}
			}

			entries, err := os.ReadDir(path)
			if err != nil {
				return nil
			}

			bestMarker := ""
			bestSpecificity := 0
			for _, entry := range entries {
				for _, marker := range p.config.ProjectMarkers {
					if entry.Name() == marker {
						specificity := markerSpecificity[marker]
						if specificity == 0 {
							specificity = 5
						}
						if specificity > bestSpecificity {
							bestSpecificity = specificity
							bestMarker = marker
						}
					}
				}
			}

			if bestMarker != "" {
				if !seen[path] {
					name, err := p.home.ShortenHome(path)
					if err != nil {
						name = path
					}
					sessions = append(sessions, model.SeshSession{
						Src:         "projects",
						Name:        name,
						Path:        path,
						ProjectType: bestMarker,
					})
					seen[path] = true
				}
			}

			return nil
		})
	}

	return sessions, nil
}
