#!/bin/bash

# Test script without preview to isolate connection issues
export PATH="/Users/rumen.d/.local/bin:$PATH"

debug_log() {
    echo "$(date): $1" >> /tmp/sesh-fzf-debug.log
}

debug_log "Starting sesh FZF script WITHOUT preview for testing"

# Run fzf without preview to test basic functionality
selected=$(sesh list --icons | fzf-tmux -p 80%,70% \
  --no-sort --ansi --border-label ' sesh ' --prompt '⚡  ' \
  --header '  ^a all ^t tmux ^g configs ^x zoxide ^d tmux kill ^m marked' \
  --bind 'ctrl-a:change-prompt(⚡  )+reload(/Users/rumen.d/.local/bin/sesh list --icons)' \
  --bind 'ctrl-t:change-prompt(🪟  )+reload(/Users/rumen.d/.local/bin/sesh list --tmux --icons)' \
  --bind 'ctrl-g:change-prompt(⚙️  )+reload(/Users/rumen.d/.local/bin/sesh list --config --icons)' \
  --bind 'ctrl-x:change-prompt(📁  )+reload(/Users/rumen.d/.local/bin/sesh list --zoxide --icons)' \
  --bind 'ctrl-d:execute(/Users/rumen.d/.local/bin/sesh kill {})' \
  --bind 'ctrl-m:change-prompt(📌  )+reload(/Users/rumen.d/.local/bin/sesh list --marked --icons)')

debug_log "Selected: $selected"

if [ -n "$selected" ]; then
  # Strip ANSI codes and icons, but preserve the marked window format for connection
  cleaned=$(echo "$selected" | sed 's/\[[0-9;]*m//g' | sed 's/^[📌🟢⚙️📁⚡[:space:]]*//' | sed 's/^[[:space:]]*//')
  debug_log "Cleaned: $cleaned"
  
  /Users/rumen.d/.local/bin/sesh connect "$cleaned" 2>&1 | tee -a /tmp/sesh-fzf-debug.log
fi

debug_log "Script completed"