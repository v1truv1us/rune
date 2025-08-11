package notifications

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

// NotificationType represents different types of notifications
type NotificationType int

const (
	BreakReminder NotificationType = iota
	EndOfDayReminder
	SessionComplete
	IdleDetected
	Custom
)

// Priority represents notification priority levels
type Priority int

const (
	Low Priority = iota
	Normal
	High
	Critical
)

// Notification represents a system notification
type Notification struct {
	Title    string
	Message  string
	Type     NotificationType
	Priority Priority
	Sound    bool
	Icon     string
}

// NotificationManager handles cross-platform notifications
type NotificationManager struct {
	enabled bool
}

// NewNotificationManager creates a new notification manager
func NewNotificationManager(enabled bool) *NotificationManager {
	return &NotificationManager{
		enabled: enabled,
	}
}

// Send sends a notification to the OS
func (nm *NotificationManager) Send(notification Notification) error {
	if !nm.enabled {
		return nil // Silently skip if notifications are disabled
	}

	switch runtime.GOOS {
	case "darwin":
		return nm.sendMacOS(notification)
	case "linux":
		return nm.sendLinux(notification)
	case "windows":
		return nm.sendWindows(notification)
	default:
		return fmt.Errorf("notifications not supported on %s", runtime.GOOS)
	}
}

// SendBreakReminder sends a break reminder notification
func (nm *NotificationManager) SendBreakReminder(duration time.Duration) error {
	notification := Notification{
		Title:    "🧘 Time for a Break",
		Message:  fmt.Sprintf("You've been working for %v. Take a short break to recharge!", formatDuration(duration)),
		Type:     BreakReminder,
		Priority: Normal,
		Sound:    true,
		Icon:     "break",
	}
	return nm.Send(notification)
}

// SendEndOfDayReminder sends an end-of-day reminder notification
func (nm *NotificationManager) SendEndOfDayReminder(totalTime time.Duration, targetHours float64) error {
	var message string
	if totalTime.Hours() >= targetHours {
		message = fmt.Sprintf("Great work! You've completed %v today. Time to wrap up and enjoy your evening!", formatDuration(totalTime))
	} else {
		remaining := time.Duration(targetHours*float64(time.Hour)) - totalTime
		message = fmt.Sprintf("You've worked %v today. Consider wrapping up soon - %v remaining to reach your target.", formatDuration(totalTime), formatDuration(remaining))
	}

	notification := Notification{
		Title:    "🌅 End of Workday",
		Message:  message,
		Type:     EndOfDayReminder,
		Priority: High,
		Sound:    true,
		Icon:     "workday",
	}
	return nm.Send(notification)
}

// SendSessionComplete sends a session completion notification
func (nm *NotificationManager) SendSessionComplete(duration time.Duration, project string) error {
	notification := Notification{
		Title:    "✅ Session Complete",
		Message:  fmt.Sprintf("Finished working on %s for %v. Great job!", project, formatDuration(duration)),
		Type:     SessionComplete,
		Priority: Normal,
		Sound:    false,
		Icon:     "complete",
	}
	return nm.Send(notification)
}

// SendIdleDetected sends an idle detection notification
func (nm *NotificationManager) SendIdleDetected(idleDuration time.Duration) error {
	notification := Notification{
		Title:    "💤 Idle Time Detected",
		Message:  fmt.Sprintf("You've been idle for %v. Should I pause your session?", formatDuration(idleDuration)),
		Type:     IdleDetected,
		Priority: Normal,
		Sound:    false,
		Icon:     "idle",
	}
	return nm.Send(notification)
}

// macOS implementation using terminal-notifier (fallback to osascript)
func (nm *NotificationManager) sendMacOS(notification Notification) error {
	// Try terminal-notifier first (more reliable)
	if err := nm.tryTerminalNotifier(notification); err == nil {
		return nil
	}

	// Fallback to osascript
	script := fmt.Sprintf(`
 display notification "%s" with title "%s" sound name "%s"
`, notification.Message, notification.Title, nm.getSoundName(notification))

	cmd := exec.Command("osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to send macOS notification. Consider installing terminal-notifier via Homebrew: 'brew install terminal-notifier'. Original error: %w", err)
	}
	return nil
}

// tryTerminalNotifier attempts to use terminal-notifier for macOS notifications
func (nm *NotificationManager) tryTerminalNotifier(notification Notification) error {
	// Ensure terminal-notifier is available
	if _, err := exec.LookPath("terminal-notifier"); err != nil {
		return fmt.Errorf("terminal-notifier not found")
	}

	args := []string{
		"-title", notification.Title,
		"-message", notification.Message,
	}

	// Add priority-based arguments to bypass DND for important notifications
	switch notification.Type {
	case BreakReminder, EndOfDayReminder:
		// Critical notifications that should bypass DND
		args = append(args, "-timeout", "30") // Stay visible longer
		args = append(args, "-ignoreDnD")     // Bypass Do Not Disturb
		if notification.Sound {
			args = append(args, "-sound", "Basso") // More attention-grabbing sound
		}
	case IdleDetected:
		// Important but less urgent
		args = append(args, "-timeout", "15")
		if notification.Sound {
			args = append(args, "-sound", nm.getSoundName(notification))
		}
	default:
		// Normal notifications
		if notification.Sound {
			args = append(args, "-sound", nm.getSoundName(notification))
		}
	}

	cmd := exec.Command("terminal-notifier", args...)
	return cmd.Run()
}

// Linux implementation using notify-send
func (nm *NotificationManager) sendLinux(notification Notification) error {
	args := []string{
		"notify-send",
		"--urgency=" + nm.getUrgencyLevel(notification.Priority),
		"--expire-time=5000", // 5 seconds
	}

	if notification.Icon != "" {
		args = append(args, "--icon="+nm.getIconPath(notification.Icon))
	}

	args = append(args, notification.Title, notification.Message)

	cmd := exec.Command(args[0], args[1:]...)
	return cmd.Run()
}

// Windows implementation using PowerShell
func (nm *NotificationManager) sendWindows(notification Notification) error {
	script := fmt.Sprintf(`
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.UI.Notifications.ToastNotification, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null

$template = @"
<toast>
    <visual>
        <binding template="ToastGeneric">
            <text>%s</text>
            <text>%s</text>
        </binding>
    </visual>
</toast>
"@

$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($template)
$toast = New-Object Windows.UI.Notifications.ToastNotification $xml
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("Rune CLI").Show($toast)
`, notification.Title, notification.Message)

	cmd := exec.Command("powershell", "-Command", script)
	return cmd.Run()
}

// Helper functions

func (nm *NotificationManager) getSoundName(notification Notification) string {
	if !notification.Sound {
		return ""
	}

	switch notification.Priority {
	case Critical:
		return "Basso"
	case High:
		return "Ping"
	default:
		return "default"
	}
}

func (nm *NotificationManager) getUrgencyLevel(priority Priority) string {
	switch priority {
	case Critical:
		return "critical"
	case High:
		return "normal"
	case Low:
		return "low"
	default:
		return "normal"
	}
}

func (nm *NotificationManager) getIconPath(iconName string) string {
	// Map icon names to system icons or custom paths
	iconMap := map[string]string{
		"break":    "appointment-soon",
		"workday":  "appointment-missed",
		"complete": "emblem-default",
		"idle":     "appointment-soon",
	}

	if icon, exists := iconMap[iconName]; exists {
		return icon
	}
	return "dialog-information"
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%d seconds", int(d.Seconds()))
	}
	if d < time.Hour {
		minutes := int(d.Minutes())
		return fmt.Sprintf("%d minutes", minutes)
	}
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	if minutes == 0 {
		return fmt.Sprintf("%d hours", hours)
	}
	return fmt.Sprintf("%d hours %d minutes", hours, minutes)
}

// IsSupported returns true if notifications are supported on the current platform
func IsSupported() bool {
	switch runtime.GOOS {
	case "darwin", "linux", "windows":
		return true
	default:
		return false
	}
}

// TestNotification sends a test notification to verify the system is working
func (nm *NotificationManager) TestNotification() error {
	notification := Notification{
		Title:    "🧪 Rune Test Notification",
		Message:  "If you can see this, notifications are working correctly!",
		Type:     Custom,
		Priority: Normal,
		Sound:    true,
		Icon:     "complete",
	}
	return nm.Send(notification)
}
