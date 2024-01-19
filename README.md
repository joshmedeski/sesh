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

To install sesh, run the following [homebrew](https://brew.sh/) command:

```sh
brew install joshmedeski/sesh/sesh
```

### Using Go

Alternatively, you can install Sesh using Go's go install command:

```sh
go install github.com/joshmedeski/sesh@latest
```

This will download and install the latest version of Sesh. Make sure that your Go environment is properly set up.

**Note:** Do you want this on another package manager? [Create an issue](https://github.com/joshmedeski/sesh/issues/new) and let me know!

## How to use

`sesh list` will list all your sessions, and `sesh connect {session}` will connect to a session (automatically creating it if it doesn't exist yet). It is best used by integrating it into your sehll and tmux.

### fzf

The easiest way to integrate sesh into your workflow is to use [fzf](https://github.com/junegunn/fzf). You can use it to select a session to connect to:

```sh
sesh connect $(sesh list | fzf)
```

### tmux + fzf

In order to integrate with tmux, you can add a binding to your tmux config (`tmux.conf`). For example, the following will bind `ctrl-a T` to open a fzf prompt as a tmux popup (using `fzf-tmux`) and using different commands to list sessions (`sesh list -t`), zoxide directories (`sesh list -z`), and find directories (`fd...`).

```sh
 bind-key "T" run-shell "sesh connect $(
	sesh list -tz | fzf-tmux -p 55%,60% \
		--no-sort --border-label ' sesh ' --prompt '⚡  ' \
		--header '  ^a all ^t tmux ^x zoxide ^f find' \
		--bind 'tab:down,btab:up' \
		--bind 'ctrl-a:change-prompt(⚡  )+reload(sesh list)' \
		--bind 'ctrl-t:change-prompt(🪟  )+reload(sesh list -t)' \
		--bind 'ctrl-x:change-prompt(📁  )+reload(sesh list -z)' \
		--bind 'ctrl-f:change-prompt(🔎  )+reload(fd -H -d 2 -t d -E .Trash . ~)'
)"
```

You can customize this however you want, see `man fzf` for more info on the different options.

### zf

[zf](https://github.com/natecraddock/zf) is an alternative fuzzy finder designed for filtering filepaths, I've found it to be more accurate than fzf. It doesn't have as many options, but it's still a great tool. You can use it to select a session to connect to:

```sh
sesh connect (sesh list | zf --height 24)
```

## Background (the "t" script)

Sesh is a predecessor to my popular [t-smart-tmux-session-manager](https://github.com/joshmedeski/t-smart-tmux-session-manager) tmux plugin. After a year of development and over 250 stars, it's clear that people enjoy the idea of a smart session manager. However, I've always felt that the tmux plugin was a bit of a hack. It's a bash script that runs in the background and parses the output of tmux commands. It works, but it's not ideal and isn't flexible enough to support other terminal multiplexers.

I've decided to start over and build a session manager from the ground up. This time, I'm using a language that's more suited for the task: Go. Go is a compiled language that's fast, statically typed, and has a great standard library. It's perfect for a project like this. I've also decided to make this session manager multiplexer agnostic. It will be able to work with any terminal multiplexer, including tmux, zellij, Wezterm, and more.

The first step is to build a CLI that can interact with tmux and be a drop-in replacement for my previous tmux plugin. Once that's complete, I'll extend it to support other terminal multiplexers.

## Contributing

Pull requests are welcome, please help me make this project better! For major changes, please open an issue first to open a discussion.
