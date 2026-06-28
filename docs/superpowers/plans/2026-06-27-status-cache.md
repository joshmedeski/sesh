# `sesh status` issue cache Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `sesh status` read GitHub issue title/state from a local, branch-keyed disk cache so every tmux status render is a sub-millisecond file read, while a detached background process keeps it fresh via `gh`.

**Architecture:** A new `statuscache` package stores one gob file per `repoRoot+branch` key. A new `refresher` package spawns a detached `sesh status --refresh` child that performs the only live `gh` call and writes the cache. `sesh status` reads the cache and renders instantly, spawning a refresh only when the entry is missing or stale; `sesh connect` warms the cache on switch. Config gains a `[github]` section with `issue_ttl`.

**Tech Stack:** Go 1.25, cobra, testify mocks (mockery), gob, `gh` CLI, tmux format markup.

## Global Constraints

- Module path: `github.com/joshmedeski/sesh/v2`, Go 1.25.
- External tools are wrapped behind interfaces depending on `shell.Shell`; dependencies are wired in `seshcli/deps.go`. Config-free deps go on `BaseDeps`.
- Mocks are gitignored (`mock_*`) and generated via `just mock` — never committed. Regenerate after any interface change.
- Run `just test` (regenerates mocks, then `go test -cover -race ./...`) before considering work complete.
- Cache key is `sha256(repoRoot + "\x00" + branch)` hex-encoded — branch-keyed, never issue-number-keyed.
- `Entry` carries both `PR *Ref` and `Issue *Ref`; the MVP only ever fills `Issue`. The refresh **always** writes an entry (both refs nil = negative entry) so the hot path respects the TTL instead of re-spawning every tick.
- No-result behavior: `sesh status` prints nothing and exits `0` for every nothing-to-show case (not a repo, no number, gh missing/unauthenticated, issue not found, malformed JSON, negative cache entry).
- `issue_ttl` is a `*int`: `nil` → effective TTL 60; explicit `0` → cache disabled (live fetch, today's behavior); `N>0` → N seconds. Use `config.Github.EffectiveTTL()`; `== 0` means disabled.
- Status output uses tmux format markup (`#[fg=…]`), never ANSI. `formatStatus` is unchanged from PR #401.
- `statuscache` must not import `github` (the cache→issue conversion lives in `seshcli`).
- Schema (`sesh.schema.json`) must stay in lockstep with the config struct (see the `config-schema-sync` skill).

---

### Task 1: Config `[github]` section + `EffectiveTTL` + schema

**Files:**
- Modify: `model/config.go`
- Create: `model/config_test.go`
- Modify: `sesh.schema.json`

**Interfaces:**
- Produces: `model.GithubConfig` with `IssueTTL *int` (`toml:"issue_ttl"`); `Config.Github GithubConfig` (`toml:"github"`); method `(GithubConfig) EffectiveTTL() int` — returns 60 when `IssueTTL` is nil, else `*IssueTTL`.

- [ ] **Step 1: Write the failing test**

Create `model/config_test.go`:

```go
package model

import "testing"

func TestGithubConfigEffectiveTTL(t *testing.T) {
	thirty := 30
	zero := 0
	cases := []struct {
		name string
		ttl  *int
		want int
	}{
		{"nil defaults to 60", nil, 60},
		{"explicit zero disables", &zero, 0},
		{"explicit value", &thirty, 30},
	}
	for _, c := range cases {
		got := GithubConfig{IssueTTL: c.ttl}.EffectiveTTL()
		if got != c.want {
			t.Errorf("%s: EffectiveTTL() = %d, want %d", c.name, got, c.want)
		}
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./model/... -run TestGithubConfigEffectiveTTL`
Expected: compile error — `GithubConfig` undefined.

- [ ] **Step 3: Add the config struct, field, and method**

In `model/config.go`, add `Github GithubConfig` to the `Config` struct (after the `TUI TUIConfig` field):

```go
	TUI                  TUIConfig            `toml:"tui"`
	Github               GithubConfig         `toml:"github"`
```

And add (near the other config sub-structs, e.g. after `TUIConfig`):

```go
// GithubConfig holds settings for GitHub integration in the status bar.
type GithubConfig struct {
	// IssueTTL is the status cache lifetime in seconds. A pointer so an absent
	// section (nil → default 60) is distinguishable from an explicit 0 (disable
	// caching, always fetch live).
	IssueTTL *int `toml:"issue_ttl"`
}

// EffectiveTTL returns the cache TTL in seconds: 60 when unset, otherwise the
// configured value (0 means caching is disabled).
func (g GithubConfig) EffectiveTTL() int {
	if g.IssueTTL == nil {
		return 60
	}
	return *g.IssueTTL
}
```

- [ ] **Step 4: Run the test to verify it passes**

Run: `go test ./model/... -run TestGithubConfigEffectiveTTL -v`
Expected: PASS.

- [ ] **Step 5: Update the JSON schema**

In `sesh.schema.json`, insert a `github` property into the top-level `properties` object, immediately after the `cache` block. Replace:

```json
    "cache": {
      "type": "boolean",
      "description": "Enable caching for improved performance",
      "default": false
    },
```

with:

```json
    "cache": {
      "type": "boolean",
      "description": "Enable caching for improved performance",
      "default": false
    },
    "github": {
      "type": "object",
      "description": "GitHub integration settings for the status bar",
      "properties": {
        "issue_ttl": {
          "type": "integer",
          "description": "Status cache lifetime in seconds. 0 disables caching (always fetch live).",
          "minimum": 0,
          "default": 60
        }
      },
      "additionalProperties": false
    },
```

- [ ] **Step 6: Verify the schema is valid JSON**

Run: `python3 -c "import json; json.load(open('sesh.schema.json')); print('valid')"`
Expected: `valid`

- [ ] **Step 7: Commit**

```bash
git add model/config.go model/config_test.go sesh.schema.json
git commit -m "feat(config): add [github] issue_ttl setting"
```

---

### Task 2: `statuscache` package

**Files:**
- Create: `statuscache/statuscache.go`
- Test: `statuscache/statuscache_test.go`
- Regenerate: `statuscache/mock_StatusCache.go` (via `just mock`)

**Interfaces:**
- Produces:
  - `statuscache.Ref{ Number int; Title string; State string }`
  - `statuscache.Entry{ PR *Ref; Issue *Ref; Timestamp time.Time }` with method `(Entry) Preferred() (*Ref, bool)` — returns `PR` if non-nil, else `Issue`, else `(nil, false)`.
  - `statuscache.StatusCache` interface: `Read(key string) (Entry, bool, error)` and `Write(key string, entry Entry) error`.
  - `statuscache.NewFileStatusCache() *FileStatusCache` implementing `StatusCache`, storing `<dir>/<key>.gob` under `$XDG_CACHE_HOME/sesh/status/` (fallback `~/.cache/sesh/status/`).
  - `statuscache.Key(repoRoot, branch string) string` — `sha256(repoRoot + "\x00" + branch)` hex.

- [ ] **Step 1: Write the failing tests**

Create `statuscache/statuscache_test.go`:

```go
package statuscache

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestCache(t *testing.T) *FileStatusCache {
	t.Helper()
	return NewFileStatusCacheWithDir(filepath.Join(t.TempDir(), "status"))
}

func TestKeyStableAndDistinct(t *testing.T) {
	a := Key("/repo/one", "400")
	assert.Equal(t, a, Key("/repo/one", "400"), "same inputs => same key")
	assert.NotEqual(t, a, Key("/repo/two", "400"), "different repo => different key")
	assert.NotEqual(t, a, Key("/repo/one", "401"), "different branch => different key")
}

func TestWriteReadRoundTrip(t *testing.T) {
	c := newTestCache(t)
	entry := Entry{Issue: &Ref{Number: 400, Title: "Dynamic tmux status bar", State: "OPEN"}}
	require.NoError(t, c.Write("k1", entry))

	got, found, err := c.Read("k1")
	require.NoError(t, err)
	assert.True(t, found)
	assert.Nil(t, got.PR)
	assert.Equal(t, &Ref{Number: 400, Title: "Dynamic tmux status bar", State: "OPEN"}, got.Issue)
}

func TestNegativeEntryRoundTrips(t *testing.T) {
	c := newTestCache(t)
	require.NoError(t, c.Write("k2", Entry{}))

	got, found, err := c.Read("k2")
	require.NoError(t, err)
	assert.True(t, found, "negative entry is still a hit")
	assert.Nil(t, got.PR)
	assert.Nil(t, got.Issue)
}

func TestMissingFileIsMiss(t *testing.T) {
	c := newTestCache(t)
	_, found, err := c.Read("does-not-exist")
	assert.NoError(t, err)
	assert.False(t, found)
}

func TestCorruptFileIsMiss(t *testing.T) {
	c := newTestCache(t)
	require.NoError(t, c.Write("k3", Entry{Issue: &Ref{Number: 1}}))
	// Corrupt the file on disk.
	require.NoError(t, writeRaw(c, "k3", []byte("not gob data")))

	_, found, err := c.Read("k3")
	assert.NoError(t, err)
	assert.False(t, found)
}

func TestPreferred(t *testing.T) {
	pr := &Ref{Number: 401, Title: "pr", State: "OPEN"}
	iss := &Ref{Number: 400, Title: "iss", State: "OPEN"}

	got, ok := Entry{PR: pr, Issue: iss}.Preferred()
	assert.True(t, ok)
	assert.Equal(t, pr, got, "PR preferred over issue")

	got, ok = Entry{Issue: iss}.Preferred()
	assert.True(t, ok)
	assert.Equal(t, iss, got)

	_, ok = Entry{}.Preferred()
	assert.False(t, ok)
}
```

- [ ] **Step 2: Run the tests to verify they fail**

Run: `go test ./statuscache/...`
Expected: compile error — package `statuscache` does not exist.

- [ ] **Step 3: Write the implementation**

Create `statuscache/statuscache.go`:

```go
package statuscache

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"os"
	"path/filepath"
	"time"
)

// Ref is one GitHub entity (issue or PR) as rendered in the status bar.
type Ref struct {
	Number int
	Title  string
	State  string
}

// Entry is the cached status for one branch. Both pointers may be nil (a
// "negative" entry: the branch has nothing to show), which still counts as a
// cache hit so the reader respects the TTL instead of refreshing every tick.
type Entry struct {
	PR        *Ref
	Issue     *Ref
	Timestamp time.Time
}

// Preferred returns the entity to render: the PR if present, otherwise the
// issue. ok is false for a negative entry.
func (e Entry) Preferred() (*Ref, bool) {
	if e.PR != nil {
		return e.PR, true
	}
	if e.Issue != nil {
		return e.Issue, true
	}
	return nil, false
}

// StatusCache reads and writes per-branch status entries.
type StatusCache interface {
	Read(key string) (Entry, bool, error) // bool=found; false (nil err) on miss/corrupt
	Write(key string, entry Entry) error
}

// Key builds the cache key (and filename stem) for a repo root + branch.
func Key(repoRoot, branch string) string {
	sum := sha256.Sum256([]byte(repoRoot + "\x00" + branch))
	return hex.EncodeToString(sum[:])
}

// FileStatusCache stores one gob file per key under a directory.
type FileStatusCache struct {
	dir string
}

// NewFileStatusCache stores entries under $XDG_CACHE_HOME/sesh/status
// (falling back to ~/.cache/sesh/status).
func NewFileStatusCache() *FileStatusCache {
	dir := os.Getenv("XDG_CACHE_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "."
		}
		dir = filepath.Join(home, ".cache")
	}
	return &FileStatusCache{dir: filepath.Join(dir, "sesh", "status")}
}

// NewFileStatusCacheWithDir stores entries under an explicit directory (tests).
func NewFileStatusCacheWithDir(dir string) *FileStatusCache {
	return &FileStatusCache{dir: dir}
}

func (c *FileStatusCache) path(key string) string {
	return filepath.Join(c.dir, key+".gob")
}

func (c *FileStatusCache) Read(key string) (Entry, bool, error) {
	data, err := os.ReadFile(c.path(key))
	if err != nil {
		return Entry{}, false, nil // missing → miss
	}
	var entry Entry
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&entry); err != nil {
		return Entry{}, false, nil // corrupt → miss
	}
	return entry, true, nil
}

func (c *FileStatusCache) Write(key string, entry Entry) error {
	if err := os.MkdirAll(c.dir, 0o755); err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(entry); err != nil {
		return err
	}
	tmp := c.path(key) + ".tmp"
	if err := os.WriteFile(tmp, buf.Bytes(), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, c.path(key))
}
```

Add a test-only raw-write helper at the bottom of `statuscache_test.go` (it needs package-internal access to `path`):

```go
func writeRaw(c *FileStatusCache, key string, data []byte) error {
	return os.WriteFile(c.path(key), data, 0o644)
}
```

And add `"os"` to that test file's imports.

- [ ] **Step 4: Regenerate mocks and run the tests**

Run: `just mock && go test ./statuscache/... -v`
Expected: PASS (all subtests).

- [ ] **Step 5: Commit**

```bash
git add statuscache/statuscache.go statuscache/statuscache_test.go
git commit -m "feat(statuscache): add branch-keyed gob cache for status entries"
```

---

### Task 3: `github.Resolve` + refactor `Issue`

**Files:**
- Modify: `github/github.go`
- Modify: `github/github_test.go`
- Regenerate: `github/mock_Github.go` (via `just mock`)

**Interfaces:**
- Consumes: `git.Git.CurrentBranch(path) (bool, string, error)`, `git.Git.ShowTopLevel(path) (bool, string, error)` (both existing).
- Produces:
  - `github.BranchRef{ RepoRoot string; Branch string; Number int; HasNumber bool }`
  - `github.Github.Resolve(path string) (BranchRef, bool)` — `ok` is true whenever `path` is in a git repo (number optional). Added to the `Github` interface.
  - `github.Issue` refactored to call `Resolve` (unchanged external behavior).

- [ ] **Step 1: Write the failing test for `Resolve`**

Add to `github/github_test.go`:

```go
func TestResolve(t *testing.T) {
	t.Run("numeric branch resolves repo, branch, and number", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockGit := new(git.MockGit)
		gh := NewGithub(mockShell, mockGit)
		path := "/Users/josh/c/sesh"
		mockGit.On("CurrentBranch", path).Return(true, "feat/400-status-bar", nil)
		mockGit.On("ShowTopLevel", path).Return(true, "/Users/josh/c/sesh", nil)

		ref, ok := gh.Resolve(path)

		assert.True(t, ok)
		assert.Equal(t, BranchRef{RepoRoot: "/Users/josh/c/sesh", Branch: "feat/400-status-bar", Number: 400, HasNumber: true}, ref)
	})

	t.Run("non-numeric branch resolves with HasNumber false", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockGit := new(git.MockGit)
		gh := NewGithub(mockShell, mockGit)
		path := "/Users/josh/c/sesh"
		mockGit.On("CurrentBranch", path).Return(true, "main", nil)
		mockGit.On("ShowTopLevel", path).Return(true, "/Users/josh/c/sesh", nil)

		ref, ok := gh.Resolve(path)

		assert.True(t, ok)
		assert.Equal(t, "main", ref.Branch)
		assert.False(t, ref.HasNumber)
		assert.Equal(t, 0, ref.Number)
	})

	t.Run("not a git repo => ok false", func(t *testing.T) {
		mockShell := new(shell.MockShell)
		mockGit := new(git.MockGit)
		gh := NewGithub(mockShell, mockGit)
		path := "/tmp/x"
		mockGit.On("CurrentBranch", path).Return(false, "", fmt.Errorf("not a git repo"))

		_, ok := gh.Resolve(path)
		assert.False(t, ok)
	})
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `just mock && go test ./github/... -run TestResolve`
Expected: FAIL — `gh.Resolve` undefined.

- [ ] **Step 3: Add `Resolve` and refactor `Issue`**

In `github/github.go`, add `strconv` to imports, add the type + method to the interface and impl, and refactor `Issue`:

```go
// BranchRef identifies the repo, branch, and (optional) issue number for a path.
type BranchRef struct {
	RepoRoot  string
	Branch    string
	Number    int
	HasNumber bool
}
```

Add to the `Github` interface:

```go
type Github interface {
	Issue(path string) (Issue, bool, error)
	// Resolve returns repo root, branch, and the issue number parsed from the
	// branch (HasNumber=false if none). ok is false only when path is not a repo.
	Resolve(path string) (BranchRef, bool)
}
```

Add the implementation and refactor `Issue` (replace the existing `Issue` method):

```go
func (g *RealGithub) Resolve(path string) (BranchRef, bool) {
	ok, branch, err := g.git.CurrentBranch(path)
	if err != nil || !ok {
		return BranchRef{}, false
	}
	topOk, repoRoot, err := g.git.ShowTopLevel(path)
	if err != nil || !topOk {
		return BranchRef{}, false
	}
	ref := BranchRef{RepoRoot: repoRoot, Branch: branch}
	if numStr, has := parseIssueNumber(branch); has {
		if n, err := strconv.Atoi(numStr); err == nil {
			ref.Number = n
			ref.HasNumber = true
		}
	}
	return ref, true
}

func (g *RealGithub) Issue(path string) (Issue, bool, error) {
	ref, ok := g.Resolve(path)
	if !ok || !ref.HasNumber {
		return Issue{}, false, nil
	}

	out, err := g.shell.Cmd("gh", "issue", "view", strconv.Itoa(ref.Number), "--json", "number,title,state")
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

- [ ] **Step 4: Update existing `Issue` tests to mock `ShowTopLevel`**

`Issue` now calls `Resolve`, which calls `ShowTopLevel` whenever `CurrentBranch` succeeds. In `github/github_test.go`'s `TestIssue`, add a `ShowTopLevel` expectation to every subtest where `CurrentBranch` returns `true`. For each such subtest add, right after its `CurrentBranch` mock line:

```go
		mockGit.On("ShowTopLevel", path).Return(true, "/Users/josh/c/sesh", nil)
```

This applies to: "returns the issue on success", "not found when branch has no number", "not found when gh errors", and "not found when gh returns malformed json". The "not found when not a git repo" subtest (where `CurrentBranch` returns `false`) needs no `ShowTopLevel` mock — `Resolve` returns before calling it.

- [ ] **Step 5: Regenerate mocks and run the package tests**

Run: `just mock && go test ./github/... -v`
Expected: PASS (`TestResolve` + `TestParseIssueNumber` + all `TestIssue` subtests).

- [ ] **Step 6: Commit**

```bash
git add github/github.go github/github_test.go
git commit -m "feat(github): add Resolve and route Issue through it"
```

---

### Task 4: `refresher` package

**Files:**
- Create: `refresher/refresher.go`
- Test: `refresher/refresher_test.go`
- Regenerate: `refresher/mock_Refresher.go` (via `just mock`)

**Interfaces:**
- Produces:
  - `refresher.Refresher` interface: `Spawn(path string) error`.
  - `refresher.NewRefresher() Refresher` returning `*RealRefresher`.
  - Internal `refreshArgs(path string) []string` — `["status", "--refresh"]`, with `path` appended only when non-empty.

`RealRefresher.Spawn` launches `<self> status --refresh [path]` detached (`Setsid`) and does not wait, so it outlives the parent (tmux waits for the foreground `#(sesh status)` to exit). `Spawn` is thin glue over `os.Executable` + `exec`; the testable logic is `refreshArgs`.

- [ ] **Step 1: Write the failing test**

Create `refresher/refresher_test.go`:

```go
package refresher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefreshArgs(t *testing.T) {
	assert.Equal(t, []string{"status", "--refresh"}, refreshArgs(""))
	assert.Equal(t, []string{"status", "--refresh", "/repo"}, refreshArgs("/repo"))
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./refresher/...`
Expected: compile error — package `refresher` does not exist.

- [ ] **Step 3: Write the implementation**

Create `refresher/refresher.go`:

```go
package refresher

import (
	"log/slog"
	"os"
	"os/exec"
	"syscall"
)

// Refresher launches a detached `sesh status --refresh` to repopulate the
// status cache without blocking the caller.
type Refresher interface {
	// Spawn launches a detached refresh for path. An empty path lets the child
	// resolve the directory itself (attached tmux session, then cwd).
	Spawn(path string) error
}

type RealRefresher struct{}

func NewRefresher() Refresher {
	return &RealRefresher{}
}

func refreshArgs(path string) []string {
	args := []string{"status", "--refresh"}
	if path != "" {
		args = append(args, path)
	}
	return args
}

func (r *RealRefresher) Spawn(path string) error {
	self, err := os.Executable()
	if err != nil {
		slog.Debug("refresher: os.Executable failed", "error", err)
		return err
	}
	cmd := exec.Command(self, refreshArgs(path)...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	// Detach into its own session so it outlives this (foreground) process.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		slog.Debug("refresher: start failed", "error", err)
		return err
	}
	// Do not Wait — release the child to run independently.
	return cmd.Process.Release()
}
```

- [ ] **Step 4: Regenerate mocks and run the test**

Run: `just mock && go test ./refresher/... -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add refresher/refresher.go refresher/refresher_test.go
git commit -m "feat(refresher): add detached sesh status --refresh spawner"
```

---

### Task 5: Wire `StatusCache` and `Refresher` into DI

**Files:**
- Modify: `seshcli/deps.go`

**Interfaces:**
- Consumes: `statuscache.NewFileStatusCache()` (Task 2), `refresher.NewRefresher()` (Task 4).
- Produces: `BaseDeps.StatusCache statuscache.StatusCache` and `BaseDeps.Refresher refresher.Refresher`, reachable as `deps.StatusCache` / `deps.Refresher`.

> No new test — wiring is covered by the build and Tasks 6–7. Deliverable is a compiling dependency graph.

- [ ] **Step 1: Add imports**

In `seshcli/deps.go`, add to the import block (alphabetical):

```go
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/joshmedeski/sesh/v2/picker"
	"github.com/joshmedeski/sesh/v2/previewer"
	"github.com/joshmedeski/sesh/v2/refresher"
	"github.com/joshmedeski/sesh/v2/replacer"
	"github.com/joshmedeski/sesh/v2/runtimewrap"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/joshmedeski/sesh/v2/startup"
	"github.com/joshmedeski/sesh/v2/statuscache"
```

(Place `refresher` after `previewer`/before `replacer`, and `statuscache` after `startup`, to keep alphabetical order — adjust to match the existing block.)

- [ ] **Step 2: Add fields to `BaseDeps`**

In the `BaseDeps` struct, after `Github github.Github`:

```go
	Github      github.Github
	StatusCache statuscache.StatusCache
	Refresher   refresher.Refresher
```

(Re-gofmt the struct so the field alignment is consistent.)

- [ ] **Step 3: Construct and assign in `NewBaseDeps`**

After `gh := github.NewGithub(sh, g)`:

```go
	gh := github.NewGithub(sh, g)
	sc := statuscache.NewFileStatusCache()
	rf := refresher.NewRefresher()
```

And in the returned `&BaseDeps{...}` literal, after `Github: gh,`:

```go
		Github:      gh,
		StatusCache: sc,
		Refresher:   rf,
```

- [ ] **Step 4: Verify it compiles**

Run: `go build ./...`
Expected: builds with no errors.

- [ ] **Step 5: Commit**

```bash
git add seshcli/deps.go
git commit -m "feat(seshcli): wire StatusCache and Refresher into BaseDeps"
```

---

### Task 6: `sesh status` cache logic + `--refresh` flag

**Files:**
- Modify: `seshcli/status.go`
- Modify: `seshcli/status_test.go`

**Interfaces:**
- Consumes: `config.Github.EffectiveTTL()`, `deps.Github.Issue`/`Resolve`, `deps.StatusCache.Read`/`Write`, `deps.Refresher.Spawn`, `statuscache.Key`, `statuscache.Entry`/`Ref`, `formatStatus` (existing).
- Produces (within `seshcli`):
  - `computeStatus(deps *Deps, ttl int, path string) (output string, spawn bool)` — the pure decision: what to print and whether a refresh should be spawned.
  - `runRefresh(deps *Deps, path string) error` — the `--refresh` handler: live fetch + cache write (always writes, even a negative entry).
  - `toIssue(ref statuscache.Ref) github.Issue` — adapts a cache ref to the value `formatStatus` consumes.
  - A `--refresh` bool flag and optional path arg on the `status` command.

- [ ] **Step 1: Write the failing tests**

Replace the body of `seshcli/status_test.go`'s imports/add new tests. Add these tests (keep the existing `TestStatusPath` and `TestFormatStatus`):

```go
func TestToIssue(t *testing.T) {
	got := toIssue(statuscache.Ref{Number: 400, Title: "x", State: "OPEN"})
	assert.Equal(t, github.Issue{Number: 400, Title: "x", State: "OPEN"}, got)
}

func TestComputeStatus(t *testing.T) {
	path := "/repo"
	ref := github.BranchRef{RepoRoot: "/repo", Branch: "400", Number: 400, HasNumber: true}
	key := statuscache.Key("/repo", "400")

	t.Run("fresh issue hit prints and does not spawn", func(t *testing.T) {
		gh := new(github.MockGithub)
		sc := new(statuscache.MockStatusCache)
		gh.On("Resolve", path).Return(ref, true)
		sc.On("Read", key).Return(statuscache.Entry{
			Issue:     &statuscache.Ref{Number: 400, Title: "Dynamic tmux status bar", State: "OPEN"},
			Timestamp: time.Now(),
		}, true, nil)
		deps := &Deps{}
		deps.Github = gh
		deps.StatusCache = sc

		out, spawn := computeStatus(deps, 60, path)

		assert.Equal(t, "#[fg=green,bold]OPEN#[default] #[fg=magenta]Issue #400#[default] Dynamic tmux status bar", out)
		assert.False(t, spawn)
	})

	t.Run("fresh negative hit prints nothing and does not spawn", func(t *testing.T) {
		gh := new(github.MockGithub)
		sc := new(statuscache.MockStatusCache)
		gh.On("Resolve", path).Return(ref, true)
		sc.On("Read", key).Return(statuscache.Entry{Timestamp: time.Now()}, true, nil)
		deps := &Deps{}
		deps.Github = gh
		deps.StatusCache = sc

		out, spawn := computeStatus(deps, 60, path)

		assert.Equal(t, "", out)
		assert.False(t, spawn)
	})

	t.Run("stale hit prints and spawns", func(t *testing.T) {
		gh := new(github.MockGithub)
		sc := new(statuscache.MockStatusCache)
		gh.On("Resolve", path).Return(ref, true)
		sc.On("Read", key).Return(statuscache.Entry{
			Issue:     &statuscache.Ref{Number: 400, Title: "T", State: "OPEN"},
			Timestamp: time.Now().Add(-2 * time.Minute),
		}, true, nil)
		deps := &Deps{}
		deps.Github = gh
		deps.StatusCache = sc

		out, spawn := computeStatus(deps, 60, path)

		assert.NotEqual(t, "", out)
		assert.True(t, spawn)
	})

	t.Run("miss prints nothing and spawns", func(t *testing.T) {
		gh := new(github.MockGithub)
		sc := new(statuscache.MockStatusCache)
		gh.On("Resolve", path).Return(ref, true)
		sc.On("Read", key).Return(statuscache.Entry{}, false, nil)
		deps := &Deps{}
		deps.Github = gh
		deps.StatusCache = sc

		out, spawn := computeStatus(deps, 60, path)

		assert.Equal(t, "", out)
		assert.True(t, spawn)
	})

	t.Run("ttl zero uses live Issue and never spawns", func(t *testing.T) {
		gh := new(github.MockGithub)
		gh.On("Issue", path).Return(github.Issue{Number: 400, Title: "T", State: "OPEN"}, true, nil)
		deps := &Deps{}
		deps.Github = gh

		out, spawn := computeStatus(deps, 0, path)

		assert.NotEqual(t, "", out)
		assert.False(t, spawn)
		gh.AssertNotCalled(t, "Resolve", path)
	})
}

func TestRunRefresh(t *testing.T) {
	path := "/repo"
	ref := github.BranchRef{RepoRoot: "/repo", Branch: "400", Number: 400, HasNumber: true}
	key := statuscache.Key("/repo", "400")

	t.Run("writes issue entry on success", func(t *testing.T) {
		gh := new(github.MockGithub)
		sc := new(statuscache.MockStatusCache)
		gh.On("Resolve", path).Return(ref, true)
		gh.On("Issue", path).Return(github.Issue{Number: 400, Title: "T", State: "OPEN"}, true, nil)
		sc.On("Write", key, mock.MatchedBy(func(e statuscache.Entry) bool {
			return e.Issue != nil && e.Issue.Number == 400 && e.PR == nil
		})).Return(nil)
		deps := &Deps{}
		deps.Github = gh
		deps.StatusCache = sc

		assert.NoError(t, runRefresh(deps, path))
		sc.AssertExpectations(t)
	})

	t.Run("writes negative entry when no issue", func(t *testing.T) {
		gh := new(github.MockGithub)
		sc := new(statuscache.MockStatusCache)
		gh.On("Resolve", path).Return(github.BranchRef{RepoRoot: "/repo", Branch: "main"}, true)
		sc.On("Write", statuscache.Key("/repo", "main"), mock.MatchedBy(func(e statuscache.Entry) bool {
			return e.Issue == nil && e.PR == nil
		})).Return(nil)
		deps := &Deps{}
		deps.Github = gh
		deps.StatusCache = sc

		assert.NoError(t, runRefresh(deps, path))
		sc.AssertExpectations(t)
	})

	t.Run("not a repo writes nothing", func(t *testing.T) {
		gh := new(github.MockGithub)
		sc := new(statuscache.MockStatusCache)
		gh.On("Resolve", path).Return(github.BranchRef{}, false)
		deps := &Deps{}
		deps.Github = gh
		deps.StatusCache = sc

		assert.NoError(t, runRefresh(deps, path))
		sc.AssertNotCalled(t, "Write", mock.Anything, mock.Anything)
	})
}
```

Update the imports of `seshcli/status_test.go` to include `"time"`, `"github.com/joshmedeski/sesh/v2/statuscache"`, and `"github.com/stretchr/testify/mock"` (alongside the existing `github`, `assert`, etc.).

- [ ] **Step 2: Run the tests to verify they fail**

Run: `just mock && go test ./seshcli/... -run 'TestToIssue|TestComputeStatus|TestRunRefresh'`
Expected: FAIL — `computeStatus`, `runRefresh`, `toIssue` undefined.

- [ ] **Step 3: Implement the logic and flag in `seshcli/status.go`**

Rewrite `seshcli/status.go` to:

```go
package seshcli

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/statuscache"
)

func NewStatusCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show contextual status for the current session (for the tmux status bar)",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			path := statusPath(deps)
			if len(args) > 0 && args[0] != "" {
				path = args[0]
			}
			if path == "" {
				return nil
			}

			refresh, _ := cmd.Flags().GetBool("refresh")
			if refresh {
				return runRefresh(deps, path)
			}

			ttl := deps.Config.Github.EffectiveTTL()
			out, spawn := computeStatus(deps, ttl, path)
			if out != "" {
				fmt.Print(out)
			}
			if spawn {
				_ = deps.Refresher.Spawn(path)
			}
			return nil
		},
	}
	cmd.Flags().Bool("refresh", false, "Internal: fetch live data and update the status cache (used for background refresh)")
	_ = cmd.Flags().MarkHidden("refresh")
	return cmd
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

// computeStatus decides what to print and whether a background refresh should
// be spawned. With ttl==0 the cache is bypassed and gh is queried live.
func computeStatus(deps *Deps, ttl int, path string) (string, bool) {
	if ttl == 0 {
		issue, found, _ := deps.Github.Issue(path)
		if found {
			return formatStatus(issue), false
		}
		return "", false
	}

	ref, ok := deps.Github.Resolve(path)
	if !ok {
		return "", false
	}

	key := statuscache.Key(ref.RepoRoot, ref.Branch)
	entry, found, _ := deps.StatusCache.Read(key)

	out := ""
	if found {
		if r, ok := entry.Preferred(); ok {
			out = formatStatus(toIssue(*r))
		}
	}
	stale := !found || time.Since(entry.Timestamp) > time.Duration(ttl)*time.Second
	return out, stale
}

// runRefresh performs the live gh fetch and always writes a cache entry
// (a negative entry when there is nothing to show), so the reader respects
// the TTL instead of re-spawning every tick.
func runRefresh(deps *Deps, path string) error {
	ref, ok := deps.Github.Resolve(path)
	if !ok {
		return nil // not a repo — nothing to cache
	}

	var entry statuscache.Entry
	if issue, found, _ := deps.Github.Issue(path); found {
		entry.Issue = &statuscache.Ref{Number: issue.Number, Title: issue.Title, State: issue.State}
	}
	entry.Timestamp = time.Now()
	return deps.StatusCache.Write(statuscache.Key(ref.RepoRoot, ref.Branch), entry)
}

// formatStatus renders an issue as a tmux-styled status line.
func formatStatus(issue github.Issue) string {
	color := "green"
	if issue.State != "OPEN" {
		color = "red"
	}
	return fmt.Sprintf("#[fg=%s,bold]%s#[default] #[fg=magenta]Issue #%d#[default] %s", color, issue.State, issue.Number, issue.Title)
}

// toIssue adapts a cache Ref into the github.Issue value formatStatus consumes.
func toIssue(ref statuscache.Ref) github.Issue {
	return github.Issue{Number: ref.Number, Title: ref.Title, State: ref.State}
}
```

(The `formatStatus` body is unchanged from PR #401 — it is reproduced here only because the file is being rewritten.)

- [ ] **Step 4: Run the tests to verify they pass**

Run: `go test ./seshcli/... -run 'TestToIssue|TestComputeStatus|TestRunRefresh|TestFormatStatus|TestStatusPath' -v`
Expected: PASS (all subtests).

- [ ] **Step 5: Verify the command still builds and registers**

Run: `go build ./... && go run . status --help`
Expected: builds; help shows the short description. `--refresh` is hidden (not listed).

- [ ] **Step 6: Commit**

```bash
git add seshcli/status.go seshcli/status_test.go
git commit -m "feat(seshcli): read status from cache and refresh in the background"
```

---

### Task 7: Warm the cache on `sesh connect`

**Files:**
- Modify: `seshcli/connect.go`
- Modify/Create: `seshcli/connect_test.go`

**Interfaces:**
- Consumes: `deps.Config.Github.EffectiveTTL()`, `deps.Refresher.Spawn`.
- Produces (within `seshcli`): `maybeWarmStatus(deps *Deps)` — spawns a path-less background refresh when caching is enabled (the child resolves the now-attached session via `statusPath`).

- [ ] **Step 1: Write the failing test**

Create `seshcli/connect_test.go`:

```go
package seshcli

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/refresher"
)

func TestMaybeWarmStatus(t *testing.T) {
	t.Run("spawns when caching enabled", func(t *testing.T) {
		rf := new(refresher.MockRefresher)
		rf.On("Spawn", "").Return(nil)
		deps := &Deps{}
		deps.Refresher = rf
		deps.Config = model.Config{} // IssueTTL nil → EffectiveTTL 60

		maybeWarmStatus(deps)

		rf.AssertExpectations(t)
	})

	t.Run("does not spawn when issue_ttl is zero", func(t *testing.T) {
		rf := new(refresher.MockRefresher)
		deps := &Deps{}
		deps.Refresher = rf
		zero := 0
		deps.Config = model.Config{Github: model.GithubConfig{IssueTTL: &zero}}

		maybeWarmStatus(deps)

		rf.AssertNotCalled(t, "Spawn", "")
	})
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `just mock && go test ./seshcli/... -run TestMaybeWarmStatus`
Expected: FAIL — `maybeWarmStatus` undefined.

- [ ] **Step 3: Add the helper and call it from connect**

In `seshcli/connect.go`, add the helper (above or below `NewConnectCommand`):

```go
// maybeWarmStatus spawns a background status refresh after a connect so the
// status bar is warm before its first render. The child resolves the
// (now-attached) session path itself. No-op when caching is disabled.
func maybeWarmStatus(deps *Deps) {
	if deps.Config.Github.EffectiveTTL() == 0 {
		return
	}
	_ = deps.Refresher.Spawn("")
}
```

And call it after the existing post-connect cache refresh block (right before `return nil` at the end of `RunE`):

```go
		// Refresh cache in background so next sesh list has fresh data
		if deps.CachingLister != nil {
			deps.CachingLister.RefreshCache(lister.ListOptions{})
			deps.CachingLister.Wait()
		}
		maybeWarmStatus(deps)
		return nil
```

- [ ] **Step 4: Run the test to verify it passes**

Run: `go test ./seshcli/... -run TestMaybeWarmStatus -v`
Expected: PASS (both subtests).

- [ ] **Step 5: Commit**

```bash
git add seshcli/connect.go seshcli/connect_test.go
git commit -m "feat(seshcli): warm the status cache on connect"
```

---

### Final verification

- [ ] **Run the full suite**

Run: `just test`
Expected: mocks regenerate; all packages pass with `-race`, including `model`, `statuscache`, `github`, `refresher`, and `seshcli`.

- [ ] **Manual smoke test (this repo is on branch `400`)**

```bash
go build -o /tmp/sesh-cache-test .
rm -rf "${XDG_CACHE_HOME:-$HOME/.cache}/sesh/status"
/tmp/sesh-cache-test status        # cold: prints nothing, spawns a refresh
sleep 2
/tmp/sesh-cache-test status        # warm: prints the magenta "Issue #400" badge instantly
```

Expected: the first call is silent (cold miss, background refresh kicked off); after a moment the cache file exists under `…/sesh/status/` and the second call prints the formatted badge with no network wait. (Requires `gh` authenticated; otherwise both print nothing.)
