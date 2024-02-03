package tmux

import "time"

type TmuxSession struct {
	// Time of session last activity
	Activity *time.Time

	// Time session created
	Created *time.Time

	// Time session last attached
	LastAttached *time.Time

	// List of window indexes with alerts
	Alerts []int

	// Window indexes in most recent order
	Stack []int

	// List of clients session is attached to
	AttachedList []string

	// List of clients sessions in group are attached to
	GroupAttachedList []string

	// List of sessions in group
	GroupList []string

	// Name of session group
	Group string

	// Unique session ID
	ID string

	// Name of session
	Name string

	// Working directory of session
	Path string

	// Number of clients session is attached to
	Attached int

	// Number of clients sessions in group are attached to
	GroupAttached int

	// Size of session group
	GroupSize int

	// Number of windows in session
	Windows int

	// 1 if format is for a session
	Format bool

	// 1 if multiple clients attached to sessions in group
	GroupManyAttached bool

	// 1 if session in a group
	Grouped bool

	// 1 if multiple clients attached
	ManyAttached bool

	// 1 if this session contains the marked pane
	Marked bool
}
