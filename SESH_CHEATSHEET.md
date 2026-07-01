# Sesh Codebase Explained

## What is Sesh?

Sesh is a terminal session manager. It helps you create, find, and jump between tmux sessions quickly. It integrates with **zoxide** (a smart directory jumper) and **tmuxinator** (a tmux session launcher).

---

## Repo Layout — The Big Picture

```
sesh/
├── main.go               # App startup
├── seshcli/               # CLI commands + wiring
├── model/                 # Data types (session, config, etc.)
├── lister/                # Finds & lists sessions
├── connector/             # Connects to sessions
├── namer/                 # Generates session names
├── picker/                # Interactive session chooser (TUI)
├── previewer/             # Shows session previews (for fzf)
├── dashboard/             # Full-screen dashboard (TUI)
├── tmux/                  # Talks to tmux
├── zoxide/                # Talks to zoxide
├── tmuxinator/            # Talks to tmuxinator
├── git/                   # Talks to git
├── shell/                 # Runs shell commands
├── configurator/          # Reads sesh.toml config file
├── startup/               # Runs startup commands for new sessions
├── icon/                  # Nerd font icons for session types
├── cloner/                # Git clone + auto-connect
├── dir/                   # Checks if paths are directories / git repos
├── home/                  # Handles ~/ in paths
├── ls/                    # Lists directory contents
├── cache/                 # Caches session lists to disk
├── execwrap/              # Wraps os/exec (for testing)
├── oswrap/                # Wraps os (for testing)
├── pathwrap/              # Wraps path/filepath (for testing)
├── runtimewrap/           # Wraps runtime (for testing)
├── json/                  # JSON encoding for sessions
├── replacer/              # Replaces text using Aho-Corasick
└── convert/               # Parses strings to ints, bools, times
```

---

## `main.go` — The Starting Point

**What it does:** Sets up logging, then starts the app.

1. `init()` creates a log file at `/tmp/.seshtmp/YYYY-MM-DD.log` and sets up structured JSON logging. If `ENV=debug` is set, logs go to both the file and the terminal. Otherwise they go to the file only (so logging doesn't clutter your terminal).
2. `main()` creates the root `sesh` command and runs it via `fang.Execute()` — a Charm framework that handles CLI execution.

---

## `seshcli/` — All the CLI Commands

Each file in this package is one subcommand that a user can run. They're built using the [cobra](https://github.com/spf13/cobra) library.

### `root_command.go` — The `sesh` command itself

Creates the root command with the `--config` flag and registers all subcommands: `list`, `last`, `connect`, `clone`, `root`, `preview`, `picker`, `window`, `dashboard`.

### `list.go` — `sesh list` (alias: `l`)

**What it does:** Lists all sessions from whatever sources you choose (tmux, zoxide, config files, tmuxinator). Supports JSON output for scripting. This is what powers the picker and other tools.

Key flags:
- `--tmux` / `--zoxide` / `--config` / `--tmuxinator` — which sources to show
- `--json` — output as JSON
- `--icons` — show nerd font icons
- `--hide-attached` — skip the session you're currently in
- `--panes` — show individual tmux panes instead of sessions

### `connect.go` — `sesh connect <name>` (alias: `cn`)

**What it does:** Attaches to a session by name. If the session doesn't exist, it tries to create it (from a directory, zoxide result, etc.).

Key flags:
- `--switch` — use `tmux switch-client` instead of `attach-session` (for use from outside tmux)
- `--command` — run a command when creating a new session
- `--root` — resolve to the git root directory first

### `picker.go` — `sesh picker` (alias: `pick`, `pk`)

**What it does:** Opens an interactive TUI that lets you type to filter sessions and select one. Uses fuzzy matching. This is the main way users interact with sesh.

Key flags:
- Same source flags as `list` (`--tmux`, `--zoxide`, etc.)
- `--prompt` — change the prompt text
- `--separator-aware` — match `-`, `_`, `/`, `\` as spaces in searches

### `dashboard.go` — `sesh dashboard` (alias: `dash`, `d`)

**What it does:** Opens a full-screen TUI showing your tmux sessions grouped by categories you define in config. Shows git branches, window counts, and attached status. Requires being inside tmux.

No flags. All configuration is in `sesh.toml`.

### `root.go` — `sesh root` (alias: `r`)

**What it does:** Prints the root directory name of your current tmux session (e.g., the git repo name if you're in one).

### `last.go` — `sesh last` (alias: `L`)

**What it does:** Switches to your second-most-recently-used tmux session (the one you were just in).

### `preview.go` — `sesh preview <name>` (alias: `p`)

**What it does:** Shows a preview of a session or directory. Used by fzf as a preview command when you're selecting sessions. Shows the terminal output of a tmux session, or a file listing of a directory.

### `window.go` — `sesh window [name]` (alias: `w`)

**What it does:** Manages tmux windows.
- With **no arguments** — lists windows in the current tmux session
- With a **window name** — switches to that window
- With a **directory path** — creates a new window in that directory

### `clone.go` — `sesh clone <url>` (alias: `cl`)

**What it does:** Clones a git repo, then immediately connects to it as a new tmux session. One command to go from "I want to work on this project" to "I'm inside a tmux session for it."

### `deps.go` — The Wiring (Dependency Injection)

**What it does:** Builds all the objects the app needs, in order. Think of it like a factory that assembles everything.

**Layer 1 (no config needed):** Creates wrappers for `os`, `exec`, `path`, `git`, `shell`, etc. These are all immediately available.

**Layer 2 (needs config):** Reads `sesh.toml`, then creates:
- `Tmux` — uses the configured tmux command
- `Lister` — combines tmux + zoxide + config sessions
- `Namer` — names sessions (git repo name, directory name, etc.)
- `Connector` — connects to sessions
- `Startup` — runs startup commands
- `Icon` — adds icons to names
- `Previewer` — generates previews
- `Picker` — the interactive picker
- `Cloner` — git clone + connect

---

## `model/` — The Data Types

This package defines all the data structures that flow through the app.

### `config.go`

**What it does:** Holds your `sesh.toml` settings. Key configs:
- `Cache` — whether to cache session lists
- `TmuxCommand` — custom tmux binary path
- `TUI` — picker appearance (prompt, placeholder, icons)
- `SortOrder` — which session sources to show first
- `SessionConfigs` — named sessions you define
- `DefaultSessionConfig` — default settings for all sessions
- `Dashboard` — dashboard sections and groups

### `sesh_session.go`

**What it does:** The core session type. Every session has:
- `Src` — where it came from (`tmux`, `config`, `zoxide`, `tmuxinator`)
- `Name` — the display name
- `Path` — the directory path
- `Attached` — how many clients are attached
- `Windows` — how many windows the session has
- `Branch` — the git branch (used by the dashboard)
- `Score` — priority from zoxide

### `connect_opts.go` & `connection.go`

**What they do:** `ConnectOpts` is what you tell the connector when connecting (switch mode, run a command, use tmuxinator). `Connection` is what the connector returns (did it find a session? was it new? should it be added to zoxide?).

### `tmux_session.go`, `tmux_window.go`, `tmux_pane.go`

**What they do:** Raw data from tmux — exactly what `tmux list-sessions -F` returns. Created date, last attached time, activity time, group info, path, name, client count, etc.

### `zoxide_result.go`

**What it does:** A zoxide database entry — just a path and a score (how often you visit it).

### `tmuxinator_config.go`

**What it does:** A tmuxinator project — just a name.

---

## `lister/` — Finding Sessions

**What it does:** This is the "search engine" of sesh. It finds sessions from every possible source.

**How it works when you run `sesh list`:**
1. Determines which sources to check (tmux? config? zoxide? tmuxinator? panes?)
2. Runs them all **at the same time** (concurrently)
3. Collects their results
4. Applies filters (blacklist, hide-attached, hide-duplicates)
5. Returns the final list

Key functions:
- `List()` — the main function. Runs source strategies in parallel.
- `listTmux()` — runs `tmux list-sessions` and parses the output
- `listZoxide()` — runs `zoxide query --list --score` and parses the output
- `listConfig()` — iterates over the sessions defined in `sesh.toml`
- `listTmuxinator()` — runs `tmuxinator list -n` and parses the output
- `listTmuxPanes()` — runs `tmux list-panes -s` and parses the output
- `applyDedup()` — when the same session appears from multiple sources (e.g., tmux and zoxide), keeps only the best one. Priority: tmux > config > zoxide.
- `isBlacklisted()` — checks if a session name matches a regex pattern from the blacklist config

**Caching:** If `Cache: true` in config, results are cached to disk with a 5-second TTL. Reading is instant; stale caches are refreshed in the background.

---

## `connector/` — Connecting to Sessions

**What it does:** Takes a session name and figures out how to connect to it. It tries each strategy in order:

1. **Pane** — if the name looks like `windowName/paneTitle` and you're in tmux, select that pane
2. **Tmux** — if an existing tmux session has this name, attach to it. If not, create one.
3. **Tmuxinator** — if a tmuxinator project has this name, start it
4. **Config** — if a config session has this name, create a tmux session from it
5. **Wildcard** — if the path matches a wildcard pattern in config, use that
6. **Directory** — if the name is a directory path, create a session there
7. **Zoxide** — look it up in zoxide and create a session

Each strategy either succeeds (returns a connection) or says "not me, try the next one."

---

## `namer/` — Naming Sessions

**What it does:** Given a directory path, figures out what to call the session.

**Strategy chain (tries in order):**
1. **Git name** — if it's inside a git repo, use `<repo-name>/<relative-path>`. Handles worktrees and bare repos.
2. **Directory name** — use the last part of the path, or join multiple levels based on `DirLength` config.

Also has `RootName()` — collapses nested directories to their git root for the `sesh root` command.

---

## `picker/` — Interactive TUI

**What it does:** A bubbletea-based interactive session picker. You type to filter, the list narrows, you press enter to select.

**How it works:**
1. Fetches sessions asynchronously (shows "Loading..." until ready)
2. Filters via fuzzy matching (`sahilm/fuzzy` library)
3. Shows nerd font icons colored by source type
4. Supports `separator-aware` mode where `-`, `_`, `/`, `\` match as spaces

**Controls:** Arrow keys / ctrl+n/p to navigate, enter to select, esc to cancel, ctrl+d/u to page.

---

## `previewer/` — Session Previews

**What it does:** Generates a text preview of a session or directory. Used by fzf when you're selecting sessions.

**Strategy chain:**
1. If it's a tmux session, capture its current terminal output (`tmux capture-pane`)
2. If it's a config session with a `preview_command`, run that
3. If it's a config session without a preview command, list its directory
4. If it's any directory, list its contents

---

## `dashboard/` — Full-Screen Dashboard

**What it does:** A full-screen TUI built with bubbletea that shows all your tmux sessions in an organized, visually-rich view.

**How it's built:**
- Uses a **pluggable section system** — new widget types can be added by registering factories
- Currently has one section type: `sessions` — shows tmux sessions grouped by config-defined path patterns
- Sessions are loaded asynchronously; git branches are fetched in the background

**How the Sessions Section works:**
1. **Load** — fetches all tmux sessions via `Lister.List(Tmux: true)`
2. **Group** — matches each session's path against config patterns (e.g., `~/work/*`). Unmatched sessions go into an "Other" group.
3. **Enrich** — for each unique directory, fetches the git branch name asynchronously
4. **Render** — draws groups as expandable trees with session names, branches, paths, and window counts

**Controls:** `j`/`k` to navigate, `t` to collapse/expand a group, `enter` to attach to a session, `q` to quit.

---

## `tmux/` — Tmux Communication

**What it does:** Talks to tmux by running `tmux` commands and parsing the output. Every function runs a specific tmux command:

- `ListSessions()` — `tmux list-sessions -F` with 21 format variables
- `ListWindows()` — `tmux list-windows -F`
- `ListTmuxPanes()` — `tmux list-panes -s -F`
- `NewSession(name, dir)` — `tmux new-session -d -s <name> -c <dir>`
- `AttachSession(name)` — `tmux attach-session -t <name>`
- `SwitchClient(name)` — `tmux switch-client -t <name>`
- `SwitchOrAttach(name, opts)` — decides between switch and attach based on context
- `CapturePane(name)` — `tmux capture-pane -e -p -t <name>` (gets terminal contents)
- `SendKeys(name, command)` — types keys into a pane
- `IsAttached()` — checks if `$TMUX` env var is set

---

## `zoxide/` — Zoxide Communication

**What it does:** Talks to zoxide to find directories you visit often.

- `ListResults()` — `zoxide query --list --score` — gets all your frequently-used directories
- `Query(path)` — `zoxide query <path>` — finds the best match for a path
- `Add(path)` — `zoxide add <path>` — adds a directory to zoxide's database

---

## `tmuxinator/` — Tmuxinator Communication

**What it does:** Talks to tmuxinator to use its project templates.

- `List()` — `tmuxinator list -n` — lists all projects
- `Start(name)` — `tmuxinator start --no-attach --name <name> <name>` — starts a project

---

## `git/` — Git Communication

**What it does:** Runs git commands to get repo information.

- `ShowTopLevel(path)` — `git -C <path> rev-parse --show-toplevel` — finds the repo root
- `GitCommonDir(path)` — `git -C <path> rev-parse --git-common-dir` — finds the .git directory
- `WorktreeList(path)` — `git -C <path> worktree list --porcelain` — lists all worktrees
- `CurrentBranch(path)` — `git -C <path> rev-parse --abbrev-ref HEAD` — gets the current branch
- `Clone(url, cmdDir, dir)` — `git clone <url>` — clones a repo

---

## `configurator/` — Reading Config Files

**What it does:** Finds and parses `sesh.toml`.

**Where it looks:**
1. `$XDG_CONFIG_HOME/sesh/sesh.toml` (or `~/.config/sesh/sesh.toml`)
2. Or the path specified by `--config` flag

**What it handles:**
- Regular TOML parsing via `pelletier/go-toml/v2`
- `StrictMode` — rejects unknown config fields
- `ImportPaths` — imports additional TOML files and merges their sessions/windows
- Human-friendly error messages when TOML is malformed
- Sets defaults: `DirLength=1`, `TUI.Prompt="> "`, `TUI.Placeholder="Filter Sessions..."`

---

## `startup/` — Startup Commands

**What it does:** When creating a new session, runs any startup commands you've configured. Creates windows first, then runs the command (with `{}` replaced by the session path).

**Strategy chain:**
1. Check the session's own config for a startup command
2. Check wildcard patterns for a matching command
3. Fall back to the default session config

---

## `shell/` — Running Shell Commands

**What it does:** The central place for running external commands. All other packages (tmux, zoxide, git, etc.) use this instead of calling `os/exec` directly.

- `Cmd(cmd, args...)` — runs a command, returns stdout as a string
- `CmdWithOutput(cmd, args...)` — runs a command, lets stdout show on the terminal
- `ListCmd(cmd, args...)` — runs a command, returns each line of output as a list
- `PrepareCmd(cmd, replacements)` — takes a command string like `ls -la {}`, replaces `{}` with values, splits into args

---

## `icon/` — Nerd Font Icons

**What it does:** Adds nerd font icons to session names based on their source:

| Source | Icon | Color |
|--------|------|-------|
| tmux |  | Blue |
| config |  | Gray |
| zoxide |  | Cyan |
| tmuxinator |  | Yellow |
| tmux-pane |  | Green |

---

## `cloner/` — Clone + Connect

**What it does:** Runs `git clone`, then immediately connects to the cloned repo as a tmux session. Combines `git.Clone()` + `connector.Connect()`.

---

## `dir/` — Directory Checks

**What it does:** Answers two questions:
- `Dir(path)` — is this path an existing directory?
- `RootDir(path)` — does this path have a git root? Handles regular repos, bare repos (with `.bare` folder), and worktrees.

---

## `home/` — Home Directory Handling

**What it does:** Helps with `~` in paths.
- `ShortenHome(path)` — replaces `/home/you/...` with `~/...`
- `ExpandPath(path)` — expands `~` and `$VARS` in paths

---

## `ls/` — Directory Listing

**What it does:** Lists files in a directory, either via a configured preview command or the system `ls`.

---

## `cache/` — Session Caching

**What it does:** Saves session lists to disk as gob-encoded files at `~/.cache/sesh/sessions.gob`. Uses atomic writes (writes to a `.tmp` file first, then renames). Used by `CachingLister` when `Cache: true` in config.

---

## Wrapper Packages — `execwrap/`, `oswrap/`, `pathwrap/`, `runtimewrap/`

**What they do:** Thin wrappers around Go standard library packages (`os/exec`, `os`, `path/filepath`, `runtime`). They exist so the code can be tested — in tests, these wrappers are swapped with mock implementations that return fake data without actually running commands.

- `execwrap.Exec` — wraps `os/exec` (running commands)
- `oswrap.Os` — wraps `os` (file system, env vars, home directory)
- `pathwrap.Path` — wraps `path/filepath` (path manipulation)
- `runtimewrap.Runtime` — wraps `runtime` (OS detection)

---

## Utility Packages

- **`json/`** — encodes session lists to JSON for the `sesh list --json` command
- **`replacer/`** — replaces text in strings using Aho-Corasick algorithm (fast, case-insensitive, whole-word matching)
- **`convert/`** — helper functions to parse strings into `time.Time`, `int`, `bool`, `float64`

---

## Build & Test

```bash
just build          # Build sesh, installs to $GOPATH/bin/sesh
just build v2.0.0   # Build with a specific version string
just test           # Generate mocks, then run all tests with coverage
just mock           # Generate mock files for testing only
go build -o ./sesh . && ./sesh <cmd>  # Build and run locally
```
