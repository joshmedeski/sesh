package model

type (
	SeshWildcards struct {
		// catalog of the sessions
		Directory SeshWildcardMap
		// unique identifiers of the sessions ordered
		OrderedIndex []string
	}

	SeshWildcardMap map[string]SeshWildcard

	SeshWildcard struct {
		Src                   string // The source of the session (config, tmux, zoxide, tmuxinator)
		Pattern               string // The absolute directory wildcard
		StartupCommand        string // The command to run when the session is started
		PreviewCommand        string // The command to run when the session is previewed
		DisableStartupCommand bool   // Ignore the default startup command if present
	}
)
