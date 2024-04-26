package model

type SeshSession struct {
	Src  string // The source of the session (config, tmux, zoxide)
	Name string // The display name
	Path string // The absolute directory path

	Attached int     // Whether the session is currently attached
	Windows  int     // The number of windows in the session
	Score    float64 // The score of the session (from Zoxide)
}

type SeshSrcs struct {
	Config bool
	Tmux   bool
	Zoxide bool
}