package model

type (
	SeshSessions struct {
		// catalog of the sessions
		Directory SeshSessionMap
		// unique identifiers of the sessions ordered
		OrderedIndex []string
	}

	SeshSessionMap map[string]SeshSession

	SeshSession struct {
		Src  string // The source of the session (config, tmux, tmuxinator, zoxide)
		Name string // The display name
		Path string // The absolute directory path

		StartupCommand string  // The command to run when the session is started
		Attached       int     // Whether the session is currently attached
		Windows        int     // The number of windows in the session
		Score          float64 // The score of the session (from Zoxide)
	}

	SeshSrcs struct {
		Config     bool
		Tmux       bool
		Tmuxinator bool
		Zoxide     bool
	}
)
