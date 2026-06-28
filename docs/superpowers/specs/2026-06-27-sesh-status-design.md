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
Only the GitHub issue title is in scope here.

## Behavior

- Resolve the directory for the current session.
- Determine the git branch of that directory.
- Parse the first run of digits out of the branch name as an issue number
  (e.g. `feat/400-status-bar` → `400`, `400` → `400`, `bugfix/issue-12` → `12`).
- Look up the issue title via the `gh` CLI: `gh issue view <n> --json title -q .title`.
- Print the title to stdout.

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
type Github interface {
    IssueTitle(path string) (string, error)
}

type RealGithub struct {
    shell shell.Shell
    git   git.Git
}

func NewGithub(shell shell.Shell, git git.Git) Github
```

`IssueTitle` flow:

1. `git.CurrentBranch(path)` → branch.
2. `parseIssueNumber(branch) (string, bool)` — pure helper, extracts the first
   contiguous run of digits. Unit-testable in isolation.
3. `gh issue view <n> --json title -q .title` via `shell.Cmd`.
4. Return the trimmed title.

Returns `("", nil)` for all no-result cases above — graceful, never an error the
caller has to suppress.

### 3. `sesh status` command (`seshcli/status.go`)

```go
func NewStatusCommand(base *BaseDeps) *cobra.Command
```

- Build deps via `buildDeps`.
- Resolve path: `deps.Lister.GetAttachedTmuxSession()` (same approach as `sesh root`);
  fall back to `os.Getwd()` when not attached to a tmux session.
- `title, _ := deps.Github.IssueTitle(path)` — error intentionally ignored.
- `fmt.Print(title)`; return `nil`.

Registered in `seshcli/root_command.go` alongside the other commands.

### Dependency injection

`github` is config-free, so it joins `BaseDeps`:

- Add `Github github.Github` field to `BaseDeps`.
- Wire `github.NewGithub(sh, g)` in `NewBaseDeps()`.

## Testing

- `parseIssueNumber` — table test: plain number, prefixed/suffixed branch, no number,
  multiple number groups (takes the first).
- `github.IssueTitle` — against mocked `shell.Shell` + `git.Git`: success path,
  no-number branch (empty, no shell call), `gh` error (empty), not-a-repo (empty).
- Run `just test` (regenerates mocks, runs with coverage + race) before completion.

## README

Add a "tmux status bar" section under usage documenting the `status-left` snippet,
the `gh` install + auth requirement, and the current behavior (issue title matching
the branch number).

## Out of scope (deferred, tracked in #400)

- Pull request titles and draft/open/merged/closed state.
- GitHub issue open/closed state.
- Claude Code / Pi conversation titles and state.
- User configuration of priority/format.
