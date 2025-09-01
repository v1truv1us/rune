package commands

import (
	"fmt"

	"github.com/ferg-cod3s/rune/internal/config"
	"github.com/ferg-cod3s/rune/internal/dnd"
	"github.com/ferg-cod3s/rune/internal/notifications"
	"github.com/ferg-cod3s/rune/internal/rituals"
	"github.com/ferg-cod3s/rune/internal/telemetry"
	"github.com/ferg-cod3s/rune/internal/tracking"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [project]",
	Short: "Start your workday and run start rituals",
	Long: `Start your workday timer and execute your configured start rituals.

This command will:
- Begin time tracking for your work session
- Execute global start rituals
- Execute project-specific start rituals (if detected)
- Launch interactive development environments (tmux sessions, terminals)
- Enable focus mode (Do Not Disturb) if configured

Interactive rituals create tmux sessions with multi-pane layouts or launch
interactive terminals for development work. Use 'rune ritual test start' to
preview what will be executed without running the commands.

If no project is specified, it will be auto-detected from the current directory.`,
	RunE: runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Wrap command with telemetry
	telemetry.WrapCommand(startCmd, runStart)
}

func runStart(cmd *cobra.Command, args []string) error {
	fmt.Println("🔮 Casting your start ritual...")

	// Load configuration to get idle threshold
	var tracker *tracking.Tracker
	cfg, configErr := config.Load()
	if configErr != nil {
		// Use default tracker if config fails to load
		var err error
		tracker, err = tracking.NewTracker()
		if err != nil {
			return fmt.Errorf("failed to initialize tracker: %w", err)
		}
	} else {
		// Use configured idle threshold
		var err error
		tracker, err = tracking.NewTrackerWithIdleThreshold(cfg.Settings.IdleThreshold)
		if err != nil {
			return fmt.Errorf("failed to initialize tracker: %w", err)
		}
	}
	defer tracker.Close()

	// Determine project name
	var project string
	if len(args) > 0 {
		// Use provided project name
		project = args[0]
	} else {
		// Auto-detect project
		detector := tracking.NewProjectDetector()
		project = detector.SanitizeProjectName(detector.DetectProject())
	}

	// Start time tracking
	session, err := tracker.Start(project)
	if err != nil {
		telemetry.TrackError(err, "start", map[string]interface{}{
			"project": project,
			"step":    "tracker_start",
		})
		return fmt.Errorf("failed to start session: %w", err)
	}

	// Track successful start
	telemetry.Track("session_started", map[string]interface{}{
		"project":       project,
		"auto_detected": len(args) == 0,
	})

	// Load configuration and execute start rituals (reuse cfg if already loaded)
	if cfg == nil {
		var err error
		cfg, err = config.Load()
		if err != nil {
			fmt.Printf("⚠ Could not load config for rituals: %v\n", err)
		}
	}
	if cfg != nil {
		engine := rituals.NewEngine(cfg)
		if err := engine.ExecuteStartRituals(project); err != nil {
			fmt.Printf("⚠ Start rituals failed: %v\n", err)
		}
	}

	// Check and enable Do Not Disturb if configured
	// Create notification manager based on config
	var notificationEnabled bool
	if cfg != nil {
		notificationEnabled = cfg.Settings.Notifications.Enabled
	}

	nm := notifications.NewNotificationManager(notificationEnabled)
	dndManager := dnd.NewDNDManager(nm)

	// Check if shortcuts are properly set up
	shortcutsReady, shortcutsErr := dndManager.CheckShortcutsSetup()
	if shortcutsErr != nil {
		fmt.Printf("⚠ Could not check Focus mode shortcuts: %v\n", shortcutsErr)
	} else if !shortcutsReady {
		fmt.Println("⚠ Focus mode shortcuts not set up")
		fmt.Println("💡 To enable automatic Focus mode control:")
		fmt.Println("   1. Open Shortcuts app")
		fmt.Println("   2. Create a new shortcut named 'Turn On Do Not Disturb'")
		fmt.Println("   3. Add action: 'Set Focus' → 'Do Not Disturb'")
		fmt.Println("   4. Create another shortcut named 'Turn Off Do Not Disturb'")
		fmt.Println("   5. Add action: 'Set Focus' → 'Turn Off Focus'")
		fmt.Println("   📖 See FOCUS_SETUP.md for detailed instructions")
	} else {
		// Shortcuts are set up, try to enable Focus mode
		if err := dndManager.Enable(); err != nil {
			fmt.Printf("⚠ Could not enable Do Not Disturb: %v\n", err)
		} else {
			fmt.Println("🎯 Focus mode enabled")
		}
	}

	fmt.Println("✓ Start ritual complete")
	fmt.Printf("⏰ Work timer started for project: %s\n", session.Project)

	return nil
}
