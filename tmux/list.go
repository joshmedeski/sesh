package tmux

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type TmuxSession struct {
	// Time of session last activity
	Activity time.Time

	// Time session created
	Created time.Time

	// Time session last attached
	LastAttached time.Time

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
		"session_activity",
		"session_alerts",
		"session_attached",
		"session_attached_list",
		"session_created",
		"session_format",
		"session_group",
		"session_group_attached",
		"session_group_attached_list",
		"session_group_list",
		"session_group_many_attached",
		"session_group_size",
		"session_grouped",
		"session_id",
		"session_last_attached",
		"session_many_attached",
		"session_marked",
		"session_name",
		"session_path",
		"session_stack",
		"session_windows",
	}
	variablesStr := ""
	for i, variable := range variables {
		variablesStr += "#{" + variable + "}"
		if i != len(variables)-1 {
			variablesStr += " "
		}
	}
	return variablesStr
}

func stringToTime(s string) time.Time {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	t := time.Unix(i, 0)
	return t
}

func stringToIntSlice(s string) []int {
	split := strings.Split(s, ",") // Or another delimiter if not ","
	ints := make([]int, 0, len(split))
	for _, str := range split {
		if i, err := strconv.Atoi(str); err == nil {
			ints = append(ints, i)
		}
	}
	return ints
}

func stringToBool(s string) bool {
	return s == "1"
}

func List() ([]*TmuxSession, error) {
	format := format()
	output, err := tmuxCmd([]string{"list-sessions", "-F", format})
	if err != nil {
		return nil, err
	}

	sessionList := strings.TrimSpace(string(output))
	lines := strings.Split(sessionList, "\n")

	sessions := make([]*TmuxSession, 0, len(lines))
	for _, line := range lines {
		fields := strings.Split(line, " ") // Strings split by single space
		if len(fields) == 21 {
			activity := stringToTime(fields[0])
			alerts := stringToIntSlice(fields[1])
			attached, _ := strconv.Atoi(fields[2])
			attachedList := strings.Split(fields[3], ",") // replace "," with appropriate delimiter
			created := stringToTime(fields[4])
			format := stringToBool(fields[5])
			group := fields[6]
			groupAttached, _ := strconv.Atoi(fields[7])
			groupAttachedList := strings.Split(fields[8], ",") // replace "," with appropriate delimiter
			groupList := strings.Split(fields[9], ",")         // replace "," with appropriate delimiter
			groupManyAttached := stringToBool(fields[10])
			groupSize, _ := strconv.Atoi(fields[11])
			grouped := stringToBool(fields[12])
			id := fields[13]
			lastAttached := stringToTime(fields[14])
			manyAttached := stringToBool(fields[15])
			marked := stringToBool(fields[16])
			name := fields[17]
			path := fields[18]
			stack := stringToIntSlice(fields[19])
			windows, _ := strconv.Atoi(fields[20])

			sessions = append(sessions, &TmuxSession{
				Activity:          activity,
				Alerts:            alerts,
				Attached:          attached,
				AttachedList:      attachedList,
				Created:           created,
				Format:            format,
				Group:             group,
				GroupAttached:     groupAttached,
				GroupAttachedList: groupAttachedList,
				GroupList:         groupList,
				GroupManyAttached: groupManyAttached,
				GroupSize:         groupSize,
				Grouped:           grouped,
				ID:                id,
				LastAttached:      lastAttached,
				ManyAttached:      manyAttached,
				Marked:            marked,
				Name:              name,
				Path:              path,
				Stack:             stack,
				Windows:           windows,
			})
		}
	}

	return sessions, nil
}
