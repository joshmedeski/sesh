package model

type TmuxPane struct {
	WindowIndex int
	WindowName  string
	PaneIndex   int
	PaneTitle   string
	PaneCommand string
	PanePath    string
	PaneID      string
}

// TmuxPaneAcrossSessions 跨所有 session 列出 pane 的精简视图。
// 让 cc-sesh 把 Claude 实例按 pane 当前 cwd 归到对应的 session。
type TmuxPaneAcrossSessions struct {
	SessionName     string
	PaneID          string
	PaneCurrentPath string
	PanePID         int
}
