package dashboard

import (
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/joshmedeski/sesh/v2/model"
)

// ListState holds the shared list-navigation state used by the sessions and
// configured sections: cursor/offset/view-height, the type-to-filter state
// (filtering, filterQuery, filtered), and the cursor movement/clamping logic.
// Sections embed it and supply their own master list and match predicate.
type ListState struct {
	cursor      int
	offset      int
	viewHeight  int
	filtering   bool
	filterQuery string
	filtered    []model.SeshSession // filtered view when filtering
}

// visible returns the currently displayed list (the filtered view while
// filtering, the full list otherwise).
func (ls *ListState) visible(all []model.SeshSession) []model.SeshSession {
	if ls.filtering && ls.filtered != nil {
		return ls.filtered
	}
	return all
}

// applyFilter rebuilds the filtered view from the master list and clamps the
// cursor. match reports whether a session matches the lower-cased query.
func (ls *ListState) applyFilter(all []model.SeshSession, match func(sess model.SeshSession, q string) bool) {
	if !ls.filtering || ls.filterQuery == "" {
		ls.filtered = nil
		ls.clampCursor(len(all))
		return
	}
	q := strings.ToLower(ls.filterQuery)
	out := make([]model.SeshSession, 0, len(all))
	for _, sess := range all {
		if match(sess, q) {
			out = append(out, sess)
		}
	}
	ls.filtered = out
	ls.clampCursor(len(out))
}

// clampCursor clamps the cursor and offset to a list of n items.
func (ls *ListState) clampCursor(n int) {
	if ls.cursor >= n {
		ls.cursor = max(n-1, 0)
	}
	if ls.offset >= n {
		ls.offset = 0
	}
}

func (ls *ListState) cursorUp(n int) {
	ls.cursor -= n
	if ls.cursor < 0 {
		ls.cursor = 0
	}
	if ls.cursor < ls.offset {
		ls.offset = ls.cursor
	}
}

func (ls *ListState) cursorDown(n, count int) {
	ls.cursor += n
	maxIdx := max(count-1, 0)
	if ls.cursor > maxIdx {
		ls.cursor = maxIdx
	}
	visible := ls.visibleCount()
	if ls.cursor >= ls.offset+visible {
		ls.offset = ls.cursor - visible + 1
	}
}

// visibleCount returns the number of rows that fit in the view, defaulting to
// 20 when the view height is unknown.
func (ls *ListState) visibleCount() int {
	if ls.viewHeight <= 0 {
		return 20
	}
	return max(ls.viewHeight, 1)
}

// ClickAt moves the cursor to the clicked view row, scrolling to reveal it.
func (ls *ListState) ClickAt(row, count int) {
	if count == 0 {
		return
	}
	ls.cursor = min(max(ls.offset+row, 0), count-1)
	if ls.cursor < ls.offset {
		ls.offset = ls.cursor
	}
	if visible := ls.visibleCount(); ls.cursor >= ls.offset+visible {
		ls.offset = ls.cursor - visible + 1
	}
}

// handleFilterKey consumes keys while type-to-filter is active: printable
// characters append to the query, backspace (and its ctrl+h / ctrl+backspace
// aliases) delete the last rune, j/k and the arrow keys move the cursor
// through the filtered results, enter selects the highlighted filtered item
// and exits filtering, and esc cancels filtering without selecting.
func (ls *ListState) handleFilterKey(msg tea.KeyPressMsg, all []model.SeshSession, match func(sess model.SeshSession, q string) bool, selectItem func()) {
	if isBackspaceKey(msg) {
		if ls.filterQuery != "" {
			r := []rune(ls.filterQuery)
			ls.filterQuery = string(r[:len(r)-1])
		}
		ls.applyFilter(all, match)
		return
	}
	switch msg.String() {
	case "esc":
		ls.filtering = false
		ls.filterQuery = ""
		ls.applyFilter(all, match)
	case "enter":
		ls.applyFilter(all, match)
		selectItem()
		ls.filtering = false
		ls.filterQuery = ""
		ls.applyFilter(all, match)
	case "j", "down":
		ls.cursorDown(1, len(ls.visible(all)))
	case "k", "up":
		ls.cursorUp(1)
	default:
		if msg.Text != "" {
			ls.filterQuery += msg.Text
			ls.applyFilter(all, match)
		}
	}
}
