package session

type Session struct {
	Src      string  // tmux or zoxide
	Name     string  // The display name
	Path     string  // The absolute directory path
	Score    float64 // The score of the session (from Zoxide)
	Attached int     // Whether the session is currently attached
	Windows  int     // The number of windows in the session
}

type Srcs struct {
	Tmux   bool
	Zoxide bool
}
