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

- `sesh status` resolves path → branch → issue number (all local, cheap), then
  reads the cache entry for that issue and prints the formatted badge
  immediately. It never calls `gh` itself.
- Whenever the entry is missing or older than `issue_ttl`, `sesh status` spawns
  a detached `sesh status --refresh <path>` child that performs the live `gh`
  fetch and writes the cache, then the parent exits without waiting.
- `sesh connect` spawns the same detached refresh for the connected session's
  path, so the entry is warm before the first redraw after a switch.
- This is stale-while-revalidate: renders are always instant; data is at most
  `issue_ttl` seconds plus one refresh stale.

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
type Entry struct {
    Number    int
    Title     string
    State     string
    Timestamp time.Time
}

type StatusCache interface {
    Read(key string) (Entry, bool, error) // bool=found (false on missing/corrupt)
    Write(key string, entry Entry) error
}

// Key builds the cache key (and on-disk filename stem) for a repo+issue.
func Key(repoRoot string, number int) string
```

- **Key:** `sha256(repoRoot + "#" + strconv.Itoa(number))` hex-encoded, used as
  the filename (`<hash>.gob`). Keying on repo root + issue number (not branch)
  means two branches pointing at the same issue share an entry and the same
  branch name across different repos/worktrees never collides.
- **One file per key** so concurrent refreshes from different sessions never
  write the same file — no locking, no last-writer-wins clobbering.
- A read error (missing file, decode failure) returns `(Entry{}, false, nil)` —
  treated as a cache miss, never surfaced as an error to the status bar.
- The directory is created on first write (`os.MkdirAll`).

`repoRoot` is obtained from the existing `git` package (`ShowTopLevel` /
`GitCommonDir`), so the key is stable across subdirectories of a repo.

## The refresh entrypoint

`sesh status` gains a hidden boolean flag `--refresh`. When set, the command:

1. Resolves path → branch → issue number (reusing the existing logic).
2. Calls `deps.Github.Issue(path)` (the live `gh` path).
3. On success, writes the `statuscache` entry for the key and prints nothing.
4. On any failure (no number, gh missing/unauthenticated, not found, etc.),
   writes nothing and prints nothing — the existing graceful-empty contract.

This is the only code path that invokes `gh`. It runs synchronously *within the
detached child*, so it can take as long as the network needs without affecting
any tmux redraw.

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

ref, ok := deps.Github.Resolve(path)   // local: {RepoRoot, Number}; ok=false if none
if !ok { return nil }                  // nothing to show

key := statuscache.Key(ref.RepoRoot, ref.Number)
entry, found, _ := deps.StatusCache.Read(key)
if found {
    print(formatStatus(github.Issue{
        Number: entry.Number, Title: entry.Title, State: entry.State,
    }))
}
if !found || time.Since(entry.Timestamp) > time.Duration(ttl)*time.Second {
    deps.Refresher.Spawn(path)         // detached; never blocks
}
return nil
```

To avoid duplicating the branch→number→repoRoot parsing that currently lives
inside `github.Issue`, the `github` package exposes it as a reusable method:

```go
type IssueRef struct {
    RepoRoot string
    Number   int
}

// Resolve returns the repo root and issue number for the branch at path.
// ok is false when path is not a repo or the branch has no issue number.
Resolve(path string) (ref IssueRef, ok bool)
```

`github.Issue` is refactored to call `Resolve` internally (then `gh issue
view`), so the parsing exists in exactly one place. `statuscache.Key(repoRoot
string, number int) string` builds the hashed key. `formatStatus` is unchanged
(the magenta `Issue #<n>` format from PR #401); the cache `Entry` is converted to
a `github.Issue` inline in the `seshcli` command (which already imports both
packages), keeping `statuscache` free of any dependency on `github`.

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

- **`statuscache`:** write→read round-trip; missing file → `(_, false, nil)`;
  corrupt file → `(_, false, nil)`; key hashing is stable and collision-free for
  distinct repo/number pairs; write creates the directory.
- **`github.Resolve`:** mocked `git.Git` → returns `{RepoRoot, Number}` for a
  numeric branch; `ok=false` for a non-repo path or a branch with no number.
  (`github.Issue` tests continue to pass via the refactor.)
- **`config.EffectiveTTL`:** nil → 60; explicit 0 → 0; explicit N → N.
- **`--refresh` path:** with mocked `github.Github` + `statuscache`, a successful
  `Issue` writes the expected `Entry`; a not-found/error `Issue` writes nothing.
- **`sesh status` decision logic** (mocked `StatusCache` + `Refresher`):
  - fresh hit → prints, no spawn
  - stale hit → prints, spawns
  - miss → prints nothing, spawns
  - `issue_ttl = 0` → no cache read, live `Github.Issue` used
- **`sesh connect`:** spawns a refresh for the connected path when caching is
  enabled; no spawn when `issue_ttl = 0`.
- Run `just test` (regenerates mocks, `-race`) before completion. New interfaces
  (`StatusCache`, `Refresher`) get mockery mocks.

## Out of scope

- Caching pull-request data or AI-conversation state (those features don't exist
  yet — tracked in #400).
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
