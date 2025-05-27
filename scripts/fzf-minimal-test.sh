#!/bin/bash

# Using system sesh via homebrew symlink

debug_log() {
    echo "$(date): $1" >> /tmp/sesh-fzf-debug.log
}

debug_log "Starting progressive test - Step 7: Add ctrl-f find directories"

# Step 7: Add ctrl-f for finding directories
selected=$(sesh list --icons | fzf-tmux -p 80%,70% \
  --no-sort --ansi --border-label ' sesh ' --prompt '⚡  ' \
  --header '  ^a all ^t tmux ^g configs ^x zoxide ^d kill ^f find ^r marked' \
  --bind 'ctrl-a:change-prompt(⚡  )+reload(sesh list --icons)' \
  --bind 'ctrl-t:change-prompt(🪟  )+reload(sesh list --tmux --icons)' \
  --bind 'ctrl-g:change-prompt(⚙️  )+reload(sesh list --config --icons)' \
  --bind 'ctrl-x:change-prompt(📁  )+reload(sesh list --zoxide --icons)' \
  --bind 'ctrl-d:execute(tmux kill-session -t {2..})+change-prompt(⚡  )+reload(sesh list --icons)' \
  --bind 'ctrl-f:change-prompt(🔎  )+reload(fd -H -d 2 -t d -E .Trash . ~)' \
  --bind 'ctrl-r:change-prompt(📌  )+reload(sesh list --marked --icons)')

debug_log "Selected: $selected"

if [ -n "$selected" ]; then
  cleaned=$(echo "$selected" | sed 's/\[[0-9;]*m//g' | sed 's/^[📌🟢⚙️📁⚡[:space:]]*//' | sed 's/^[[:space:]]*//')
  debug_log "Cleaned: $cleaned"
  
  # Try original cleaned name first
  if sesh connect "$cleaned" 2>/dev/null; then
    debug_log "Connected successfully with: $cleaned"
  else
    # If connection fails and name contains underscores, try replacing with spaces
    if [[ "$cleaned" == *"_"* ]] && [[ "$cleaned" != *"/"* ]] && [[ "$cleaned" != *"~"* ]]; then
      spaced_name=$(echo "$cleaned" | tr '_' ' ')
      debug_log "Trying with spaces: $spaced_name"
      sesh connect "$spaced_name" 2>&1 | tee -a /tmp/sesh-fzf-debug.log
    else
      debug_log "Connection failed for: $cleaned"
      exit 1
    fi
  fi
else
  debug_log "Failed - no selection"
fi

debug_log "Step 1 completed"