package dashboard

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/joshmedeski/sesh/v2/model"
)

type aiAgentsLoadedMsg struct {
	agents []aiAgent
}

type aiAgent struct {
	Name   string
	Status string
	PID    string
}

type AIAgentSection struct {
	config model.DashboardSectionConfig
	deps   SectionDeps
	agents []aiAgent
	cursor int
	chosen string
	loading bool
}

func NewAIAgentSection(cfg model.DashboardSectionConfig, deps SectionDeps) Section {
	return &AIAgentSection{
		config:  cfg,
		deps:    deps,
		loading: true,
	}
}

func (s *AIAgentSection) Name() string    { return s.config.Title }
func (s *AIAgentSection) TotalItems() int { return len(s.agents) }
func (s *AIAgentSection) Width() float64  { return s.config.Width }
func (s *AIAgentSection) Chosen() string  { return s.chosen }

func (s *AIAgentSection) Init() tea.Cmd {
	return s.fetchAgents
}

func (s *AIAgentSection) fetchAgents() tea.Msg {
	out, err := runCommand("sh", "-c",
		`ps aux | grep -E '(claude|copilot|aider|cursor|chatgpt|openai|ollama|llm)' | grep -v grep || true`)
	if err != nil {
		return aiAgentsLoadedMsg{agents: nil}
	}

	var agents []aiAgent
	seen := make(map[string]bool)

	for line := range strings.SplitSeq(strings.TrimSpace(out), "\n") {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 11 {
			continue
		}

		pid := fields[1]
		cmd := strings.Join(fields[10:], " ")

		name := "unknown"
		switch {
		case strings.Contains(cmd, "claude"):
			name = "Claude"
		case strings.Contains(cmd, "copilot"):
			name = "Copilot"
		case strings.Contains(cmd, "aider"):
			name = "Aider"
		case strings.Contains(cmd, "cursor"):
			name = "Cursor"
		case strings.Contains(cmd, "chatgpt"):
			name = "ChatGPT"
		case strings.Contains(cmd, "ollama"):
			name = "Ollama"
		default:
			name = "AI Agent"
		}

		if seen[name] {
			continue
		}
		seen[name] = true

		agents = append(agents, aiAgent{
			Name:   name,
			Status: "running",
			PID:    pid,
		})
	}

	return aiAgentsLoadedMsg{agents: agents}
}

func (s *AIAgentSection) Update(msg tea.Msg) (Section, tea.Cmd) {
	switch msg := msg.(type) {
	case aiAgentsLoadedMsg:
		s.loading = false
		s.agents = msg.agents
	case tea.KeyPressMsg:
		switch msg.String() {
		case "j", "down":
			if s.cursor < len(s.agents)-1 {
				s.cursor++
			}
		case "k", "up":
			if s.cursor > 0 {
				s.cursor--
			}
		case "enter":
			if len(s.agents) > 0 {
				s.chosen = s.agents[s.cursor].Name
			}
		case "r":
			s.loading = true
			return s, s.fetchAgents
		}
	}
	return s, nil
}

func (s *AIAgentSection) View(width, height int) string {
	const minWidth = 20
	if width < minWidth {
		return lipgloss.NewStyle().Faint(true).Width(width).Height(height).Render("  AI")
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
		b.WriteString(lipgloss.NewStyle().Faint(true).Render("  Scanning..."))
		return NewStyleBorder(width, width, height, height, 15, false, []int{0, 0, 0, 1}).
			Render(b.String())
	}

	if len(s.agents) == 0 {
		b.WriteString(lipgloss.NewStyle().Faint(true).Render("  No agents detected"))
		return NewStyleBorder(width, width, height, height, 15, false, []int{0, 0, 0, 1}).
			Render(b.String())
	}

	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(2)).Bold(true)
	runningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15)).Bold(true)
	pidStyle := lipgloss.NewStyle().Faint(true)

	end := min(s.cursor+1, len(s.agents))
	start := max(end-available, 0)
	if len(s.agents) <= available {
		start = 0
		end = len(s.agents)
	}

	for i := start; i < end; i++ {
		a := s.agents[i]

		var prefix string
		if i == s.cursor {
			prefix = cursorStyle.Render("▸ ")
		} else {
			prefix = "  "
		}

		b.WriteString(fmt.Sprintf("%s %s %-16s %s\n",
			prefix,
			runningStyle.Render("●"),
			nameStyle.Render(a.Name),
			pidStyle.Render("PID:"+a.PID),
		))
	}

	return NewStyleBorder(width, width, height, height, 15, false, []int{0, 0, 0, 1}).
		Render(b.String())
}
