package marker

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/joshmedeski/sesh/v2/home"
)

type Marker interface {
	Mark(session, window string) error
	Unmark(session, window string) error
	IsMarked(session, window string) bool
	GetMarkedSessions() ([]MarkedSession, error)
	GetMarkedSessionsForSession(session string) ([]MarkedSession, error)
	UpdateActivity(session, window string) error
	GetAlertLevel(session, window string) int
	ResetAlertForWindow(session, window string) error
}

type RealMarker struct {
	home home.Home
}

type MarkedSession struct {
	Session           string `json:"session"`
	Window            string `json:"window"`
	Timestamp         int64  `json:"timestamp"`
	Marked            bool   `json:"marked"`
	LastActivity      int64  `json:"last_activity"`
	AlertStartTime    int64  `json:"alert_start_time"`
	LastNavigated     int64  `json:"last_navigated"`
}

type MarkedSessionMap map[string]MarkedSession

func NewMarker(home home.Home) Marker {
	return &RealMarker{home: home}
}

func (m *RealMarker) getMarkerFilePath() string {
	return filepath.Join(m.home.SeshDir(), "marked.json")
}

func (m *RealMarker) loadMarkedSessions() (MarkedSessionMap, error) {
	filePath := m.getMarkerFilePath()
	
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return make(MarkedSessionMap), nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read marked sessions file: %w", err)
	}

	var sessions MarkedSessionMap
	if err := json.Unmarshal(data, &sessions); err != nil {
		return nil, fmt.Errorf("failed to parse marked sessions: %w", err)
	}

	return sessions, nil
}

func (m *RealMarker) saveMarkedSessions(sessions MarkedSessionMap) error {
	filePath := m.getMarkerFilePath()
	
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal marked sessions: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write marked sessions file: %w", err)
	}

	return nil
}

func (m *RealMarker) Mark(session, window string) error {
	sessions, err := m.loadMarkedSessions()
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s", session, window)
	now := time.Now().Unix()
	sessions[key] = MarkedSession{
		Session:      session,
		Window:       window,
		Timestamp:    now,
		Marked:       true,
		LastActivity: now,
	}

	return m.saveMarkedSessions(sessions)
}

func (m *RealMarker) Unmark(session, window string) error {
	sessions, err := m.loadMarkedSessions()
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s", session, window)
	delete(sessions, key)

	return m.saveMarkedSessions(sessions)
}

func (m *RealMarker) IsMarked(session, window string) bool {
	sessions, err := m.loadMarkedSessions()
	if err != nil {
		return false
	}

	key := fmt.Sprintf("%s:%s", session, window)
	marked, exists := sessions[key]
	return exists && marked.Marked
}

func (m *RealMarker) GetMarkedSessions() ([]MarkedSession, error) {
	sessions, err := m.loadMarkedSessions()
	if err != nil {
		return nil, err
	}

	var result []MarkedSession
	for _, session := range sessions {
		if session.Marked {
			result = append(result, session)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp > result[j].Timestamp
	})

	return result, nil
}

func (m *RealMarker) GetMarkedSessionsForSession(session string) ([]MarkedSession, error) {
	allMarked, err := m.GetMarkedSessions()
	if err != nil {
		return nil, err
	}

	var result []MarkedSession
	for _, marked := range allMarked {
		if marked.Session == session {
			result = append(result, marked)
		}
	}

	return result, nil
}

func (m *RealMarker) UpdateActivity(session, window string) error {
	sessions, err := m.loadMarkedSessions()
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s", session, window)
	if marked, exists := sessions[key]; exists && marked.Marked {
		marked.LastActivity = time.Now().Unix()
		sessions[key] = marked
		return m.saveMarkedSessions(sessions)
	}

	return nil
}

func (m *RealMarker) getTmuxWindowActivity(session, window string) (int64, error) {
	cmd := exec.Command("tmux", "list-windows", "-t", session, "-F", "#{window_index}:#{window_activity}")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get tmux window activity: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) == 2 && parts[0] == window {
			activity, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse activity timestamp: %w", err)
			}
			return activity, nil
		}
	}

	return 0, fmt.Errorf("window %s not found in session %s", window, session)
}

func (m *RealMarker) ResetAlertForWindow(session, window string) error {
	sessions, err := m.loadMarkedSessions()
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s", session, window)
	if marked, exists := sessions[key]; exists && marked.Marked {
		now := time.Now().Unix()
		marked.LastNavigated = now
		marked.AlertStartTime = 0 // Reset alert timer
		sessions[key] = marked
		return m.saveMarkedSessions(sessions)
	}

	return nil
}

func (m *RealMarker) GetAlertLevel(session, window string) int {
	sessions, err := m.loadMarkedSessions()
	if err != nil {
		return 0
	}

	key := fmt.Sprintf("%s:%s", session, window)
	marked, exists := sessions[key]
	if !exists || !marked.Marked {
		return 0
	}

	// Get actual tmux window activity
	tmuxActivity, err := m.getTmuxWindowActivity(session, window)
	if err != nil {
		// Fallback to stored activity if tmux query fails
		tmuxActivity = marked.LastActivity
	}

	now := time.Now().Unix()
	const inactivityThreshold = 10

	// Check if window has recent activity or user recently navigated to it
	lastRelevantTime := tmuxActivity
	if marked.LastNavigated > lastRelevantTime {
		lastRelevantTime = marked.LastNavigated
	}

	inactiveTime := now - lastRelevantTime

	// If recently active or recently navigated, reset alert
	if inactiveTime <= inactivityThreshold {
		if marked.AlertStartTime > 0 {
			marked.AlertStartTime = 0
			sessions[key] = marked
			m.saveMarkedSessions(sessions)
		}
		return 0
	}

	// Start alert timer if not already started
	if marked.AlertStartTime == 0 {
		// Alert will start counting from when inactivity period ends
		marked.AlertStartTime = lastRelevantTime + inactivityThreshold
		sessions[key] = marked
		m.saveMarkedSessions(sessions)
	}

	// Calculate how long we've been in alert state (time since inactivity threshold was reached)
	alertDuration := now - marked.AlertStartTime

	// Only show alert if we're past the alert start time
	if alertDuration <= 0 {
		return 0
	} else if alertDuration <= 60 {
		return 1 // 0-1 minute after becoming inactive
	} else if alertDuration <= 300 {
		return 2 // 1-5 minutes after becoming inactive  
	} else {
		return 3 // 5+ minutes after becoming inactive
	}
}