# Browser-driven worktree creation

**Date:** 2026-07-14
**Status:** Approved
**Branch:** 409-worktree-support

## Summary

Let `sesh worktree create` read the active browser tab's URL, extract the
GitHub `org/repo` and issue/PR number from it, and create (or reconnect to) a
worktree for that reference — without the user typing a number or being inside
the target repository.

Invoked via a new `--browser`/`-b` flag on the existing `create` command. macOS
only for now (uses `osascript`), mirroring the existing `focuser` package.

## Motivation

The worktree feature (issue #409) creates a worktree from an issue/PR number the
user supplies while inside the matching repo. In practice the user is often
looking at the issue/PR in their browser. Reading the active tab's URL removes
two friction points: typing the number, and having to be in the right working
directory (the URL already carries `org/repo`).

Issues are the first-class focus of the worktree feature, so they are the
primary tested path. PR URLs are supported at zero extra cost because the URL
disambiguates `/issues/` from `/pull/` and `worktree.Create` already handles PRs.

## Configuration

New top-level `[browser]` table:

```toml
[browser]
application = "Helium"
# url_command = "URL of active tab of front window"  # optional override
```

- `application` — the browser app name passed to AppleScript. Empty (no
  `[browser]` table) means the feature is unavailable.
- `url_command` — the AppleScript fragment that yields the URL. Defaults to
  `URL of active tab of front window` (Chrome-family: Helium, Chrome, Arc,
  Brave, Edge). Safari would use `URL of current tab of front window`.

Model change (`model/config.go`):

```go
Config struct {
    // ...existing...
    Browser BrowserConfig `toml:"browser"`
}

BrowserConfig struct {
    Application string `toml:"application"`
    URLCommand  string `toml:"url_command"`
}
```

`sesh.schema.json` gains the `browser` object with `application` and
`url_command` string properties (required by the repo's config-schema-sync
convention whenever a TOML struct changes).

## Components

### `browser` package (new)

Mirrors `focuser`: macOS-only osascript, skip semantics elsewhere.

```go
type Browser interface {
    // ActiveTabURL returns the front window's active-tab URL.
    // (\"\", false, nil) when skipped: non-macOS or no application configured.
    ActiveTabURL() (url string, ok bool, err error)
}

type RealBrowser struct {
    runtime runtimewrap.Runtime
    shell   shell.Shell
    config  model.BrowserConfig
}

func NewBrowser(runtime runtimewrap.Runtime, shell shell.Shell, config model.BrowserConfig) Browser
```

Behavior:

- `runtime.GOOS() != "darwin"` → `("", false, nil)`.
- `config.Application == ""` → `("", false, nil)`.
- Otherwise build the script:
  `tell application "<app>" to return <url_command> of front window`
  where `<url_command>` defaults to `URL of active tab of front window`, and run
  `osascript -e <script>` via `shell.Cmd`. Return the trimmed URL.
- A shell/osascript error (browser not running, no window) is returned as an
  error.

### `parseGitHubRef` (pure function, `worktree` package)

```go
// parseGitHubRef extracts repo + number + kind from a GitHub issue/PR URL.
//   https://github.com/joshmedeski/sesh/issues/409 → (\"joshmedeski/sesh\", 409, false, true)
//   https://github.com/joshmedeski/sesh/pull/409   → (\"joshmedeski/sesh\", 409, true,  true)
func parseGitHubRef(rawURL string) (repo string, number int, isPR bool, ok bool)
```

- Regex-based, matching `github.com/<org>/<repo>/(issues|pull)/<number>`.
- Tolerant of trailing path segments, query strings, and fragments
  (`/issues/409/`, `?foo=bar`, `#issuecomment-123`).
- Unrecognized URL → `ok == false`.

### Flow wiring

- `model.WorktreeCreateOpts` gains `FromBrowser bool`.
- `worktree.RealWorktree` gains a `browser browser.Browser` field (added to
  `NewWorktree`'s signature).
- `RealWorktree.Create`: when `opts.FromBrowser` is set, before the existing
  `resolveConfig` step:
  1. `url, ok, err := w.browser.ActiveTabURL()`. Error → return it. `!ok` →
     return a clear "browser lookup unavailable / not configured" error.
  2. `repo, number, isPR, ok := parseGitHubRef(url)`. `!ok` → return an error
     naming the URL that was read.
  3. Populate `opts.Repo = repo`, `opts.Number = number`,
     `opts.Pr = opts.Pr || isPR`.
  The rest of `Create` is unchanged. Because `resolveConfig` already matches a
  `[[worktree]]` config by `Name` when `opts.Repo` is non-empty, the URL's
  `org/repo` flows straight into the existing "match by name" branch — no new
  matching logic.

### CLI (`seshcli/worktree.go`)

- Add `--browser`/`-b` bool flag.
- Change `Args` from `cobra.ExactArgs(1)` to `cobra.MaximumNArgs(1)`.
- Require exactly one of {a number arg, `--browser`}:
  - neither → usage error ("provide an issue/PR number or use --browser").
  - both → usage error.
- When `--browser`, set `FromBrowser: true` and skip the `strconv.Atoi` parse of
  the (absent) number arg.
- `--repo`, `--pr`, `--switch` continue to work; `--repo` overrides the URL's
  repo, `--pr` forces the PR path.

### DI (`seshcli/deps.go`)

- Add `Browser browser.Browser` to `BaseDeps`.
- Construct in `NewBaseDeps`: `br := browser.NewBrowser(runtime, sh, ...)`.
  Because the browser config is config-dependent, construction moves to
  `BuildAll` (like `Zoxide`), or `NewBrowser` takes config at build time in
  `BuildAll`. Decision: construct the `Browser` in `BuildAll` (it needs
  `config.Browser`) and store it on `Deps`, then pass it into
  `worktree.NewWorktree`. `BaseDeps` keeps only config-free deps, consistent
  with the existing split.
- Pass the browser into `worktree.NewWorktree(config, git, github, connector,
  browser, home, os, path)` in `BuildAll`.

## Error handling

| Situation | Result |
|-----------|--------|
| Not macOS, or no `[browser].application` | Clear error: browser lookup unavailable / not configured |
| osascript fails (browser closed, no window) | Surface the shell error |
| URL isn't a GitHub issue/PR URL | Error naming the URL that was read |
| URL parses but no matching `[[worktree]]` config | Existing `resolveConfig` "no config found for repo %q" error |
| Both a number arg and `--browser` given | CLI usage error |
| Neither given | CLI usage error |

## Testing

- `browser` package: table tests with mocked `runtimewrap.Runtime` (darwin vs
  linux) and mocked `shell.Shell`, asserting:
  - exact `osascript -e <script>` args for a configured app,
  - default vs overridden `url_command`,
  - skip (`"", false, nil`) on non-darwin and empty application,
  - error propagation from the shell.
- `parseGitHubRef`: table of issue URL, PR URL, trailing slash, query string,
  fragment, non-GitHub URL, and garbage → asserting repo/number/isPR/ok.
- `worktree.Create` with `FromBrowser`: mocked `Browser` returning an issue URL
  and a PR URL, asserting the correct `resolveConfig` match by name and the
  right downstream git/github calls; plus the unparseable-URL error and the
  non-macOS skip error.
- Regenerate mocks (`just mock`) for the new `Browser` interface and the updated
  `NewWorktree` signature.

## Out of scope

- Non-macOS browser reading (Linux/Windows).
- Browsers whose AppleScript dictionary differs beyond a single `url_command`
  fragment (handled later if requested).
- Non-GitHub hosts (GitLab, Bitbucket).
- Reading anything other than the front window's active tab.
