package session

type Session struct {
	Src  string // tmux or zoxide
	Name string // The display name
	Path string // The absolute directory path
}

type Srcs struct {
	Tmux   bool
	Zoxide bool
}
