package live

// Logical 是单个 Claude 实例归类后的逻辑状态。
// 排序即严重度：值越大越严重，用于 session 聚合时取 max。
type Logical int

const (
	LogicalIdle Logical = iota
	LogicalSubagent
	LogicalBusy
	LogicalNeedsInput
)

func (l Logical) String() string {
	switch l {
	case LogicalNeedsInput:
		return "needs-input"
	case LogicalBusy:
		return "busy"
	case LogicalSubagent:
		return "subagent"
	default:
		return "idle"
	}
}

// Status 是一个 SeshSession（按 cwd 聚合后）下所有活实例的统计。
// Idle 由 Total - Busy - Subagent - Needing 推得。
type Status struct {
	Total    int
	Busy     int
	Subagent int
	Needing  int
}

func (s Status) Idle() int {
	idle := s.Total - s.Busy - s.Subagent - s.Needing
	if idle < 0 {
		return 0
	}
	return idle
}

func (s Status) IsEmpty() bool { return s.Total == 0 }

// Severity 返回该 session 的最高严重度，用于决定主徽章字符。
func (s Status) Severity() Logical {
	switch {
	case s.Needing > 0:
		return LogicalNeedsInput
	case s.Busy > 0:
		return LogicalBusy
	case s.Subagent > 0:
		return LogicalSubagent
	default:
		return LogicalIdle
	}
}

// classify 把 Claude 写到 json 的 raw status / kind 归到 Logical。
// 未知 status 一律当 Idle，避免穷举遗漏导致显示异常。
func classify(rawStatus, kind string) Logical {
	switch rawStatus {
	case "auth_url", "pending":
		return LogicalNeedsInput
	case "busy", "running", "in_progress", "async_launched", "compacting":
		if kind == "subagent" {
			return LogicalSubagent
		}
		return LogicalBusy
	default:
		// idle / completed / 未知值 / 缺失
		return LogicalIdle
	}
}
