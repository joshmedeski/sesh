#!/bin/bash

# Toggle mark/unmark for current tmux session:window
# Usage: Called from tmux binding to toggle marking state

export PATH="/Users/rumen.d/.local/bin:$PATH"

# Get current session and window
SESSION=$(tmux display-message -p '#S')
WINDOW=$(tmux display-message -p '#I')

# Check if currently marked by trying to list marked sessions and checking for our session:window
if sesh list --marked | grep -q "${SESSION}:.*\(${WINDOW}\)"; then
    # Currently marked - unmark it
    if sesh unmark; then
        tmux display-message "❌ Unmarked window ${SESSION}:${WINDOW}"
    else
        tmux display-message "❌ Failed to unmark window ${SESSION}:${WINDOW}"
    fi
else
    # Not marked - mark it
    if sesh mark; then
        tmux display-message "📌 Marked window ${SESSION}:${WINDOW}"
    else
        tmux display-message "❌ Failed to mark window ${SESSION}:${WINDOW}"
    fi
fi