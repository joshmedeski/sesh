#!/bin/bash

# Enhanced sesh FZF script with working preview functionality
export PATH="/Users/rumen.d/.local/bin:$PATH"

# Debug logging function
debug_log() {
    echo "$(date): $1" >> /tmp/sesh-fzf-debug.log
}

debug_log "Starting sesh FZF script with preview"

# Test preview command and fallback if it fails
PREVIEW_CMD='/Users/rumen.d/.local/bin/sesh preview {}'
if ! timeout 2s /Users/rumen.d/.local/bin/sesh preview "sesh config" >/dev/null 2>&1; then
    debug_log "Preview command failed, using simple fallback"
    PREVIEW_CMD='echo "Preview for: {}"'
else
    debug_log "Preview command works, using enhanced preview"
fi

debug_log "Using preview command: $PREVIEW_CMD"

# Run fzf with preview
selected=$(sesh list --icons | fzf-tmux -p 80%,70% \
  --no-sort --ansi --border-label ' sesh ' --prompt '⚡  ' \
  --header '  ^a all ^t tmux ^g configs ^x zoxide ^d tmux kill ^m marked' \
  --preview "$PREVIEW_CMD" \
  --preview-window "right:50%:wrap" \
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