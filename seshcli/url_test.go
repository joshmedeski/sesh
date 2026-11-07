package seshcli

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/joshmedeski/sesh/v2/browser"
	"github.com/joshmedeski/sesh/v2/execwrap"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/picker"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/worktree"
)

func TestPickURLsOpensEveryChosenURL(t *testing.T) {
	mockPicker := new(picker.MockPicker)
	mockBrowser := new(browser.MockBrowser)
	mockExec := new(execwrap.MockExec)
	mockExec.On("LookPath", "defuddle").Return("", errors.New("not found"))
	mockPicker.On("PickURLs", []string{"https://a.dev", "https://b.dev"}, mock.Anything).
		Return([]string{"https://a.dev", "https://b.dev"}, nil)
	mockBrowser.On("Open", "https://a.dev").Return(nil)
	mockBrowser.On("Open", "https://b.dev").Return(nil)

	deps := &Deps{BaseDeps: BaseDeps{Exec: mockExec}, Picker: mockPicker, Browser: mockBrowser}
	assert.NoError(t, pickURLs(deps, []string{"https://a.dev", "https://b.dev"}))
	mockBrowser.AssertExpectations(t)
}

func TestPickURLsSkipsThePickerWithNoURLs(t *testing.T) {
	mockPicker := new(picker.MockPicker)
	assert.NoError(t, pickURLs(&Deps{Picker: mockPicker}, nil))
	mockPicker.AssertNotCalled(t, "PickURLs", mock.Anything, mock.Anything)
}

func TestIssueRepo(t *testing.T) {
	t.Run("prefers the session's worktree repo", func(t *testing.T) {
		mockTmux := new(tmux.MockTmux)
		mockWorktree := new(worktree.MockWorktree)
		mockTmux.On("FocusedSessionPath").Return("/repo/.wk/12", nil)
		mockWorktree.On("RepoAt", "/repo/.wk/12").Return("joshmedeski/sesh", true)

		assert.Equal(t, "joshmedeski/sesh", issueRepo(&Deps{Tmux: mockTmux, Worktree: mockWorktree}))
	})

	t.Run("uses the pane's worktree repo", func(t *testing.T) {
		mockTmux := new(tmux.MockTmux)
		mockWorktree := new(worktree.MockWorktree)
		mockTmux.On("FocusedSessionPath").Return("/code", nil)
		mockWorktree.On("RepoAt", "/code").Return("", false)
		mockTmux.On("CurrentPanePath").Return("/repo", nil)
		mockWorktree.On("RepoAt", "/repo").Return("joshmedeski/sesh", true)

		assert.Equal(t, "joshmedeski/sesh", issueRepo(&Deps{Tmux: mockTmux, Worktree: mockWorktree}))
	})

	t.Run("falls back to the pane's GitHub origin", func(t *testing.T) {
		mockTmux := new(tmux.MockTmux)
		mockWorktree := new(worktree.MockWorktree)
		mockGit := new(git.MockGit)
		mockTmux.On("FocusedSessionPath").Return("/code", nil)
		mockWorktree.On("RepoAt", "/code").Return("", false)
		mockTmux.On("CurrentPanePath").Return("/code/app", nil)
		mockWorktree.On("RepoAt", "/code/app").Return("", false)
		mockGit.On("OriginURL", "/code/app").Return("git@github.com:org/app.git", nil)

		deps := &Deps{BaseDeps: BaseDeps{Git: mockGit}, Tmux: mockTmux, Worktree: mockWorktree}
		assert.Equal(t, "org/app", issueRepo(deps))
	})
}

func TestURLPreview(t *testing.T) {
	newDeps := func(command string) (*Deps, *github.MockGithub, *shell.MockShell) {
		mockGithub := new(github.MockGithub)
		mockShell := new(shell.MockShell)
		mockExec := new(execwrap.MockExec)
		mockExec.On("LookPath", "defuddle").Return("", errors.New("not found"))
		deps := &Deps{
			BaseDeps: BaseDeps{Github: mockGithub, Shell: mockShell, Exec: mockExec},
			Config:   model.Config{URL: model.URLConfig{PreviewCommand: command}},
		}
		return deps, mockGithub, mockShell
	}

	t.Run("renders GitHub issues", func(t *testing.T) {
		deps, mockGithub, _ := newDeps("")
		mockGithub.On("IssueDetail", "org/app", 7).
			Return(github.IssueDetail{Title: "Bug", State: "OPEN", Body: "details"}, nil)
		out, err := urlPreview(deps)("https://github.com/org/app/issues/7")
		assert.NoError(t, err)
		assert.Equal(t, "# Bug\n\n**OPEN**\n\ndetails", out)
	})

	t.Run("runs the preview command", func(t *testing.T) {
		deps, _, mockShell := newDeps("fetch {url}")
		mockShell.On("ShellCmd", "fetch {url}", map[string]string{"{url}": "https://a.dev"}).Return("# A", nil)
		out, _ := urlPreview(deps)("https://a.dev")
		assert.Equal(t, "# A", out)
	})

	stubFetch := func(t *testing.T, out string, err error) {
		original := fetchMarkdown
		fetchMarkdown = func(string) (string, error) { return out, err }
		t.Cleanup(func() { fetchMarkdown = original })
	}

	t.Run("fetches the page itself with no command", func(t *testing.T) {
		stubFetch(t, "# A page", nil)
		deps, _, _ := newDeps("")
		out, _ := urlPreview(deps)("https://a.dev")
		assert.Equal(t, "# A page", out)
	})

	t.Run("falls back to the built-in fetch when the command fails", func(t *testing.T) {
		stubFetch(t, "# A page", nil)
		deps, _, mockShell := newDeps("fetch {url}")
		mockShell.On("ShellCmd", mock.Anything, mock.Anything).Return("", errors.New("no content"))
		out, _ := urlPreview(deps)("https://a.dev")
		assert.Equal(t, "# A page", out)
	})

	t.Run("shows the bare URL when the command and the built-in fetch fail", func(t *testing.T) {
		stubFetch(t, "", errors.New("offline"))
		deps, _, mockShell := newDeps("fetch {url}")
		mockShell.On("ShellCmd", mock.Anything, mock.Anything).Return("", errors.New("no content"))
		out, err := urlPreview(deps)("https://a.dev")
		assert.NoError(t, err)
		assert.Equal(t, "https://a.dev", out)
	})

	t.Run("shows the bare URL when the built-in fetch fails", func(t *testing.T) {
		stubFetch(t, "", errors.New("offline"))
		deps, _, _ := newDeps("")
		out, _ := urlPreview(deps)("https://a.dev")
		assert.Equal(t, "https://a.dev", out)
	})
}

func TestAppendUnique(t *testing.T) {
	assert.Equal(t, []string{"a", "b", "c"}, appendUnique([]string{"a", "b"}, []string{"b", "c", "c"}))
}
