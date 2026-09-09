# `sesh status` — Dynamic tmux status bar (MVP)

Issue: [#400](https://github.com/joshmedeski/sesh/issues/400)
Date: 2026-06-27

## Goal

Add a `sesh status` command that prints the GitHub issue title associated with the
current session's git branch. It is designed to be embedded in a tmux status bar:

```sh
set -g status-left "#[fg=blue,bold]#S #[fg=white,nobold]#(sesh status)"
```

This is the MVP slice of the larger "dynamic tmux status bar" vision (gitmux-inspired).
The GitHub issue **title and open/closed state** are in scope here.

## Behavior

- Resolve the directory for the current session.
- Determine the git branch of that directory.
- Parse the first run of digits out of the branch name as an issue number
  (e.g. `feat/400-status-bar` → `400`, `400` → `400`, `bugfix/issue-12` → `12`).
- Look up the issue via the `gh` CLI: `gh issue view <n> --json number,title,state`.
- Format and print a tmux-styled status line to stdout (see Output format).

### Output format

The status line uses **tmux format markup** (`#[fg=…]` … `#[default]`), which tmux
re-interprets in `#()` output — the same mechanism gitmux relies on. (ANSI escape
codes, as used by the `icon` package for the picker, do **not** render in the status
bar, so they are not used here.)

A leading colored state badge is always shown, followed by `#<number> <title>`:

```
OPEN:    #[fg=green,bold]OPEN#[default] #400 Dynamic tmux status bar
CLOSED:  #[fg=red,bold]CLOSED#[default] #400 Dynamic tmux status bar
```

### No-result behavior

Print nothing and exit `0` for **every** "nothing to show" case:

- branch name contains no number
- directory is not a git repository
- `gh` is not installed or not authenticated
- the issue does not exist

This keeps `#(sesh status)` visually blank rather than leaking error text into the
status bar.

## Architecture

Follows existing repo conventions: each external tool is wrapped behind an interface
that depends on `shell.Shell`, and dependencies are wired in `seshcli/deps.go`.

### 1. `git` package — add `CurrentBranch`

Add to the existing `git.Git` interface:

```go
CurrentBranch(path string) (bool, string, error)
```

Implementation runs `git -C <path> rev-parse --abbrev-ref HEAD`. Mirrors the existing
`ShowTopLevel` / `GitCommonDir` methods. Mock regenerated via `just mock`.

### 2. New `github` package (`github/github.go`)

```go
type Issue struct {
    Number int    `json:"number"`
    Title  string `json:"title"`
    State  string `json:"state"` // "OPEN" | "CLOSED"
}

type Github interface {
    Issue(path string) (Issue, bool, error) // bool = found
}

type RealGithub struct {
    shell shell.Shell
    git   git.Git
}

func NewGithub(shell shell.Shell, git git.Git) Github
```

`Issue` flow:

1. `git.CurrentBranch(path)` → branch.
2. `parseIssueNumber(branch) (string, bool)` — pure helper, extracts the first
   contiguous run of digits. Unit-testable in isolation.
3. `gh issue view <n> --json number,title,state` via `shell.Cmd`.
4. Unmarshal the JSON into `Issue` with `encoding/json` and return `(issue, true, nil)`.

Returns `(Issue{}, false, nil)` for all no-result cases above — graceful, never an
error the caller has to suppress. Returning a struct (rather than a bare title) keeps
the API open to additional fields later (labels, assignee, PR data) without reshaping
the signature.

### 3. `sesh status` command (`seshcli/status.go`)

```go
func NewStatusCommand(base *BaseDeps) *cobra.Command
```

- Build deps via `buildDeps`.
- Resolve path: `deps.Lister.GetAttachedTmuxSession()` (same approach as `sesh root`);
  fall back to `os.Getwd()` when not attached to a tmux session.
- `issue, found, _ := deps.Github.Issue(path)` — error intentionally ignored.
- If `!found`, print nothing and return `nil`.
- Otherwise `fmt.Print(formatStatus(issue))`; return `nil`.

`formatStatus(issue Issue) string` is a small pure helper in the `seshcli` package
(unit-testable) that builds the tmux-styled line. The badge color is green for
`OPEN`, red for `CLOSED` (red for any non-`OPEN` state, defensively).

Registered in `seshcli/root_command.go` alongside the other commands.

### Dependency injection

`github` is config-free, so it joins `BaseDeps`:

- Add `Github github.Github` field to `BaseDeps`.
- Wire `github.NewGithub(sh, g)` in `NewBaseDeps()`.

## Testing

- `parseIssueNumber` — table test: plain number, prefixed/suffixed branch, no number,
  multiple number groups (takes the first).
- `github.Issue` — against mocked `shell.Shell` + `git.Git`: success path (JSON →
  struct), no-number branch (`found=false`, no shell call), `gh` error
  (`found=false`), not-a-repo (`found=false`), malformed JSON (`found=false`).
- `formatStatus` — open → green `OPEN` badge; closed → red `CLOSED` badge; non-`OPEN`
  state → red badge.
- Run `just test` (regenerates mocks, runs with coverage + race) before completion.

## README

Add a "tmux status bar" section under usage documenting the `status-left` snippet,
the `gh` install + auth requirement, and the current behavior (state badge + issue
number + title matching the branch number).

## Out of scope (deferred, tracked in #400)

- Pull request titles and draft/open/merged/closed state.
- Claude Code / Pi conversation titles and state.
- User configuration of priority/format.
