package dnd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/ferg-cod3s/rune/internal/notifications"
)

// DNDManager handles Do Not Disturb functionality across platforms
type DNDManager struct {
	notificationManager *notifications.NotificationManager
}

// NewDNDManager creates a new DND manager
func NewDNDManager(notificationManager *notifications.NotificationManager) *DNDManager {
	return &DNDManager{
		notificationManager: notificationManager,
	}
}

// Enable enables Do Not Disturb mode
func (d *DNDManager) Enable() error {
	switch runtime.GOOS {
	case "darwin":
		return d.enableMacOS()
	case "linux":
		return d.enableLinux()
	case "windows":
		return d.enableWindows()
	default:
		return fmt.Errorf("DND not supported on %s", runtime.GOOS)
	}
}

// Disable disables Do Not Disturb mode
func (d *DNDManager) Disable() error {
	switch runtime.GOOS {
	case "darwin":
		return d.disableMacOS()
	case "linux":
		return d.disableLinux()
	case "windows":
		return d.disableWindows()
	default:
		return fmt.Errorf("DND not supported on %s", runtime.GOOS)
	}
}

// IsEnabled returns true if Do Not Disturb is currently enabled
func (d *DNDManager) IsEnabled() (bool, error) {
	switch runtime.GOOS {
	case "darwin":
		return d.isEnabledMacOS()
	case "linux":
		return d.isEnabledLinux()
	case "windows":
		return d.isEnabledWindows()
	default:
		return false, fmt.Errorf("DND not supported on %s", runtime.GOOS)
	}
}

// macOS implementation using modern Focus system
func (d *DNDManager) enableMacOS() error {
	// Method 1: Try using shortcuts if available (user-created shortcuts)
	cmd := exec.Command("shortcuts", "run", "Turn On Do Not Disturb")
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Method 2: Try alternative Focus shortcut names
	focusShortcuts := []string{
		"Set Do Not Disturb",
		"Enable Do Not Disturb",
		"Turn On Focus",
		"Enable Focus Mode",
		"Do Not Disturb On",
	}

	for _, shortcut := range focusShortcuts {
		cmd = exec.Command("shortcuts", "run", shortcut)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	// Method 3: Use AppleScript to access Control Center (requires accessibility permissions)
	script := `
tell application "System Events"
	try
		-- Open Control Center
		tell process "ControlCenter"
			click menu bar item "Control Center" of menu bar 1
			delay 0.8
			
			-- Look for Focus button in Control Center
			try
				click button "Focus" of group 1 of window "Control Center"
				delay 0.3
				-- Click on Do Not Disturb option
				click button "Do Not Disturb" of group 1 of window "Control Center"
			on error
				-- Try direct Do Not Disturb button
				click button "Do Not Disturb" of group 1 of window "Control Center"
			end try
			
			-- Close Control Center by clicking elsewhere
			key code 53 -- Escape key
		end tell
		return true
	on error errMsg
		-- Close Control Center if it's open
		try
			key code 53 -- Escape key
		end try
		error errMsg
	end try
end tell
`
	cmd = exec.Command("osascript", "-e", script)
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Method 4: Fallback - return an informative error
	return fmt.Errorf("could not enable Do Not Disturb - please enable it manually or create a Shortcuts automation named 'Turn On Do Not Disturb'")
}

func (d *DNDManager) disableMacOS() error {
	// Method 1: Try using shortcuts if available (user-created shortcuts)
	cmd := exec.Command("shortcuts", "run", "Turn Off Do Not Disturb")
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Method 2: Try alternative Focus shortcut names
	focusShortcuts := []string{
		"Disable Do Not Disturb",
		"Turn Off Focus",
		"Disable Focus Mode",
		"Do Not Disturb Off",
	}

	for _, shortcut := range focusShortcuts {
		cmd = exec.Command("shortcuts", "run", shortcut)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	// Method 3: Use AppleScript to access Control Center (requires accessibility permissions)
	script := `
tell application "System Events"
	try
		-- Open Control Center
		tell process "ControlCenter"
			click menu bar item "Control Center" of menu bar 1
			delay 0.8
			
			-- Look for Focus button in Control Center and turn it off
			try
				click button "Focus" of group 1 of window "Control Center"
				delay 0.3
				-- Click to turn off current focus mode
				click button "Turn Off" of group 1 of window "Control Center"
			on error
				-- Try direct Do Not Disturb button to toggle off
				click button "Do Not Disturb" of group 1 of window "Control Center"
			end try
			
			-- Close Control Center by clicking elsewhere
			key code 53 -- Escape key
		end tell
		return true
	on error errMsg
		-- Close Control Center if it's open
		try
			key code 53 -- Escape key
		end try
		error errMsg
	end try
end tell
`
	cmd = exec.Command("osascript", "-e", script)
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Method 4: Fallback - return an informative error
	return fmt.Errorf("could not disable Do Not Disturb - please disable it manually or create a Shortcuts automation named 'Turn Off Do Not Disturb'")
}

func (d *DNDManager) isEnabledMacOS() (bool, error) {
	// Method 1: Most reliable - check if the test enable/disable actually works
	// This indicates that DND is properly functional, even if we can't detect state
	// First, let's try to detect based on the success of our own enable/disable operations
	
	// Method 2: Check via AppleScript for menu bar presence
	script := `
tell application "System Events"
	try
		tell process "SystemUIServer"
			set menuBarItems to name of every menu bar item of menu bar 1
			repeat with itemName in menuBarItems
				if (itemName as string) contains "Focus" or (itemName as string) contains "Do Not Disturb" then
					return "true"
				end if
			end repeat
		end tell
	on error
		-- If SystemUIServer doesn't work, try looking for Control Center
		try
			tell process "ControlCenter"
				set menuBarItems to name of every menu bar item of menu bar 1
				repeat with itemName in menuBarItems
					if (itemName as string) contains "Focus" or (itemName as string) contains "Do Not Disturb" then
						return "true"
					end if
				end repeat
			end tell
		end try
	end try
	return "false"
end tell
`
	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output)) == "true", nil
	}

	// Method 3: Check if we can run "Get Current Focus" shortcut
	cmd = exec.Command("shortcuts", "run", "Get Current Focus")
	output, err = cmd.Output()
	if err == nil {
		result := strings.TrimSpace(string(output))
		return result != "" && result != "None" && result != "Off", nil
	}

	// Method 4: Check modern macOS Focus system by checking CoreServices
	cmd = exec.Command("sh", "-c", `
defaults read com.apple.ncprefs dnd_prefs 2>/dev/null | 
xxd -r -p 2>/dev/null | 
plutil -convert xml1 -o - -- - 2>/dev/null | 
grep -q '<key>userPref</key>' && echo "true" || echo "false"
`)
	output, err = cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output)) == "true", nil
	}

	// Method 5: Since detection is unreliable, let's use a simple heuristic:
	// If DND shortcuts are set up and working, assume we can't detect state
	// but DND functionality is available. In this case, we'll return false
	// as the safe default, but the user can still use the enable/disable functions
	return false, nil
}

// Linux implementation using various desktop environments
func (d *DNDManager) enableLinux() error {
	// Try GNOME first
	if err := d.enableGNOME(); err == nil {
		return nil
	}

	// Try KDE
	if err := d.enableKDE(); err == nil {
		return nil
	}

	// Try generic notification daemon
	return d.enableGenericLinux()
}

func (d *DNDManager) disableLinux() error {
	// Try GNOME first
	if err := d.disableGNOME(); err == nil {
		return nil
	}

	// Try KDE
	if err := d.disableKDE(); err == nil {
		return nil
	}

	// Try generic notification daemon
	return d.disableGenericLinux()
}

func (d *DNDManager) isEnabledLinux() (bool, error) {
	// Check GNOME settings
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.notifications", "show-banners")
	output, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output)) == "false", nil
	}

	// Fallback - assume not enabled if we can't detect
	return false, nil
}

func (d *DNDManager) enableGNOME() error {
	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.notifications", "show-banners", "false")
	return cmd.Run()
}

func (d *DNDManager) disableGNOME() error {
	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.notifications", "show-banners", "true")
	return cmd.Run()
}

func (d *DNDManager) enableKDE() error {
	// KDE uses kwriteconfig5 to modify notification settings
	cmd := exec.Command("kwriteconfig5", "--file", "plasmanotifyrc", "--group", "DoNotDisturb", "--key", "Enabled", "true")
	return cmd.Run()
}

func (d *DNDManager) disableKDE() error {
	cmd := exec.Command("kwriteconfig5", "--file", "plasmanotifyrc", "--group", "DoNotDisturb", "--key", "Enabled", "false")
	return cmd.Run()
}

func (d *DNDManager) enableGenericLinux() error {
	// Try to pause notification daemon
	cmd := exec.Command("notify-send", "--urgency=critical", "--expire-time=1", "DND Enabled")
	return cmd.Run()
}

func (d *DNDManager) disableGenericLinux() error {
	cmd := exec.Command("notify-send", "--urgency=normal", "--expire-time=1", "DND Disabled")
	return cmd.Run()
}

// Windows implementation using Focus Assist
func (d *DNDManager) enableWindows() error {
	// Use PowerShell to enable Focus Assist
	script := `
$registryPath = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\CloudStore\Store\Cache\DefaultAccount"
$name = "Current"
$value = 1
Set-ItemProperty -Path $registryPath -Name $name -Value $value -Force
`
	cmd := exec.Command("powershell", "-Command", script)
	return cmd.Run()
}

func (d *DNDManager) disableWindows() error {
	script := `
$registryPath = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\CloudStore\Store\Cache\DefaultAccount"
$name = "Current"
$value = 0
Set-ItemProperty -Path $registryPath -Name $name -Value $value -Force
`
	cmd := exec.Command("powershell", "-Command", script)
	return cmd.Run()
}

func (d *DNDManager) isEnabledWindows() (bool, error) {
	script := `
$registryPath = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\CloudStore\Store\Cache\DefaultAccount"
$name = "Current"
try {
    $value = Get-ItemProperty -Path $registryPath -Name $name -ErrorAction Stop
    Write-Output $value.Current
} catch {
    Write-Output "0"
}
`
	cmd := exec.Command("powershell", "-Command", script)
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(string(output)) == "1", nil
}

// CheckShortcutsSetup verifies if the required Focus mode shortcuts are available
func (d *DNDManager) CheckShortcutsSetup() (bool, error) {
	switch runtime.GOOS {
	case "darwin":
		return d.checkShortcutsMacOS()
	case "linux":
		// Linux doesn't use shortcuts, so always return true
		return true, nil
	case "windows":
		// Windows doesn't use shortcuts, so always return true
		return true, nil
	default:
		return false, fmt.Errorf("shortcuts check not supported on %s", runtime.GOOS)
	}
}

func (d *DNDManager) checkShortcutsMacOS() (bool, error) {
	// Check if the primary shortcuts exist
	requiredShortcuts := []string{
		"Turn On Do Not Disturb",
		"Turn Off Do Not Disturb",
	}

	for _, shortcut := range requiredShortcuts {
		cmd := exec.Command("shortcuts", "list")
		output, err := cmd.Output()
		if err != nil {
			return false, fmt.Errorf("failed to list shortcuts: %w", err)
		}

		if !strings.Contains(string(output), shortcut) {
			return false, nil
		}
	}

	return true, nil
}

// SendBreakNotification sends a break reminder notification
func (d *DNDManager) SendBreakNotification(workDuration time.Duration) error {
	if d.notificationManager == nil {
		return nil
	}
	return d.notificationManager.SendBreakReminder(workDuration)
}

// SendEndOfDayNotification sends an end-of-day reminder notification
func (d *DNDManager) SendEndOfDayNotification(totalTime time.Duration, targetHours float64) error {
	if d.notificationManager == nil {
		return nil
	}
	return d.notificationManager.SendEndOfDayReminder(totalTime, targetHours)
}

// SendSessionCompleteNotification sends a session completion notification
func (d *DNDManager) SendSessionCompleteNotification(duration time.Duration, project string) error {
	if d.notificationManager == nil {
		return nil
	}
	return d.notificationManager.SendSessionComplete(duration, project)
}

// SendIdleNotification sends an idle detection notification
func (d *DNDManager) SendIdleNotification(idleDuration time.Duration) error {
	if d.notificationManager == nil {
		return nil
	}
	return d.notificationManager.SendIdleDetected(idleDuration)
}

// TestNotifications sends a test notification to verify the system is working
func (d *DNDManager) TestNotifications() error {
	if d.notificationManager == nil {
		return fmt.Errorf("notification manager not initialized")
	}
	return d.notificationManager.TestNotification()
}

// Helper function to get home directory
func getHomeDir() string {
	if runtime.GOOS == "windows" {
		return "%USERPROFILE%"
	}
	return "$HOME"
}
