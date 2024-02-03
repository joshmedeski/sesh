package tmux

import (
	"sort"
	"strings"

	"github.com/joshmedeski/sesh/convert"
)

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

type Options struct {
	HideAttached bool
}

func processSessions(o Options, sessionList []string) []Session {
	sessions := make([]Session, 0, len(sessionList))
	for _, line := range sessionList {
		fields := strings.Split(line, " ") // Strings split by single space

		if len(fields) != 21 {
			continue
		}
		if o.HideAttached && fields[2] == "1" {
			continue
		}

		session := Session{
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
			name:              fields[17],
			path:              fields[18],
			Stack:             convert.StringToIntSlice(fields[19]),
			Windows:           convert.StringToInt(fields[20]),
		}
		sessions = append(sessions, session)
	}

	return sessions
}

func sortSessions(sessions []Session) []Session {
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[j].LastAttached.Before(*sessions[i].LastAttached)
	})

	return sessions
}

func (c *Command) List(o Options) ([]Session, error) {
	format := format()
	output, err := command.Run([]string{"list-sessions", "-F", format})
	if err != nil {
		return nil, err
	}

	sessionList := output
	lines := strings.Split(sessionList, "\n")
	sessions := processSessions(o, lines)

	return sortSessions(sessions), nil
}
