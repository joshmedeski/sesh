# Worktree Support Design

**Issue:** [#409](https://github.com/joshmedeski/sesh/issues/409)
**Date:** 2026-07-13

## Goal

Add native git worktree support to sesh so a single command creates and connects to a
per-issue/PR worktree, driven by declarative per-repo configuration in `sesh.toml`. This
ports the logic from a personal `browser-dispatch` bash script into sesh, using a
`[[worktree]]` config mapping inspired by [treekanga](https://github.com/garrettkrohn/treekanga).

Headline command:

```
sesh worktree create 2345
```

Resolves the current repo's config, creates a worktree for issue/PR `2345`, connects to it
as a tmux session, and runs the configured startup command.

## Config

New `[[worktree]]` array-of-tables section, mirroring the existing `[[session]]` /
`[[wildcard]]` blocks in `model/config.go`.

```toml
[[worktree]]
name = "nutiliti/nutiliti"          # org/repo — used for gh + --repo matching
path = "~/c/nu"                      # local repo root (main worktree)
worktree_dir = "w"                   # where worktrees go; relative to path, or absolute (default ".wk")
branch_template = "jam/{number}-1"   # {number} token; prefix baked in (default "{number}")
base_branch = "origin/main"          # branch new worktrees off this (default "origin/main")
fetch = true                         # git fetch before creating (default true)
startup_command = "nu_setup"         # runs on connect (reuses existing startup mechanism)
```

Field defaults: `worktree_dir` → `.wk`, `branch_template` → `{number}`,
`base_branch` → `origin/main`, `fetch` → `true`. `path` supports `~` expansion.

### Go model

```go
WorktreeConfig struct {
	Name           string `toml:"name"`            // org/repo
	Path           string `toml:"path"`            // local repo root
	WorktreeDir    string `toml:"worktree_dir"`    // default ".wk"
	BranchTemplate string `toml:"branch_template"` // default "{number}"
	BaseBranch     string `toml:"base_branch"`     // default "origin/main"
	Fetch          *bool  `toml:"fetch"`           // default true (pointer to distinguish unset)
	StartupCommand string `toml:"startup_command"`
}
```

Added to `Config` as `WorktreeConfigs []WorktreeConfig \`toml:"worktree"\``. Mirrored into
`sesh.schema.json` per the `config-schema-sync` skill.

## Command

```
sesh worktree create <number> [--repo org/repo] [--pr] [--switch]
```

- **Repo resolution:** default = detect the git root of cwd via
  `git rev-parse --git-common-dir` (works from inside an existing worktree too), then match
  to a `[[worktree]]` entry by expanded `path`. `--repo org/repo` overrides by `name` (used
  by the browser dispatcher, which has no cwd context).
- **Issue vs PR:** auto-detect — try `gh pr view <n>`; if it resolves as a PR, run the PR
  path; otherwise treat as an issue. `--pr` forces the PR path and skips the probe.
- **Branch/dir:** worktree lands at `<worktree_dir>/<number>`; branch name from
  `branch_template`. If the target dir already exists, reconnect instead of recreating.
- **Connect:** derive the new worktree path and call the existing connector with the
  configured `startup_command`, so worktree sessions are named `repo/branch` and behave like
  any other sesh session.

### PR path (full dispatch parity)

Ported from `browser-dispatch`:

1. `gh pr view <n>` → author login; `gh api user` → current login.
2. **Foreign PR** (author ≠ me): `git worktree add --detach <target> <base>` then
   `gh pr checkout <n>` inside it.
3. **Own PR:** resolve the closing-issue reference
   (`gh pr view --json closingIssuesReferences`), fall back to scanning title/body for
   `#N`. If an issue is found, create the worktree keyed on that issue number; if none,
   fall back to a PR-branch checkout keyed on the PR number.
4. Existing target dir + PR path → `gh pr checkout` + `git pull --ff-only` before connecting.

## Architecture

Follows the existing `cloner` pattern (external git op → derive path → `Connect`).

| Layer | File | Responsibility |
|---|---|---|
| Config | `model/config.go`, `sesh.schema.json` | `[[worktree]]` model + schema |
| Git wrapper | `git/git.go` | `WorktreeAdd`, `Fetch` (+ regen mock) |
| GitHub wrapper | `github/github.go` | `PrView`, `PrCheckout`, `CurrentUser` (+ regen mock) |
| Service | `worktree/worktree.go` | `Create(opts)` — orchestrates resolve → gh → git → connect |
| Command | `seshcli/worktree.go` | cobra `worktree create` |
| Wiring | `seshcli/deps.go`, `seshcli/root_command.go` | build + register the service/command |
| Docs | `README.md` | `[[worktree]]` + `sesh worktree` docs |

The `worktree` service depends on `Config`, `Git`, `Github`, `Connector`, and the path/os
wrappers — mirroring how `cloner.NewCloner(config, git)` is wired in `BuildAll`.

## Data flow

```
sesh worktree create <n> [--repo] [--pr]
  → resolve WorktreeConfig (cwd git-common-dir → path match, or --repo → name match)
  → (optional) git fetch
  → determine kind (issue | PR) via gh probe or --pr
  → compute target = <worktree_dir>/<n>, branch = render(branch_template, n)
  → if target exists: (PR path may `gh pr checkout` + pull) → Connect
    else: git worktree add [-b branch --no-track base | --detach base + gh pr checkout] → Connect
  → connector.Connect(target, ConnectOpts{Switch, Command: startup_command})
```

## Error handling

- No `[[worktree]]` entry matches the resolved repo → clear error naming the repo path and
  pointing at config.
- `git`/`gh` failures surface the underlying command output (via the existing `shell`
  wrappers), consistent with `cloner`/`github`.
- Missing `gh` auth or a non-existent number → propagate gh's error message.

## Testing

- `worktree` service: table-driven tests with mocked `Git`, `Github`, `Connector`, `Config`
  covering: cwd resolution, `--repo` override, issue path, own-PR-with-issue,
  own-PR-no-issue fallback, foreign-PR detach, existing-dir reconnect, branch-template
  rendering, default fill-in.
- `git.WorktreeAdd`/`Fetch` and `github.Pr*`: shell-mocked tests following existing
  `git_test.go` / `github_test.go` patterns.
- Regenerate mocks via `just mock`; run `just test`.

## Scope

**In (v1):** everything above — full dispatch parity for the worktree/gh logic, plus the
three documented example configs (nutiliti, joshmedeski.com, sesh).

**Out (v1):** `sesh worktree list` / `remove` lifecycle commands; macOS/Helium browser
activation (stays in the thin dispatcher wrapper that calls `sesh worktree create --repo … <n>`).

## Housekeeping

Implement on a fresh branch off `main` (current branch `400-tmux-status-bar` has unrelated
uncommitted edits to `github/github.go`).
