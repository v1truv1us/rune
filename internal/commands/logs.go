package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ferg-cod3s/rune/internal/colors"
	"github.com/ferg-cod3s/rune/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Time      time.Time              `json:"time"`
	Level     string                 `json:"level"`
	Message   string                 `json:"msg"`
	Component string                 `json:"component,omitempty"`
	Command   string                 `json:"command,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

var logsCmd = &cobra.Command{
	Use:     "logs",
	Aliases: []string{"log"},
	Short:   "Display recent logs",
	Long: `Display recent application logs, including errors that would be sent to Sentry.
This is useful for troubleshooting issues without connecting to external services.`,
	RunE: runLogs,
}

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().IntP("lines", "n", 50, "number of log lines to display")
	logsCmd.Flags().StringP("level", "l", "", "filter by log level (debug, info, warn, error)")
	logsCmd.Flags().Bool("json", false, "output logs in JSON format")
	logsCmd.Flags().Bool("follow", false, "follow log output (like tail -f)")
	logsCmd.Flags().StringP("component", "c", "", "filter by component")
	logsCmd.Flags().StringP("file", "f", "", "specific log file to read (default: most recent session)")
	logsCmd.Flags().Bool("list", false, "list available log files")
}

func runLogs(cmd *cobra.Command, args []string) error {
	lines, _ := cmd.Flags().GetInt("lines")
	levelFilter, _ := cmd.Flags().GetString("level")
	jsonOutput, _ := cmd.Flags().GetBool("json")
	follow, _ := cmd.Flags().GetBool("follow")
	componentFilter, _ := cmd.Flags().GetString("component")
	specificFile, _ := cmd.Flags().GetString("file")
	listFiles, _ := cmd.Flags().GetBool("list")

	// Handle --list flag
	if listFiles {
		return listLogFiles()
	}

	// Check if --log flag was used at root level
	if viper.GetBool("log") && len(args) == 0 && !cmd.Flags().Changed("lines") {
		lines = 20 // Show fewer lines for quick --log usage
	}

	var logPath string
	if specificFile != "" {
		logPath = specificFile
	} else {
		logPath = getLogFilePath()
	}

	if follow {
		return followLogs(logPath, levelFilter, componentFilter, jsonOutput)
	}

	return displayLogs(logPath, lines, levelFilter, componentFilter, jsonOutput)
}

func listLogFiles() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	logsDir := filepath.Join(home, ".rune", "logs")

	// Look for all log files
	files, err := filepath.Glob(filepath.Join(logsDir, "*.log"))
	if err != nil {
		return fmt.Errorf("failed to list log files: %w", err)
	}

	if len(files) == 0 {
		fmt.Printf("No log files found in %s\n", logsDir)
		return nil
	}

	fmt.Printf("📂 Available log files in %s:\n\n", logsDir)

	// Sort files by modification time (newest first)
	sort.Slice(files, func(i, j int) bool {
		info1, err1 := os.Stat(files[i])
		info2, err2 := os.Stat(files[j])
		if err1 != nil || err2 != nil {
			return false
		}
		return info1.ModTime().After(info2.ModTime())
	})

	for i, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		basename := filepath.Base(file)
		size := info.Size()
		modTime := info.ModTime().Format("2006-01-02 15:04:05")

		marker := "  "
		if i == 0 {
			marker = colors.Success("▶ ") // Mark the newest/current file
		}

		fmt.Printf("%s%s\n", marker, colors.Accent(basename))
		fmt.Printf("   %s │ %s │ %d bytes\n", colors.Muted(modTime), colors.Muted("Modified"), size)
		if i == 0 {
			fmt.Printf("   %s\n", colors.Muted("(most recent - used by default)"))
		}
		fmt.Println()
	}

	fmt.Printf("Use %s to read a specific file\n", colors.Accent("rune logs --file <filename>"))

	return nil
}

func getLogFilePath() string {
	// Try to get from config first
	if cfg, err := config.Load(); err == nil {
		if cfg.Logging.Output != "" && cfg.Logging.Output != "stdout" && cfg.Logging.Output != "stderr" {
			return cfg.Logging.Output
		}
	}

	// Find the most recent session log file
	home, err := os.UserHomeDir()
	if err != nil {
		return "./rune-logs.json" // fallback
	}

	logsDir := filepath.Join(home, ".rune", "logs")

	// Look for the most recent session log file
	files, err := filepath.Glob(filepath.Join(logsDir, "rune-session-*.log"))
	if err == nil && len(files) > 0 {
		// Sort files by name (which includes timestamp) to get the latest
		sort.Strings(files)
		return files[len(files)-1] // Return the most recent file
	}

	// Fallback to the old filename for backward compatibility
	return filepath.Join(logsDir, "rune.log")
}

func displayLogs(logPath string, lines int, levelFilter, componentFilter string, jsonOutput bool) error {
	file, err := os.Open(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("No log file found at %s\n", logPath)
			fmt.Println("Logs will be created when the application runs with JSON logging enabled.")
			return nil
		}
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			// Skip malformed lines
			continue
		}

		// Apply filters
		if levelFilter != "" && !strings.EqualFold(entry.Level, levelFilter) {
			continue
		}

		if componentFilter != "" && !strings.Contains(entry.Component, componentFilter) {
			continue
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log file: %w", err)
	}

	// Sort by time (most recent last)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Time.Before(entries[j].Time)
	})

	// Take only the last N entries
	start := 0
	if len(entries) > lines {
		start = len(entries) - lines
	}
	entries = entries[start:]

	// Display entries
	for _, entry := range entries {
		if jsonOutput {
			jsonBytes, _ := json.Marshal(entry)
			fmt.Println(string(jsonBytes))
		} else {
			displayLogEntry(entry)
		}
	}

	return nil
}

func followLogs(logPath string, levelFilter, componentFilter string, jsonOutput bool) error {
	// Simple tail -f implementation
	// In a real implementation, you might want to use a proper file watching library
	fmt.Printf("Following logs at %s (Ctrl+C to exit)\n", logPath)

	file, err := os.Open(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Waiting for log file to be created at %s...\n", logPath)
			// In a real implementation, you'd watch for file creation
			return nil
		}
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// Seek to end of file
	_, err = file.Seek(0, 2)
	if err != nil {
		return fmt.Errorf("failed to seek to end of file: %w", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		// Apply filters
		if levelFilter != "" && !strings.EqualFold(entry.Level, levelFilter) {
			continue
		}

		if componentFilter != "" && !strings.Contains(entry.Component, componentFilter) {
			continue
		}

		if jsonOutput {
			jsonBytes, _ := json.Marshal(entry)
			fmt.Println(string(jsonBytes))
		} else {
			displayLogEntry(entry)
		}
	}

	return scanner.Err()
}

func displayLogEntry(entry LogEntry) {
	// Format timestamp
	timeStr := entry.Time.Format("2006-01-02 15:04:05")

	// Color-code log levels
	var levelStr string
	switch strings.ToLower(entry.Level) {
	case "debug":
		levelStr = colors.Muted("DEBUG")
	case "info":
		levelStr = colors.Success("INFO ")
	case "warn":
		levelStr = colors.Warning("WARN ")
	case "error":
		levelStr = colors.Error("ERROR")
	default:
		levelStr = entry.Level
	}

	// Build the log line
	var parts []string
	parts = append(parts, colors.Muted(timeStr))
	parts = append(parts, levelStr)

	if entry.Component != "" {
		parts = append(parts, colors.Accent("["+entry.Component+"]"))
	}

	if entry.Command != "" {
		parts = append(parts, colors.Accent("cmd:"+entry.Command))
	}

	parts = append(parts, entry.Message)

	fmt.Println(strings.Join(parts, " "))

	// Display error details if present
	if entry.Error != "" {
		fmt.Printf("  %s %s\n", colors.Error("└─ Error:"), entry.Error)
	}

	// Display additional context if present and not in JSON mode
	if len(entry.Context) > 0 {
		for key, value := range entry.Context {
			if key != "error" { // Don't duplicate error info
				fmt.Printf("  %s %s: %v\n", colors.Muted("├─"), key, value)
			}
		}
	}
}
