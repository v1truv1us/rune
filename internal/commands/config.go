package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ferg-cod3s/rune/internal/config"
	"github.com/ferg-cod3s/rune/internal/logger"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Rune configuration",
	Long: `Manage your Rune configuration file.

This command provides subcommands to edit, validate, and manage
your Rune configuration.`,
	RunE: runConfigRoot,
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the configuration file",
	Long:  `Open the configuration file in your default editor.`,
	RunE:  runConfigEdit,
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the configuration file",
	Long:  `Validate the syntax and content of your configuration file.`,
	RunE:  runConfigValidate,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current configuration",
	Long:  `Display the current configuration with resolved values.`,
	RunE:  runConfigShow,
}

var configSetupTelemetryCmd = &cobra.Command{
	Use:   "setup-telemetry",
	Short: "Setup telemetry integration (Sentry)",
	Long: `Configure Sentry DSN for telemetry integration.

This command allows you to easily set up error tracking:
- Sentry: For error tracking, performance monitoring, and crash reporting

Future updates will include OpenTelemetry integration for logs and traces.`,
	RunE: runConfigSetupTelemetry,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configValidateCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetupTelemetryCmd)

	// Add flags for telemetry setup
	configSetupTelemetryCmd.Flags().String("sentry-dsn", "", "Sentry DSN for error tracking")
	configSetupTelemetryCmd.Flags().Bool("enable", true, "Enable telemetry after setting keys")
	configSetupTelemetryCmd.Flags().Bool("disable", false, "Disable telemetry")
	configSetupTelemetryCmd.Flags().Bool("interactive", false, "Interactive setup mode")
	configSetupTelemetryCmd.Flags().Bool("show-examples", false, "Show example values and setup instructions")

	// Add convenient aliases for the main flags
	configCmd.PersistentFlags().String("add-sentry-dsn", "", "Quick setup: Add Sentry DSN")
}

func runConfigRoot(cmd *cobra.Command, args []string) error {
	// Handle quick setup flags
	addSentryDSN, _ := cmd.Flags().GetString("add-sentry-dsn")

	if addSentryDSN != "" {
		// Quick setup mode - delegate to setup-telemetry command
		setupCmd := configSetupTelemetryCmd
		_ = setupCmd.Flags().Set("sentry-dsn", addSentryDSN)
		return runConfigSetupTelemetry(setupCmd, args)
	}

	// Default behavior - show help
	return cmd.Help()
}

func runConfigEdit(cmd *cobra.Command, args []string) error {
	configPath, err := config.GetConfigPath()
	if err != nil {
		return err
	}

	exists, err := config.Exists()
	if err != nil {
		return err
	}
	if !exists {
		fmt.Println("⚠ Configuration file does not exist.")
		fmt.Println("Run 'rune init' to create a new configuration.")
		return nil
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi" // fallback to vi
	}

	cmd_exec := exec.Command(editor, configPath)
	cmd_exec.Stdin = os.Stdin
	cmd_exec.Stdout = os.Stdout
	cmd_exec.Stderr = os.Stderr

	return cmd_exec.Run()
}

func runConfigValidate(cmd *cobra.Command, args []string) error {
	exists, err := config.Exists()
	if err != nil {
		return err
	}
	if !exists {
		fmt.Println("⚠ Configuration file does not exist.")
		fmt.Println("Run 'rune init' to create a new configuration.")
		return nil
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Configuration validation failed: %v\n", err)
		return nil // Don't return error to avoid double error message
	}

	fmt.Println("✅ Configuration is valid!")
	fmt.Printf("   Version: %d\n", cfg.Version)
	fmt.Printf("   Projects: %d\n", len(cfg.Projects))
	fmt.Printf("   Work hours: %.1f\n", cfg.Settings.WorkHours)

	return nil
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	exists, err := config.Exists()
	if err != nil {
		return err
	}
	if !exists {
		fmt.Println("⚠ Configuration file does not exist.")
		fmt.Println("Run 'rune init' to create a new configuration.")
		return nil
	}

	configPath, err := config.GetConfigPath()
	if err != nil {
		return err
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	fmt.Printf("Configuration file: %s\n", configPath)
	fmt.Println("=" + string(make([]byte, len(configPath)+20)))
	fmt.Print(string(content))

	return nil
}

func runConfigSetupTelemetry(cmd *cobra.Command, args []string) error {
	sentryDSN, _ := cmd.Flags().GetString("sentry-dsn")
	enableTelemetry, _ := cmd.Flags().GetBool("enable")
	disableTelemetry, _ := cmd.Flags().GetBool("disable")
	interactive, _ := cmd.Flags().GetBool("interactive")
	showExamples, _ := cmd.Flags().GetBool("show-examples")

	if showExamples {
		return showTelemetryExamples()
	}

	// Load existing config or create new one
	cfg, err := config.Load()
	if err != nil {
		// If config doesn't exist, create a default one
		cfg = &config.Config{
			Version: 1,
			Settings: config.Settings{
				WorkHours:     8.0,
				BreakInterval: 25 * 60 * 1000000000, // 25 minutes in nanoseconds
				IdleThreshold: 5 * 60 * 1000000000,  // 5 minutes in nanoseconds
			},
			Integrations: config.Integrations{
				Telemetry: config.TelemetryIntegration{
					Enabled: true,
				},
			},
		}
	}

	if interactive {
		err := runInteractiveTelemetrySetup(cfg)
		if err != nil {
			return err
		}
		// Save the config after interactive setup
		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}
		fmt.Println("💾 Configuration saved successfully!")
		return nil
	}

	if disableTelemetry {
		cfg.Integrations.Telemetry.Enabled = false
		fmt.Println("🔕 Telemetry disabled")
	} else {
		// Update keys if provided
		updated := false

		if sentryDSN != "" {
			if !isValidSentryDSN(sentryDSN) {
				return fmt.Errorf("invalid Sentry DSN format. Expected format: https://public_key@sentry.io/project_id")
			}
			cfg.Integrations.Telemetry.SentryDSN = sentryDSN
			fmt.Println("✅ Sentry DSN configured")
			updated = true
		}

		if updated {
			if enableTelemetry {
				cfg.Integrations.Telemetry.Enabled = true
				fmt.Println("📊 Telemetry enabled")
			}
		} else if !cmd.Flags().Changed("enable") && !cmd.Flags().Changed("disable") {
			// Show current status if no changes were made
			return showTelemetryStatus(cfg)
		}
	}

	// Save the updated configuration
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Println("💾 Configuration saved successfully!")

	// Log the telemetry setup
	logger.LogStructuredEvent("info", "Telemetry configuration updated", "config", "setup-telemetry", map[string]interface{}{
		"sentry_configured": cfg.Integrations.Telemetry.SentryDSN != "",
		"telemetry_enabled": cfg.Integrations.Telemetry.Enabled,
	})

	return nil
}

func runInteractiveTelemetrySetup(cfg *config.Config) error {
	fmt.Println("🔧 Interactive Telemetry Setup")
	fmt.Println("==============================")
	fmt.Println()

	// Show current status
	fmt.Println("Current Configuration:")
	fmt.Printf("  Telemetry: %s\n", map[bool]string{true: "✅ Enabled", false: "❌ Disabled"}[cfg.Integrations.Telemetry.Enabled])
	fmt.Printf("  Sentry DSN: %s\n", maskTelemetryKey(cfg.Integrations.Telemetry.SentryDSN))
	fmt.Println()

	// Ask about Sentry DSN
	fmt.Println("📍 Sentry Setup (Error Tracking)")
	fmt.Println("Sentry helps track errors and performance issues.")
	fmt.Println("Get your DSN from: https://sentry.io/settings/projects/")
	fmt.Print("Enter Sentry DSN (or press Enter to skip): ")

	var sentryDSN string
	_, _ = fmt.Scanln(&sentryDSN)
	sentryDSN = strings.TrimSpace(sentryDSN)

	if sentryDSN != "" {
		if !isValidSentryDSN(sentryDSN) {
			fmt.Println("⚠️  Warning: DSN format looks invalid, but continuing...")
		}
		cfg.Integrations.Telemetry.SentryDSN = sentryDSN
		fmt.Println("✅ Sentry DSN configured")
	}

	fmt.Println()

	// Ask about enabling telemetry
	hasAnyKey := cfg.Integrations.Telemetry.SentryDSN != ""
	if hasAnyKey {
		fmt.Print("Enable telemetry? [Y/n]: ")
		var enable string
		_, _ = fmt.Scanln(&enable)
		enable = strings.ToLower(strings.TrimSpace(enable))

		cfg.Integrations.Telemetry.Enabled = enable != "n" && enable != "no"
		fmt.Printf("📊 Telemetry %s\n", map[bool]string{true: "enabled", false: "disabled"}[cfg.Integrations.Telemetry.Enabled])
	} else {
		fmt.Println("ℹ️  No telemetry keys configured, telemetry will remain disabled")
		cfg.Integrations.Telemetry.Enabled = false
	}

	return nil
}

func showTelemetryStatus(cfg *config.Config) error {
	fmt.Println("📊 Current Telemetry Status")
	fmt.Println("==========================")
	fmt.Printf("Enabled: %s\n", map[bool]string{true: "✅ Yes", false: "❌ No"}[cfg.Integrations.Telemetry.Enabled])
	fmt.Printf("Sentry DSN: %s\n", maskTelemetryKey(cfg.Integrations.Telemetry.SentryDSN))
	fmt.Println()
	fmt.Println("Use flags to update:")
	fmt.Println("  --sentry-dsn <dsn>     Set Sentry DSN")
	fmt.Println("  --enable               Enable telemetry")
	fmt.Println("  --disable              Disable telemetry")
	fmt.Println("  --interactive          Interactive setup")

	return nil
}

func isValidSentryDSN(dsn string) bool {
	// Basic validation for Sentry DSN format
	// Expected: https://public_key@sentry.io/project_id or similar
	return strings.HasPrefix(dsn, "https://") && strings.Contains(dsn, "@") && strings.Contains(dsn, "/")
}

func maskTelemetryKey(key string) string {
	if key == "" {
		return "❌ Not configured"
	}
	if len(key) < 8 {
		return "✅ Configured"
	}
	return "✅ " + key[:4] + "****" + key[len(key)-4:]
}

func showTelemetryExamples() error {
	fmt.Println("🔧 Telemetry Setup Examples & Instructions")
	fmt.Println("==========================================")
	fmt.Println()

	fmt.Println("📍 Sentry (Error Tracking)")
	fmt.Println("---------------------------")
	fmt.Println("Sentry helps track errors, performance issues, and crashes.")
	fmt.Println("Future updates will integrate OpenTelemetry for comprehensive observability.")
	fmt.Println()
	fmt.Println("Setup:")
	fmt.Println("1. Go to https://sentry.io/")
	fmt.Println("2. Create a new project or select an existing one")
	fmt.Println("3. Go to Settings → Projects → [Your Project] → Client Keys (DSN)")
	fmt.Println("4. Copy your DSN")
	fmt.Println()
	fmt.Println("DSN Format:")
	fmt.Println("  https://public_key@sentry.io/project_id")
	fmt.Println("  https://abc123def456@sentry.io/123456")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Printf("  %s\n", "rune config --add-sentry-dsn https://your_key@sentry.io/project_id")
	fmt.Printf("  %s\n", "rune config setup-telemetry --sentry-dsn https://your_key@sentry.io/project_id")
	fmt.Println()

	fmt.Println("🚀 Setup Examples")
	fmt.Println("-----------------")
	fmt.Printf("Interactive setup:\n")
	fmt.Printf("  %s\n", "rune config setup-telemetry --interactive")
	fmt.Println()
	fmt.Printf("Quick setup shortcut:\n")
	fmt.Printf("  %s\n", "rune config --add-sentry-dsn https://your_key@sentry.io/project_id")
	fmt.Println()

	fmt.Println("🔮 Coming Soon: OpenTelemetry Integration")
	fmt.Println("-----------------------------------------")
	fmt.Println("Future versions will include:")
	fmt.Println("• OpenTelemetry logging, tracing, and metrics")
	fmt.Println("• OTLP exporter for sending logs to any OpenTelemetry collector")
	fmt.Println("• Distributed tracing correlation")
	fmt.Println("• Integration with popular observability platforms")
	fmt.Println()

	fmt.Println("💡 Security Notes")
	fmt.Println("-----------------")
	fmt.Println("• Keys are stored in your local config file (~/.rune/config.yaml)")
	fmt.Println("• Keys are never embedded in public binaries")
	fmt.Println("• Environment variables (RUNE_SENTRY_DSN) override config")
	fmt.Println("• You can disable telemetry anytime with: rune config setup-telemetry --disable")
	fmt.Println("• For developers: Use .env files or environment variables")
	fmt.Println("• For production: Keys can be embedded at build-time for official releases")
	fmt.Println()

	fmt.Println("📋 Current Status")
	fmt.Println("-----------------")

	if cfg, err := config.Load(); err == nil {
		fmt.Printf("Telemetry: %s\n", map[bool]string{true: "✅ Enabled", false: "❌ Disabled"}[cfg.Integrations.Telemetry.Enabled])
		fmt.Printf("Sentry DSN: %s\n", maskTelemetryKey(cfg.Integrations.Telemetry.SentryDSN))
	} else {
		fmt.Println("❌ No configuration file found. Run 'rune init' first.")
	}

	return nil
}
