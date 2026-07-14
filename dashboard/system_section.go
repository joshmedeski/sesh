package dashboard

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type systemMetricsMsg struct {
	cpuPercent float64
	memPercent float64
	memTotal   uint64
}

type SystemSection struct {
	config     model.DashboardSectionConfig
	cpuUsage   float64
	memUsage   float64
	memTotal   float64
	lastUpdate time.Time
}

func NewSystemSection(cfg model.DashboardSectionConfig, deps SectionDeps) Section {
	return &SystemSection{
		config: cfg,
	}
}

func (s *SystemSection) Name() string    { return s.config.Title }
func (s *SystemSection) Chosen() string  { return "" }
func (s *SystemSection) TotalItems() int { return 0 }
func (s *SystemSection) Width() float64  { return s.config.Width }

func fetchSystemMetrics() tea.Msg {
	// Get virtual memory usage
	vMem, err := mem.VirtualMemory()
	if err != nil {
		return systemMetricsMsg{}
	}

	cpuPercent, err := cpu.Percent(0, false) // Using false aggregates all cores into one overall percentage
	var cpuUsage float64
	if err == nil && len(cpuPercent) > 0 {
		cpuUsage = cpuPercent[0]
	}

	return systemMetricsMsg{
		cpuPercent: cpuUsage,
		memPercent: vMem.UsedPercent, // Raw percent (e.g. 45.12345...)
		memTotal:   vMem.Total,       // Raw bytes (e.g. 17179869184)
	}
}

func systemTick() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return fetchSystemMetrics()
	})
}

func (s *SystemSection) Init() tea.Cmd {
	return fetchSystemMetrics
}

func (s *SystemSection) Update(msg tea.Msg) (Section, tea.Cmd) {
	switch msg := msg.(type) {
	case systemMetricsMsg:
		s.cpuUsage = msg.cpuPercent
		s.memUsage = msg.memPercent
		s.memTotal = float64(msg.memTotal)
		s.lastUpdate = time.Now()

		return s, systemTick()
	}
	return s, nil
}

func (s *SystemSection) View(width, height int) string {
	b := strings.Builder{}

	labelStyle := lipgloss.NewStyle().Bold(true).Width(6)
	// valueStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15))

	// RAM
	totalGB := float64(s.memTotal) / (1024 * 1024 * 1024)
	usedGB := totalGB * (s.memUsage / 100.0)

	cpuRow := fmt.Sprintf("%s %.1f%%\n", labelStyle.Render("CPU:"), s.cpuUsage)

	// Displays clean readings like: "RAM:  45.2% (7.2 / 16.0 GB)"
	ramRow := fmt.Sprintf("%s %.1f%% (%.1f / %.1f GB)\n",
		labelStyle.Render("RAM:"),
		s.memUsage,
		usedGB,
		totalGB,
	)

	b.WriteString(cpuRow)
	b.WriteString(ramRow)

	return b.String()
}
