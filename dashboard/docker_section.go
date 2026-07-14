package dashboard

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/joshmedeski/sesh/v2/model"
)

type dockerContainersLoadedMsg struct {
	containers []dockerContainer
}

type dockerContainer struct {
	ID     string
	Name   string
	Image  string
	Status string
	State  string
}

type DockerSection struct {
	config     model.DashboardSectionConfig
	deps       SectionDeps
	containers []dockerContainer
	cursor     int
	chosen     string
	loading    bool
}

func NewDockerSection(cfg model.DashboardSectionConfig, deps SectionDeps) Section {
	return &DockerSection{
		config:  cfg,
		deps:    deps,
		loading: true,
	}
}

func (s *DockerSection) Name() string    { return s.config.Title }
func (s *DockerSection) TotalItems() int { return len(s.containers) }
func (s *DockerSection) Width() float64  { return s.config.Width }
func (s *DockerSection) Chosen() string  { return s.chosen }

func (s *DockerSection) Init() tea.Cmd {
	return s.fetchContainers
}

func (s *DockerSection) fetchContainers() tea.Msg {
	args := []string{"ps", "--format", "{{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}\t{{.State}}"}
	if s.config.Docker.All {
		args = append(args, "-a")
	}
	for _, f := range s.config.Docker.Filters {
		args = append(args, "--filter", f)
	}

	out, err := runCommand("docker", args...)
	if err != nil {
		return dockerContainersLoadedMsg{containers: nil}
	}

	var containers []dockerContainer
	for line := range strings.SplitSeq(strings.TrimSpace(out), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 5)
		if len(parts) < 5 {
			continue
		}
		containers = append(containers, dockerContainer{
			ID:     parts[0],
			Name:   parts[1],
			Image:  parts[2],
			Status: parts[3],
			State:  parts[4],
		})
	}

	return dockerContainersLoadedMsg{containers: containers}
}

func (s *DockerSection) Update(msg tea.Msg) (Section, tea.Cmd) {
	switch msg := msg.(type) {
	case dockerContainersLoadedMsg:
		s.loading = false
		s.containers = msg.containers
	case tea.KeyPressMsg:
		switch msg.String() {
		case "j", "down":
			if s.cursor < len(s.containers)-1 {
				s.cursor++
			}
		case "k", "up":
			if s.cursor > 0 {
				s.cursor--
			}
		case "enter":
			if len(s.containers) > 0 {
				s.chosen = s.containers[s.cursor].Name
			}
		case "r":
			s.loading = true
			return s, s.fetchContainers
		}
	}
	return s, nil
}

func (s *DockerSection) View(width, height int) string {
	const minWidth = 24
	if width < minWidth {
		return lipgloss.NewStyle().Faint(true).Width(width).Height(height).Render("  Docker")
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

	if len(s.containers) == 0 {
		b.WriteString(lipgloss.NewStyle().Faint(true).Render("  No containers found"))
		return NewStyleBorder(width, width, height, height, 15, false, []int{0, 0, 0, 1}).
			Render(b.String())
	}

	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(2)).Bold(true)
	runningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	exitedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15)).Bold(true)
	statusStyle := lipgloss.NewStyle().Faint(true)

	end := min(s.cursor+1, len(s.containers))
	start := max(end-available, 0)
	if len(s.containers) <= available {
		start = 0
		end = len(s.containers)
	}

	for i := start; i < end; i++ {
		c := s.containers[i]

		var prefix string
		if i == s.cursor {
			prefix = cursorStyle.Render("▸ ")
		} else {
			prefix = "  "
		}

		var stateIndicator string
		if c.State == "running" {
			stateIndicator = runningStyle.Render("●")
		} else {
			stateIndicator = exitedStyle.Render("●")
		}

		nameWidth := max(width-30, 10)
		if nameWidth > 20 {
			nameWidth = 20
		}

		b.WriteString(fmt.Sprintf("%s %s %-20s %s\n",
			prefix,
			stateIndicator,
			nameStyle.Render(truncateString(c.Name, nameWidth)),
			statusStyle.Render(truncateString(c.Status, width-nameWidth-10)),
		))
	}

	return NewStyleBorder(width, width, height, height, 15, false, []int{0, 0, 0, 1}).
		Render(b.String())
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
