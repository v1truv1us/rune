package commands

import (
	"fmt"

	"github.com/ferg-cod3s/rune/internal/colors"
	"github.com/ferg-cod3s/rune/internal/config"
	"github.com/ferg-cod3s/rune/internal/dnd"
	"github.com/ferg-cod3s/rune/internal/notifications"
	"github.com/ferg-cod3s/rune/internal/telemetry"
	"github.com/ferg-cod3s/rune/internal/tracking"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current session status",
	Long: `Display the current status of your work session.

This command shows:
- Current timer state (running, paused, stopped)
- Active project (if detected)
- Session duration
- Today's total work time
- Focus mode status`,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Wrap command with telemetry
	telemetry.WrapCommand(statusCmd, runStatus)
}

func runStatus(cmd *cobra.Command, args []string) error {
	fmt.Println(colors.Header("📊 Current Session Status"))
	fmt.Println(colors.Secondary("========================"))
	fmt.Println()

	// Initialize tracker
	tracker, err := tracking.NewTracker()
	if err != nil {
		return fmt.Errorf("failed to initialize tracker: %w", err)
	}
	defer tracker.Close()

	// Get current session
	session, err := tracker.GetCurrentSession()
	if err != nil {
		return fmt.Errorf("failed to get current session: %w", err)
	}

	if session == nil {
		fmt.Printf("Timer:        %s\n", colors.StatusStopped("Stopped"))
		fmt.Printf("Project:      %s\n", colors.Muted("Not detected"))
		fmt.Printf("Session:      %s\n", colors.Duration("0h 0m"))
	} else {
		duration, err := tracker.GetSessionDuration()
		if err != nil {
			return fmt.Errorf("failed to get session duration: %w", err)
		}

		var timerStatus string
		switch session.State {
		case tracking.StateRunning:
			timerStatus = colors.StatusRunning("Running")
		case tracking.StatePaused:
			timerStatus = colors.StatusPaused("Paused")
		case tracking.StateStopped:
			timerStatus = colors.StatusStopped("Stopped")
		default:
			timerStatus = colors.Muted(session.State.String())
		}

		fmt.Printf("Timer:        %s\n", timerStatus)
		fmt.Printf("Project:      %s\n", colors.Project(session.Project))
		fmt.Printf("Session:      %s %s\n", colors.Duration(formatDuration(duration)), colors.RelativeTime("started "+formatRelativeTime(session.StartTime)))
	}

	// Get daily total
	dailyTotal, err := tracker.GetDailyTotal()
	if err != nil {
		return fmt.Errorf("failed to get daily total: %w", err)
	}
	fmt.Printf("Today Total:  %s\n", colors.Duration(formatDuration(dailyTotal)))

	// Get idle status
	isIdle, err := tracker.IsIdle()
	if err != nil {
		fmt.Printf("Idle Status:  %s\n", colors.Warning("Unknown (detection failed)"))
	} else {
		if isIdle {
			idleTime, err := tracker.GetIdleTime()
			if err == nil {
				fmt.Printf("Idle Status:  %s for %s\n", colors.Warning("Idle"), colors.Duration(formatDuration(idleTime)))
			} else {
				fmt.Printf("Idle Status:  %s\n", colors.Warning("Idle"))
			}
		} else {
			fmt.Printf("Idle Status:  %s\n", colors.Success("Active"))
		}
	}

	// Check DND status
	cfg, _ := config.Load()
	var notificationEnabled bool
	if cfg != nil {
		notificationEnabled = cfg.Settings.Notifications.Enabled
	}

	nm := notifications.NewNotificationManager(notificationEnabled)
	dndManager := dnd.NewDNDManager(nm)

	// Check if shortcuts are set up
	shortcutsOK, shortcutsErr := dndManager.CheckShortcutsSetup()

	if shortcutsErr != nil {
		fmt.Printf("Focus Mode:   %s\n", colors.Muted("Not supported on this platform"))
	} else if !shortcutsOK {
		fmt.Printf("Focus Mode:   %s\n", colors.Warning("Not configured (run 'rune start' for setup)"))
	} else {
		// Try to detect current status
		dndEnabled, err := dndManager.IsEnabled()
		if err != nil {
			fmt.Printf("Focus Mode:   %s\n", colors.Glow("Available (detection unavailable)"))
		} else {
			if dndEnabled {
				fmt.Printf("Focus Mode:   %s\n", colors.Success("Enabled"))
			} else {
				fmt.Printf("Focus Mode:   %s\n", colors.Glow("Available (currently off)"))
			}
		}
	}

	return nil
}
