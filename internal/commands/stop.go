package commands

import (
	"fmt"

	"github.com/ferg-cod3s/rune/internal/dnd"
	"github.com/ferg-cod3s/rune/internal/notifications"
	"github.com/ferg-cod3s/rune/internal/rituals"
	"github.com/ferg-cod3s/rune/internal/telemetry"
	"github.com/ferg-cod3s/rune/internal/tracking"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "End your workday and run stop rituals",
	Long: `End your workday timer and execute your configured stop rituals.

This command will:
- Stop time tracking for your work session
- Execute global stop rituals
- Execute project-specific stop rituals (if detected)
- Disable focus mode (Do Not Disturb) if configured
- Generate a summary of your work session`,
	RunE: runStop,
}

func init() {
	rootCmd.AddCommand(stopCmd)

	// Wrap command with telemetry
	telemetry.WrapCommand(stopCmd, runStop)
}

func runStop(cmd *cobra.Command, args []string) error {
	fmt.Println("🔮 Casting your stop ritual...")

	// Initialize tracker
	tracker, err := tracking.NewTracker()
	if err != nil {
		return fmt.Errorf("failed to initialize tracker: %w", err)
	}
	defer tracker.Close()

	// Stop time tracking
	session, err := tracker.Stop()
	if err != nil {
		telemetry.TrackError(err, "stop", map[string]interface{}{
			"step": "tracker_stop",
		})
		return fmt.Errorf("failed to stop session: %w", err)
	}

	// Track successful stop
	telemetry.Track("session_stopped", map[string]interface{}{
		"project":  session.Project,
		"duration": session.Duration.Milliseconds(),
	})

	// Load configuration and execute stop rituals
	cfg, err := loadConfigWithProfile()
	if err != nil {
		fmt.Printf("⚠ Could not load config for rituals: %v\n", err)
	} else {
		engine := rituals.NewEngine(cfg)
		if err := engine.ExecuteStopRituals(session.Project); err != nil {
			fmt.Printf("⚠ Stop rituals failed: %v\n", err)
		}
	}

	// Disable Do Not Disturb
	var notificationEnabled bool
	if cfg != nil {
		notificationEnabled = cfg.Settings.Notifications.Enabled
	}

	nm := notifications.NewNotificationManager(notificationEnabled)
	dndManager := dnd.NewDNDManager(nm)
	if err := dndManager.Disable(); err != nil {
		fmt.Printf("⚠ Could not disable Do Not Disturb: %v\n", err)
	} else {
		fmt.Println("🎯 Focus mode disabled")
	}

	fmt.Println("✓ Stop ritual complete")
	fmt.Println("⏰ Work timer stopped")
	fmt.Printf("📊 Session summary: %s (project: %s)\n",
		formatDuration(session.Duration), session.Project)

	return nil
}
