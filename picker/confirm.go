package picker

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// zoxideSrc is the source name the frecency backend lists sessions under, and
// tmuxSrc the one live sessions are listed under. They are the only two sources
// ctrl+x owns: a config entry and a tmuxinator config describe a session that
// might exist, and deleting either means editing a file.
const (
	zoxideSrc = "zoxide"
	tmuxSrc   = "tmux"
)

// RemoveFunc drops a path from the frecency backend. Nil means removal is
// unavailable.
type RemoveFunc func(path string) error

// KillFunc kills a live tmux session by name. Nil means killing is unavailable.
type KillFunc func(name string) error

// confirmAction is which of the two things ctrl+x does on the highlighted row.
// It is decided by the row's source, and carried through the dialog and the
// result message so the wording and the dropped row agree with what ran.
type confirmAction int

const (
	actionRemove confirmAction = iota
	actionKill
)

// confirmState is the confirmation dialog: the row it was opened on, what would
// happen to it, and which of the two buttons is focused. It is nil whenever the
// dialog is closed, so its presence is the mode.
type confirmState struct {
	action confirmAction
	name   string
	path   string
	// yes is the focused button. The dialog opens on Yes, so enter alone
	// confirms — the row was already chosen by pressing ctrl+x on it.
	yes bool
}

// entryRemovedMsg carries the result back to Update. The row is dropped from
// the list only once this arrives without an error, so a backend that refuses
// the removal leaves the list telling the truth.
type entryRemovedMsg struct {
	action confirmAction
	name   string
	path   string
	err    error
}

// startRemove handles ctrl+x: it opens the confirmation dialog for the
// highlighted row, or explains why it can't. It is unavailable while the list
// is still loading — there is nothing highlighted to act on — and on any row
// that is neither a live tmux session nor a frecency entry.
func (m Model) startRemove() (tea.Model, tea.Cmd) {
	if m.loading {
		return m, nil
	}
	if m.cursor < 0 || m.cursor >= len(m.filtered) {
		return m, nil
	}
	item := m.filtered[m.cursor].item
	action, ok := m.actionFor(item.src)
	if !ok {
		m.status = "Only tmux sessions and zoxide entries can be removed"
		return m, nil
	}
	m.confirm = &confirmState{action: action, name: item.name, path: item.session.Path, yes: true}
	return m, nil
}

// actionFor maps a row's source to what ctrl+x does to it, reporting false when
// the source has no action or the backend behind it wasn't wired up.
func (m Model) actionFor(src string) (confirmAction, bool) {
	switch src {
	case zoxideSrc:
		return actionRemove, m.removeFunc != nil
	case tmuxSrc:
		return actionKill, m.killFunc != nil
	}
	return actionRemove, false
}

// updateConfirm handles a keypress while the dialog is open. It swallows every
// key: the filter input keeps its value but sees nothing typed, so the list is
// exactly as it was when the dialog closes either way.
func (m Model) updateConfirm(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		// ctrl+c means "get me out of here" everywhere else in the picker, and
		// a dialog is a poor place to start making it mean something narrower.
		m.confirm = nil
		m.quit = true
		return m, tea.Quit
	case "y", "Y":
		return m.applyConfirm(true)
	case "n", "N", "q", "esc":
		return m.applyConfirm(false)
	case "enter":
		return m.applyConfirm(m.confirm.yes)
	case "left", "right", "tab", "shift+tab", "h", "l":
		flipped := *m.confirm
		flipped.yes = !flipped.yes
		m.confirm = &flipped
		return m, nil
	}
	return m, nil
}

// applyConfirm closes the dialog, starting the action when it was confirmed.
func (m Model) applyConfirm(confirmed bool) (tea.Model, tea.Cmd) {
	target := *m.confirm
	m.confirm = nil
	if !confirmed {
		return m, nil
	}
	return m, m.removeEntry(target)
}

// removeEntry runs the backend off the event loop: it shells out, and blocking
// Update on it would freeze the picker for as long as the command takes.
func (m Model) removeEntry(target confirmState) tea.Cmd {
	removeFunc, killFunc := m.removeFunc, m.killFunc
	return func() tea.Msg {
		msg := entryRemovedMsg{action: target.action, name: target.name, path: target.path}
		if target.action == actionKill {
			msg.err = killFunc(target.name)
		} else {
			msg.err = removeFunc(target.path)
		}
		return msg
	}
}

// dropItem removes an entry from the loaded list and refilters, keeping the
// cursor where it was — pulled back by one when it sat on the last row, so it
// never points past the end.
//
// The row is matched on source as well as name: a zoxide directory and a live
// tmux session can share a name, and only the one the action ran against is
// gone.
func (m *Model) dropItem(action confirmAction, name, path string) {
	src, matchPath := zoxideSrc, true
	if action == actionKill {
		// A tmux session is identified by its name alone; its path is the
		// working directory it was created in, which says nothing about which
		// session was killed.
		src, matchPath = tmuxSrc, false
	}
	for i, item := range m.allItems {
		if item.src != src || item.name != name {
			continue
		}
		if matchPath && item.session.Path != path {
			continue
		}
		m.allItems = append(m.allItems[:i], m.allItems[i+1:]...)
		break
	}
	m.applyFilter()
	if m.cursor >= len(m.filtered) {
		m.cursor = max(len(m.filtered)-1, 0)
	}
	// The list lost a line, so the viewport is re-derived from the cursor in
	// row space rather than nudged: a group separator may have gone with it.
	m.scrollToCursor()
}

const (
	// confirmMinWidth keeps a short prompt from rendering as a cramped box,
	// and confirmMaxWidth keeps a long path from pushing the dialog wider than
	// the list it sits over.
	confirmMinWidth = 40
	confirmMaxWidth = 60
)

// confirmView renders the dialog itself, sized to fit within the frame.
func (m Model) confirmView() string {
	c := m.confirm

	inner := confirmMaxWidth
	if m.width > 0 && m.width-6 < inner {
		inner = m.width - 6
	}
	if inner < confirmMinWidth {
		inner = confirmMinWidth
	}

	center := lipgloss.NewStyle().Width(inner).MaxWidth(inner).Align(lipgloss.Center)
	faint := center.Faint(true)

	lines := []string{
		center.Render(confirmPrompt(c.action)),
		"",
		center.Bold(true).Render(truncate(c.name, inner)),
	}
	// The path is what actually gets removed from zoxide. It is only worth a
	// line of its own when the row is displayed under some other name, and
	// nothing to a kill, which targets the name.
	if c.action == actionRemove && c.path != "" && c.path != c.name {
		lines = append(lines, faint.Render(truncate(c.path, inner)))
	}
	lines = append(lines, "", center.Render(confirmButtons(c.yes)))

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.ANSIColor(4)).
		Padding(1, 2).
		Render(strings.Join(lines, "\n"))
}

// confirmPrompt is the question the dialog asks, which has to name what is
// about to happen: one forgets a directory, the other ends a running session.
func confirmPrompt(action confirmAction) string {
	if action == actionKill {
		return "Do you want to kill this tmux session?"
	}
	return "Do you want to remove this directory from zoxide?"
}

// confirmButtons renders the Yes/No pair with the focused one filled in.
func confirmButtons(yes bool) string {
	focused := lipgloss.NewStyle().Reverse(true).Bold(true)
	blurred := lipgloss.NewStyle().Faint(true)
	yesStyle, noStyle := blurred, focused
	if yes {
		yesStyle, noStyle = focused, blurred
	}
	return yesStyle.Render(" Yes ") + "   " + noStyle.Render(" No ")
}

// truncate cuts a string to width columns, marking the cut with an ellipsis so
// a clipped path doesn't read as a shorter one that also exists.
func truncate(s string, width int) string {
	if width < 2 || lipgloss.Width(s) <= width {
		return s
	}
	runes := []rune(s)
	for len(runes) > 0 && lipgloss.Width(string(runes))+1 > width {
		runes = runes[:len(runes)-1]
	}
	return string(runes) + "…"
}

// overlayCentered composes the dialog over the frame, centered. The frame is
// left visible around it: the dialog is about a row in the list, and hiding the
// list to ask about it would take away the context the question needs.
//
// A frame too small to hold the dialog — or one rendered before the terminal
// size is known — gets the dialog on its own instead of a clipped overlay.
func overlayCentered(base, dialog string, width, height int) string {
	dw, dh := lipgloss.Width(dialog), lipgloss.Height(dialog)
	if width < dw || height < dh {
		return dialog
	}
	// The layers go through a compositor rather than being composed onto the
	// canvas one at a time: a layer drawn directly fills the whole canvas, and
	// only the compositor honors its x/y offset.
	canvas := lipgloss.NewCanvas(width, height)
	canvas.Compose(lipgloss.NewCompositor(
		lipgloss.NewLayer(base),
		lipgloss.NewLayer(dialog).X((width-dw)/2).Y((height-dh)/2).Z(1),
	))
	return canvas.Render()
}

// removalFailed is the status shown when the backend refuses. The row stays in
// the list, so the message has to say why it is still there.
func removalFailed(action confirmAction, err error) string {
	if action == actionKill {
		return fmt.Sprintf("Couldn't kill session: %v", err)
	}
	return fmt.Sprintf("Couldn't remove entry: %v", err)
}
