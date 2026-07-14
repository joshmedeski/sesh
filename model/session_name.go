package model

// SessionNameSeparator joins a session's base name and its enriched suffix
// (e.g. a GitHub issue title). The em dash is legal in tmux session names,
// unlike ':' and '.'.
const SessionNameSeparator = " — "
