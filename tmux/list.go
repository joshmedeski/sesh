package tmux

import (
	"sort"
	"strings"
	"time"

	"github.com/joshmedeski/sesh/convert"
)

type TmuxSession struct {
	// Time of session last activity
	Activity *time.Time

	// Time session created
	Created *time.Time

	// Time session last attached
	LastAttached *time.Time

	// List of window indexes with alerts
	Alerts []int

	// Window indexes in most recent order
	Stack []int

	// List of clients session is attached to
	AttachedList []string

	// List of clients sessions in group are attached to
	GroupAttachedList []string

	// List of sessions in group
	GroupList []string

	// Name of session group
	Group string

	// Unique session ID
	ID string

	// Name of session
	Name string

	// Working directory of session
	Path string

	// Number of clients session is attached to
	Attached int

	// Number of clients sessions in group are attached to
	GroupAttached int

	// Size of session group
	GroupSize int

	// Number of windows in session
	Windows int

	// 1 if format is for a session
	Format bool

	// 1 if multiple clients attached to sessions in group
	GroupManyAttached bool

	// 1 if session in a group
	Grouped bool

	// 1 if multiple clients attached
	ManyAttached bool

	// 1 if this session contains the marked pane
	Marked bool
}

func format() string {
	variables := []string{
		"#{session_activity}",
		"#{session_alerts}",
		"#{session_attached}",
		"#{session_attached_list}",
		"#{session_created}",
		"#{session_format}",
		"#{session_group}",
		"#{session_group_attached}",
		"#{session_group_attached_list}",
		"#{session_group_list}",
		"#{session_group_many_attached}",
		"#{session_group_size}",
		"#{session_grouped}",
		"#{session_id}",
		"#{session_last_attached}",
		"#{session_many_attached}",
		"#{session_marked}",
		"#{session_name}",
		"#{session_path}",
		"#{session_stack}",
		"#{session_windows}",
	}

	return strings.Join(variables, " ")
}

func processSessions(sessionList []string) []*TmuxSession {
	sessions := make([]*TmuxSession, 0, len(sessionList))
	for _, line := range sessionList {
		fields := strings.Split(line, " ") // Strings split by single space

		if len(fields) != 21 {
			continue
		}

		session := &TmuxSession{
			Activity:          convert.StringToTime(fields[0]),
			Alerts:            convert.StringToIntSlice(fields[1]),
			Attached:          convert.StringToInt(fields[2]),
			AttachedList:      strings.Split(fields[3], ","),
			Created:           convert.StringToTime(fields[4]),
			Format:            convert.StringToBool(fields[5]),
			Group:             fields[6],
			GroupAttached:     convert.StringToInt(fields[7]),
			GroupAttachedList: strings.Split(fields[8], ","),
			GroupList:         strings.Split(fields[9], ","),
			GroupManyAttached: convert.StringToBool(fields[10]),
			GroupSize:         convert.StringToInt(fields[11]),
			Grouped:           convert.StringToBool(fields[12]),
			ID:                fields[13],
			LastAttached:      convert.StringToTime(fields[14]),
			ManyAttached:      convert.StringToBool(fields[15]),
			Marked:            convert.StringToBool(fields[16]),
			Name:              fields[17],
			Path:              fields[18],
			Stack:             convert.StringToIntSlice(fields[19]),
			Windows:           convert.StringToInt(fields[20]),
		}
		sessions = append(sessions, session)
	}

	return sessions
}

func sortSessions(sessions []*TmuxSession) []*TmuxSession {
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[j].LastAttached.Before(*sessions[i].LastAttached)
	})

	return sessions
}

func List() ([]*TmuxSession, error) {
	format := format()
	output, err := tmuxCmd([]string{"list-sessions", "-F", format})
	cleanOutput := strings.TrimSpace(output)
	if err != nil || strings.HasPrefix(cleanOutput, "no server running on") {
		return nil, nil
	}
	sessionList := strings.TrimSpace(string(output))
	lines := strings.Split(sessionList, "\n")
	sessions := processSessions(lines)

	return sortSessions(sessions), nil
}
