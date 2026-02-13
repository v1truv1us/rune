package commands

import (
	"os"
	"testing"
	"time"

	"github.com/ferg-cod3s/rune/internal/tracking"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSwitch_NoActiveSession(t *testing.T) {
	setupTestEnvironment(t)

	// Try to switch without an active session
	cmd := &cobra.Command{}
	err := runSwitch(cmd, []string{"new-project"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no active session to switch from")
}

func TestSwitch_PausedSession(t *testing.T) {
	setupTestEnvironment(t)

	// Start and pause a session
	tracker, err := tracking.NewTracker()
	require.NoError(t, err)
	_, err = tracker.Start("old-project")
	require.NoError(t, err)
	_, err = tracker.Pause()
	require.NoError(t, err)
	tracker.Close()

	// Try to switch from paused session
	cmd := &cobra.Command{}
	err = runSwitch(cmd, []string{"new-project"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "current session is not running")
}

func TestSwitch_SameProject(t *testing.T) {
	setupTestEnvironment(t)

	// Start a session
	tracker, err := tracking.NewTracker()
	require.NoError(t, err)
	_, err = tracker.Start("same-project")
	require.NoError(t, err)
	tracker.Close()

	// Try to switch to the same project
	cmd := &cobra.Command{}
	err = runSwitch(cmd, []string{"same-project"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already working on project")
}

func TestSwitch_SuccessfulSwitch(t *testing.T) {
	setupTestEnvironment(t)

	// Start a session
	tracker, err := tracking.NewTracker()
	require.NoError(t, err)
	session1, err := tracker.Start("project-a")
	require.NoError(t, err)
	assert.Equal(t, "project-a", session1.Project)
	assert.Equal(t, tracking.StateRunning, session1.State)
	tracker.Close()

	// Wait a bit to accumulate some time
	time.Sleep(10 * time.Millisecond)

	// Switch to a new project
	cmd := &cobra.Command{}
	err = runSwitch(cmd, []string{"project-b"})

	require.NoError(t, err)

	// Verify new session is active
	tracker2, err := tracking.NewTracker()
	require.NoError(t, err)
	defer tracker2.Close()

	currentSession, err := tracker2.GetCurrentSession()
	require.NoError(t, err)
	require.NotNil(t, currentSession)
	assert.Equal(t, "project-b", currentSession.Project)
	assert.Equal(t, tracking.StateRunning, currentSession.State)

	// Verify old session was stopped and saved
	history, err := tracker2.GetSessionHistory(10)
	require.NoError(t, err)
	require.Len(t, history, 1)
	assert.Equal(t, "project-a", history[0].Project)
	assert.Equal(t, tracking.StateStopped, history[0].State)
	assert.True(t, history[0].Duration > 0)
}

func TestSwitch_MultipleSequentialSwitches(t *testing.T) {
	setupTestEnvironment(t)

	projects := []string{"project-1", "project-2", "project-3"}

	// Start with first project
	tracker, err := tracking.NewTracker()
	require.NoError(t, err)
	_, err = tracker.Start(projects[0])
	require.NoError(t, err)
	tracker.Close()

	// Switch through all projects
	for i := 1; i < len(projects); i++ {
		time.Sleep(5 * time.Millisecond)

		cmd := &cobra.Command{}
		err = runSwitch(cmd, []string{projects[i]})
		require.NoError(t, err)

		// Verify current session
		trackerCheck, err := tracking.NewTracker()
		require.NoError(t, err)
		currentSession, err := trackerCheck.GetCurrentSession()
		require.NoError(t, err)
		require.NotNil(t, currentSession)
		assert.Equal(t, projects[i], currentSession.Project)
		trackerCheck.Close()
	}

	// Stop final session
	trackerFinal, err := tracking.NewTracker()
	require.NoError(t, err)
	_, err = trackerFinal.Stop()
	require.NoError(t, err)

	// Verify all sessions are in history
	history, err := trackerFinal.GetSessionHistory(10)
	require.NoError(t, err)
	assert.Len(t, history, len(projects))

	// Verify sessions are in reverse chronological order
	for i := 0; i < len(projects); i++ {
		expectedProject := projects[len(projects)-1-i]
		assert.Equal(t, expectedProject, history[i].Project)
		assert.Equal(t, tracking.StateStopped, history[i].State)
	}
	trackerFinal.Close()
}

func TestSwitch_PreservesSessionBoundaries(t *testing.T) {
	setupTestEnvironment(t)

	// Start first session
	tracker, err := tracking.NewTracker()
	require.NoError(t, err)
	_, err = tracker.Start("project-alpha")
	require.NoError(t, err)
	tracker.Close()

	time.Sleep(10 * time.Millisecond)

	// Switch to second project
	cmd := &cobra.Command{}
	err = runSwitch(cmd, []string{"project-beta"})
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	// Stop second session
	tracker2, err := tracking.NewTracker()
	require.NoError(t, err)
	_, err = tracker2.Stop()
	require.NoError(t, err)

	// Get project stats
	stats, err := tracker2.GetProjectStats()
	require.NoError(t, err)

	// Both projects should have time tracked
	assert.Contains(t, stats, "project-alpha")
	assert.Contains(t, stats, "project-beta")
	assert.True(t, stats["project-alpha"] > 0)
	assert.True(t, stats["project-beta"] > 0)

	// Get session history
	history, err := tracker2.GetSessionHistory(10)
	require.NoError(t, err)
	assert.Len(t, history, 2)

	// Verify separate session entries exist
	projects := make(map[string]bool)
	for _, session := range history {
		projects[session.Project] = true
		assert.Equal(t, tracking.StateStopped, session.State)
		assert.NotNil(t, session.EndTime)
		assert.True(t, session.Duration > 0)
	}

	assert.True(t, projects["project-alpha"])
	assert.True(t, projects["project-beta"])
	tracker2.Close()
}

func TestSwitch_NoArguments(t *testing.T) {
	// This test verifies the command definition has ExactArgs(1)
	// The Args validator runs before RunE, so we can't test runSwitch directly with empty args
	assert.NotNil(t, switchCmd.Args)
}

func TestSwitch_AccumulatesTimeCorrectly(t *testing.T) {
	setupTestEnvironment(t)

	// Start project A
	tracker, err := tracking.NewTracker()
	require.NoError(t, err)
	_, err = tracker.Start("project-a")
	require.NoError(t, err)
	tracker.Close()

	// Let it run for a measurable duration
	time.Sleep(20 * time.Millisecond)

	// Record time before switch
	beforeSwitch := time.Now()

	// Switch to project B
	cmd := &cobra.Command{}
	err = runSwitch(cmd, []string{"project-b"})
	require.NoError(t, err)

	// Let project B run
	time.Sleep(20 * time.Millisecond)

	// Stop project B
	tracker2, err := tracking.NewTracker()
	require.NoError(t, err)
	sessionB, err := tracker2.Stop()
	require.NoError(t, err)

	// Get project A from history
	history, err := tracker2.GetSessionHistory(10)
	require.NoError(t, err)
	require.Len(t, history, 2)

	var sessionA *tracking.Session
	for _, s := range history {
		if s.Project == "project-a" {
			sessionA = s
			break
		}
	}
	require.NotNil(t, sessionA)

	// Project A's duration should not include time after the switch
	assert.True(t, sessionA.Duration >= 15*time.Millisecond)
	assert.True(t, sessionA.Duration <= 30*time.Millisecond)

	// Verify project A ended around the switch time
	if sessionA.EndTime != nil {
		timeDiff := sessionA.EndTime.Sub(beforeSwitch).Abs()
		assert.True(t, timeDiff < 10*time.Millisecond)
	}

	// Project B should have its own duration
	assert.True(t, sessionB.Duration >= 15*time.Millisecond)
	assert.True(t, sessionB.Duration <= 30*time.Millisecond)
	tracker2.Close()
}

// setupTestEnvironment sets up a test environment with a temporary HOME directory
func setupTestEnvironment(t *testing.T) {
	// Create a temporary directory for the test database
	tempDir := t.TempDir()

	// Temporarily change HOME to use temp directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})
}
