<p align="center">
  <img width="256" height="256" src="https://github.com/joshmedeski/sesh/blob/main/sesh-icon.png" />
</p>

<h1 align="center">Sesh, the smart terminal session manager</h1>

<p align="center">
  <a href="https://github.com/joshmedeski/sesh/actions/workflows/ci-cd.yml">
    <img alt="tests" src="https://github.com/joshmedeski/sesh/actions/workflows/ci-cd.yml/badge.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/joshmedeski/sesh">
    <img alt="goreport" src="https://goreportcard.com/badge/github.com/joshmedeski/sesh" />
  </a>
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
</p>

**BETA:** This project is in active development. Please check out the issues and contribute if you're interested in helping out.

## How to install

### Homebrew

To install sesh, run the following [homebrew](https://brew.sh/) command:

```sh
brew install joshmedeski/sesh/sesh
```

### Go

Alternatively, you can install Sesh using Go's go install command:

```sh
go install github.com/joshmedeski/sesh@latest
```

This will download and install the latest version of Sesh. Make sure that your Go environment is properly set up.

**Note:** Do you want this on another package manager? [Create an issue](https://github.com/joshmedeski/sesh/issues/new) and let me know!

## Raycast Extension

The [sesh companion extension](https://www.raycast.com/joshmedeski/sesh) for [Raycast](https://www.raycast.com/) makes it easy to use sesh outside of the terminal.

Here are limitations to keep in mind:

- tmux has to be running before you can use the extension
- The extension caches results for a few seconds, so it may not always be up to date

## How to use

### tmux for sessions

[tmux](https://github.com/tmux/tmux) is a powerful terminal multiplexer that allows you to create and manage multiple terminal sessions. Sesh is designed to make managing tmux sessions easier.

### zoxide for directories

[zoxide](https://github.com/ajeetdsouza/zoxide) is a blazing fast alternative to `cd` that tracks your most used directories. Sesh uses zoxide to manage your projects. You'll have to set up zoxide first, but once you do, you can use it to quickly jump to your most used directories.

### Basic usage

Once tmux and zoxide are setup, `sesh list` will list all your tmux sessions and zoxide results, and `sesh connect {session}` will connect to a session (automatically creating it if it doesn't exist yet). It is best used by integrating it into your shell and tmux.

#### fzf

The easiest way to integrate sesh into your workflow is to use [fzf](https://github.com/junegunn/fzf). You can use it to select a session to connect to:

```sh
sesh connect $(sesh list | fzf)
```

#### tmux + fzf

In order to integrate with tmux, you can add a binding to your tmux config (`tmux.conf`). For example, the following will bind `ctrl-a T` to open a fzf prompt as a tmux popup (using `fzf-tmux`) and using different commands to list sessions (`sesh list -t`), zoxide directories (`sesh list -z`), and find directories (`fd...`).

```sh
bind-key "T" run-shell "sesh connect \"$(
	sesh list -tz | fzf-tmux -p 55%,60% \
		--no-sort --border-label ' sesh ' --prompt '⚡  ' \
		--header '  ^a all ^t tmux ^x zoxide ^f find' \
		--bind 'tab:down,btab:up' \
		--bind 'ctrl-a:change-prompt(⚡  )+reload(sesh list)' \
		--bind 'ctrl-t:change-prompt(🪟  )+reload(sesh list -t)' \
		--bind 'ctrl-x:change-prompt(📁  )+reload(sesh list -z)' \
		--bind 'ctrl-f:change-prompt(🔎  )+reload(fd -H -d 2 -t d -E .Trash . ~)'
)\""
```

You can customize this however you want, see `man fzf` for more info on the different options.

See my video, [Top 4 Fuzzy CLIs](https://www.youtube.com/watch?v=T0O2qrOhauY) for more inspiration for tooling that can be integrated with sesh.

## Recommended tmux Settings

I recommend you add these settings to your `tmux.conf` to have a better experience with this plugin.

```sh
bind-key x kill-pane # skip "kill-pane 1? (y/n)" prompt
set -g detach-on-destroy off  # don't exit from tmux when closing a session
```

## Configuration

You can configure sesh by creating a `sesh.toml` file in your `$XDG_CONFIG_HOME/sesh` or `$HOME/.config/sesh` directory.

```sh
mkdir -p ~/.config/sesh
touch sesh.toml
```

### Startup Scripts

A startup script is a script that is run when a session is created. It is useful for setting up your environment for a given project. For example, you may want to run `npm run dev` to automatically start a dev server.

**Note:** If you use the `--command/-c` flag, then the startup script will not be run.

The default startup will run on every project that doesn't have a specific startup script. You can configure the default startup script by setting the `default_startup_script` property in your `sesh.toml` file.

```toml
default_startup_script = "~/.config/sesh/scripts/default"
```

I like to use a script that opens nvim on session startup:

```sh
#!/usr/bin/env bash
tmux send-keys "nvim" Enter
```

The set a specific startup script for a project, you can add a `startup_script` property to your `sesh.toml` file.

```toml
[[startup_scripts]]
session_path = "~/code/sesh"
script_path = "~/.config/sesh/scripts/go"

[[startup_scripts]]
session_path = "~/code/joshmedeski.com"
script_path = "~/.config/sesh/scripts/node_dev"
```

The script can execute tmux commands to create panes, additional windows and trigger commands. Here is an example of a script (node_dev) that creates a pane for a dev server and opens a new pane with nvim:

```sh
#!/usr/bin/env bash
tmux split-window -v -p 30 "npm run dev"
tmux select-pane -t :.+
tmux send-keys "nvim" Enter
```

Set the file as an executable and it will be run when you connect to the specified session.

## Background (the "t" script)

Sesh is the successor to my popular [t-smart-tmux-session-manager](https://github.com/joshmedeski/t-smart-tmux-session-manager) tmux plugin. After a year of development and over 250 stars, it's clear that people enjoy the idea of a smart session manager. However, I've always felt that the tmux plugin was a bit of a hack. It's a bash script that runs in the background and parses the output of tmux commands. It works, but it's not ideal and isn't flexible enough to support other terminal multiplexers.

I've decided to start over and build a session manager from the ground up. This time, I'm using a language that's more suited for the task: Go. Go is a compiled language that's fast, statically typed, and has a great standard library. It's perfect for a project like this. I've also decided to make this session manager multiplexer agnostic. It will be able to work with any terminal multiplexer, including tmux, zellij, Wezterm, and more.

The first step is to build a CLI that can interact with tmux and be a drop-in replacement for my previous tmux plugin. Once that's complete, I'll extend it to support other terminal multiplexers.

## Contributors

<a href="https://github.com/joshmedeski/sesh/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=joshmedeski/sesh" />
</a>

Made with [contrib.rocks](https://contrib.rocks).
