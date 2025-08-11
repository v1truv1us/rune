package commands

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/ferg-cod3s/rune/internal/config"
	"github.com/ferg-cod3s/rune/internal/telemetry"
	"github.com/getsentry/sentry-go"
	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug and diagnostic commands",
	Long:  "Debug and diagnostic commands for troubleshooting Rune issues",
}

var debugTelemetryCmd = &cobra.Command{
	Use:   "telemetry",
	Short: "Debug telemetry configuration and connectivity",
	Long:  "Show telemetry configuration, test connectivity, and verify event delivery",
	RunE:  runDebugTelemetry,
}

var debugKeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Show telemetry key configuration (masked for security)",
	Long:  "Display telemetry API keys and DSNs with masking for security",
	RunE:  runDebugKeys,
}

var debugNotificationsCmd = &cobra.Command{
	Use:   "notifications",
	Short: "Debug OS notifications setup",
	Long:  "Diagnose notification tooling on your OS and provide setup guidance.",
	RunE:  runDebugNotifications,
}

func init() {
	rootCmd.AddCommand(debugCmd)
	debugCmd.AddCommand(debugTelemetryCmd)
	debugCmd.AddCommand(debugKeysCmd)
	debugCmd.AddCommand(debugNotificationsCmd)

	// Wrap debug commands with telemetry
	telemetry.WrapCommand(debugTelemetryCmd, runDebugTelemetry)
	telemetry.WrapCommand(debugKeysCmd, runDebugKeys)
	telemetry.WrapCommand(debugNotificationsCmd, runDebugNotifications)
}

func runDebugTelemetry(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Rune Telemetry Debug Report")
	fmt.Println("=" + strings.Repeat("=", 40))

	// System Information
	fmt.Printf("\n📊 System Information:\n")
	fmt.Printf("  OS: %s\n", runtime.GOOS)
	fmt.Printf("  Architecture: %s\n", runtime.GOARCH)
	fmt.Printf("  Go Version: %s\n", runtime.Version())

	// Environment Variables
	fmt.Printf("\n🌍 Environment Variables:\n")
	fmt.Printf("  RUNE_TELEMETRY_DISABLED: %s\n", getEnvOrDefault("RUNE_TELEMETRY_DISABLED", "not set"))
	fmt.Printf("  RUNE_DEBUG: %s\n", getEnvOrDefault("RUNE_DEBUG", "not set"))
	fmt.Printf("  RUNE_ENV: %s\n", getEnvOrDefault("RUNE_ENV", "not set"))
	fmt.Printf("  RUNE_OTLP_ENDPOINT: %s\n", getEnvOrDefault("RUNE_OTLP_ENDPOINT", "not set"))
	fmt.Printf("  RUNE_SENTRY_DSN: %s\n", maskDSN(os.Getenv("RUNE_SENTRY_DSN")))

	// Configuration File
	fmt.Printf("\n📄 Configuration:\n")
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("  Config Load Error: %v\n", err)
	} else {
		fmt.Printf("  Config File Found: ✅\n")
		fmt.Printf("  Telemetry Enabled: %t\n", cfg.Integrations.Telemetry.Enabled)
		fmt.Printf("  Sentry DSN (config): %s\n", maskDSN(cfg.Integrations.Telemetry.SentryDSN))
		fmt.Printf("  User ID: %s\n", cfg.UserID)
	}

	// Build-time Keys (check if embedded)
	fmt.Printf("\n🔧 Build-time Configuration:\n")
	fmt.Printf("  Build-time keys embedded: %s\n", checkBuildTimeKeys())

	// Network Connectivity Tests
	fmt.Printf("\n🌐 Network Connectivity:\n")
	if os.Getenv("RUNE_OTLP_ENDPOINT") != "" {
		testConnectivity("OTLP endpoint", os.Getenv("RUNE_OTLP_ENDPOINT"))
	} else {
		fmt.Printf("  OTLP endpoint: not set (set RUNE_OTLP_ENDPOINT to enable OTLP logging)\n")
	}
	testConnectivity("Sentry API", "https://sentry.io/api/")

	// Test Event Sending
	fmt.Printf("\n📡 Test Event Sending:\n")
	testEventSending()

	// Send a Sentry test message if DSN is configured
	sentryDSN := os.Getenv("RUNE_SENTRY_DSN")
	if sentryDSN == "" && cfg != nil {
		sentryDSN = cfg.Integrations.Telemetry.SentryDSN
	}
	if sentryDSN != "" {
		fmt.Printf("\n🧪 Sending Sentry test message...\n")
		telemetry.CaptureMessage("rune debug telemetry test message", sentry.LevelInfo, map[string]string{
			"component": "debug",
			"command":   "debug telemetry",
		})
		fmt.Printf("  ✅ Sentry test message enqueued (check your Sentry project)\n")
	} else {
		fmt.Printf("\nℹ️  Sentry DSN not configured; skipping Sentry test message.\n")
	}

	return nil
}

func runDebugKeys(cmd *cobra.Command, args []string) error {
	fmt.Println("🔑 Rune API Keys Debug")
	fmt.Println("=" + strings.Repeat("=", 30))

	// Environment Variables
	sentryEnv := os.Getenv("RUNE_SENTRY_DSN")

	fmt.Printf("\n🌍 Environment Variables:\n")
	fmt.Printf("  RUNE_SENTRY_DSN: %s\n", maskDSN(sentryEnv))

	// Configuration File
	cfg, err := config.Load()
	if err == nil {
		fmt.Printf("\n📄 Configuration File:\n")
		fmt.Printf("  Sentry DSN: %s\n", maskDSN(cfg.Integrations.Telemetry.SentryDSN))
	}

	finalSentry := sentryEnv
	if finalSentry == "" && cfg != nil {
		finalSentry = cfg.Integrations.Telemetry.SentryDSN
	}

	fmt.Printf("\n🎯 Final Resolution:\n")
	fmt.Printf("  Active Sentry DSN: %s\n", maskDSN(finalSentry))

	// Validation
	fmt.Printf("\n🔍 Validation:\n")

	fmt.Printf("\n✅ Final Resolution:\n")
	fmt.Printf("  Active Sentry DSN: %s\n", maskDSN(finalSentry))

	if finalSentry == "" {
		fmt.Printf("  ❌ No Sentry DSN configured\n")
	} else {
		fmt.Printf("  ✅ Sentry DSN configured (%d chars)\n", len(finalSentry))
	}

	return nil
}

func runDebugNotifications(cmd *cobra.Command, args []string) error {
	fmt.Println("🔔 Rune Notifications Debug")
	fmt.Println("=" + strings.Repeat("=", 32))

	fmt.Printf("\n📊 System Information:\n")
	fmt.Printf("  OS: %s\n", runtime.GOOS)

	switch runtime.GOOS {
	case "darwin":
		fmt.Printf("\n🍎 macOS Checks:\n")
		// Check terminal-notifier
		if _, err := exec.LookPath("terminal-notifier"); err != nil {
			fmt.Printf("  ❌ terminal-notifier not found\n")
			fmt.Println("  ➜ Install with: brew install terminal-notifier")
		} else {
			fmt.Printf("  ✅ terminal-notifier found\n")
		}
		// Check Shortcuts CLI
		if _, err := exec.LookPath("shortcuts"); err != nil {
			fmt.Printf("  ❌ shortcuts CLI not found\n")
			fmt.Println("  ➜ Ensure macOS Shortcuts app is installed (CLI is built-in on recent macOS)")
		} else {
			fmt.Printf("  ✅ shortcuts CLI found\n")
		}
		fmt.Println("\n  Focus tips:")
		fmt.Println("  - Add your terminal (Terminal, iTerm2, Ghostty) to Allowed Apps in System Settings > Focus")
		fmt.Println("  - Some Focus modes may still silence notifications; use -ignoreDnD when critical")
	case "linux":
		fmt.Printf("\n🐧 Linux Checks:\n")
		// notify-send
		if _, err := exec.LookPath("notify-send"); err != nil {
			fmt.Printf("  ❌ notify-send not found\n")
			fmt.Println("  ➜ Install libnotify-bin (Debian/Ubuntu) or libnotify (Arch/Fedora)")
		} else {
			fmt.Printf("  ✅ notify-send found\n")
		}
		// dunstctl
		if _, err := exec.LookPath("dunstctl"); err == nil {
			fmt.Printf("  ✅ dunstctl found (dunst daemon)\n")
		}
		// makoctl
		if _, err := exec.LookPath("makoctl"); err == nil {
			fmt.Printf("  ✅ makoctl found (mako daemon)\n")
		}
		fmt.Println("\n  Notes:")
		fmt.Println("  - GNOME Shell ignores expire-time; urgency critical may be required for persistent alerts")
		fmt.Println("  - Ensure a notification daemon (dunst, mako, notification-daemon) is running")
	case "windows":
		fmt.Printf("\n🪟 Windows Checks:\n")
		// PowerShell presence
		if _, err := exec.LookPath("powershell"); err != nil {
			fmt.Printf("  ❌ PowerShell not found in PATH\n")
			fmt.Println("  ➜ Ensure Windows PowerShell is available (default on Windows 10/11)")
		} else {
			fmt.Printf("  ✅ PowerShell found\n")
		}
		fmt.Println("\n  Focus Assist tips:")
		fmt.Println("  - Settings > System > Focus assist: set Priority only and configure priority senders")
		fmt.Println("  - Registry keys under HKCU\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\QuietHours control status")
	default:
		fmt.Println("Unsupported OS for notifications debug")
	}

	fmt.Println("\n📎 Next steps:")
	fmt.Println("  - Run: rune test notifications")
	fmt.Println("  - Run: rune test dnd")
	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func maskKey(key string) string {
	if key == "" {
		return "not set"
	}
	if len(key) <= 8 {
		return strings.Repeat("*", len(key))
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}

func maskDSN(dsn string) string {
	if dsn == "" {
		return "not set"
	}
	// For Sentry DSN format: https://public_key@sentry.io/project_id
	if strings.Contains(dsn, "@") {
		parts := strings.Split(dsn, "@")
		if len(parts) >= 2 {
			return maskKey(parts[0]) + "@" + parts[1]
		}
	}
	return maskKey(dsn)
}

func checkBuildTimeKeys() string {
	// This is a simple check - in a real implementation, you'd check if build-time
	// variables were properly injected during the build process
	// For now, we'll indicate if the binary likely has embedded keys
	return "checking..." // This would need actual implementation
}

func testConnectivity(name, url string) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Head(url)
	if err != nil {
		fmt.Printf("  %s: ❌ Failed (%v)\n", name, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 400 {
		fmt.Printf("  %s: ✅ Connected (HTTP %d)\n", name, resp.StatusCode)
	} else {
		fmt.Printf("  %s: ⚠️  HTTP %d\n", name, resp.StatusCode)
	}
}

func testEventSending() {
	fmt.Printf("  Sending test event...\n")

	// Send a test event
	telemetry.Track("debug_test_event", map[string]interface{}{
		"test":      true,
		"timestamp": time.Now().Unix(),
		"source":    "debug_command",
	})

	fmt.Printf("  ✅ Test event sent (check your OTLP collector or Sentry project if configured)\n")
	fmt.Printf("  💡 Enable RUNE_DEBUG=true for detailed telemetry logs\n")
}
