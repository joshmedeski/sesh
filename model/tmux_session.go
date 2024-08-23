package model

import "time"

type TmuxSession struct {
	Created           *time.Time
	LastAttached      *time.Time
	Activity          *time.Time
	Group             string
	Path              string
	Name              string
	ID                string
	AttachedList      []string
	GroupList         []string
	GroupAttachedList []string
	Stack             []int
	Alerts            []int
	GroupSize         int
	GroupAttached     int
	Attached          int
	Windows           int
	Format            bool
	GroupManyAttached bool
	Grouped           bool
	ManyAttached      bool
	Marked            bool
}
