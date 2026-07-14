# Enrich tmux session name with issue title — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the tmux status-bar feature with a `sesh rename --enrich` command that folds a branch's GitHub issue title into the tmux session name (`base — issue title`), triggered by an opt-in tmux hook.

**Architecture:** A session is created fast with the namer name (unchanged). A tmux `session-created` hook runs `sesh rename --enrich` in the background; it recomputes the base name from the session path, looks up the issue via `gh`, and renames the live session. Reconnection stays correct because `dirStrategy`/`zoxideStrategy` match an existing session by base name (exact or `base — …` prefix) and reattach instead of creating a duplicate.

**Tech Stack:** Go 1.25, cobra (CLI), testify + mockery (`testify` template) for mocks, `just` task runner.

## Global Constraints

- Module path: `github.com/joshmedeski/sesh/v2`.
- Session-name separator is the string `" — "` (space, em dash U+2014, space). Define it **once** as `model.SessionNameSeparator` and reference it everywhere.
- tmux session names may not contain `.` or `:`; titles must have those two characters replaced with a space before use.
- Mocks are generated, never hand-edited. After changing any interface, run `just mock` (which runs `mockery` with `all: true`, recursive).
- Tests use `testify` (`assert`, `mock`); mock expectations use either `.On(...)` or `.EXPECT()` styles already present in the package — match the neighboring file.
- Run the full suite with `just test` (regenerates mocks, then `go test -cover -race ./...`). Run a single package with `go test ./<pkg>/...`.
- Follow existing error handling and `slog` conventions; a "nothing to do" outcome is a nil error, not a failure.

---

### Task 1: `SessionNameSeparator` constant + `SanitizeTitle` helper

**Files:**
- Create: `model/session_name.go`
- Create: `namer/sanitize.go`
- Test: `namer/sanitize_test.go`

**Interfaces:**
- Produces: `model.SessionNameSeparator` (untyped string const `" — "`); `namer.SanitizeTitle(title string) string`.

- [ ] **Step 1: Create the shared separator constant**

Create `model/session_name.go`:

```go
package model

// SessionNameSeparator joins a session's base name and its enriched suffix
// (e.g. a GitHub issue title). The em dash is legal in tmux session names,
// unlike ':' and '.'.
const SessionNameSeparator = " — "
```

- [ ] **Step 2: Write the failing test for `SanitizeTitle`**

Create `namer/sanitize_test.go`:

```go
package namer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeTitle(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"plain title kept", "warm the status cache", "warm the status cache"},
		{"casing preserved", "Warm The Cache", "Warm The Cache"},
		{"colon replaced with space", "fix: crash on start", "fix crash on start"},
		{"dot replaced with space", "bump v2.0 release", "bump v2 0 release"},
		{"collapses resulting double spaces", "a:  b", "a b"},
		{"trims leading and trailing space", "  hello  ", "hello"},
		{"colon then space stays single space", "feat: thing", "feat thing"},
		{"empty stays empty", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, SanitizeTitle(c.input))
		})
	}
}
```

- [ ] **Step 3: Run the test to verify it fails**

Run: `go test ./namer/ -run TestSanitizeTitle -v`
Expected: FAIL — `undefined: SanitizeTitle`.

- [ ] **Step 4: Implement `SanitizeTitle`**

Create `namer/sanitize.go`:

```go
package namer

import "strings"

// SanitizeTitle prepares a GitHub issue title for use inside a tmux session
// name: it keeps spaces and original casing but replaces the two characters
// tmux forbids in session names ('.' and ':') with a space, then collapses
// the resulting runs of whitespace and trims.
func SanitizeTitle(title string) string {
	replaced := strings.NewReplacer(".", " ", ":", " ").Replace(title)
	return strings.Join(strings.Fields(replaced), " ")
}
```

- [ ] **Step 5: Run the test to verify it passes**

Run: `go test ./namer/ -run TestSanitizeTitle -v`
Expected: PASS (all sub-tests).

- [ ] **Step 6: Commit**

```bash
git add model/session_name.go namer/sanitize.go namer/sanitize_test.go
git commit -m "feat(namer): add SessionNameSeparator and SanitizeTitle"
```

---

### Task 2: `Tmux.RenameSession` primitive

**Files:**
- Modify: `tmux/tmux.go` (add method to `Tmux` interface + `RealTmux`)
- Create: `tmux/rename.go`
- Test: `tmux/rename_test.go`
- Regenerate: `tmux/mock_Tmux.go` (via `just mock`)

**Interfaces:**
- Produces: `Tmux.RenameSession(target, newName string) (string, error)`.

- [ ] **Step 1: Add `RenameSession` to the `Tmux` interface**

In `tmux/tmux.go`, add this line to the `Tmux` interface block (alongside the other method signatures, e.g. after `GetCurrentSession()`):

```go
	RenameSession(target string, newName string) (string, error)
```

- [ ] **Step 2: Write the failing test**

Create `tmux/rename_test.go`:

```go
package tmux

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestRenameSession(t *testing.T) {
	t.Run("calls tmux rename-session with target and new name", func(t *testing.T) {
		mockShell := &shell.MockShell{}
		tmux := &RealTmux{shell: mockShell, bin: "tmux"}
		mockShell.EXPECT().
			Cmd("tmux", "rename-session", "-t", "400-status", "400-status — warm the cache").
			Return("", nil)

		result, err := tmux.RenameSession("400-status", "400-status — warm the cache")

		assert.Nil(t, err)
		assert.Equal(t, "", result)
	})
}
```

- [ ] **Step 3: Run the test to verify it fails**

Run: `go test ./tmux/ -run TestRenameSession -v`
Expected: FAIL — `RealTmux` has no field/method `RenameSession` (compile error).

- [ ] **Step 4: Implement `RenameSession`**

Create `tmux/rename.go`:

```go
package tmux

// RenameSession renames the tmux session identified by target to newName.
// newName may contain spaces (it is passed as a single argv element, so no
// shell quoting is required).
func (t *RealTmux) RenameSession(target string, newName string) (string, error) {
	return t.shell.Cmd(t.bin, "rename-session", "-t", target, newName)
}
```

- [ ] **Step 5: Regenerate mocks**

Run: `just mock`
Expected: `tmux/mock_Tmux.go` now contains a `RenameSession` method; no errors.

- [ ] **Step 6: Run the tests to verify they pass**

Run: `go test ./tmux/ -run TestRenameSession -v`
Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add tmux/tmux.go tmux/rename.go tmux/rename_test.go tmux/mock_Tmux.go
git commit -m "feat(tmux): add RenameSession"
```

---

### Task 3: `Lister.FindTmuxSessionByBase`

**Files:**
- Modify: `lister/lister.go` (add method to `Lister` interface)
- Modify: `lister/tmux.go` (implement on `RealLister`)
- Modify: `lister/caching_lister.go` (delegate to inner)
- Test: `lister/tmux_test.go` (create if absent, else append)
- Regenerate: `lister/mock_Lister.go` (via `just mock`)

**Interfaces:**
- Consumes: `model.SessionNameSeparator` (Task 1).
- Produces: `Lister.FindTmuxSessionByBase(base string) (model.SeshSession, bool)` — returns the tmux session whose name equals `base`, or (failing that) the first whose name starts with `base + model.SessionNameSeparator`.

- [ ] **Step 1: Add `FindTmuxSessionByBase` to the `Lister` interface**

In `lister/lister.go`, add to the `Lister` interface block (after `FindTmuxSession`):

```go
	FindTmuxSessionByBase(base string) (model.SeshSession, bool)
```

- [ ] **Step 2: Write the failing test**

Create `lister/find_by_base_test.go`:

```go
package lister

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/stretchr/testify/assert"
)

func TestFindTmuxSessionByBase(t *testing.T) {
	newListerWith := func(names ...string) *RealLister {
		sessions := make([]*model.TmuxSession, 0, len(names))
		for _, n := range names {
			sessions = append(sessions, &model.TmuxSession{Name: n, Path: "/p/" + n})
		}
		mockTmux := new(tmux.MockTmux)
		mockTmux.On("ListSessions").Return(sessions, nil)
		return &RealLister{tmux: mockTmux}
	}

	t.Run("exact match wins", func(t *testing.T) {
		l := newListerWith("sesh", "other")
		got, ok := l.FindTmuxSessionByBase("sesh")
		assert.True(t, ok)
		assert.Equal(t, "sesh", got.Name)
	})

	t.Run("matches enriched prefix", func(t *testing.T) {
		l := newListerWith("400-status — warm the cache", "other")
		got, ok := l.FindTmuxSessionByBase("400-status")
		assert.True(t, ok)
		assert.Equal(t, "400-status — warm the cache", got.Name)
	})

	t.Run("prefers exact over enriched prefix", func(t *testing.T) {
		l := newListerWith("400-status — warm the cache", "400-status")
		got, ok := l.FindTmuxSessionByBase("400-status")
		assert.True(t, ok)
		assert.Equal(t, "400-status", got.Name)
	})

	t.Run("does not match a bare prefix without the separator", func(t *testing.T) {
		l := newListerWith("sesh-ui")
		_, ok := l.FindTmuxSessionByBase("sesh")
		assert.False(t, ok)
	})

	t.Run("matches base containing a slash", func(t *testing.T) {
		l := newListerWith("w/400 — warm the cache")
		got, ok := l.FindTmuxSessionByBase("w/400")
		assert.True(t, ok)
		assert.Equal(t, "w/400 — warm the cache", got.Name)
	})

	t.Run("no match returns false", func(t *testing.T) {
		l := newListerWith("alpha", "beta")
		_, ok := l.FindTmuxSessionByBase("gamma")
		assert.False(t, ok)
	})
}
```

- [ ] **Step 3: Run the test to verify it fails**

Run: `go test ./lister/ -run TestFindTmuxSessionByBase -v`
Expected: FAIL — `RealLister` has no method `FindTmuxSessionByBase` (compile error).

- [ ] **Step 4: Implement `FindTmuxSessionByBase` on `RealLister`**

Add to `lister/tmux.go` (it already imports `fmt` and `model`; add `strings` to the import block):

```go
// FindTmuxSessionByBase finds a live tmux session for a namer-produced base
// name. It prefers an exact name match; failing that it returns the first
// session whose name is the base followed by the enrichment separator
// (e.g. "base — issue title"). This keeps reconnection working after a
// session has been renamed to include its issue title.
func (l *RealLister) FindTmuxSessionByBase(base string) (model.SeshSession, bool) {
	sessions, err := listTmux(l)
	if err != nil {
		return model.SeshSession{}, false
	}
	if session, exists := sessions.Directory[tmuxKey(base)]; exists {
		return session, true
	}
	prefix := base + model.SessionNameSeparator
	for _, key := range sessions.OrderedIndex {
		session := sessions.Directory[key]
		if strings.HasPrefix(session.Name, prefix) {
			return session, true
		}
	}
	return model.SeshSession{}, false
}
```

Update the import block at the top of `lister/tmux.go` to:

```go
import (
	"fmt"
	"strings"

	"github.com/joshmedeski/sesh/v2/model"
)
```

- [ ] **Step 5: Delegate from `CachingLister`**

Add to `lister/caching_lister.go`, in the "Delegate all other Lister methods to inner" section (next to `FindTmuxSession`):

```go
func (cl *CachingLister) FindTmuxSessionByBase(base string) (model.SeshSession, bool) {
	return cl.inner.FindTmuxSessionByBase(base)
}
```

- [ ] **Step 6: Regenerate mocks**

Run: `just mock`
Expected: `lister/mock_Lister.go` now has `FindTmuxSessionByBase`; no errors.

- [ ] **Step 7: Run the tests to verify they pass**

Run: `go test ./lister/ -run TestFindTmuxSessionByBase -v`
Expected: PASS (all sub-tests).

- [ ] **Step 8: Commit**

```bash
git add lister/lister.go lister/tmux.go lister/caching_lister.go lister/find_by_base_test.go lister/mock_Lister.go
git commit -m "feat(lister): add FindTmuxSessionByBase for enriched-name matching"
```

---

### Task 4: `sesh rename --enrich` command

**Files:**
- Create: `seshcli/rename.go`
- Modify: `seshcli/root_command.go` (register the command)
- Test: `seshcli/rename_test.go`

**Interfaces:**
- Consumes: `model.SessionNameSeparator` (Task 1), `namer.SanitizeTitle` (Task 1), `Tmux.RenameSession` (Task 2), `Lister.GetAttachedTmuxSession`/`FindTmuxSession` (existing), `Namer.Name` (existing), `Github.Issue` (existing).
- Produces: `NewRenameCommand(base *BaseDeps) *cobra.Command`; internal helpers `renameTarget(deps *Deps, args []string) (model.SeshSession, bool)` and `enrichedName(deps *Deps, path string) string`.

- [ ] **Step 1: Write the failing tests**

Create `seshcli/rename_test.go`:

```go
package seshcli

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/namer"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/stretchr/testify/assert"
)

// tmuxRenamer embeds the generated tmux mock (so it satisfies the full
// tmux.Tmux interface) and overrides RenameSession to capture its arguments.
type tmuxRenamer struct {
	*tmux.MockTmux
	called     bool
	gotTarget  string
	gotNewName string
}

func newTmuxRenamer() *tmuxRenamer {
	return &tmuxRenamer{MockTmux: new(tmux.MockTmux)}
}

func (t *tmuxRenamer) RenameSession(target string, newName string) (string, error) {
	t.called = true
	t.gotTarget = target
	t.gotNewName = newName
	return "", nil
}

func TestEnrichedName(t *testing.T) {
	t.Run("appends sanitized issue title to base", func(t *testing.T) {
		mockNamer := new(namer.MockNamer)
		mockGithub := new(github.MockGithub)
		mockNamer.On("Name", "/p").Return("400-status", nil)
		mockGithub.On("Issue", "/p").
			Return(github.Issue{Number: 400, Title: "warm: the cache", State: "OPEN"}, true, nil)

		deps := &Deps{Namer: mockNamer}
		deps.Github = mockGithub

		assert.Equal(t, "400-status — warm the cache", enrichedName(deps, "/p"))
	})

	t.Run("returns bare base when no issue", func(t *testing.T) {
		mockNamer := new(namer.MockNamer)
		mockGithub := new(github.MockGithub)
		mockNamer.On("Name", "/p").Return("400-status", nil)
		mockGithub.On("Issue", "/p").Return(github.Issue{}, false, nil)

		deps := &Deps{Namer: mockNamer}
		deps.Github = mockGithub

		assert.Equal(t, "400-status", enrichedName(deps, "/p"))
	})

	t.Run("returns empty when namer fails", func(t *testing.T) {
		mockNamer := new(namer.MockNamer)
		mockNamer.On("Name", "/p").Return("", assert.AnError)

		deps := &Deps{Namer: mockNamer}

		assert.Equal(t, "", enrichedName(deps, "/p"))
	})
}

func TestRenameTarget(t *testing.T) {
	t.Run("uses the named session when an arg is given", func(t *testing.T) {
		mockLister := new(lister.MockLister)
		mockLister.On("FindTmuxSession", "foo").
			Return(model.SeshSession{Name: "foo", Path: "/p/foo"}, true)

		deps := &Deps{Lister: mockLister}
		got, ok := renameTarget(deps, []string{"foo"})

		assert.True(t, ok)
		assert.Equal(t, "foo", got.Name)
	})

	t.Run("falls back to the attached session", func(t *testing.T) {
		mockLister := new(lister.MockLister)
		mockLister.On("GetAttachedTmuxSession").
			Return(model.SeshSession{Name: "bar", Path: "/p/bar"}, true)

		deps := &Deps{Lister: mockLister}
		got, ok := renameTarget(deps, nil)

		assert.True(t, ok)
		assert.Equal(t, "bar", got.Name)
	})
}

func TestRunEnrich(t *testing.T) {
	t.Run("renames when the enriched name differs", func(t *testing.T) {
		mockLister := new(lister.MockLister)
		mockNamer := new(namer.MockNamer)
		mockGithub := new(github.MockGithub)
		mockTmux := newTmuxRenamer()
		mockLister.On("GetAttachedTmuxSession").
			Return(model.SeshSession{Name: "400-status", Path: "/p"}, true)
		mockNamer.On("Name", "/p").Return("400-status", nil)
		mockGithub.On("Issue", "/p").
			Return(github.Issue{Number: 400, Title: "warm the cache", State: "OPEN"}, true, nil)

		deps := &Deps{Lister: mockLister, Namer: mockNamer}
		deps.Github = mockGithub
		deps.Tmux = mockTmux

		err := runEnrich(deps, nil)

		assert.NoError(t, err)
		assert.Equal(t, "400-status", mockTmux.gotTarget)
		assert.Equal(t, "400-status — warm the cache", mockTmux.gotNewName)
	})

	t.Run("does not rename when the name is unchanged (idempotent)", func(t *testing.T) {
		mockLister := new(lister.MockLister)
		mockNamer := new(namer.MockNamer)
		mockGithub := new(github.MockGithub)
		mockTmux := newTmuxRenamer()
		mockLister.On("GetAttachedTmuxSession").
			Return(model.SeshSession{Name: "400-status — warm the cache", Path: "/p"}, true)
		mockNamer.On("Name", "/p").Return("400-status", nil)
		mockGithub.On("Issue", "/p").
			Return(github.Issue{Number: 400, Title: "warm the cache", State: "OPEN"}, true, nil)

		deps := &Deps{Lister: mockLister, Namer: mockNamer}
		deps.Github = mockGithub
		deps.Tmux = mockTmux

		err := runEnrich(deps, nil)

		assert.NoError(t, err)
		assert.False(t, mockTmux.called)
	})
}
```

- [ ] **Step 2: Run the tests to verify they fail**

Run: `go test ./seshcli/ -run 'TestEnrichedName|TestRenameTarget|TestRunEnrich' -v`
Expected: FAIL — `undefined: enrichedName`, `renameTarget`, `runEnrich` (compile error).

- [ ] **Step 3: Implement the command**

Create `seshcli/rename.go`:

```go
package seshcli

import (
	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/namer"
)

func NewRenameCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename [name]",
		Short: "Rename a tmux session, optionally enriching it with its GitHub issue title",
		Long: "Rename a tmux session. With --enrich, the session named by [name] " +
			"(or the attached session) is renamed to '<namer name> — <issue title>' " +
			"when its branch resolves to a GitHub issue, and back to the bare namer " +
			"name when it does not. Intended to be run from a tmux session-created hook.",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}
			enrich, _ := cmd.Flags().GetBool("enrich")
			if !enrich {
				return nil // only the enrich mode is implemented today
			}
			return runEnrich(deps, args)
		},
	}
	cmd.Flags().Bool("enrich", false, "Rename the session to include its GitHub issue title")
	return cmd
}

// runEnrich resolves the target session, computes its enriched name, and
// renames it when that differs from the current name. Every "nothing to do"
// path returns nil.
func runEnrich(deps *Deps, args []string) error {
	target, ok := renameTarget(deps, args)
	if !ok || target.Name == "" {
		return nil
	}
	newName := enrichedName(deps, target.Path)
	if newName == "" || newName == target.Name {
		return nil
	}
	_, err := deps.Tmux.RenameSession(target.Name, newName)
	return err
}

// renameTarget resolves which session to rename: the one named by the first
// arg, or the attached session when no arg is given.
func renameTarget(deps *Deps, args []string) (model.SeshSession, bool) {
	if len(args) > 0 && args[0] != "" {
		return deps.Lister.FindTmuxSession(args[0])
	}
	return deps.Lister.GetAttachedTmuxSession()
}

// enrichedName recomputes the base name from path (the deterministic source of
// truth, which makes re-runs idempotent) and appends the sanitized issue title
// when the branch resolves to an issue. With no issue it returns the bare base,
// which self-heals a stale suffix. Returns "" when the base cannot be computed.
func enrichedName(deps *Deps, path string) string {
	baseName, err := deps.Namer.Name(path)
	if err != nil || baseName == "" {
		return ""
	}
	issue, found, _ := deps.Github.Issue(path)
	if !found {
		return baseName
	}
	title := namer.SanitizeTitle(issue.Title)
	if title == "" {
		return baseName
	}
	return baseName + model.SessionNameSeparator + title
}
```

- [ ] **Step 4: Register the command**

In `seshcli/root_command.go`, add `NewRenameCommand(base),` to the `rootCmd.AddCommand(...)` list (leave `NewStatusCommand(base)` in place for now — Task 6 removes it):

```go
	rootCmd.AddCommand(
		NewListCommand(base),
		NewLastCommand(base),
		NewConnectCommand(base),
		NewCloneCommand(base),
		NewRootSessionCommand(base),
		NewStatusCommand(base),
		NewRenameCommand(base),
		NewPreviewCommand(base),
		NewPickerCommand(base),
		NewWindowCommand(base),
	)
```

- [ ] **Step 5: Run the tests to verify they pass**

Run: `go test ./seshcli/ -run 'TestEnrichedName|TestRenameTarget|TestRunEnrich' -v`
Expected: PASS.

- [ ] **Step 6: Build to confirm the command compiles and registers**

Run: `go build ./... && go run . rename --help`
Expected: build succeeds; help text for `rename` shows the `--enrich` flag.

- [ ] **Step 7: Commit**

```bash
git add seshcli/rename.go seshcli/rename_test.go seshcli/root_command.go
git commit -m "feat(seshcli): add sesh rename --enrich"
```

---

### Task 5: Reattach to enriched sessions in dir/zoxide strategies

**Files:**
- Modify: `connector/dir.go`
- Modify: `connector/zoxide.go`
- Test: `connector/dir_test.go` (create), `connector/zoxide_test.go` (create)

**Interfaces:**
- Consumes: `Lister.FindTmuxSessionByBase` (Task 3).
- Produces: no new exported symbols; changes the behavior of `dirStrategy`/`zoxideStrategy` to return `New: false` with the existing session when a base match is found.

- [ ] **Step 1: Write the failing test for `dirStrategy`**

Create `connector/dir_test.go`:

```go
package connector

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/dir"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/namer"
	"github.com/stretchr/testify/assert"
)

func TestDirStrategy(t *testing.T) {
	newConnector := func(l lister.Lister, n namer.Namer, h home.Home, d dir.Dir) *RealConnector {
		return &RealConnector{lister: l, namer: n, home: h, dir: d}
	}

	t.Run("reattaches to an enriched session instead of creating a new one", func(t *testing.T) {
		mockHome := new(home.MockHome)
		mockDir := new(dir.MockDir)
		mockNamer := new(namer.MockNamer)
		mockLister := new(lister.MockLister)
		mockHome.On("ExpandPath", "/p/repo").Return("/p/repo", nil)
		mockDir.On("Dir", "/p/repo").Return(true, "/p/repo")
		mockNamer.On("Name", "/p/repo").Return("repo", nil)
		mockLister.On("FindTmuxSessionByBase", "repo").
			Return(model.SeshSession{Name: "repo — fix the bug", Path: "/p/repo"}, true)

		c := newConnector(mockLister, mockNamer, mockHome, mockDir)
		conn, err := dirStrategy(c, "/p/repo")

		assert.NoError(t, err)
		assert.True(t, conn.Found)
		assert.False(t, conn.New)
		assert.Equal(t, "repo — fix the bug", conn.Session.Name)
	})

	t.Run("creates a new session when no base match exists", func(t *testing.T) {
		mockHome := new(home.MockHome)
		mockDir := new(dir.MockDir)
		mockNamer := new(namer.MockNamer)
		mockLister := new(lister.MockLister)
		mockHome.On("ExpandPath", "/p/repo").Return("/p/repo", nil)
		mockDir.On("Dir", "/p/repo").Return(true, "/p/repo")
		mockNamer.On("Name", "/p/repo").Return("repo", nil)
		mockLister.On("FindTmuxSessionByBase", "repo").Return(model.SeshSession{}, false)

		c := newConnector(mockLister, mockNamer, mockHome, mockDir)
		conn, err := dirStrategy(c, "/p/repo")

		assert.NoError(t, err)
		assert.True(t, conn.Found)
		assert.True(t, conn.New)
		assert.Equal(t, "repo", conn.Session.Name)
	})

	t.Run("not a directory returns not found", func(t *testing.T) {
		mockHome := new(home.MockHome)
		mockDir := new(dir.MockDir)
		mockHome.On("ExpandPath", "/p/x").Return("/p/x", nil)
		mockDir.On("Dir", "/p/x").Return(false, "")

		c := newConnector(nil, nil, mockHome, mockDir)
		conn, err := dirStrategy(c, "/p/x")

		assert.NoError(t, err)
		assert.False(t, conn.Found)
	})
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./connector/ -run TestDirStrategy -v`
Expected: FAIL — the "reattaches" case gets `New: true` (current behavior), and the mock `FindTmuxSessionByBase` is unexpected.

- [ ] **Step 3: Add base-match to `dirStrategy`**

In `connector/dir.go`, after computing `nameFromPath` and before the final `return`, insert the reattach check:

```go
	nameFromPath, err := c.namer.Name(absPath)
	if err != nil {
		return model.Connection{}, err
	}
	if existing, ok := c.lister.FindTmuxSessionByBase(nameFromPath); ok {
		return model.Connection{
			Found:       true,
			New:         false,
			AddToZoxide: true,
			Session:     existing,
		}, nil
	}
	return model.Connection{
		Found:       true,
		New:         true,
		AddToZoxide: true,
		Session: model.SeshSession{
			Src:  "dir",
			Name: nameFromPath,
			Path: absPath,
		},
	}, nil
```

- [ ] **Step 4: Write the failing test for `zoxideStrategy`**

Create `connector/zoxide_test.go`:

```go
package connector

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/namer"
	"github.com/stretchr/testify/assert"
)

func TestZoxideStrategy(t *testing.T) {
	t.Run("reattaches to an enriched session", func(t *testing.T) {
		mockLister := new(lister.MockLister)
		mockNamer := new(namer.MockNamer)
		mockLister.On("FindZoxideSession", "repo").
			Return(model.SeshSession{Name: "repo", Path: "/p/repo"}, true)
		mockNamer.On("Name", "/p/repo").Return("repo", nil)
		mockLister.On("FindTmuxSessionByBase", "repo").
			Return(model.SeshSession{Name: "repo — fix the bug", Path: "/p/repo"}, true)

		c := &RealConnector{lister: mockLister, namer: mockNamer}
		conn, err := zoxideStrategy(c, "repo")

		assert.NoError(t, err)
		assert.True(t, conn.Found)
		assert.False(t, conn.New)
		assert.Equal(t, "repo — fix the bug", conn.Session.Name)
	})

	t.Run("creates a new session when no base match", func(t *testing.T) {
		mockLister := new(lister.MockLister)
		mockNamer := new(namer.MockNamer)
		mockLister.On("FindZoxideSession", "repo").
			Return(model.SeshSession{Name: "repo", Path: "/p/repo"}, true)
		mockNamer.On("Name", "/p/repo").Return("repo", nil)
		mockLister.On("FindTmuxSessionByBase", "repo").Return(model.SeshSession{}, false)

		c := &RealConnector{lister: mockLister, namer: mockNamer}
		conn, err := zoxideStrategy(c, "repo")

		assert.NoError(t, err)
		assert.True(t, conn.Found)
		assert.True(t, conn.New)
		assert.Equal(t, "repo", conn.Session.Name)
	})

	t.Run("zoxide miss returns not found", func(t *testing.T) {
		mockLister := new(lister.MockLister)
		mockLister.On("FindZoxideSession", "nope").Return(model.SeshSession{}, false)

		c := &RealConnector{lister: mockLister}
		conn, err := zoxideStrategy(c, "nope")

		assert.NoError(t, err)
		assert.False(t, conn.Found)
	})
}
```

- [ ] **Step 5: Run the test to verify it fails**

Run: `go test ./connector/ -run TestZoxideStrategy -v`
Expected: FAIL — the "reattaches" case gets `New: true`; `FindTmuxSessionByBase` unexpected.

- [ ] **Step 6: Add base-match to `zoxideStrategy`**

Rewrite `zoxideStrategy` in `connector/zoxide.go` (keep `zoxideToTmuxName` unchanged):

```go
func zoxideStrategy(c *RealConnector, path string) (model.Connection, error) {
	session, exists := c.lister.FindZoxideSession(path)
	if !exists {
		return model.Connection{Found: false}, nil
	}
	name, err := c.namer.Name(session.Path)
	if err != nil {
		return model.Connection{}, err
	}
	if existing, ok := c.lister.FindTmuxSessionByBase(name); ok {
		return model.Connection{
			Found:       true,
			Session:     existing,
			New:         false,
			AddToZoxide: true,
		}, nil
	}
	session.Name = name
	return model.Connection{
		Found:       true,
		Session:     session,
		New:         true,
		AddToZoxide: true,
	}, nil
}
```

- [ ] **Step 7: Run both strategy tests to verify they pass**

Run: `go test ./connector/ -run 'TestDirStrategy|TestZoxideStrategy' -v`
Expected: PASS.

- [ ] **Step 8: Run the full connector package to check nothing regressed**

Run: `go test ./connector/...`
Expected: PASS (including the existing `TestEstablishTmuxConnection`).

- [ ] **Step 9: Commit**

```bash
git add connector/dir.go connector/zoxide.go connector/dir_test.go connector/zoxide_test.go
git commit -m "feat(connector): reattach to enriched sessions by base name"
```

---

### Task 6: Remove the status-bar feature

**Files:**
- Delete: `seshcli/status.go`, `seshcli/status_test.go`, `seshcli/connect_test.go`
- Delete: `statuscache/` (entire package), `refresher/` (entire package)
- Modify: `seshcli/connect.go` (remove `maybeWarmStatus` + its call)
- Modify: `seshcli/deps.go` (remove `StatusCache`/`Refresher` fields, construction, imports)
- Modify: `seshcli/root_command.go` (remove `NewStatusCommand`)
- Modify: `model/config.go` (remove `GithubConfig` type, `Github` field, `EffectiveTTL`)
- Modify: `sesh.schema.json` (remove the `github` property)

**Interfaces:**
- Consumes: nothing new.
- Produces: nothing new — this task only removes code. After it, `github.Issue`/`github.Resolve` remain (used by Task 4); `github.NewGithub` wiring in `deps.go` stays.

- [ ] **Step 1: Delete the status command and its now-orphaned tests**

```bash
git rm seshcli/status.go seshcli/status_test.go seshcli/connect_test.go
```

(`connect_test.go` contained only `TestMaybeWarmStatus`, removed next.)

- [ ] **Step 2: Remove `maybeWarmStatus` from `connect.go`**

In `seshcli/connect.go`, delete the entire `maybeWarmStatus` function (lines defining `// maybeWarmStatus ...` through its closing brace) and delete the call `maybeWarmStatus(deps)` in the `RunE` body. The `RunE` tail should read:

```go
			// Refresh cache in background so next sesh list has fresh data
			if deps.CachingLister != nil {
				deps.CachingLister.RefreshCache(lister.ListOptions{})
				deps.CachingLister.Wait()
			}
			return nil
```

Leave the `lister` and `model` imports (still used).

- [ ] **Step 3: Remove `NewStatusCommand` registration**

In `seshcli/root_command.go`, delete the `NewStatusCommand(base),` line from `rootCmd.AddCommand(...)`. `NewRenameCommand(base),` stays.

- [ ] **Step 4: Remove status cache + refresher wiring from `deps.go`**

In `seshcli/deps.go`:
- Delete the `StatusCache statuscache.StatusCache` and `Refresher refresher.Refresher` fields from `BaseDeps`.
- Delete `sc := statuscache.NewFileStatusCache()` and `rf := refresher.NewRefresher()`.
- Delete the `StatusCache: sc,` and `Refresher: rf,` entries in the returned `&BaseDeps{...}`.
- Remove the now-unused imports `"github.com/joshmedeski/sesh/v2/refresher"` and `"github.com/joshmedeski/sesh/v2/statuscache"`.

- [ ] **Step 5: Delete the packages**

```bash
git rm -r statuscache refresher
```

- [ ] **Step 6: Remove the dead GitHub config**

In `model/config.go`:
- Delete the `Github GithubConfig \`toml:"github"\`` field from the `Config` struct.
- Delete the entire `GithubConfig` type definition.
- Delete the `EffectiveTTL` method.

- [ ] **Step 7: Sync the JSON schema**

In `sesh.schema.json`, delete the `"github": { ... }` property block (the object with `issue_ttl`) from `properties`. Ensure the surrounding JSON stays valid (no trailing comma left dangling).

Run: `python3 -m json.tool sesh.schema.json > /dev/null && echo VALID`
Expected: `VALID`.

- [ ] **Step 8: Build to confirm nothing references the removed code**

Run: `go build ./...`
Expected: builds clean. If the compiler reports an unused import or missing symbol, remove that reference (it is dead status-bar code).

- [ ] **Step 9: Regenerate mocks (drops MockStatusCache / MockRefresher)**

Run: `just mock`
Expected: no `statuscache`/`refresher` mocks remain; no errors.

- [ ] **Step 10: Run the full suite**

Run: `just test`
Expected: PASS across all packages.

- [ ] **Step 11: Commit**

```bash
git add -A
git commit -m "refactor: remove tmux status-bar feature in favor of session rename"
```

---

### Task 7: Document the tmux hook

**Files:**
- Modify: `README.md` (or the docs page where tmux integration is described)

**Interfaces:** none.

- [ ] **Step 1: Find the right doc location**

Run: `grep -rn "status" README.md docs/ 2>/dev/null | grep -i "tmux\|status bar" | head`
Expected: locate any existing status-bar documentation to replace; if none, add a new "Enrich session names with issue titles" subsection near the tmux integration docs.

- [ ] **Step 2: Add the hook documentation**

Add a subsection with this content (adapt heading level to the surrounding doc):

````markdown
### Enrich session names with the GitHub issue title

`sesh` can rename a session to include its branch's GitHub issue title, e.g.
`400-status` → `400-status — warm the status cache`. It parses the issue number
from the branch name and looks it up with the `gh` CLI (which must be installed
and authenticated).

Add an opt-in tmux hook so every new session is enriched in the background:

```tmux
set-hook -g session-created 'run-shell -b "sesh rename --enrich"'
```

The command is a no-op when the branch has no resolvable issue (the session
keeps its plain name), and it self-heals: switching to a branch without an issue
renames the session back to its base name. Reconnecting to the directory
reattaches to the enriched session rather than creating a duplicate.
````

- [ ] **Step 3: Commit**

```bash
git add README.md
git commit -m "docs: document sesh rename --enrich tmux hook"
```

---

## Self-Review Notes

- **Spec coverage:** two-phase behavior (Tasks 4 + hook docs in 7); `sesh rename --enrich` with attached/arg resolution + recompute-base idempotency + force-to-base (Task 4); `SanitizeTitle` replacing `:`/`.` with space (Task 1); ` — ` separator via `model.SessionNameSeparator` (Task 1, used in 3/4); `Tmux.RenameSession` (Task 2); base-match reconnection in dir/zoxide (Tasks 3+5); removal of status.go/statuscache/refresher + deps wiring + config/schema (Task 6); hook docs (Task 7). All spec sections map to a task.
- **Type consistency:** `enrichedName`, `renameTarget`, `runEnrich`, `FindTmuxSessionByBase`, `RenameSession`, `SanitizeTitle`, `SessionNameSeparator` are used with identical signatures across the tasks that define and consume them.
- **Known caveat (from spec, out of scope):** two distinct repos producing the identical base name remain ambiguous; not addressed here.
