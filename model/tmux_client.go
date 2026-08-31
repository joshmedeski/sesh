package model

// TmuxClient is a terminal attached to the tmux server. Sesh needs the whole
// record rather than just the name: when several clients are attached, which
// session each one is on and when it was last active is what decides which
// client a switch should move.
type TmuxClient struct {
	Name      string
	TTY       string
	SessionID string // e.g. "$3"
	Activity  int64  // unix seconds; higher is more recent
}
