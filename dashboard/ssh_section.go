package dashboard

import (
	"fmt"
	"os"
	"os/exec"

	tea "charm.land/bubbletea/v2"
	"github.com/joshmedeski/sesh/v2/model"
)

type SSHStatusMsg struct {
	index  int
	status string
	err    error
}

type SSHHostConfig struct {
	Name     string
	Host     string
	Port     int
	Username string
	Status   SSHStatusMsg
}

type SSHSection struct {
	config model.DashboardSectionConfig
	deps   SectionDeps
	hosts  []SSHHostConfig
	cursor int
	chosen string
}

func NewSSHSection(cfg model.DashboardSectionConfig, deps SectionDeps) Section {
	return &SSHSection{
		config: cfg,
	}
}

func (s *SSHSection) Name() string    { return s.config.Title }
func (s *SSHSection) Chosen() string  { return "" }
func (s *SSHSection) TotalItems() int { return 0 }
func (s *SSHSection) Width() float64  { return s.config.Width }

func (s *SSHSection) Init() tea.Cmd {
	return func() tea.Msg {
		for i, h := range s.hosts {
			cmd := checkSSHHost(i, h)
			return tea.Msg(cmd)
		}
		return nil
	}
}

func (s *SSHSection) Update(msg tea.Msg) (Section, tea.Cmd) {
	return s, nil
}

// TODO: checks only on launch of the dashboard. Need to check more often.
func checkSSHHost(index int, h SSHHostConfig) tea.Msg {
	user := h.Username
	if user == "" {
		user = os.Getenv("USER")
	}
	port := h.Port
	if port == 0 {
		port = 22
	}
	target := fmt.Sprintf("%s@%s", user, h.Host)
	cmd := exec.Command("ssh", "-p", fmt.Sprintf("%d", port),
		"-o", "ConnectTimeout=3",
		"-o", "BatchMode=yes",
		"-O", "check",
		target,
	)
	err := cmd.Run()
	status := "online"
	if err != nil {
		status = "offline"
	}
	return SSHStatusMsg{index: index, status: status}
}

func (s *SSHSection) View(width, height int) string {
	return ""
}
