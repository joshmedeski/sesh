# `sesh status` Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a `sesh status` command that prints the current branch's GitHub issue (state badge + number + title) as a tmux-styled string for the status bar.

**Architecture:** A new `git.CurrentBranch` method feeds a new `github` package that wraps the `gh` CLI and returns a structured `Issue`. A thin `sesh status` cobra command resolves the session path, calls `github.Issue`, and formats a tmux-markup line. The `github` dependency is config-free and wired into `BaseDeps`.

**Tech Stack:** Go 1.25, cobra, testify mocks (mockery), `gh` CLI, tmux format markup.

## Global Constraints

- Module path: `github.com/joshmedeski/sesh/v2`.
- All external tools are wrapped behind interfaces that depend on `shell.Shell`; dependencies are wired in `seshcli/deps.go`.
- Mocks are generated, not committed (`.gitignore` ignores `mock_*`). Regenerate with `just mock` after any interface change.
- Run `just test` (regenerates mocks, then `go test -cover -race ./...`) before considering work complete.
- No-result behavior: `sesh status` prints nothing and exits `0` for every "nothing to show" case (no number in branch, not a git repo, `gh` missing/unauthenticated, issue not found, malformed output).
- Status output uses tmux format markup (`#[fg=…,bold]…#[default]`), never ANSI escape codes.

---

### Task 1: Add `CurrentBranch` to the `git` package

**Files:**
- Modify: `git/git.go` (interface + implementation)
- Test: `git/git_test.go` (create)
- Regenerate: `git/mock_Git.go` (via `just mock`)

**Interfaces:**
- Consumes: `shell.Shell.Cmd(cmd string, arg ...string) (string, error)` (existing).
- Produces: `git.Git.CurrentBranch(path string) (bool, string, error)` — `(true, branchName, nil)` on success; `(false, "", err)` when the directory is not a git repo or git fails.

- [ ] **Step 1: Write the failing test**

Create `git/git_test.go`:

```go
package git

import (
	"fmt"
	"testing"

	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestCurrentBranch(t *testing.T) {
	t.Run("returns the branch name on success", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		g := NewGit(mockShell)
		path := "/Users/josh/c/sesh"
		mockShell.On("Cmd", "git", "-C", path, "rev-parse", "--abbrev-ref", "HEAD").
			Return("400", nil)

		ok, branch, err := g.CurrentBranch(path)

		assert.True(t, ok)
		assert.Equal(t, "400", branch)
		assert.NoError(t, err)
	})

	t.Run("returns false when not a git repo", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		g := NewGit(mockShell)
		path := "/tmp/not-a-repo"
		mockShell.On("Cmd", "git", "-C", path, "rev-parse", "--abbrev-ref", "HEAD").
			Return("", fmt.Errorf("fatal: not a git repository"))

		ok, branch, err := g.CurrentBranch(path)

		assert.False(t, ok)
		assert.Equal(t, "", branch)
		assert.Error(t, err)
	})
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `just mock && go test ./git/... -run TestCurrentBranch -v`
Expected: compile error / FAIL — `g.CurrentBranch` undefined.

- [ ] **Step 3: Add the method to the interface and implementation**

In `git/git.go`, add to the `Git` interface (after `WorktreeList`):

```go
	WorktreeList(name string) (bool, string, error)
	CurrentBranch(path string) (bool, string, error)
```

And add the implementation at the end of the file:

```go
func (g *RealGit) CurrentBranch(path string) (bool, string, error) {
	out, err := g.shell.Cmd("git", "-C", path, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return false, "", err
	}
	return true, out, nil
}
```

- [ ] **Step 4: Regenerate mocks and run the test**

Run: `just mock && go test ./git/... -run TestCurrentBranch -v`
Expected: PASS (both subtests).

- [ ] **Step 5: Commit**

```bash
git add git/git.go git/git_test.go
git commit -m "feat(git): add CurrentBranch to read the branch for a path"
```

---

### Task 2: Create the `github` package

**Files:**
- Create: `github/github.go`
- Test: `github/github_test.go`
- Regenerate: `github/mock_Github.go` (via `just mock`)

**Interfaces:**
- Consumes: `git.Git.CurrentBranch(path string) (bool, string, error)` (Task 1); `shell.Shell.Cmd(cmd string, arg ...string) (string, error)`.
- Produces:
  - `github.Issue` struct with fields `Number int`, `Title string`, `State string` (json tags `number`/`title`/`state`).
  - `github.Github.Issue(path string) (Issue, bool, error)` — `(issue, true, nil)` when found; `(Issue{}, false, nil)` for every no-result case.
  - `github.NewGithub(shell shell.Shell, git git.Git) Github`.

- [ ] **Step 1: Write the failing tests**

Create `github/github_test.go`:

```go
package github

import (
	"fmt"
	"testing"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestParseIssueNumber(t *testing.T) {
	cases := []struct {
		branch string
		want   string
		ok     bool
	}{
		{"400", "400", true},
		{"feat/400-status-bar", "400", true},
		{"bugfix/issue-12", "12", true},
		{"feat/400-then-401", "400", true},
		{"main", "", false},
		{"", "", false},
	}
	for _, c := range cases {
		got, ok := parseIssueNumber(c.branch)
		assert.Equal(t, c.ok, ok, c.branch)
		assert.Equal(t, c.want, got, c.branch)
	}
}

func TestIssue(t *testing.T) {
	t.Run("returns the issue on success", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockGit := new(git.MockGit)
		gh := NewGithub(mockShell, mockGit)
		path := "/Users/josh/c/sesh"
		mockGit.On("CurrentBranch", path).Return(true, "feat/400-status-bar", nil)
		mockShell.On("Cmd", "gh", "issue", "view", "400", "--json", "number,title,state").
			Return(`{"number":400,"state":"OPEN","title":"Dynamic tmux status bar"}`, nil)

		issue, found, err := gh.Issue(path)

		assert.True(t, found)
		assert.NoError(t, err)
		assert.Equal(t, Issue{Number: 400, Title: "Dynamic tmux status bar", State: "OPEN"}, issue)
	})

	t.Run("not found when branch has no number", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockGit := new(git.MockGit)
		gh := NewGithub(mockShell, mockGit)
		path := "/Users/josh/c/sesh"
		mockGit.On("CurrentBranch", path).Return(true, "main", nil)

		issue, found, err := gh.Issue(path)

		assert.False(t, found)
		assert.NoError(t, err)
		assert.Equal(t, Issue{}, issue)
		// No gh call is set up on mockShell; testify panics on an
		// unexpected call, so reaching gh would fail this test.
	})

	t.Run("not found when not a git repo", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockGit := new(git.MockGit)
		gh := NewGithub(mockShell, mockGit)
		path := "/tmp/x"
		mockGit.On("CurrentBranch", path).Return(false, "", fmt.Errorf("not a git repo"))

		_, found, err := gh.Issue(path)

		assert.False(t, found)
		assert.NoError(t, err)
	})

	t.Run("not found when gh errors", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockGit := new(git.MockGit)
		gh := NewGithub(mockShell, mockGit)
		path := "/Users/josh/c/sesh"
		mockGit.On("CurrentBranch", path).Return(true, "400", nil)
		mockShell.On("Cmd", "gh", "issue", "view", "400", "--json", "number,title,state").
			Return("", fmt.Errorf("gh: not found"))

		_, found, err := gh.Issue(path)

		assert.False(t, found)
		assert.NoError(t, err)
	})

	t.Run("not found when gh returns malformed json", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockGit := new(git.MockGit)
		gh := NewGithub(mockShell, mockGit)
		path := "/Users/josh/c/sesh"
		mockGit.On("CurrentBranch", path).Return(true, "400", nil)
		mockShell.On("Cmd", "gh", "issue", "view", "400", "--json", "number,title,state").
			Return("not json", nil)

		_, found, err := gh.Issue(path)

		assert.False(t, found)
		assert.NoError(t, err)
	})
}
```

- [ ] **Step 2: Run the tests to verify they fail**

Run: `go test ./github/... -v`
Expected: compile error — package `github` does not exist yet.

- [ ] **Step 3: Write the implementation**

Create `github/github.go`:

```go
package github

import (
	"encoding/json"
	"regexp"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/shell"
)

// Issue is the subset of GitHub issue data sesh renders in the status bar.
type Issue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"` // "OPEN" | "CLOSED"
}

type Github interface {
	// Issue returns the GitHub issue for the branch checked out at path.
	// The bool is false (with a nil error) for every "nothing to show" case.
	Issue(path string) (Issue, bool, error)
}

type RealGithub struct {
	shell shell.Shell
	git   git.Git
}

func NewGithub(shell shell.Shell, git git.Git) Github {
	return &RealGithub{shell, git}
}

var issueNumberRe = regexp.MustCompile(`\d+`)

// parseIssueNumber returns the first run of digits in a branch name.
func parseIssueNumber(branch string) (string, bool) {
	match := issueNumberRe.FindString(branch)
	if match == "" {
		return "", false
	}
	return match, true
}

func (g *RealGithub) Issue(path string) (Issue, bool, error) {
	ok, branch, err := g.git.CurrentBranch(path)
	if err != nil || !ok {
		return Issue{}, false, nil
	}

	number, ok := parseIssueNumber(branch)
	if !ok {
		return Issue{}, false, nil
	}

	out, err := g.shell.Cmd("gh", "issue", "view", number, "--json", "number,title,state")
	if err != nil || out == "" {
		return Issue{}, false, nil
	}

	var issue Issue
	if err := json.Unmarshal([]byte(out), &issue); err != nil {
		return Issue{}, false, nil
	}
	return issue, true, nil
}
```

- [ ] **Step 4: Regenerate mocks and run the tests**

Run: `just mock && go test ./github/... -v`
Expected: PASS (all subtests, including `TestParseIssueNumber`).

- [ ] **Step 5: Commit**

```bash
git add github/github.go github/github_test.go
git commit -m "feat(github): add Issue lookup via the gh CLI"
```

---

### Task 3: Wire `github` into dependency injection

**Files:**
- Modify: `seshcli/deps.go`

**Interfaces:**
- Consumes: `github.NewGithub(shell.Shell, git.Git) github.Github` (Task 2).
- Produces: `BaseDeps.Github github.Github`, available on `*Deps` (which embeds `BaseDeps`) for command code.

> No new test — this is wiring covered by the build and by Task 4's command. The deliverable is a compiling dependency graph.

- [ ] **Step 1: Add the import**

In `seshcli/deps.go`, add to the import block (keep alphabetical, after `git`):

```go
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/home"
```

- [ ] **Step 2: Add the field to `BaseDeps`**

In the `BaseDeps` struct, add after `Git git.Git`:

```go
	Git        git.Git
	Github     github.Github
	Dir        dir.Dir
```

- [ ] **Step 3: Construct and assign it in `NewBaseDeps`**

In `NewBaseDeps`, after `g := git.NewGit(sh)`:

```go
	g := git.NewGit(sh)
	gh := github.NewGithub(sh, g)
```

And add it to the returned `&BaseDeps{...}` literal, after the `Git: g,` line:

```go
		Git:        g,
		Github:     gh,
		Dir:        d,
```

- [ ] **Step 4: Verify it compiles**

Run: `go build ./...`
Expected: builds with no errors.

- [ ] **Step 5: Commit**

```bash
git add seshcli/deps.go
git commit -m "feat(seshcli): wire github dependency into BaseDeps"
```

---

### Task 4: Add the `sesh status` command

**Files:**
- Create: `seshcli/status.go`
- Test: `seshcli/status_test.go` (create)
- Modify: `seshcli/root_command.go` (register the command)

**Interfaces:**
- Consumes: `deps.Github.Issue(path string) (github.Issue, bool, error)` (Task 2/3); `deps.Lister.GetAttachedTmuxSession() (model.SeshSession, bool)` (existing); `buildDeps(cmd, base) (*Deps, error)` (existing).
- Produces: `NewStatusCommand(base *BaseDeps) *cobra.Command`; `formatStatus(issue github.Issue) string`.

- [ ] **Step 1: Write the failing test for `formatStatus`**

Create `seshcli/status_test.go`:

```go
package seshcli

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/github"
	"github.com/stretchr/testify/assert"
)

func TestFormatStatus(t *testing.T) {
	t.Run("open issue gets a green OPEN badge", func(t *testing.T) {
		got := formatStatus(github.Issue{Number: 400, Title: "Dynamic tmux status bar", State: "OPEN"})
		assert.Equal(t, "#[fg=green,bold]OPEN#[default] #400 Dynamic tmux status bar", got)
	})

	t.Run("closed issue gets a red CLOSED badge", func(t *testing.T) {
		got := formatStatus(github.Issue{Number: 400, Title: "Dynamic tmux status bar", State: "CLOSED"})
		assert.Equal(t, "#[fg=red,bold]CLOSED#[default] #400 Dynamic tmux status bar", got)
	})

	t.Run("any non-OPEN state is treated as red", func(t *testing.T) {
		got := formatStatus(github.Issue{Number: 7, Title: "x", State: "MERGED"})
		assert.Equal(t, "#[fg=red,bold]MERGED#[default] #7 x", got)
	})
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./seshcli/... -run TestFormatStatus -v`
Expected: compile error — `formatStatus` undefined.

- [ ] **Step 3: Write the command and helpers**

Create `seshcli/status.go`:

```go
package seshcli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/github"
)

func NewStatusCommand(base *BaseDeps) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show contextual status for the current session (for the tmux status bar)",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			path := statusPath(deps)
			if path == "" {
				return nil
			}

			issue, found, _ := deps.Github.Issue(path)
			if !found {
				return nil
			}

			fmt.Print(formatStatus(issue))
			return nil
		},
	}
}

// statusPath resolves the directory to inspect: the attached tmux session's
// path when running inside tmux, otherwise the current working directory.
func statusPath(deps *Deps) string {
	if session, exists := deps.Lister.GetAttachedTmuxSession(); exists {
		return session.Path
	}
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return cwd
}

// formatStatus renders an issue as a tmux-styled status line.
func formatStatus(issue github.Issue) string {
	color := "green"
	if issue.State != "OPEN" {
		color = "red"
	}
	return fmt.Sprintf("#[fg=%s,bold]%s#[default] #%d %s", color, issue.State, issue.Number, issue.Title)
}
```

- [ ] **Step 4: Run the test to verify it passes**

Run: `go test ./seshcli/... -run TestFormatStatus -v`
Expected: PASS (all three subtests).

- [ ] **Step 5: Register the command**

In `seshcli/root_command.go`, add to the `rootCmd.AddCommand(...)` call (after `NewRootSessionCommand(base),`):

```go
		NewRootSessionCommand(base),
		NewStatusCommand(base),
		NewPreviewCommand(base),
```

- [ ] **Step 6: Verify it builds and the command is registered**

Run: `go build ./... && go run . status --help`
Expected: builds; help shows "Show contextual status for the current session (for the tmux status bar)".

- [ ] **Step 7: Commit**

```bash
git add seshcli/status.go seshcli/status_test.go seshcli/root_command.go
git commit -m "feat(seshcli): add sesh status command for the tmux status bar"
```

---

### Task 5: Document `sesh status` in the README

**Files:**
- Modify: `README.md` (add a subsection under `## Bonus`)

> No automated test; deliverable is accurate user-facing docs. Verify by reading.

- [ ] **Step 1: Add the documentation section**

In `README.md`, insert the following immediately before `### Connect to root` (after the `### Last` section's `bind` block, around line 463):

````markdown
### Dynamic status bar

`sesh status` prints a tmux-styled string describing the GitHub issue that
matches the current session's branch — a state badge (green `OPEN` / red
`CLOSED`) followed by the issue number and title. It is inspired by
[gitmux](https://github.com/arl/gitmux).

Add it to your `status-left` (or `status-right`):

```sh
set -g status-left "#[fg=blue,bold]#S #[fg=white,nobold]#(sesh status)"
```

The issue number is parsed from the branch name (`400` or `feat/400-status-bar`
both resolve to issue `#400`), so it works best with a branch-per-issue
workflow.

**Requirements:** the [`gh` CLI](https://cli.github.com) must be installed and
authenticated (`gh auth login`). When there is nothing to show — no number in
the branch, no matching issue, or `gh` is unavailable — `sesh status` prints
nothing, leaving the status bar clean.
````

- [ ] **Step 2: Verify the rendered section reads correctly**

Run: `grep -n "Dynamic status bar" README.md`
Expected: one match; the surrounding section reads as written.

- [ ] **Step 3: Commit**

```bash
git add README.md
git commit -m "docs: document sesh status for the tmux status bar"
```

---

### Final verification

- [ ] **Run the full suite**

Run: `just test`
Expected: mocks regenerate; all packages pass with `-race`, including `git`, `github`, and `seshcli`.

- [ ] **Manual smoke test (this repo is on branch `400`)**

Run: `go run . status`
Expected: prints something like `#[fg=green,bold]OPEN#[default] #400 Dynamic tmux status bar` (requires `gh` authenticated; otherwise prints nothing and exits 0).
