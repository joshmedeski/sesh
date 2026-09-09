---
name: sesh
description: Use when starting work in a terminal session — opening a project or directory, connecting to a session by name, cloning a repo to work in, or running a command or coding agent somewhere other than the current shell. Sesh resolves names and directories to sessions across tmux, zoxide, config, and tmuxinator. Use instead of `tmux new-session`, `tmux new-window`, `tmux attach`, or `tmux send-keys`.
---

# Starting work with sesh

When work needs to start somewhere other than the shell you're in, call `sesh`, not `tmux`.

Reaching for `tmux new-session -s <name>` or `tmux new-window -t <session>` throws away what the
user set sesh up to do: resolve a rough name or a directory to the right session, start it from
its configured root with its startup commands, reuse what's already running instead of duplicating
it, and decide whether to move the user or leave them where they are.

If `sesh` is not on `PATH`, fall back to raw tmux. Otherwise use the commands below.

## Session or window?

This is the first decision, and the one most worth getting right.

**A session per unit of work.** That's the model — a project, a repo, a directory. If the work has
its own directory, it wants its own session:

```bash
sesh connect ~/c/sesh
sesh connect sesh          # a name works too; sesh resolves it
```

**A window when it belongs inside a context that already exists.** Reach for a window when the
work is *part of* a session that's already running rather than a thing of its own:

```bash
sesh window connect claude -t 'second brain' -c 'claude "<prompt>"'
```

Good reasons to add a window instead of a session:

- The user named a session to put it in ("in my dotfiles session", "next to what I'm working on")
- The work is a second track in the same project — tests, logs, an agent, a dev server — and the
  session for that project is already open
- You're already inside the right session and just need another place to run something

If neither applies, use a session. Don't create a window in an unrelated session just because
that session happens to be attached.

Either way, both commands **select-or-create**: they reuse what's there and only create on a miss,
so it's safe to run them without checking first.

## Sessions

```bash
sesh connect <name|directory>
```

Connects to the session, creating it if it isn't running. The argument can be a path, a session
name, a zoxide result, a configured session, or a tmuxinator config — sesh resolves across all of
them, which is the main reason to call it instead of tmux.

```bash
sesh connect ~/c/sesh              # by path
sesh connect sesh                  # by rough name
sesh connect 'second brain' -s     # names with spaces need quoting
sesh connect ~/c/api -c 'npm run dev'
```

| Flag | Meaning |
|---|---|
| `-s`, `--switch` | Switch the client rather than attach. Use when triggering sesh from outside the terminal — a hotkey, a script, another app. |
| `-c`, `--command` | Startup command. **Only runs on a cold start** — silently ignored if the session already exists. |
| `-T`, `--tmuxinator` | Start via tmuxinator if the session isn't running. |

That `-c` caveat is the sharp edge: if the user might already have the session open, `sesh connect
-c` may do nothing at all. To run a command in a session that may already be running, add a window
(next section).

Two related entry points, for when the directory doesn't exist yet:

```bash
sesh mkdir ~/c/new-project -c '<command>'   # create the directory, then connect to it
sesh clone <repo>                           # clone a git repo and connect to it
```

Use `sesh mkdir` rather than `mkdir && sesh connect` — it's one step and doesn't wait on zoxide.

## Windows

```bash
sesh window connect <name|directory> -t <session> -c '<command>'
```

Selects the window called `<name>` in `<session>`, creating it if it isn't there. Prints the
window's `session:index` target.

```bash
sesh window connect tests -c 'npm test -- --watch'   # in the current session
sesh window connect ~/c/sesh -t dotfiles             # window "sesh", rooted at that path
```

A directory argument is named after its basename and rooted there, so running it twice selects the
window the first run created rather than duplicating it.

### `--command` means different things on create and reuse

Identity is the window name:

- **Created** window — the command becomes the window's process. The window closes when it exits.
- **Existing** window — the command is typed in as keystrokes, into whatever is already running.
  If that's a shell, it runs. If it's a long-running process, you just typed at that process.

So launching an agent twice with the same window name starts one agent, then types the second
prompt into the first agent's input. When each invocation must get its own window, pass `--new`:

```bash
sesh window connect claude -t 'second brain' --new -c 'claude "<prompt>"'
```

`--new` always creates, letting windows share a name; a later connect without it reuses the
lowest-indexed one. Omitting the name also always creates — there's nothing to match on.

### Targeting

`-t <session>` resolves the same way `sesh connect` does and **starts the session if it isn't
running**, without moving the user's client. So you don't have to check whether the session is up,
and cold and warm targets take the same command.

With no `-t`, the target is the attached session. Outside tmux that fails with
`Not inside a tmux session, use --target to specify one.` — pass `-t`.

## Focus: don't yank the user out of their work

Both commands move the user to what they connected to. That's right when they asked for it and
wrong when it's a background step.

| Situation | Flag |
|---|---|
| The user asked for this and wants to land there | *(default)* |
| Background work; leave them where they are | `-b` (windows only) |
| Triggered from outside tmux — hotkey, script, another app | `-s` |

`-b` / `--background` creates the window without selecting it or moving the client. Prefer it for
anything the user didn't explicitly ask to be taken to. Sessions have no `-b`; a window connect
with `-t` is the way to set something up without moving.

## Finding the right target

Don't guess a session name and don't invent one:

```bash
sesh list -j                       # every source: tmux, zoxide, config, tmuxinator
sesh list -t                       # running tmux sessions only
sesh window list -t <session>      # windows in a session (-j for json)
```

When the user names a project or directory, pass what they said — sesh resolves it. Session names
contain spaces (`second brain`) more often than you'd expect, so quote arguments.

## Do not

- **Do not** use `tmux new-session`, `tmux attach`, `tmux new-window`, or `tmux split-window` to
  start work. Use `sesh connect` or `sesh window connect`.
- **Do not** use `tmux send-keys` to run a command. Use `-c`, which knows whether the target is
  new or already running.
- **Do not** create a second session for a directory or name that already has one — both commands
  already select-or-create.
- **Do not** put unrelated work in a window of whatever session happens to be attached. New work
  with its own directory gets its own session.
- **Do not** guess session names. Run `sesh list`.
- **Do not** move the user for background work. Use `-b`.

Panes are outside sesh's scope; `tmux split-window` and `tmux send-keys` remain correct there.

## Command reference

```
sesh connect <name|directory>          Connect to a session, creating it if it doesn't exist
  -s, --switch             Switch the client rather than attach
  -c, --command <cmd>      Startup command; ignored if the session already exists
  -T, --tmuxinator         Start via tmuxinator if the session isn't running

sesh mkdir <path>                      Create a directory and connect to it as a session
  -c, --command <cmd>      Startup command
  -s, --switch             Switch rather than attach

sesh clone <repo>                      Clone a git repo and connect to it as a session
  -d, --dir <name>         Directory git creates
  -c, --cmdDir <dir>       Directory to run git in

sesh window connect [name|directory]   Select a window, creating it if it doesn't exist
  -t, --target <session>   Target session (default: current attached session)
  -c, --command <cmd>      Command to run in the target
  -p, --path <dir>         Start directory (default: the target session's root)
  -b, --background         Create without selecting it or moving the client
  -s, --switch             Switch the client rather than attach
      --new                Always create, even when the name matches an existing window

sesh window list                       List the windows of a session
  -t, --target <session>   Target session (default: current attached session)
  -j, --json               Output as json

sesh list                              List sessions from all sources
  -j, --json               Output as json
  -t, --tmux               tmux sessions only
  -z, --zoxide             zoxide results only

sesh last                              Connect to the last tmux session
sesh root                              Print the root directory of the active session
```
