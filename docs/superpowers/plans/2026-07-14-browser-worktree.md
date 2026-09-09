# Browser-Driven Worktree Creation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Let `sesh worktree create --browser` read the active browser tab's URL, extract the GitHub `org/repo` + issue/PR number, and create (or reconnect to) a worktree for it.

**Architecture:** A new `browser` package reads the front window's active-tab URL via `osascript` (macOS only), mirroring the existing `focuser` package. A pure `parseGitHubRef` function in the `worktree` package turns a GitHub URL into `(repo, number, isPR)`. `worktree.Create`, when given `FromBrowser`, calls the browser + parser to populate the existing `WorktreeCreateOpts` before its unchanged create flow runs. A `--browser`/`-b` flag on the CLI toggles this.

**Tech Stack:** Go 1.25, cobra (CLI), mockery/testify (mocks + tests), `osascript` (macOS AppleScript).

## Global Constraints

- Module: `github.com/joshmedeski/sesh/v2`; Go 1.25.0.
- macOS-only feature: guard on `runtimewrap.Runtime.GOOS() == "darwin"`; skip (no error) otherwise.
- Follow existing wrapper/DI patterns: external calls go through `shell.Shell`; dependencies wired in `seshcli/deps.go`.
- Mocks are generated, not hand-written: `.mockery.yaml` has `all: true` recursive. Run `just mock` after adding/changing any interface; never edit `mock_*.go` by hand.
- Use `just build` / `just test` (regenerates mocks) rather than raw `go` where a recipe exists.
- Config-schema-sync: any change to a `toml:"..."` struct in `model/` requires a matching update to `sesh.schema.json`.
- Test framework: `testify` (`require` for fatal, `assert` for non-fatal), table tests where natural.

---

### Task 1: Add `BrowserConfig` to the config model + JSON schema

**Files:**
- Modify: `model/config.go` (add `Browser` field to `Config`; add `BrowserConfig` struct)
- Modify: `sesh.schema.json` (add `browser` object)

**Interfaces:**
- Produces: `model.BrowserConfig{ Application string; URLCommand string }`, reachable as `model.Config.Browser`.

- [ ] **Step 1: Add the struct and field**

In `model/config.go`, add a field to the `Config` struct (place it near `Terminal`, since both are top-level UX options):

```go
	Terminal                string               `toml:"terminal"`
	Browser                 BrowserConfig        `toml:"browser"`
```

And add a new struct inside the `type ( ... )` block (place it after `WorktreeConfig`):

```go
	// BrowserConfig configures reading the active browser tab's URL so
	// `sesh worktree create --browser` can derive the target issue/PR.
	// macOS-only (osascript). An empty Application disables the feature.
	BrowserConfig struct {
		Application string `toml:"application"` // browser app name, e.g. "Helium"
		URLCommand  string `toml:"url_command"` // AppleScript fragment; default "URL of active tab of front window"
	}
```

- [ ] **Step 2: Update the JSON schema**

Find the top-level `properties` object in `sesh.schema.json` (alongside the existing `terminal` / `worktree` entries) and add a `browser` property:

```json
    "browser": {
      "type": "object",
      "description": "Read the active browser tab URL (macOS only) for `sesh worktree create --browser`.",
      "additionalProperties": false,
      "properties": {
        "application": {
          "type": "string",
          "description": "Browser application name, e.g. \"Helium\"."
        },
        "url_command": {
          "type": "string",
          "description": "AppleScript fragment yielding the URL. Default: \"URL of active tab of front window\"."
        }
      }
    }
```

(Match the existing indentation/quoting style in the file; if `worktree` sits inside a `properties` block with `additionalProperties: false`, mirror that exactly.)

- [ ] **Step 3: Verify it compiles**

Run: `go build ./...`
Expected: builds with no errors.

- [ ] **Step 4: Verify the schema is valid JSON**

Run: `python3 -m json.tool sesh.schema.json > /dev/null && echo OK`
Expected: `OK`

- [ ] **Step 5: Commit**

```bash
git add model/config.go sesh.schema.json
git commit -m "feat(config): add [browser] application + url_command (#409)"
```

---

### Task 2: Create the `browser` package

**Files:**
- Create: `browser/browser.go`
- Create: `browser/browser_test.go`
- Generated (via `just mock`): `browser/mock_Browser.go`

**Interfaces:**
- Consumes: `runtimewrap.Runtime`, `shell.Shell`, `model.BrowserConfig` (Task 1).
- Produces:
  - `type Browser interface { ActiveTabURL() (url string, ok bool, err error) }`
  - `func NewBrowser(runtime runtimewrap.Runtime, shell shell.Shell, config model.BrowserConfig) Browser`
  - Generated mock `browser.NewMockBrowser(t)` with `.EXPECT().ActiveTabURL()`.

- [ ] **Step 1: Write the failing tests**

Create `browser/browser_test.go`:

```go
package browser

import (
	"errors"
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/runtimewrap"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActiveTabURLOnMacOSDefaultCommand(t *testing.T) {
	rt := runtimewrap.NewMockRuntime(t)
	rt.EXPECT().GOOS().Return("darwin")
	s := shell.NewMockShell(t)
	s.EXPECT().
		Cmd("osascript", "-e", `tell application "Helium" to return URL of active tab of front window`).
		Return("https://github.com/joshmedeski/sesh/issues/409", nil)

	b := NewBrowser(rt, s, model.BrowserConfig{Application: "Helium"})
	url, ok, err := b.ActiveTabURL()
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "https://github.com/joshmedeski/sesh/issues/409", url)
}

func TestActiveTabURLWithURLCommandOverride(t *testing.T) {
	rt := runtimewrap.NewMockRuntime(t)
	rt.EXPECT().GOOS().Return("darwin")
	s := shell.NewMockShell(t)
	s.EXPECT().
		Cmd("osascript", "-e", `tell application "Safari" to return URL of current tab of front window`).
		Return("https://github.com/joshmedeski/sesh/pull/410", nil)

	b := NewBrowser(rt, s, model.BrowserConfig{
		Application: "Safari",
		URLCommand:  "URL of current tab of front window",
	})
	url, ok, err := b.ActiveTabURL()
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "https://github.com/joshmedeski/sesh/pull/410", url)
}

func TestActiveTabURLSkippedOnLinux(t *testing.T) {
	rt := runtimewrap.NewMockRuntime(t)
	rt.EXPECT().GOOS().Return("linux")
	s := shell.NewMockShell(t) // no shell calls expected
	b := NewBrowser(rt, s, model.BrowserConfig{Application: "Helium"})
	url, ok, err := b.ActiveTabURL()
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, "", url)
}

func TestActiveTabURLSkippedWhenNoApplication(t *testing.T) {
	rt := runtimewrap.NewMockRuntime(t)
	s := shell.NewMockShell(t) // no runtime/shell calls expected
	b := NewBrowser(rt, s, model.BrowserConfig{})
	url, ok, err := b.ActiveTabURL()
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, "", url)
}

func TestActiveTabURLPropagatesShellError(t *testing.T) {
	rt := runtimewrap.NewMockRuntime(t)
	rt.EXPECT().GOOS().Return("darwin")
	s := shell.NewMockShell(t)
	s.EXPECT().
		Cmd("osascript", "-e", `tell application "Helium" to return URL of active tab of front window`).
		Return("", errors.New("no window"))

	b := NewBrowser(rt, s, model.BrowserConfig{Application: "Helium"})
	_, ok, err := b.ActiveTabURL()
	require.Error(t, err)
	assert.False(t, ok)
}
```

Note: the empty-application test constructs the mocks but expects **no** calls on them — check the application field before touching `runtime`, so neither mock records an expectation.

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./browser/...`
Expected: FAIL — `browser.go` doesn't exist / `NewBrowser` undefined.

- [ ] **Step 3: Write the implementation**

Create `browser/browser.go`:

```go
package browser

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/runtimewrap"
	"github.com/joshmedeski/sesh/v2/shell"
)

// Browser reads the active tab URL from a browser's front window. Currently
// macOS-only via osascript; a no-op elsewhere.
type Browser interface {
	// ActiveTabURL returns the front window's active-tab URL.
	// Returns ("", false, nil) when skipped: non-macOS or no application
	// configured.
	ActiveTabURL() (url string, ok bool, err error)
}

// defaultURLCommand is the Chrome-family AppleScript fragment (Helium, Chrome,
// Arc, Brave, Edge). Safari uses "URL of current tab of front window".
const defaultURLCommand = "URL of active tab of front window"

type RealBrowser struct {
	runtime runtimewrap.Runtime
	shell   shell.Shell
	config  model.BrowserConfig
}

func NewBrowser(runtime runtimewrap.Runtime, shell shell.Shell, config model.BrowserConfig) Browser {
	return &RealBrowser{runtime, shell, config}
}

func (b *RealBrowser) ActiveTabURL() (string, bool, error) {
	if b.config.Application == "" {
		return "", false, nil
	}
	if b.runtime.GOOS() != "darwin" {
		return "", false, nil
	}
	urlCommand := b.config.URLCommand
	if urlCommand == "" {
		urlCommand = defaultURLCommand
	}
	script := fmt.Sprintf("tell application %q to return %s of front window", b.config.Application, urlCommand)
	url, err := b.shell.Cmd("osascript", "-e", script)
	if err != nil {
		return "", false, err
	}
	return url, true, nil
}
```

- [ ] **Step 4: Generate the mock**

Run: `just mock`
Expected: creates `browser/mock_Browser.go` (untracked until committed).

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./browser/...`
Expected: PASS (all 5 tests).

- [ ] **Step 6: Commit**

```bash
git add browser/browser.go browser/browser_test.go browser/mock_Browser.go
git commit -m "feat(browser): read active tab URL via osascript (macOS) (#409)"
```

---

### Task 3: Add `parseGitHubRef` to the `worktree` package

**Files:**
- Modify: `worktree/worktree.go` (add the pure function + its regexp)
- Create: `worktree/parse_test.go`

**Interfaces:**
- Produces: `func parseGitHubRef(rawURL string) (repo string, number int, isPR bool, ok bool)` (unexported; used by Task 4).

- [ ] **Step 1: Write the failing tests**

Create `worktree/parse_test.go`:

```go
package worktree

import "testing"

func TestParseGitHubRef(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		repo    string
		number  int
		isPR    bool
		ok      bool
	}{
		{"issue", "https://github.com/joshmedeski/sesh/issues/409", "joshmedeski/sesh", 409, false, true},
		{"pull", "https://github.com/joshmedeski/sesh/pull/410", "joshmedeski/sesh", 410, true, true},
		{"trailing slash", "https://github.com/joshmedeski/sesh/issues/409/", "joshmedeski/sesh", 409, false, true},
		{"query string", "https://github.com/joshmedeski/sesh/issues/409?foo=bar", "joshmedeski/sesh", 409, false, true},
		{"fragment", "https://github.com/joshmedeski/sesh/pull/410#issuecomment-1", "joshmedeski/sesh", 410, true, true},
		{"no scheme", "github.com/joshmedeski/sesh/issues/409", "joshmedeski/sesh", 409, false, true},
		{"repo home", "https://github.com/joshmedeski/sesh", "", 0, false, false},
		{"non-github", "https://gitlab.com/joshmedeski/sesh/issues/409", "", 0, false, false},
		{"garbage", "not a url", "", 0, false, false},
		{"empty", "", "", 0, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, number, isPR, ok := parseGitHubRef(tt.url)
			if ok != tt.ok || repo != tt.repo || number != tt.number || isPR != tt.isPR {
				t.Fatalf("parseGitHubRef(%q) = (%q, %d, %v, %v); want (%q, %d, %v, %v)",
					tt.url, repo, number, isPR, ok, tt.repo, tt.number, tt.isPR, tt.ok)
			}
		})
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./worktree/ -run TestParseGitHubRef`
Expected: FAIL — `parseGitHubRef` undefined.

- [ ] **Step 3: Write the implementation**

In `worktree/worktree.go`, add near the existing `issueRefRe` block:

```go
// githubRefRe matches a GitHub issue or PR URL, capturing org/repo, the kind
// ("issues"|"pull"), and the number. Tolerant of scheme, trailing path,
// query, and fragment.
var githubRefRe = regexp.MustCompile(`github\.com/([^/]+/[^/]+)/(issues|pull)/(\d+)`)

// parseGitHubRef extracts repo, number, and kind from a GitHub issue/PR URL.
// ok is false for any URL that is not a recognizable issue/PR reference.
func parseGitHubRef(rawURL string) (repo string, number int, isPR bool, ok bool) {
	m := githubRefRe.FindStringSubmatch(rawURL)
	if len(m) < 4 {
		return "", 0, false, false
	}
	n, err := strconv.Atoi(m[3])
	if err != nil {
		return "", 0, false, false
	}
	return m[1], n, m[2] == "pull", true
}
```

(`regexp` and `strconv` are already imported in `worktree.go`.)

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./worktree/ -run TestParseGitHubRef -v`
Expected: PASS (all subtests).

- [ ] **Step 5: Commit**

```bash
git add worktree/worktree.go worktree/parse_test.go
git commit -m "feat(worktree): parse GitHub issue/PR URLs into refs (#409)"
```

---

### Task 4: Wire `FromBrowser` into `worktree.Create`

**Files:**
- Modify: `model/connect_opts.go` (add `FromBrowser` to `WorktreeCreateOpts`)
- Modify: `worktree/worktree.go` (add `browser` field, extend `NewWorktree`, add resolution step in `Create`)
- Modify: `worktree/worktree_test.go` (update existing `NewWorktree` calls with a mock browser)
- Create: `worktree/browser_create_test.go`

**Interfaces:**
- Consumes: `browser.Browser` (Task 2), `parseGitHubRef` (Task 3).
- Produces:
  - `model.WorktreeCreateOpts.FromBrowser bool`
  - New signature: `func NewWorktree(config model.Config, g git.Git, gh github.Github, c connector.Connector, b browser.Browser, h home.Home, os oswrap.Os, p pathwrap.Path) Worktree`
  - (Note the new `b browser.Browser` parameter, inserted after the connector.)

- [ ] **Step 1: Add the opts field**

In `model/connect_opts.go`, add to `WorktreeCreateOpts`:

```go
type WorktreeCreateOpts struct {
	Number      int
	Repo        string // GitHub "org/repo" override; empty => detect from cwd
	Pr          bool   // force the pull-request path
	Switch      bool
	FromBrowser bool   // resolve Number/Repo/Pr from the active browser tab URL
}
```

- [ ] **Step 2: Write the failing test**

Create `worktree/browser_create_test.go`:

```go
package worktree

import (
	"os"
	"testing"

	"github.com/joshmedeski/sesh/v2/browser"
	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateFromBrowserIssue(t *testing.T) {
	mGit := git.NewMockGit(t)
	mGh := github.NewMockGithub(t)
	mConn := connector.NewMockConnector(t)
	mBrowser := browser.NewMockBrowser(t)
	mOs := oswrap.NewMockOs(t)
	h := home.NewHome(mOs)
	p := pathwrap.NewPath()

	mOs.EXPECT().UserHomeDir().Return("/home/me", nil).Maybe()
	mOs.EXPECT().ExpandEnv("/repo").Return("/repo").Maybe()

	// Browser resolves the issue URL for the configured repo.
	mBrowser.EXPECT().ActiveTabURL().
		Return("https://github.com/nutiliti/nutiliti/issues/2345", true, nil)

	// Not a PR => issue path.
	mGh.EXPECT().PrView("nutiliti/nutiliti", 2345).Return(github.PullRequest{}, false, nil)

	mOs.EXPECT().Stat("/repo/w/2345").Return(nil, os.ErrNotExist)
	mOs.EXPECT().MkdirAll("/repo/w", mock.Anything).Return(nil)
	mGit.EXPECT().WorktreeAdd("/repo", "/repo/w/2345", "jam/2345-1", "origin/main").Return("", nil)
	mConn.EXPECT().
		Connect("/repo/w/2345", model.ConnectOpts{Switch: true, Command: "nu_setup"}).
		Return("", nil)

	w := NewWorktree(nuConfig(), mGit, mGh, mConn, mBrowser, h, mOs, p)
	_, err := w.Create(model.WorktreeCreateOpts{FromBrowser: true, Switch: true})
	require.NoError(t, err)
}

func TestCreateFromBrowserUnparseableURL(t *testing.T) {
	mBrowser := browser.NewMockBrowser(t)
	mBrowser.EXPECT().ActiveTabURL().Return("https://example.com/", true, nil)

	w := NewWorktree(
		nuConfig(),
		git.NewMockGit(t), github.NewMockGithub(t), connector.NewMockConnector(t),
		mBrowser, home.NewHome(oswrap.NewMockOs(t)), oswrap.NewMockOs(t), pathwrap.NewPath(),
	)
	_, err := w.Create(model.WorktreeCreateOpts{FromBrowser: true})
	require.Error(t, err)
}

func TestCreateFromBrowserUnavailable(t *testing.T) {
	mBrowser := browser.NewMockBrowser(t)
	// Non-macOS / no application configured => skipped.
	mBrowser.EXPECT().ActiveTabURL().Return("", false, nil)

	w := NewWorktree(
		nuConfig(),
		git.NewMockGit(t), github.NewMockGithub(t), connector.NewMockConnector(t),
		mBrowser, home.NewHome(oswrap.NewMockOs(t)), oswrap.NewMockOs(t), pathwrap.NewPath(),
	)
	_, err := w.Create(model.WorktreeCreateOpts{FromBrowser: true})
	require.Error(t, err)
}
```

- [ ] **Step 3: Run the test to verify it fails**

Run: `go test ./worktree/ -run TestCreateFromBrowser`
Expected: FAIL — `NewWorktree` takes the wrong number of args / `browser` mock unused. (May not compile yet — that counts as failing.)

- [ ] **Step 4: Update `NewWorktree` and `Create`**

In `worktree/worktree.go`:

Add the import:

```go
	"github.com/joshmedeski/sesh/v2/browser"
```

Add the field to `RealWorktree` (after `connector`):

```go
type RealWorktree struct {
	config    model.Config
	git       git.Git
	github    github.Github
	connector connector.Connector
	browser   browser.Browser
	home      home.Home
	os        oswrap.Os
	path      pathwrap.Path
}
```

Update the constructor:

```go
func NewWorktree(
	config model.Config,
	g git.Git,
	gh github.Github,
	c connector.Connector,
	b browser.Browser,
	h home.Home,
	os oswrap.Os,
	p pathwrap.Path,
) Worktree {
	return &RealWorktree{config, g, gh, c, b, h, os, p}
}
```

Add a resolution step at the very top of `Create`, before `resolveConfig`:

```go
func (w *RealWorktree) Create(opts model.WorktreeCreateOpts) (string, error) {
	if opts.FromBrowser {
		resolved, err := w.resolveFromBrowser(opts)
		if err != nil {
			return "", err
		}
		opts = resolved
	}

	cfg, err := w.resolveConfig(opts.Repo)
	// ... unchanged ...
```

Add the helper (place it after `Create`, before `planCreate`):

```go
// resolveFromBrowser reads the active browser tab URL and fills Number/Repo/Pr
// from the GitHub issue/PR it points at.
func (w *RealWorktree) resolveFromBrowser(opts model.WorktreeCreateOpts) (model.WorktreeCreateOpts, error) {
	url, ok, err := w.browser.ActiveTabURL()
	if err != nil {
		return opts, err
	}
	if !ok {
		return opts, fmt.Errorf("browser tab lookup unavailable: set [browser].application (macOS only)")
	}
	repo, number, isPR, ok := parseGitHubRef(url)
	if !ok {
		return opts, fmt.Errorf("could not parse a GitHub issue/PR from browser URL %q", url)
	}
	opts.Repo = repo
	opts.Number = number
	opts.Pr = opts.Pr || isPR
	return opts, nil
}
```

(`fmt` is already imported.)

- [ ] **Step 5: Update existing `NewWorktree` call sites in tests**

In `worktree/worktree_test.go`, add the browser import:

```go
	"github.com/joshmedeski/sesh/v2/browser"
```

Every existing `NewWorktree(nuConfig(), mGit, mGh, mConn, h, mOs, p)` call must gain a mock browser argument after `mConn`. For each test, add near the other mocks:

```go
	mBrowser := browser.NewMockBrowser(t)
```

and change the call to:

```go
	w := NewWorktree(nuConfig(), mGit, mGh, mConn, mBrowser, h, mOs, p)
```

These tests set `FromBrowser: false`, so `ActiveTabURL` is never called — constructing `NewMockBrowser(t)` with no `.EXPECT()` is correct (testify passes when no calls are expected and none occur).

- [ ] **Step 6: Regenerate mocks**

Run: `just mock`
Expected: no diff for `browser/mock_Browser.go` (already generated in Task 2); confirms interfaces unchanged.

- [ ] **Step 7: Run the worktree tests**

Run: `go test ./worktree/...`
Expected: PASS — new browser tests plus all pre-existing tests.

- [ ] **Step 8: Commit**

```bash
git add model/connect_opts.go worktree/worktree.go worktree/worktree_test.go worktree/browser_create_test.go
git commit -m "feat(worktree): resolve create opts from browser tab URL (#409)"
```

---

### Task 5: DI wiring + `--browser` CLI flag

**Files:**
- Modify: `seshcli/deps.go` (construct `Browser` in `BuildAll`, add to `Deps`, pass into `NewWorktree`)
- Modify: `seshcli/worktree.go` (add `--browser`/`-b` flag, make number arg optional, validate)

**Interfaces:**
- Consumes: `browser.NewBrowser` (Task 2), updated `worktree.NewWorktree` (Task 4), `model.WorktreeCreateOpts.FromBrowser` (Task 4).

- [ ] **Step 1: Wire the browser in `deps.go`**

Add the import:

```go
	"github.com/joshmedeski/sesh/v2/browser"
```

Add the field to the `Deps` struct (after `Worktree`):

```go
	Worktree      worktree.Worktree
	Browser       browser.Browser
```

In `BuildAll`, construct the browser (it needs `config.Browser`, so it lives here alongside the other config-dependent deps) and pass it into `NewWorktree`. Replace:

```go
	wt := worktree.NewWorktree(config, b.Git, b.Github, c, b.Home, b.Os, b.Path)
```

with:

```go
	br := browser.NewBrowser(b.Runtime, b.Shell, config.Browser)
	wt := worktree.NewWorktree(config, b.Git, b.Github, c, br, b.Home, b.Os, b.Path)
```

Add `Browser: br,` to the returned `&Deps{...}` literal (after `Worktree: wt,`):

```go
		Worktree:      wt,
		Browser:       br,
```

- [ ] **Step 2: Verify it builds**

Run: `go build ./...`
Expected: builds cleanly.

- [ ] **Step 3: Add the flag and validation in `seshcli/worktree.go`**

Replace the body of `newWorktreeCreateCommand` so the number arg is optional and `--browser` is accepted. The full updated function:

```go
func newWorktreeCreateCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [number]",
		Aliases: []string{"c"},
		Short:   "Create (or reconnect to) a worktree for an issue/PR and connect to it",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fromBrowser, _ := cmd.Flags().GetBool("browser")

			if fromBrowser && len(args) > 0 {
				return fmt.Errorf("provide a number or --browser, not both")
			}
			if !fromBrowser && len(args) == 0 {
				return fmt.Errorf("provide an issue/PR number or use --browser")
			}

			number := 0
			if !fromBrowser {
				n, err := strconv.Atoi(args[0])
				if err != nil {
					return fmt.Errorf("worktree number must be an integer: %q", args[0])
				}
				number = n
			}

			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			repo, _ := cmd.Flags().GetString("repo")
			pr, _ := cmd.Flags().GetBool("pr")
			switchFlag, _ := cmd.Flags().GetBool("switch")

			_, err = deps.Worktree.Create(model.WorktreeCreateOpts{
				Number:      number,
				Repo:        repo,
				Pr:          pr,
				Switch:      switchFlag,
				FromBrowser: fromBrowser,
			})
			return err
		},
	}

	cmd.Flags().StringP("repo", "r", "", "Target repo as org/repo (overrides cwd detection)")
	cmd.Flags().BoolP("pr", "p", false, "Treat the number as a pull request")
	cmd.Flags().BoolP("switch", "s", false, "Switch the session rather than attach (for use outside tmux)")
	cmd.Flags().BoolP("browser", "b", false, "Read the issue/PR from the active browser tab (macOS)")

	return cmd
}
```

- [ ] **Step 4: Verify it builds and the CLI wires up**

Run: `go build ./... && go run . worktree create --help`
Expected: builds; help output lists the `-b, --browser` flag and `Usage: ... create [number]`.

- [ ] **Step 5: Verify argument validation**

Run: `go run . worktree create`
Expected: errors with `provide an issue/PR number or use --browser` (non-zero exit).

Run: `go run . worktree create 1 --browser`
Expected: errors with `provide a number or --browser, not both`.

- [ ] **Step 6: Commit**

```bash
git add seshcli/deps.go seshcli/worktree.go
git commit -m "feat(seshcli): add --browser flag + DI wiring for worktree create (#409)"
```

---

### Task 6: Full test suite + documentation

**Files:**
- Modify: `README.md` and/or `docs/` worktree documentation (whichever documents `[[worktree]]` and `sesh worktree`)

**Interfaces:** none (docs + verification only).

- [ ] **Step 1: Run the full test suite**

Run: `just test`
Expected: all packages pass, including `browser`, `worktree`, `seshcli`.

- [ ] **Step 2: Locate existing worktree docs**

Run: `grep -rln "sesh worktree\|\[\[worktree\]\]" README.md docs/`
Expected: lists the file(s) documenting the worktree feature. Open the most relevant one.

- [ ] **Step 3: Document `[browser]` and `--browser`**

In the worktree documentation, add a subsection near the existing `[[worktree]]` / `sesh worktree` content:

````markdown
### Creating a worktree from the browser (macOS)

With a `[browser]` configured, `sesh worktree create --browser` reads the URL
of your browser's active tab, extracts the GitHub `org/repo` and issue/PR
number, and creates the matching worktree — no need to type the number or be
inside the repo.

```toml
[browser]
application = "Helium"
# url_command = "URL of active tab of front window"  # optional; Safari uses "URL of current tab of front window"
```

```bash
# With github.com/joshmedeski/sesh/issues/409 open in the front tab:
sesh worktree create --browser   # or: sesh wt c -b
```

The URL's `org/repo` is matched against your `[[worktree]]` entries by `name`.
Both `/issues/N` and `/pull/N` URLs are supported. macOS only.
````

Match the surrounding heading levels and prose style of the file you're editing.

- [ ] **Step 4: Commit**

```bash
git add README.md docs/
git commit -m "docs: document [browser] config and worktree create --browser (#409)"
```

---

## Self-Review

**Spec coverage:**
- `[browser]` config (application + url_command) → Task 1. ✓
- JSON schema sync → Task 1. ✓
- `browser` package with skip semantics + osascript → Task 2. ✓
- `parseGitHubRef` (issue/PR/trailing/query/fragment/non-github) → Task 3. ✓
- `FromBrowser` opt + `Create` resolution + browser dependency → Task 4. ✓
- Match `[[worktree]]` by URL org/repo → reuses existing `resolveConfig` via `opts.Repo` (Task 4). ✓
- Issue-vs-PR from URL path → `parseGitHubRef` + `opts.Pr = opts.Pr || isPR` (Tasks 3, 4). ✓
- CLI `--browser`/`-b`, optional number arg, both/neither validation → Task 5. ✓
- DI wiring in `BuildAll` → Task 5. ✓
- Error handling (unavailable, unparseable, no matching config) → Task 4 (unavailable/unparseable) + reused `resolveConfig` error; CLI both/neither → Task 5. ✓
- Tests for browser, parse, worktree.Create, plus mock regen → Tasks 2–4. ✓
- Docs → Task 6. ✓

**Placeholder scan:** No TBD/TODO; every code step shows complete code; every test step shows the assertions.

**Type consistency:** `NewWorktree(config, git, github, connector, browser, home, os, path)` is used identically in Task 4 (definition + test call sites) and Task 5 (deps.go). `ActiveTabURL() (string, bool, error)` matches between Task 2 (definition), Task 4 (mock usage), and the `browser.Browser` field. `WorktreeCreateOpts.FromBrowser` defined in Task 4, consumed in Task 5. `parseGitHubRef` signature matches between Task 3 (definition/tests) and Task 4 (usage).
