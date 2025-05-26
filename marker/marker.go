package marker

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/joshmedeski/sesh/v2/home"
)

type Marker interface {
	Mark(session, window string) error
	Unmark(session, window string) error
	IsMarked(session, window string) bool
	GetMarkedSessions() ([]MarkedSession, error)
	GetMarkedSessionsForSession(session string) ([]MarkedSession, error)
}

type RealMarker struct {
	home home.Home
}

type MarkedSession struct {
	Session   string `json:"session"`
	Window    string `json:"window"`
	Timestamp int64  `json:"timestamp"`
	Marked    bool   `json:"marked"`
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
	sessions[key] = MarkedSession{
		Session:   session,
		Window:    window,
		Timestamp: time.Now().Unix(),
		Marked:    true,
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