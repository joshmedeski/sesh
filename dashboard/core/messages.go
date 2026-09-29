// messages.go
package core

// HoveredSessionMsg carries the session under the cursor from the sessions
// list to the details widget.
type HoveredSessionMsg struct {
	Name    string
	Path    string
	Windows int
}
