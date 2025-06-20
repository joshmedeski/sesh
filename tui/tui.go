package tui

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// TODO: move to model package
// TuiModel represents the state of the TUI application
type TuiModel struct {
	connector connector.Connector
	list      list.Model
}

type Tui interface {
	NewModel() TuiModel
}

type RealTui struct {
	connector connector.Connector
	lister    lister.Lister
}

func NewTui(connector connector.Connector, lister lister.Lister) Tui {
	return &RealTui{connector, lister}
}

func (t *RealTui) NewModel() TuiModel {
	// TODO: get items from lister

	sessions, err := t.lister.List(
		lister.ListOptions{
			Tmux:   true,
			Config: true,
			Zoxide: true,
		},
	)
	if err != nil {
		slog.Error("seshcli/tui.go: NewModel", "error", err)
		panic(err)
	}

	if len(sessions.Directory) == 0 {
		slog.Info("seshcli/tui.go: NewModel", "message", "No sessions found")
		return TuiModel{
			connector: t.connector,
			list:      list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		}
	}
	items := make([]list.Item, len(sessions.Directory))
	for i, session := range sessions.OrderedIndex {
		items[i] = item{
			title: sessions.Directory[session].Name,
			desc:  sessions.Directory[session].Description,
		}
	}

	m := TuiModel{
		connector: t.connector,
		list:      list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.list.Title = "Sesh"

	return m
}

func (m TuiModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m TuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			selectedItem := m.list.SelectedItem()
			name := selectedItem.(item).title
			if name == "" {
				slog.Error("seshcli/tui.go: Update", "error", "No session selected")
				return m, tea.Quit
			}
			m.connector.Connect(selectedItem.(item).title, model.ConnectOpts{})
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m TuiModel) View() string {
	return docStyle.Render(m.list.View())
}
