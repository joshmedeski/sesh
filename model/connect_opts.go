package model

type ConnectOpts struct {
	Command    string
	Switch     bool
	Tmuxinator bool
}

type WorktreeConnectOpts struct {
	Number      int
	Repo        string // GitHub "org/repo" override; empty => detect from cwd
	Pr          bool   // force the pull-request path
	Switch      bool
	FromBrowser bool // resolve Number/Repo/Pr from the active browser tab URL
}
