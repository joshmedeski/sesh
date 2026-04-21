package model

type WeztermWorkspace struct {
	Name  string
	Panes []WeztermPane
}

type WeztermPane struct {
	PaneID    int
	TabID     int
	WindowID  int
	Workspace string
	Cwd       string
	Title     string
	IsActive  bool
}
