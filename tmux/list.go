package tmux

import (
	"sort"
	"strings"

	"github.com/joshmedeski/sesh/v2/convert"
	"github.com/joshmedeski/sesh/v2/model"
)

func (t *RealTmux) ListSessions() ([]*model.TmuxSession, error) {
	output, err := t.shell.ListCmd("tmux", "list-sessions", "-F", listsessionsformat())
	if err != nil {
		return []*model.TmuxSession{}, nil
	}
	sessions, err := parseTmuxSessionsOutput(output)
	if err != nil {
		return nil, err
	}
	sortedSessions := sortByLastAttached(sessions)
	return sortedSessions, nil
}

var separator = "::"

func listsessionsformat() string {
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

func parseTmuxSessionsOutput(rawList []string) ([]*model.TmuxSession, error) {
	sessions := make([]*model.TmuxSession, 0, len(rawList))
	for _, line := range rawList {
		fields := strings.Split(line, separator)

		if len(fields) != 21 {
			continue
		}

		session := &model.TmuxSession{
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

	return sessions, nil
}

func sortByLastAttached(sessions []*model.TmuxSession) []*model.TmuxSession {
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[j].LastAttached.Before(*sessions[i].LastAttached)
	})
	return sessions
}
