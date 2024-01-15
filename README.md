<p align="center">
  <img width="50%" height="50%" src="https://github.com/arl/gitmux/raw/readme-images/logo-transparent.png" />
</p>

<p align="center">Gitmux shows git status in your tmux status bar</p>

<hr/>

<p align="center">
  <a href="https://github.com/joshmedeski/sesh/actions/workflows/ci-cd.yaml">
    <img alt="tests" src="https://github.com/joshmedeski/sesh/actions/workflows/ci-cd.yaml/badge.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/arl/gitmux">
    <img alt="goreport" src="https://goreportcard.com/badge/github.com/arl/gitmux" />
  </a>
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
</p>

# sesh

Smart session manager tooling for the terminal.

**ALPHA:** This project is in active development and is not ready for use. I will update the README with proper installation instructions once it's ready. Please check out the issues and contribute if you're interested in helping out.

## Rewrite

- [x] List tmux sessions
- [x] List zoxide results
- [ ] Conditionally render tmux popup
- [ ] **Connect to a session**
- [ ] Keymaps to switch views in fzf (ctrl+f to find)
- [ ] Quick cloning
- [ ] Add confiugration

## Background

Sesh is a predecessor to my popular [t-smart-tmux-session-manager](https://github.com/joshmedeski/t-smart-tmux-session-manager) tmux plugin. After a year of development and over 250 stars, it's clear that people enjoy the idea of a smart session manager. However, I've always felt that the tmux plugin was a bit of a hack. It's a bash script that runs in the background and parses the output of tmux commands. It works, but it's not ideal and isn't flexible enough to support other terminal multiplexers.

I've decided to start over and build a session manager from the ground up. This time, I'm using a language that's more suited for the task: Go. Go is a compiled language that's fast, statically typed, and has a great standard library. It's perfect for a project like this. I've also decided to make this session manager multiplexer agnostic. It will be able to work with any terminal multiplexer, including tmux, zellij, Wezterm, and more.

The first step is to build a CLI that can interact with tmux and be a drop-in replacement for my previous tmux plugin. Once that's complete, I'll extend it to support other terminal multiplexers.
