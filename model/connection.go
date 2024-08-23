package model

// Connection represents an established connection to a sesh session
type Connection struct {
	Session     SeshSession
	AddToZoxide bool // Whether to add the path to Zoxide
	Switch      bool // Whether to switch to the session (otherwise attach)
	Found       bool // Whether the connection was found
	New         bool // Whether the session was new
}
