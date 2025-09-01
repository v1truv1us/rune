package commands

import (
	"fmt"

	"github.com/ferg-cod3s/rune/internal/config"
	"github.com/ferg-cod3s/rune/internal/rituals"
	"github.com/ferg-cod3s/rune/internal/tracking"
	"github.com/spf13/cobra"
)

var ritualCmd = &cobra.Command{
	Use:   "ritual",
	Short: "Manage and test rituals",
	Long: `Manage and test your configured rituals.
	
This command provides subcommands to list, run, and test rituals
without affecting your time tracking. Interactive rituals (those with
'interactive: true') create tmux sessions or terminal environments
for development work.

Use 'ritual test' to preview interactive commands before execution.`,
}

var ritualListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available rituals",
	Long:  `List all configured rituals for all projects.`,
	RunE:  runRitualList,
}

var ritualTestCmd = &cobra.Command{
	Use:   "test <start|stop> [project]",
	Short: "Test a ritual without executing it",
	Long:  `Test a ritual configuration without actually executing the commands.

This shows what commands would run, including interactive rituals that would
create tmux sessions or launch terminals. Useful for validating configuration
before execution.`,
	Args:  cobra.RangeArgs(1, 2),
	RunE:  runRitualTest,
}

var ritualRunCmd = &cobra.Command{
	Use:   "run <start|stop> [project]",
	Short: "Run a specific ritual",
	Long:  `Run a specific ritual without affecting time tracking.`,
	Args:  cobra.RangeArgs(1, 2),
	RunE:  runRitualRun,
}

func init() {
	rootCmd.AddCommand(ritualCmd)
	ritualCmd.AddCommand(ritualListCmd)
	ritualCmd.AddCommand(ritualTestCmd)
	ritualCmd.AddCommand(ritualRunCmd)
}

func runRitualList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("🔮 Configured Rituals")
	fmt.Println("====================")
	fmt.Println()

	// Show start rituals
	fmt.Println("Start Rituals:")
	if len(cfg.Rituals.Start.Global) > 0 {
		fmt.Println("  Global:")
		for _, cmd := range cfg.Rituals.Start.Global {
			fmt.Printf("    - %s: %s\n", cmd.Name, cmd.Command)
		}
	}

	if len(cfg.Rituals.Start.PerProject) > 0 {
		fmt.Println("  Per-Project:")
		for project, commands := range cfg.Rituals.Start.PerProject {
			fmt.Printf("    %s:\n", project)
			for _, cmd := range commands {
				fmt.Printf("      - %s: %s\n", cmd.Name, cmd.Command)
			}
		}
	}

	fmt.Println()

	// Show stop rituals
	fmt.Println("Stop Rituals:")
	if len(cfg.Rituals.Stop.Global) > 0 {
		fmt.Println("  Global:")
		for _, cmd := range cfg.Rituals.Stop.Global {
			fmt.Printf("    - %s: %s\n", cmd.Name, cmd.Command)
		}
	}

	if len(cfg.Rituals.Stop.PerProject) > 0 {
		fmt.Println("  Per-Project:")
		for project, commands := range cfg.Rituals.Stop.PerProject {
			fmt.Printf("    %s:\n", project)
			for _, cmd := range commands {
				fmt.Printf("      - %s: %s\n", cmd.Name, cmd.Command)
			}
		}
	}

	return nil
}

func runRitualTest(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ritualType := args[0]
	var project string

	if len(args) > 1 {
		project = args[1]
	} else {
		// Auto-detect project
		detector := tracking.NewProjectDetector()
		project = detector.SanitizeProjectName(detector.DetectProject())
	}

	engine := rituals.NewEngine(cfg)
	return engine.TestRitual(ritualType, project)
}

func runRitualRun(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ritualType := args[0]
	var project string

	if len(args) > 1 {
		project = args[1]
	} else {
		// Auto-detect project
		detector := tracking.NewProjectDetector()
		project = detector.SanitizeProjectName(detector.DetectProject())
	}

	engine := rituals.NewEngine(cfg)

	switch ritualType {
	case "start":
		return engine.ExecuteStartRituals(project)
	case "stop":
		return engine.ExecuteStopRituals(project)
	default:
		return fmt.Errorf("unknown ritual type: %s (use 'start' or 'stop')", ritualType)
	}
}
