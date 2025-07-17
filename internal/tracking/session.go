package tracking

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.etcd.io/bbolt"
)

// SessionState represents the current state of a work session
type SessionState int

const (
	StateStopped SessionState = iota
	StateRunning
	StatePaused
)

func (s SessionState) String() string {
	switch s {
	case StateStopped:
		return "Stopped"
	case StateRunning:
		return "Running"
	case StatePaused:
		return "Paused"
	default:
		return "Unknown"
	}
}

// Session represents a work session
type Session struct {
	ID        string        `json:"id"`
	Project   string        `json:"project"`
	StartTime time.Time     `json:"start_time"`
	EndTime   *time.Time    `json:"end_time,omitempty"`
	PausedAt  *time.Time    `json:"paused_at,omitempty"`
	Duration  time.Duration `json:"duration"`
	State     SessionState  `json:"state"`
}

// Tracker manages time tracking sessions
type Tracker struct {
	db           *bbolt.DB
	idleDetector *IdleDetector
	idleStop     chan struct{}
}

var (
	sessionsBucket = []byte("sessions")
	currentBucket  = []byte("current")
)

// NewTracker creates a new time tracker
func NewTracker() (*Tracker, error) {
	return NewTrackerWithIdleThreshold(5 * time.Minute)
}

// NewTrackerWithIdleThreshold creates a new time tracker with custom idle threshold
func NewTrackerWithIdleThreshold(idleThreshold time.Duration) (*Tracker, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	runeDir := filepath.Join(home, ".rune")
	if err := os.MkdirAll(runeDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create .rune directory: %w", err)
	}

	dbPath := filepath.Join(runeDir, "sessions.db")
	db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	idleDetector := NewIdleDetector(idleThreshold)

	tracker := &Tracker{
		db:           db,
		idleDetector: idleDetector,
	}
	if err := tracker.initBuckets(); err != nil {
		db.Close()
		return nil, err
	}

	return tracker, nil
}

// Close closes the tracker and database
func (t *Tracker) Close() error {
	if t.idleStop != nil {
		close(t.idleStop)
	}
	return t.db.Close()
}

// initBuckets creates the necessary database buckets
func (t *Tracker) initBuckets() error {
	return t.db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(sessionsBucket); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(currentBucket); err != nil {
			return err
		}
		return nil
	})
}

// Start starts a new work session
func (t *Tracker) Start(project string) (*Session, error) {
	// Check if there's already an active session
	current, err := t.GetCurrentSession()
	if err != nil {
		return nil, err
	}
	if current != nil && current.State != StateStopped {
		return nil, fmt.Errorf("session already active (state: %s)", current.State)
	}

	session := &Session{
		ID:        generateSessionID(),
		Project:   project,
		StartTime: time.Now(),
		State:     StateRunning,
	}

	if err := t.saveSession(session); err != nil {
		return nil, err
	}

	if err := t.setCurrentSession(session); err != nil {
		return nil, err
	}

	// Start idle monitoring when a session starts
	if err := t.StartIdleMonitoring(); err != nil {
		// Log error but don't fail the session start
		fmt.Printf("Warning: Failed to start idle monitoring: %v\n", err)
	}

	return session, nil
}

// Stop stops the current work session
func (t *Tracker) Stop() (*Session, error) {
	session, err := t.GetCurrentSession()
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, fmt.Errorf("no active session to stop")
	}

	now := time.Now()
	session.EndTime = &now
	session.State = StateStopped

	// Calculate total duration
	if session.PausedAt != nil {
		// If paused, don't include time since pause
		session.Duration = session.PausedAt.Sub(session.StartTime)
	} else {
		session.Duration = now.Sub(session.StartTime)
	}

	if err := t.saveSession(session); err != nil {
		return nil, err
	}

	if err := t.clearCurrentSession(); err != nil {
		return nil, err
	}

	// Stop idle monitoring when session stops
	t.StopIdleMonitoring()

	return session, nil
}

// Pause pauses the current work session
func (t *Tracker) Pause() (*Session, error) {
	session, err := t.GetCurrentSession()
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, fmt.Errorf("no active session to pause")
	}
	if session.State != StateRunning {
		return nil, fmt.Errorf("session is not running (state: %s)", session.State)
	}

	now := time.Now()
	session.PausedAt = &now
	session.State = StatePaused

	if err := t.saveSession(session); err != nil {
		return nil, err
	}

	if err := t.setCurrentSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

// Resume resumes a paused work session
func (t *Tracker) Resume() (*Session, error) {
	session, err := t.GetCurrentSession()
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, fmt.Errorf("no session to resume")
	}
	if session.State != StatePaused {
		return nil, fmt.Errorf("session is not paused (state: %s)", session.State)
	}

	// Calculate duration while paused and adjust start time
	pauseDuration := time.Since(*session.PausedAt)
	session.StartTime = session.StartTime.Add(pauseDuration)
	session.PausedAt = nil
	session.State = StateRunning

	if err := t.saveSession(session); err != nil {
		return nil, err
	}

	if err := t.setCurrentSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

// GetCurrentSession returns the current active session
func (t *Tracker) GetCurrentSession() (*Session, error) {
	var session *Session

	err := t.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(currentBucket)
		data := bucket.Get([]byte("session"))
		if data == nil {
			return nil
		}

		session = &Session{}
		return json.Unmarshal(data, session)
	})

	return session, err
}

// GetSessionDuration returns the current session duration
func (t *Tracker) GetSessionDuration() (time.Duration, error) {
	session, err := t.GetCurrentSession()
	if err != nil {
		return 0, err
	}
	if session == nil {
		return 0, nil
	}

	switch session.State {
	case StateRunning:
		return time.Since(session.StartTime), nil
	case StatePaused:
		if session.PausedAt != nil {
			return session.PausedAt.Sub(session.StartTime), nil
		}
		return time.Since(session.StartTime), nil
	case StateStopped:
		return session.Duration, nil
	default:
		return 0, nil
	}
}

// saveSession saves a session to the database
func (t *Tracker) saveSession(session *Session) error {
	return t.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(sessionsBucket)
		data, err := json.Marshal(session)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(session.ID), data)
	})
}

// setCurrentSession sets the current active session
func (t *Tracker) setCurrentSession(session *Session) error {
	return t.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(currentBucket)
		data, err := json.Marshal(session)
		if err != nil {
			return err
		}
		return bucket.Put([]byte("session"), data)
	})
}

// clearCurrentSession clears the current session
func (t *Tracker) clearCurrentSession() error {
	return t.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(currentBucket)
		return bucket.Delete([]byte("session"))
	})
}

// GetDailyTotal returns the total time worked today
func (t *Tracker) GetDailyTotal() (time.Duration, error) {
	var total time.Duration
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	err := t.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(sessionsBucket)
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var session Session
			if err := json.Unmarshal(v, &session); err != nil {
				continue
			}

			// Only count completed sessions from today
			if session.State == StateStopped &&
				session.StartTime.After(today) &&
				session.StartTime.Before(tomorrow) {
				total += session.Duration
			}
		}
		return nil
	})

	return total, err
}

// GetWeeklyTotal returns the total time worked this week
func (t *Tracker) GetWeeklyTotal() (time.Duration, error) {
	var total time.Duration
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	weekStart = weekStart.Truncate(24 * time.Hour)
	weekEnd := weekStart.Add(7 * 24 * time.Hour)

	err := t.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(sessionsBucket)
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var session Session
			if err := json.Unmarshal(v, &session); err != nil {
				continue
			}

			if session.State == StateStopped &&
				session.StartTime.After(weekStart) &&
				session.StartTime.Before(weekEnd) {
				total += session.Duration
			}
		}
		return nil
	})

	return total, err
}

// GetSessionHistory returns recent sessions
func (t *Tracker) GetSessionHistory(limit int) ([]*Session, error) {
	var sessions []*Session

	err := t.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(sessionsBucket)
		cursor := bucket.Cursor()

		// Collect all sessions first
		var allSessions []*Session
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var session Session
			if err := json.Unmarshal(v, &session); err != nil {
				continue
			}
			if session.State == StateStopped {
				allSessions = append(allSessions, &session)
			}
		}

		// Sort by start time (most recent first)
		for i := 0; i < len(allSessions)-1; i++ {
			for j := i + 1; j < len(allSessions); j++ {
				if allSessions[i].StartTime.Before(allSessions[j].StartTime) {
					allSessions[i], allSessions[j] = allSessions[j], allSessions[i]
				}
			}
		}

		// Take only the requested number
		if limit > 0 && len(allSessions) > limit {
			sessions = allSessions[:limit]
		} else {
			sessions = allSessions
		}

		return nil
	})

	return sessions, err
}

// GetProjectStats returns time statistics by project
func (t *Tracker) GetProjectStats() (map[string]time.Duration, error) {
	stats := make(map[string]time.Duration)

	err := t.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(sessionsBucket)
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var session Session
			if err := json.Unmarshal(v, &session); err != nil {
				continue
			}

			if session.State == StateStopped {
				stats[session.Project] += session.Duration
			}
		}
		return nil
	})

	return stats, err
}

// SetIdleThreshold sets the idle detection threshold
func (t *Tracker) SetIdleThreshold(threshold time.Duration) {
	t.idleDetector = NewIdleDetector(threshold)
}

// StartIdleMonitoring starts monitoring for idle state changes
func (t *Tracker) StartIdleMonitoring() error {
	if t.idleStop != nil {
		// Already monitoring
		return nil
	}

	t.idleStop = t.idleDetector.StartIdleMonitoring(
		func() {
			// On idle start - pause current session if running
			session, err := t.GetCurrentSession()
			if err != nil || session == nil || session.State != StateRunning {
				return
			}

			// Auto-pause due to idle
			_, _ = t.Pause()
		},
		func() {
			// On idle end - could potentially resume, but we'll leave that manual
			// to avoid accidentally resuming the wrong session
		},
	)

	return nil
}

// StopIdleMonitoring stops idle monitoring
func (t *Tracker) StopIdleMonitoring() {
	if t.idleStop != nil {
		close(t.idleStop)
		t.idleStop = nil
	}
}

// IsIdle returns true if the system is currently idle
func (t *Tracker) IsIdle() (bool, error) {
	return t.idleDetector.IsIdle()
}

// GetIdleTime returns the current system idle time
func (t *Tracker) GetIdleTime() (time.Duration, error) {
	return t.idleDetector.GetIdleTime()
}

// SaveImportedSession saves an imported session to the database
func (t *Tracker) SaveImportedSession(session *Session) error {
	return t.saveSession(session)
}

// generateSessionID generates a unique session ID
func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}
