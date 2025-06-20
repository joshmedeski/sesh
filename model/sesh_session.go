package model

type (
	SeshSessions struct {
		// catalog of the sessions
		Directory SeshSessionMap
		// unique identifiers of the sessions ordered
		OrderedIndex []string
	}

	SeshSessionMap map[string]SeshSession
	SeshWindowMap  map[string]WindowConfig

	SeshSession struct {
		Src  string // The source of the session (config, tmux, zoxide, tmuxinator)
		Name string // The display name
		Path string // The absolute directory path

		Attached              int            // Whether the session is currently attached
		Description           string         // The path with ~/ notation (e.g., "~/projects/myapp")
		DisableStartupCommand bool           // Ignore the default startup command if present
		Icon                  string         // Icon identifier or emoji for the session
		PreviewCommand        string         // The command to run when the session is previewed
		Score                 float64        // The score of the session (from Zoxide)
		StartupCommand        string         // The command to run when the session is started
		Tmuxinator            string         // Name of the tmuxinator config
		WindowConfigs         []WindowConfig // The windows used in session config
		WindowNames           []string       // The names of the windows in session config
		Windows               int            // The number of windows in the session
	}

	SeshSrcs struct {
		Config     bool
		Tmux       bool
		Tmuxinator bool
		Zoxide     bool
	}
)
