package dashboard

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/joshmedeski/sesh/v2/model"
)

type gitReposLoadedMsg struct {
	repos []gitRepo
}

type gitRepo struct {
	Path    string
	Name    string
	Branch  string
	Status  string
	IsRepo  bool
}

type GitSection struct {
	config model.DashboardSectionConfig
	deps   SectionDeps
	repos  []gitRepo
	cursor int
	chosen string
	loading bool
}

func NewGitSection(cfg model.DashboardSectionConfig, deps SectionDeps) Section {
	return &GitSection{
		config:  cfg,
		deps:    deps,
		loading: true,
	}
}

func (s *GitSection) Name() string    { return s.config.Title }
func (s *GitSection) TotalItems() int { return len(s.repos) }
func (s *GitSection) Width() float64  { return s.config.Width }
func (s *GitSection) Chosen() string  { return s.chosen }

func (s *GitSection) Init() tea.Cmd {
	return s.fetchRepos
}

func (s *GitSection) fetchRepos() tea.Msg {
	paths := s.config.Git.Paths
	if len(paths) == 0 {
		return gitReposLoadedMsg{repos: nil}
	}

	repos := make([]gitRepo, 0, len(paths))
	for _, p := range paths {
		expanded := p
		if strings.HasPrefix(p, "~/") {
			expanded = filepath.Join(s.deps.HomeDir, p[2:])
		}

		branch, err := runCommand("git", "-C", expanded, "rev-parse", "--abbrev-ref", "HEAD")
		if err != nil || strings.TrimSpace(branch) == "" {
			repos = append(repos, gitRepo{
				Path:   p,
				Name:   filepath.Base(expanded),
				IsRepo: false,
			})
			continue
		}

		statusOut, err := runCommand("git", "-C", expanded, "status", "--porcelain")
		status := ""
		if err == nil {
			lines := strings.Split(strings.TrimSpace(statusOut), "\n")
			if len(lines) > 0 && lines[0] != "" {
				parts := make([]string, 0, 4)
				for _, line := range lines {
					if strings.HasPrefix(line, "?? ") {
						parts = append(parts, fmt.Sprintf("!%d", 1))
						continue
					}
					first := line[0]
					second := line[1]
					if first != ' ' {
						parts = append(parts, fmt.Sprintf("+%d", 1))
					}
					if second == 'M' {
						parts = append(parts, fmt.Sprintf("~%d", 1))
					}
					if first == 'D' || second == 'D' {
						parts = append(parts, fmt.Sprintf("-%d", 1))
					}
				}
				status = strings.Join(parts, " ")
			}
		}

		repos = append(repos, gitRepo{
			Path:   p,
			Name:   filepath.Base(expanded),
			Branch: strings.TrimSpace(branch),
			Status: status,
			IsRepo: true,
		})
	}

	return gitReposLoadedMsg{repos: repos}
}

func (s *GitSection) Update(msg tea.Msg) (Section, tea.Cmd) {
	switch msg := msg.(type) {
	case gitReposLoadedMsg:
		s.loading = false
		s.repos = msg.repos
	case tea.KeyPressMsg:
		switch msg.String() {
		case "j", "down":
			if s.cursor < len(s.repos)-1 {
				s.cursor++
			}
		case "k", "up":
			if s.cursor > 0 {
				s.cursor--
			}
		case "enter":
			if len(s.repos) > 0 && s.repos[s.cursor].IsRepo {
				s.chosen = s.repos[s.cursor].Path
			}
		case "r":
			s.loading = true
			return s, s.fetchRepos
		}
	}
	return s, nil
}

func (s *GitSection) View(width, height int) string {
	const minWidth = 24
	if width < minWidth {
		return lipgloss.NewStyle().Faint(true).Width(width).Height(height).Render("  Git")
	}

	chrome := 4
	available := height - chrome
	if available < 1 {
		available = 5
	}

	var b strings.Builder

	titleStyle := NewStyle(width, width, 1, 1, 15, false, []int{0, 0, 0, 0})
	b.WriteString(titleStyle.Render(s.config.Title))
	b.WriteString("\n\n")

	if s.loading {
		b.WriteString(lipgloss.NewStyle().Faint(true).Render("  Loading..."))
		return NewStyleBorder(width, width, height, height, 15, false, []int{0, 0, 0, 1}).
			Render(b.String())
	}

	if len(s.repos) == 0 {
		b.WriteString(lipgloss.NewStyle().Faint(true).Render("  No repos configured"))
		return NewStyleBorder(width, width, height, height, 15, false, []int{0, 0, 0, 1}).
			Render(b.String())
	}

	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(2)).Bold(true)
	branchStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15)).Bold(true)
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	end := min(s.cursor+1, len(s.repos))
	start := max(end-available, 0)
	if len(s.repos) <= available {
		start = 0
		end = len(s.repos)
	}

	for i := start; i < end; i++ {
		repo := s.repos[i]

		var prefix string
		if i == s.cursor {
			prefix = cursorStyle.Render("▸ ")
		} else {
			prefix = "  "
		}

		if !repo.IsRepo {
			b.WriteString(fmt.Sprintf("%s%s\n", prefix, errorStyle.Render(repo.Name+" (not a git repo)")))
			continue
		}

		branch := ""
		if repo.Branch != "" {
			branch = branchStyle.Render(" ["+repo.Branch+"]")
		}

		status := ""
		if repo.Status != "" {
			status = statusStyle.Render(" "+repo.Status)
		}

		b.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, nameStyle.Render(repo.Name), branch, status))
	}

	return NewStyleBorder(width, width, height, height, 15, false, []int{0, 0, 0, 1}).
		Render(b.String())
}
