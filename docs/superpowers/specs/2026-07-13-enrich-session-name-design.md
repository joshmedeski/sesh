# Design: enrich tmux session name with issue title

**Date:** 2026-07-13
**Branch:** `400-tmux-status-bar`
**Status:** proposed

## Summary

Replace the tmux status-bar feature with a command that folds a GitHub issue
title into the **tmux session name itself**. A session is created fast with the
normal namer name, then a tmux hook asynchronously enriches the name with the
issue title, e.g.:

```
400-tmux-status-bar  ──►  400-tmux-status-bar — warm the status cache
```

## Motivation

The current branch renders issue context into the tmux **status bar** via a
custom tmux variable, backed by a status cache and a background refresher. That
machinery only exists to keep a per-tick status bar from blocking on `gh`.

Putting the context in the **session name** instead makes it visible everywhere
tmux shows session names (session list, choose-tree, `sesh` picker) without any
status-bar configuration, and collapses the feature to a single command.

## Decisions

| Decision | Choice |
| --- | --- |
| Identity model | Prefix stays matchable — one session, name carries display |
| Name format | `base — <issue title>` (space, em dash U+2014, space) |
| Title handling | Near-original: keep spaces and casing; strip `:` and `.` (illegal in tmux names) |
| Trigger | Opt-in tmux `session-created` hook (documented) |
| Scope | Replace the status-bar feature; remove cache + refresher |
| Issue vs PR | Issue only, as today (`github.Issue` — number parsed from branch) |
| Command | `sesh rename --enrich [name]` |
| No-issue case | Force the name back to the bare base (self-healing) |

## Behavior

### Phase 1 — create (unchanged)

`sesh connect` creates the session with the namer name (`base`). Instant, no
network.

### Phase 2 — enrich (new)

A tmux `session-created` hook runs `sesh rename --enrich` in the background. It
resolves the branch's issue via `gh` and renames the live session to
`base — <title>`. With no issue, the name is forced back to the bare `base`.

Non-blocking behavior comes from the hook (`run-shell -b`); the command itself
is a simple synchronous single-`gh`-call path.

## The `sesh rename --enrich [name]` command

Lives in `seshcli/rename.go` (replacing `seshcli/status.go`).

1. **Resolve target.** If `name` arg is given, use it. Otherwise use
   `Lister.GetAttachedTmuxSession()`, which yields the current session name
   **and** its path.
2. **Recompute the base from the path** via `Namer.Name(session.Path)` — not by
   splitting the current name. The path→name mapping is deterministic, which
   makes re-runs fully idempotent: whatever the session is currently called, the
   base is recomputed and a fresh suffix appended. You can never accumulate
   `base — title — title`.
3. **Resolve issue** via `Github.Issue(session.Path)` (existing, issue-only).
   - Found → `newName = base + " — " + sanitizeTitle(issue.Title)`
   - Not found → `newName = base`
4. **Rename** only if `newName != currentName`, via `Tmux.RenameSession(current, newName)`.

### `sanitizeTitle`

Keep spaces and original casing; neutralize the two characters tmux forbids in
session names by replacing each with a space:

- `.` → space, `:` → space
- collapse resulting double spaces; trim

(Distinct from `namer/convert.go:convertToValidName`, which additionally
lowercases-agnostic collapses whitespace to `_`. Add `sanitizeTitle` alongside
it in the `namer` package, or a small local helper — decide in the plan.)

## Reconnection matching (key integration point)

**Problem.** sesh dedupes sessions by tmux session name. `dirStrategy` /
`zoxideStrategy` compute the namer `base`, return `New: true`, and
`connectToTmux` calls `tmux.NewSession(base, path)` — a no-op when a session
named `base` exists, then `SwitchOrAttach(base)` reattaches. Once the session is
renamed to `base — title`, `NewSession("base", …)` no longer sees a match and
creates a **duplicate**.

**Fix.** Add a base-aware lookup in `lister`: a tmux session matches a base name
when its name `== base` **or** starts with `base + " — "`. `dirStrategy` and
`zoxideStrategy`, after computing the namer name, consult it; on a hit they
return that real session with `New: false` so `connectToTmux` reattaches instead
of creating a new session.

- Exact-match `FindTmuxSession(name)` stays for the connect-by-name path
  (picking `base — title` straight from `sesh list` still exact-matches).
- The ` — ` boundary is unambiguous even when a base name itself contains `/`
  (multi-segment dir names, `namer/dir.go:41`).
- **Caveat (out of scope):** two distinct repos that produce the *identical*
  base name are already ambiguous today; this design does not change that.

## Removals (replace the status-bar feature)

- Delete `seshcli/status.go` and the `status` command (incl. hidden `--refresh`).
- Delete the `statuscache` package and the refresher — no per-tick polling
  remains to justify them.
- Remove `StatusCache` and `Refresher` from `BaseDeps` / `Deps` wiring in
  `seshcli/deps.go`.
- Keep `github.Issue` and `github.Resolve` (still used by the rename command).

## New tmux primitive

Add to the `Tmux` interface + `RealTmux` (+ regenerated mock):

```go
RenameSession(target, newName string) (string, error)
// tmux rename-session -t <target> <newName>
```

Spaces in `newName` are safe: `shell.Cmd` uses `exec.Command(bin, args...)`
(`shell/shell.go:35`), passing an argv vector with no shell-string
interpolation.

## Documentation

Ship the opt-in hook snippet (README / docs):

```tmux
set-hook -g session-created 'run-shell -b "sesh rename --enrich"'
```

Note that it also covers sessions created outside sesh, and that it is a no-op
when the branch has no resolvable issue.

## Touch map

- **New:** `seshcli/rename.go`; `Tmux.RenameSession`; lister base-match helper;
  `sanitizeTitle`.
- **Changed:** `connector/dir.go`, `connector/zoxide.go` (consult base-match);
  `seshcli/deps.go` (drop wiring + register `rename`, unregister `status`);
  tmux interface + mock; README/docs.
- **Deleted:** `seshcli/status.go`; `statuscache/`; refresher.

## Testing

- `sanitizeTitle`: replaces `:`/`.` with space, preserves other spaces + casing, collapses/trims.
- rename command: found issue → `base — title`; no issue → bare `base`;
  idempotent on re-run (already-enriched name → recomputed identical name, no
  spurious rename); arg vs attached-session resolution.
- base-match helper: exact, prefix (`base — …`), non-match, base-containing-`/`.
- `dirStrategy` / `zoxideStrategy`: existing enriched session → `New: false`
  reattach (no duplicate).
- `Tmux.RenameSession`: correct argv, spaces preserved.

## Resolved

- **No-issue case:** force the name back to the bare base (self-healing).
- **`sanitizeTitle` on `:`/`.`:** replace each with a space, then collapse/trim.
