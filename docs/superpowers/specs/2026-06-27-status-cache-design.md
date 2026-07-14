# Cache GitHub issue title/state for `sesh status`

Follow-up to: [#400](https://github.com/joshmedeski/sesh/issues/400) / PR #401
Date: 2026-06-27

## Problem

`sesh status` runs `gh issue view` (a network round-trip, typically 200ms–1s+)
synchronously every time the tmux status line refreshes. tmux re-runs `#()`
status commands on the `status-interval` timer (default 15s; this user runs `3`)
and forces a status refresh on session switch/attach. The result is a visible
delay when switching sessions and repeated `gh` calls (≈20/min/session at
`status-interval 3`), which also risks GitHub API rate-limit throttling.

## Goal

Make `sesh status` read from a local cache so every render is a sub-millisecond
file read, while a detached background process keeps the cache fresh by calling
`gh`. The `gh` call never blocks a tmux redraw.

## Behavior overview

- `sesh status` resolves path → repo root + branch (all local, cheap), reads the
  cache entry **keyed on the branch**, and prints the formatted badge
  immediately. It never calls `gh` itself.
- Whenever the entry is missing or older than `issue_ttl`, `sesh status` spawns
  a detached `sesh status --refresh <path>` child that performs the live `gh`
  fetch and writes the cache, then the parent exits without waiting.
- `sesh connect` spawns the same detached refresh for the connected session's
  path, so the entry is warm before the first redraw after a switch.
- This is stale-while-revalidate: renders are always instant; data is at most
  `issue_ttl` seconds plus one refresh stale.

The cache is keyed on **repo root + branch**, not the issue number. A session
always knows its branch locally, and a branch maps to at most one PR — so
branch-keying is what the hot path can resolve without a network call, and it is
the foundation that makes future PR support (see *Forward compatibility* below)
purely additive.

## Configuration — new `[github]` section

Added to `sesh.toml`:

```toml
[github]
issue_ttl = 60   # seconds; default 60. 0 = disable caching (always fetch live)
```

In `model/config.go`, mirroring the existing nested-section pattern (`TUI
TUIConfig` / `DefaultSessionConfig`):

```go
type Config struct {
    // ... existing fields ...
    Github GithubConfig `toml:"github"`
}

type GithubConfig struct {
    IssueTTL *int `toml:"issue_ttl"` // pointer: distinguishes absent from explicit 0
}
```

`IssueTTL` is a `*int` to disambiguate three states that a plain `int` cannot,
because an absent section and an explicit `issue_ttl = 0` both decode to the
Go zero value:

- `nil` (section or key absent) → effective TTL **60** (default).
- `0` (explicit `issue_ttl = 0`) → cache **disabled**: `sesh status` skips the
  cache read and performs a live `gh` fetch inline (today's behavior), and
  `sesh connect` performs no warm.
- `N > 0` → effective TTL `N` seconds.

A single helper `config.Github.EffectiveTTL() int` (returns 60 for nil,
otherwise the value) centralizes this so command code never dereferences the
pointer directly. `EffectiveTTL() == 0` is the "disabled" check.
- `sesh.schema.json` gains a matching `github` object with `issue_ttl`
  (integer, default 60) so TOML editors autocomplete it. (See the
  `config-schema-sync` skill — schema must stay in lockstep with the struct.)

## Storage — new `statuscache` package

A new package `statuscache/` storing one gob file per cache key under
`$XDG_CACHE_HOME/sesh/status/` (falling back to `~/.cache/sesh/status/`),
mirroring `cache.FileCache`'s XDG resolution and atomic tmp+rename write.

```go
// Ref is one GitHub entity (issue or PR) as rendered in the status bar.
type Ref struct {
    Number int
    Title  string
    State  string // issue: OPEN|CLOSED ; PR: OPEN|DRAFT|MERGED|CLOSED
}

// Entry is the resolved status for one branch. Both pointers may be nil
// (a "negative" entry — branch has nothing to show — which still suppresses
// re-spawning a refresh until the TTL expires).
type Entry struct {
    PR        *Ref // this branch's pull request; nil if none. Unused in the MVP.
    Issue     *Ref // the linked issue; nil if none.
    Timestamp time.Time
}

type StatusCache interface {
    Read(key string) (Entry, bool, error) // bool=found (false on missing/corrupt)
    Write(key string, entry Entry) error
}

// Key builds the cache key (and on-disk filename stem) for a repo+branch.
func Key(repoRoot, branch string) string
```

- **Key:** `sha256(repoRoot + "\x00" + branch)` hex-encoded, used as the filename
  (`<hash>.gob`). The NUL separator avoids any ambiguity between repo root and
  branch. Keying on repo root + branch means each session resolves its own key
  locally, the same branch name across different repos/worktrees never collides,
  and multiple branches targeting the same issue each get their own entry — so
  the cache naturally holds the several PRs of one issue (one per branch).
- **One file per key** so concurrent refreshes from different sessions never
  write the same file — no locking, no last-writer-wins clobbering.
- **Negative entries:** the refresh always writes an `Entry` even when the branch
  has nothing to show (both pointers nil). Without this, a branch with no issue
  number (and, later, no PR) would miss on every render and re-spawn a refresh
  every tick; a negative entry makes the hot path respect the TTL instead.
- A read error (missing file, decode failure) returns `(Entry{}, false, nil)` —
  treated as a cache miss, never surfaced as an error to the status bar.
- The directory is created on first write (`os.MkdirAll`).

`repoRoot` is obtained from the existing `git` package (`ShowTopLevel` /
`GitCommonDir`) and `branch` from `git.CurrentBranch`, so the key is stable
across subdirectories of a repo.

## The refresh entrypoint

`sesh status` gains a hidden boolean flag `--refresh`. When set, the command:

1. Resolves path → repo root + branch + (optional) issue number via
   `github.Resolve`. If `path` is not a git repo, it writes nothing and returns.
2. If the branch has an issue number, calls `deps.Github.Issue(path)` (the live
   `gh` path) and, on success, sets `Entry.Issue`. On a gh failure or no number,
   `Entry.Issue` stays nil. (Future: also set `Entry.PR` from `gh pr view`.)
3. Writes the `statuscache` entry under `Key(repoRoot, branch)` — **always**,
   even when both refs are nil (a negative entry, per the storage rules), and
   prints nothing.

This is the only code path that invokes `gh`. It runs synchronously *within the
detached child*, so it can take as long as the network needs without affecting
any tmux redraw. Failures never propagate to the foreground status render.

### Spawning detached refreshes

A small injectable component spawns the refresh child detached from the parent
so it outlives the parent process (required because tmux waits for the
foreground `#(sesh status)` process to exit):

```go
type Refresher interface {
    Spawn(path string) error // launches `sesh status --refresh <path>` detached
}
```

Implementation uses `os.Executable()` for the binary path and
`exec.Command(self, "status", "--refresh", path)` with `SysProcAttr{Setsid:
true}` and no waiting (the parent does not call `Wait`). Stdout/stderr are
discarded. Spawn failures are logged via `slog` and otherwise ignored — a failed
spawn just means the cache stays stale until the next tick.

`Refresher` is wired into `BaseDeps` (config-free) so both the status command
and the connect command can use it. Behind the interface, tests assert a spawn
was requested without launching a process.

## Data flow

### `sesh status` (no `--refresh`)

```
ttl := config.Github.EffectiveTTL()   // 60 when unset
path := statusPath(deps)               // existing helper
if path == "" { return nil }

if ttl == 0 {                          // cache disabled
    issue, found, _ := deps.Github.Issue(path)   // live fetch (today's path)
    if found { print(formatStatus(issue)) }
    return nil
}

ref, ok := deps.Github.Resolve(path)   // local: {RepoRoot, Branch, Number, HasNumber}
if !ok { return nil }                  // not a git repo → nothing to show, no spawn

key := statuscache.Key(ref.RepoRoot, ref.Branch)
entry, found, _ := deps.StatusCache.Read(key)
if found {
    // PR-first, issue-fallback. PR is always nil in the MVP.
    switch {
    case entry.PR != nil:    print(formatStatus(toIssue("pr", *entry.PR)))
    case entry.Issue != nil: print(formatStatus(toIssue("issue", *entry.Issue)))
    // both nil: negative entry → print nothing
    }
}
if !found || time.Since(entry.Timestamp) > time.Duration(ttl)*time.Second {
    deps.Refresher.Spawn(path)         // detached; never blocks
}
return nil
```

`Resolve` succeeds (`ok == true`) whenever `path` is inside a git repo, so the
hot path can build the branch key even when the branch carries no issue number
(the refresh will write a negative entry, suppressing re-spawns until the TTL).

To avoid duplicating the repo-root/branch/number resolution that currently lives
inside `github.Issue`, the `github` package exposes it as a reusable method:

```go
type BranchRef struct {
    RepoRoot  string
    Branch    string
    Number    int  // issue number parsed from the branch; 0 if none
    HasNumber bool
}

// Resolve returns the repo root and branch for path, plus the issue number
// parsed from the branch name if present. ok is false only when path is not a
// git repo.
Resolve(path string) (ref BranchRef, ok bool)
```

`github.Issue` is refactored to call `Resolve` internally (then `gh issue view`
when `HasNumber`), so the parsing exists in exactly one place. `formatStatus` is
unchanged (the magenta `Issue #<n>` / future `PR #<n>` format); `toIssue(kind,
ref)` is a small `seshcli` helper that adapts a cache `Ref` to the
`github.Issue` value `formatStatus` consumes — keeping `statuscache` free of any
dependency on `github`. In the MVP `kind` is always `"issue"` and is ignored;
it is threaded through now so that adding PR support later means only extending
`formatStatus` to branch on `kind` (rendering `PR #<n>` with merged→purple),
with no change to the call sites.

### `sesh connect` (warm-on-switch)

After a successful connect/switch, and only when `issue_ttl != 0`:

```go
deps.Refresher.Spawn(connectedPath)
```

Fire-and-forget, added alongside the existing background
`CachingLister.RefreshCache` call. Never blocks the connect.

## Error handling

- gh/network/not-found failures during refresh: cache untouched, nothing
  printed. Identical to the current contract.
- Corrupt or unreadable cache file: treated as a miss.
- Spawn failure: logged at debug, ignored; cache self-heals on a later tick.
- `issue_ttl = 0`: cache bypassed entirely, restoring pre-cache live behavior.

## Testing

- **`statuscache`:** write→read round-trip (incl. a negative entry — both refs
  nil — which round-trips as `found=true`); missing file → `(_, false, nil)`;
  corrupt file → `(_, false, nil)`; `Key` is stable and collision-free for
  distinct repo/branch pairs; write creates the directory.
- **`github.Resolve`:** mocked `git.Git` → `{RepoRoot, Branch, Number,
  HasNumber=true}` for a numeric branch; `{…, HasNumber=false}` for a branch with
  no number (still `ok=true`); `ok=false` for a non-repo path.
  (`github.Issue` tests continue to pass via the refactor.)
- **`config.EffectiveTTL`:** nil → 60; explicit 0 → 0; explicit N → N.
- **`--refresh` path** (mocked `github.Github` + `statuscache`): a numeric branch
  with a successful `Issue` writes an `Entry` with `Issue` set; a gh failure or a
  no-number branch writes a negative `Entry` (both refs nil); a non-repo path
  writes nothing.
- **`sesh status` decision logic** (mocked `StatusCache` + `Refresher`):
  - fresh issue hit → prints, no spawn
  - fresh negative hit → prints nothing, no spawn
  - stale hit → prints, spawns
  - miss → prints nothing, spawns
  - `issue_ttl = 0` → no cache read, live `Github.Issue` used
- **`sesh connect`:** spawns a refresh for the connected path when caching is
  enabled; no spawn when `issue_ttl = 0`.
- Run `just test` (regenerates mocks, `-race`) before completion. New interfaces
  (`StatusCache`, `Refresher`) get mockery mocks.

## Forward compatibility: PR support

This design is shaped so that the future PR work from #400 (PR title + number +
draft/open/merged/closed state, with PRs taking priority over the issue, and one
issue tracking several PRs across branches) is **purely additive** — no key
change, no hot-path change, no cache migration:

- **Branch key** already models "branch → its one PR." Each branch's session
  reads its own entry; multiple PRs of one issue are simply multiple branch
  entries that share an `Issue.Number`. "All PRs for issue N" is reconstructable
  by scanning entries whose `Issue.Number == N` (not needed for rendering, but
  available).
- **`Entry`** already carries both `PR` and `Issue` refs; the MVP leaves `PR`
  nil. Adding PR support means the refresh additionally runs `gh pr view --json
  number,title,state,isDraft` for the branch and sets `Entry.PR`.
- **Render priority** (`PR != nil` first, else `Issue`) is already in the hot
  path. Only `formatStatus` needs extending to render `PR #<n>` and map PR
  states (e.g. merged → purple) — the `kind` is already threaded through
  `toIssue`.
- Because the cache is disposable, even unforeseen `Entry` additions cost
  nothing: old files read as misses and re-fetch.

What is **not** built now: the `gh pr view` call, the PR state→color mapping, and
the `PR #<n>` format branch. Those land with the PR feature.

## Out of scope

- Fetching pull-request or AI-conversation data (the storage model accommodates
  PRs per *Forward compatibility* above, but no PR fetching is implemented here).
- Cross-machine or shared cache.
- A manual `sesh status --clear-cache` command (cache self-expires; can be added
  later if needed).

## Expected benefit

| | Today | With cache (`issue_ttl = 60`) |
|---|---|---|
| Per-render cost | 1 `gh` network call (≈200ms–1s+) | 1 local gob read (sub-ms–few ms) |
| Session switch | Blocks on `gh` | Instant (warmed on connect) |
| `gh` calls/min/session (`status-interval 3`) | ≈20 | ≈1 |
| Staleness | Always live | ≤ `issue_ttl` + one refresh |
