package dashboard

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/joshmedeski/sesh/v2/model"
)

type sshStatusMsg struct {
	index  int
	status string
}

type SSHHost struct {
	Name     string
	Host     string
	Port     int
	Username string
	Status   string
}

type SSHSection struct {
	config  model.DashboardSectionConfig
	deps    SectionDeps
	hosts   []SSHHost
	cursor  int
	chosen  string
	loading bool
}

func NewSSHSection(cfg model.DashboardSectionConfig, deps SectionDeps) Section {
	hosts := make([]SSHHost, len(cfg.SSH))
	for i, h := range cfg.SSH {
		hosts[i] = SSHHost{
			Name:     h.Name,
			Host:     h.Host,
			Port:     h.Port,
			Username: h.Username,
			Status:   "checking",
		}
	}
	return &SSHSection{
		config:  cfg,
		deps:    deps,
		hosts:   hosts,
		loading: len(hosts) > 0,
	}
}

func (s *SSHSection) Name() string    { return s.config.Title }
func (s *SSHSection) TotalItems() int { return len(s.hosts) }
func (s *SSHSection) Width() float64  { return s.config.Width }
func (s *SSHSection) Chosen() string  { return s.chosen }

func (s *SSHSection) Init() tea.Cmd {
	if len(s.hosts) == 0 {
		return nil
	}
	cmds := make([]tea.Cmd, len(s.hosts))
	for i, h := range s.hosts {
		idx := i
		host := h
		cmds[i] = s.checkHost(idx, host)
	}
	return tea.Batch(cmds...)
}

func (s *SSHSection) checkHost(index int, host SSHHost) tea.Cmd {
	return func() tea.Msg {
		user := host.Username
		if user == "" {
			user = os.Getenv("USER")
		}
		port := host.Port
		if port == 0 {
			port = 22
		}
		target := fmt.Sprintf("%s@%s", user, host.Host)
		_, err := runCommand("ssh",
			"-p", fmt.Sprintf("%d", port),
			"-o", "ConnectTimeout=3",
			"-o", "BatchMode=yes",
			"-O", "check",
			target,
		)
		status := "online"
		if err != nil {
			status = "offline"
		}
		return sshStatusMsg{index: index, status: status}
	}
}

func (s *SSHSection) Update(msg tea.Msg) (Section, tea.Cmd) {
	switch msg := msg.(type) {
	case sshStatusMsg:
		if msg.index >= 0 && msg.index < len(s.hosts) {
			s.hosts[msg.index].Status = msg.status
		}
	case tea.KeyPressMsg:
		switch msg.String() {
		case "j", "down":
			if s.cursor < len(s.hosts)-1 {
				s.cursor++
			}
		case "k", "up":
			if s.cursor > 0 {
				s.cursor--
			}
		case "enter":
			if len(s.hosts) > 0 {
				s.chosen = s.hosts[s.cursor].Host
			}
		case "r":
			cmds := make([]tea.Cmd, len(s.hosts))
			for i, h := range s.hosts {
				h.Status = "checking"
				idx := i
				host := h
				cmds[i] = s.checkHost(idx, host)
			}
			return s, tea.Batch(cmds...)
		}
	}
	return s, nil
}

func (s *SSHSection) View(width, height int, focused bool) string {
	var b strings.Builder

	const minWidth = 20
	if width < minWidth {
		return lipgloss.NewStyle().Faint(true).Width(width).Height(height).Render("  SSH")
	}

	if len(s.hosts) == 0 {
		return NewStyleBorder(width, width, height, height, 15, false, []int{0, 0, 0, 1}, focused).
			Render(s.config.Title + "\n\n  No hosts configured")
	}

	chrome := 4
	available := height - chrome
	if available < 1 {
		available = 5
	}

	titleStyle := NewStyle(width, width, 1, 1, 15, false, []int{0, 0, 0, 0})
	b.WriteString(titleStyle.Render(s.config.Title))
	b.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15)).Bold(true)
	onlineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	offlineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	checkingStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(8)).Faint(true)
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(2)).Bold(true)

	end := min(s.cursor+1, len(s.hosts))
	start := max(end-available, 0)
	if len(s.hosts) <= available {
		start = 0
		end = len(s.hosts)
	}

	for i := start; i < end; i++ {
		h := s.hosts[i]

		var prefix string
		if i == s.cursor {
			prefix = cursorStyle.Render("▸ ")
		} else {
			prefix = "  "
		}

		var statusRendered string
		switch h.Status {
		case "online":
			statusRendered = onlineStyle.Render("● online")
		case "offline":
			statusRendered = offlineStyle.Render("● offline")
		default:
			statusRendered = checkingStyle.Render("○ checking")
		}

		nameDisplay := h.Name
		if nameDisplay == "" {
			nameDisplay = h.Host
		}

		line := fmt.Sprintf("%s%-20s %s", prefix, labelStyle.Render(nameDisplay), statusRendered)
		b.WriteString(line)
		b.WriteString("\n")
	}

	return NewStyleBorder(width, width, height, height, 15, false, []int{0, 0, 0, 1}, focused).
		Render(b.String())
}
