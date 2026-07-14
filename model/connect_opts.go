package model

type ConnectOpts struct {
	Command    string
	Switch     bool
	Tmuxinator bool
}

type WorktreeCreateOpts struct {
	Number int
	Repo   string // GitHub "org/repo" override; empty => detect from cwd
	Pr     bool   // force the pull-request path
	Switch bool
}
