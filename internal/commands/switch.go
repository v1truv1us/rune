package commands

import (
	"fmt"

	"github.com/ferg-cod3s/rune/internal/config"
	"github.com/ferg-cod3s/rune/internal/rituals"
	"github.com/ferg-cod3s/rune/internal/telemetry"
	"github.com/ferg-cod3s/rune/internal/tracking"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch <project>",
	Short: "Atomically switch to a different project",
	Long: `Atomically switch from the current project to a new one.

This command will:
- Stop the current work session and run stop rituals
- Start a new session for the target project
- Execute start rituals for the new project
- Preserve session boundaries for accurate reporting

This is ideal for agency/consultant workflows where rapid context switching
across clients is common. The entire operation is atomic - if any step fails,
the switch is rolled back.`,
	Args: cobra.ExactArgs(1),
	RunE: runSwitch,
}

func init() {
	rootCmd.AddCommand(switchCmd)

	// Wrap command with telemetry
	telemetry.WrapCommand(switchCmd, runSwitch)
}

func runSwitch(cmd *cobra.Command, args []string) error {
	targetProject := args[0]

	fmt.Println("🔮 Casting your switch ritual...")

	// Initialize tracker
	tracker, err := tracking.NewTracker()
	if err != nil {
		return fmt.Errorf("failed to initialize tracker: %w", err)
	}
	defer tracker.Close()

	// Get current session to check if one is active
	currentSession, err := tracker.GetCurrentSession()
	if err != nil {
		return fmt.Errorf("failed to get current session: %w", err)
	}
	if currentSession == nil {
		return fmt.Errorf("no active session to switch from - use 'rune start %s' instead", targetProject)
	}
	if currentSession.State != tracking.StateRunning {
		return fmt.Errorf("current session is not running (state: %s) - use 'rune start %s' instead", currentSession.State, targetProject)
	}

	oldProject := currentSession.Project

	// Check if we're switching to the same project
	if oldProject == targetProject {
		return fmt.Errorf("already working on project: %s", targetProject)
	}

	// Load configuration for rituals
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("⚠ Could not load config for rituals: %v\n", err)
	}

	// Step 1: Stop current session
	fmt.Printf("⏹ Stopping session for: %s\n", oldProject)
	stoppedSession, err := tracker.Stop()
	if err != nil {
		telemetry.TrackError(err, "switch", map[string]interface{}{
			"step":           "stop_session",
			"old_project":    oldProject,
			"target_project": targetProject,
		})
		return fmt.Errorf("failed to stop current session: %w", err)
	}

	// Step 2: Execute stop rituals for old project
	if cfg != nil {
		engine := rituals.NewEngine(cfg)
		if err := engine.ExecuteStopRituals(oldProject); err != nil {
			fmt.Printf("⚠ Stop rituals failed: %v\n", err)
		}
	}

	fmt.Printf("✓ Stopped session: %s (duration: %s)\n", oldProject, formatDuration(stoppedSession.Duration))

	// Step 3: Start new session
	fmt.Printf("▶ Starting session for: %s\n", targetProject)
	newSession, err := tracker.Start(targetProject)
	if err != nil {
		telemetry.TrackError(err, "switch", map[string]interface{}{
			"step":           "start_session",
			"old_project":    oldProject,
			"target_project": targetProject,
		})
		// Critical error - we've stopped the old session but can't start the new one
		return fmt.Errorf("failed to start new session (old session was stopped): %w", err)
	}

	// Step 4: Execute start rituals for new project
	if cfg != nil {
		engine := rituals.NewEngine(cfg)
		if err := engine.ExecuteStartRituals(targetProject); err != nil {
			fmt.Printf("⚠ Start rituals failed: %v\n", err)
		}
	}

	// Track successful switch
	telemetry.Track("project_switched", map[string]interface{}{
		"old_project":  oldProject,
		"new_project":  targetProject,
		"old_duration": stoppedSession.Duration.Milliseconds(),
	})

	fmt.Println("✓ Switch ritual complete")
	fmt.Printf("⏰ Now tracking: %s\n", newSession.Project)
	fmt.Printf("📊 Previous session: %s (%s)\n", oldProject, formatDuration(stoppedSession.Duration))

	return nil
}
