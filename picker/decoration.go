package picker

import (
	"time"

	"github.com/Wingsdh/cc-sesh/v2/model"
)

// Decoration 是 picker 渲染单条 session 时额外要展示的信息。
// 由调用方在 fetch 阶段构造（通常基于 claude/live + claude/attention 的查询），
// picker 包本身不知道这些数据从哪来 —— 仅按字段渲染。
type Decoration struct {
	Live      LiveBadge
	Attention AttentionBadge
}

// LiveBadge 是该 session 当前实时的 Claude 状态聚合。Total=0 表示该 session 内没 Claude。
type LiveBadge struct {
	Total    int
	Busy     int
	Subagent int
	Needing  int
}

func (l LiveBadge) IsEmpty() bool { return l.Total == 0 }

func (l LiveBadge) Idle() int {
	idle := l.Total - l.Busy - l.Subagent - l.Needing
	if idle < 0 {
		return 0
	}
	return idle
}

// Severity 决定主徽章字符的归类。
type Severity int

const (
	SevIdle Severity = iota
	SevSubagent
	SevBusy
	SevNeeding
)

func (l LiveBadge) Severity() Severity {
	switch {
	case l.Needing > 0:
		return SevNeeding
	case l.Busy > 0:
		return SevBusy
	case l.Subagent > 0:
		return SevSubagent
	default:
		return SevIdle
	}
}

// AttentionBadge 是粘性的「需介入」标记。Triggered=false 表示无标记。
type AttentionBadge struct {
	Triggered bool
	FirstAt   time.Time
	Reason    string
}

// Decorator 把 SeshSession 映射为 Decoration。picker 渲染时按需调用。
// 实现可以是空操作（NoDecoration），picker 此时退化为原 sesh 行为。
type Decorator interface {
	Decorate(s model.SeshSession) Decoration
}

// NoDecoration 在调用方不需要装饰时使用，避免 picker 内部判 nil。
type NoDecoration struct{}

func (NoDecoration) Decorate(model.SeshSession) Decoration { return Decoration{} }
