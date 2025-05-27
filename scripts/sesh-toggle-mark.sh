#!/bin/bash

# Toggle mark/unmark for current tmux session:window
# Usage: Called from tmux binding to toggle marking state

# Using system sesh via homebrew symlink

# Get current session and window
SESSION=$(tmux display-message -p '#S')
WINDOW=$(tmux display-message -p '#I')

# Check if currently marked by trying to list marked sessions and checking for our session:window
# Convert session name spaces to underscores to match display format
SESSION_DISPLAY=$(echo "$SESSION" | tr ' ' '_')
if sesh list --marked | grep -q "${SESSION_DISPLAY}:.*\(${WINDOW}\)"; then
    # Currently marked - unmark it
    if sesh unmark >/dev/null 2>&1; then
        tmux display-message "❌ Unmarked window ${SESSION}:${WINDOW}"
    else
        tmux display-message "❌ Failed to unmark window ${SESSION}:${WINDOW}"
    fi
else
    # Not marked - mark it
    if sesh mark >/dev/null 2>&1; then
        tmux display-message "📌 Marked window ${SESSION}:${WINDOW}"
    else
        tmux display-message "❌ Failed to mark window ${SESSION}:${WINDOW}"
    fi
fi
