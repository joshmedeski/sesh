package dashboard

import (
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/joshmedeski/sesh/v2/dashboard/render"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/model"
)

// branchLoadedMsg carries the current git branch for a session path, fetched
// asynchronously by every list section.
type branchLoadedMsg struct {
	path   string
	branch string
}

// statusLoadedMsg carries the formatted git status for a session path, fetched
// asynchronously by every list section.
type statusLoadedMsg struct {
	path   string
	status string
}

func uniquePaths(sessions []model.SeshSession) map[string]bool {
	paths := make(map[string]bool)
	for _, sess := range sessions {
		if sess.Path != "" {
			paths[sess.Path] = true
		}
	}
	return paths
}

func fetchBranches(g git.Git, sessions []model.SeshSession) tea.Cmd {
	paths := uniquePaths(sessions)
	cmds := make([]tea.Cmd, 0, len(paths))
	for path := range paths {
		cmds = append(cmds, func() tea.Msg {
			found, branch, err := g.CurrentBranch(path)
			if err != nil || !found {
				return branchLoadedMsg{path: path, branch: ""}
			}
			return branchLoadedMsg{path: path, branch: strings.TrimSpace(branch)}
		})
	}
	return tea.Batch(cmds...)
}

func fetchStatuses(g git.Git, sessions []model.SeshSession) tea.Cmd {
	paths := uniquePaths(sessions)
	cmds := make([]tea.Cmd, 0, len(paths))
	for path := range paths {
		cmds = append(cmds, func() tea.Msg {
			status, err := g.StatusSummary(path)
			if err != nil {
				return statusLoadedMsg{path: path, status: ""}
			}
			return statusLoadedMsg{path: path, status: render.FormatGitStatus(status)}
		})
	}
	return tea.Batch(cmds...)
}

func applyBranch(sessions []model.SeshSession, path, branch string) {
	for i := range sessions {
		if sessions[i].Path == path {
			sessions[i].Branch = branch
		}
	}
}

func applyStatus(sessions []model.SeshSession, path, status string) {
	for i := range sessions {
		if sessions[i].Path == path {
			sessions[i].GitStatus = status
		}
	}
}
