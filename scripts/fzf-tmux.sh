#!/bin/bash

# Enhanced sesh FZF script with preview functionality
export PATH="/Users/rumen.d/.local/bin:$PATH"

# Debug logging function
debug_log() {
    echo "$(date): $1" >> /tmp/sesh-fzf-debug.log
}

debug_log "Starting sesh FZF script"

# Try different preview approaches
PREVIEW_SIMPLE='echo "Preview for: {}"'
PREVIEW_SESH='/Users/rumen.d/.local/bin/sesh preview {}'
PREVIEW_SESH_CLEAN='/Users/rumen.d/.local/bin/sesh preview $(echo {} | sed "s/\[[0-9;]*m//g" | sed "s/^[📌[:space:]]*//" | sed "s/^[[:space:]]*//")'

debug_log "Testing preview functionality..."

# Test which preview works
test_preview() {
    local preview_cmd="$1"
    local test_session="sesh config"
    
    debug_log "Testing preview: $preview_cmd"
    if timeout 2s bash -c "$preview_cmd" <<< "$test_session" >/dev/null 2>&1; then
        debug_log "Preview test successful: $preview_cmd"
        return 0
    else
        debug_log "Preview test failed: $preview_cmd"
        return 1
    fi
}

# Determine best preview
PREVIEW_CMD=""
if test_preview "$PREVIEW_SESH_CLEAN"; then
    PREVIEW_CMD="$PREVIEW_SESH_CLEAN"
elif test_preview "$PREVIEW_SESH"; then
    PREVIEW_CMD="$PREVIEW_SESH"
else
    PREVIEW_CMD="$PREVIEW_SIMPLE"
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
  # Strip ANSI codes, icons, and extract session name
  cleaned=$(echo "$selected" | sed 's/\[[0-9;]*m//g' | sed 's/^[📌[:space:]]*//' | sed 's/^[[:space:]]*//')
  debug_log "Cleaned: $cleaned"
  
  /Users/rumen.d/.local/bin/sesh connect "$cleaned" 2>&1 | tee -a /tmp/sesh-fzf-debug.log
fi

debug_log "Script completed"