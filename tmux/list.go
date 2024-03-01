package tmux

import (
	"sort"
	"strings"
	"time"

	"github.com/joshmedeski/sesh/convert"
)

type TmuxSession struct {
	Activity          *time.Time // Time of session last activity
	Created           *time.Time // Time session created
	LastAttached      *time.Time // Time session last attached
	Alerts            []int      // List of window indexes with alerts
	Stack             []int      // Window indexes in most recent order
	AttachedList      []string   // List of clients session is attached to
	GroupAttachedList []string   // List of clients sessions in group are attached to
	GroupList         []string   // List of sessions in group
	Group             string     // Name of session group
	ID                string     // Unique session ID
	Name              string     // Name of session
	Path              string     // Working directory of session
	Attached          int        // Number of clients session is attached to
	GroupAttached     int        // Number of clients sessions in group are attached to
	GroupSize         int        // Size of session group
	Windows           int        // Number of windows in session
	Format            bool       // 1 if format is for a session
	GroupManyAttached bool       // 1 if multiple clients attached to sessions in group
	Grouped           bool       // 1 if session in a group
	ManyAttached      bool       // 1 if multiple clients attached
	Marked            bool       // 1 if this session contains the marked pane
}

var separator = "::"

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

	return strings.Join(variables, separator)
}

type Options struct {
	HideAttached bool
}

func processSessions(o Options, sessionList []string) []*TmuxSession {
	sessions := make([]*TmuxSession, 0, len(sessionList))
	for _, line := range sessionList {
		fields := strings.Split(line, separator) // Strings split by single space

		if len(fields) != 21 {
			continue
		}
		if o.HideAttached && fields[2] == "1" {
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

func List(o Options) ([]*TmuxSession, error) {
	format := format()
	output, err := tmuxCmd([]string{"list-sessions", "-F", format})
	cleanOutput := strings.TrimSpace(output)
	if err != nil || strings.HasPrefix(cleanOutput, "no server running on") {
		return nil, nil
	}
	sessionList := strings.TrimSpace(string(output))
	lines := strings.Split(sessionList, "\n")
	sessions := processSessions(o, lines)

	return sortSessions(sessions), nil
}
