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
	// Detect desktop environment and use appropriate method
	de := d.detectDesktopEnvironment()

	switch de {
	case "gnome", "ubuntu", "pop":
		return d.enableGNOME()
	case "kde", "plasma":
		return d.enableKDE()
	case "xfce":
		return d.enableXFCE()
	case "mate":
		return d.enableMATE()
	case "cinnamon":
		return d.enableCinnamon()
	case "i3", "sway", "dwm", "awesome", "bspwm":
		return d.enableTilingWM()
	default:
		// Try multiple methods in order of preference
		if err := d.enableGNOME(); err == nil {
			return nil
		}
		if err := d.enableKDE(); err == nil {
			return nil
		}
		if err := d.enableXFCE(); err == nil {
			return nil
		}
		if err := d.enableMATE(); err == nil {
			return nil
		}
		if err := d.enableCinnamon(); err == nil {
			return nil
		}
		if err := d.enableTilingWM(); err == nil {
			return nil
		}
		return d.enableGenericLinux()
	}
}

func (d *DNDManager) disableLinux() error {
	// Detect desktop environment and use appropriate method
	de := d.detectDesktopEnvironment()

	switch de {
	case "gnome", "ubuntu", "pop":
		return d.disableGNOME()
	case "kde", "plasma":
		return d.disableKDE()
	case "xfce":
		return d.disableXFCE()
	case "mate":
		return d.disableMATE()
	case "cinnamon":
		return d.disableCinnamon()
	case "i3", "sway", "dwm", "awesome", "bspwm":
		return d.disableTilingWM()
	default:
		// Try multiple methods in order of preference
		if err := d.disableGNOME(); err == nil {
			return nil
		}
		if err := d.disableKDE(); err == nil {
			return nil
		}
		if err := d.disableXFCE(); err == nil {
			return nil
		}
		if err := d.disableMATE(); err == nil {
			return nil
		}
		if err := d.disableCinnamon(); err == nil {
			return nil
		}
		if err := d.disableTilingWM(); err == nil {
			return nil
		}
		return d.disableGenericLinux()
	}
}

func (d *DNDManager) isEnabledLinux() (bool, error) {
	// Detect desktop environment and check appropriate settings
	de := d.detectDesktopEnvironment()

	switch de {
	case "gnome", "ubuntu", "pop":
		return d.isEnabledGNOME()
	case "kde", "plasma":
		return d.isEnabledKDE()
	case "xfce":
		return d.isEnabledXFCE()
	case "mate":
		return d.isEnabledMATE()
	case "cinnamon":
		return d.isEnabledCinnamon()
	case "i3", "sway", "dwm", "awesome", "bspwm":
		return d.isEnabledTilingWM()
	default:
		// Try multiple detection methods
		if enabled, err := d.isEnabledGNOME(); err == nil {
			return enabled, nil
		}
		if enabled, err := d.isEnabledKDE(); err == nil {
			return enabled, nil
		}
		if enabled, err := d.isEnabledXFCE(); err == nil {
			return enabled, nil
		}
		if enabled, err := d.isEnabledMATE(); err == nil {
			return enabled, nil
		}
		if enabled, err := d.isEnabledCinnamon(); err == nil {
			return enabled, nil
		}
		if enabled, err := d.isEnabledTilingWM(); err == nil {
			return enabled, nil
		}
		return false, nil
	}
}

// detectDesktopEnvironment detects the current Linux desktop environment
func (d *DNDManager) detectDesktopEnvironment() string {
	// Check environment variables first
	envVars := []string{
		"XDG_CURRENT_DESKTOP",
		"DESKTOP_SESSION",
		"XDG_SESSION_DESKTOP",
		"GDMSESSION",
	}

	for _, envVar := range envVars {
		cmd := exec.Command("sh", "-c", "echo $"+envVar)
		output, err := cmd.Output()
		if err == nil {
			value := strings.ToLower(strings.TrimSpace(string(output)))
			if value != "" {
				if strings.Contains(value, "gnome") || strings.Contains(value, "ubuntu") || strings.Contains(value, "pop") {
					return "gnome"
				}
				if strings.Contains(value, "kde") || strings.Contains(value, "plasma") {
					return "kde"
				}
				if strings.Contains(value, "xfce") {
					return "xfce"
				}
				if strings.Contains(value, "mate") {
					return "mate"
				}
				if strings.Contains(value, "cinnamon") {
					return "cinnamon"
				}
				if strings.Contains(value, "i3") {
					return "i3"
				}
				if strings.Contains(value, "sway") {
					return "sway"
				}
			}
		}
	}

	// Check for running processes
	processes := []struct {
		process string
		de      string
	}{
		{"gnome-shell", "gnome"},
		{"plasmashell", "kde"},
		{"xfce4-panel", "xfce"},
		{"mate-panel", "mate"},
		{"cinnamon", "cinnamon"},
		{"i3", "i3"},
		{"sway", "sway"},
		{"dwm", "dwm"},
		{"awesome", "awesome"},
		{"bspwm", "bspwm"},
	}

	for _, p := range processes {
		cmd := exec.Command("pgrep", "-x", p.process)
		if err := cmd.Run(); err == nil {
			return p.de
		}
	}

	return "unknown"
}

// GNOME implementation with comprehensive settings
func (d *DNDManager) enableGNOME() error {
	commands := [][]string{
		// Disable notification banners
		{"gsettings", "set", "org.gnome.desktop.notifications", "show-banners", "false"},
		// Disable notification sounds
		{"gsettings", "set", "org.gnome.desktop.notifications", "show-in-lock-screen", "false"},
		// Set to priority mode if available (GNOME 40+)
		{"gsettings", "set", "org.gnome.desktop.notifications", "application-children", "[]"},
	}

	var lastErr error
	successCount := 0

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err == nil {
			successCount++
		} else {
			lastErr = err
		}
	}

	// If at least one command succeeded, consider it successful
	if successCount > 0 {
		return nil
	}

	return fmt.Errorf("failed to enable GNOME DND: %w", lastErr)
}

func (d *DNDManager) disableGNOME() error {
	commands := [][]string{
		// Enable notification banners
		{"gsettings", "set", "org.gnome.desktop.notifications", "show-banners", "true"},
		// Enable notification sounds
		{"gsettings", "set", "org.gnome.desktop.notifications", "show-in-lock-screen", "true"},
	}

	var lastErr error
	successCount := 0

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err == nil {
			successCount++
		} else {
			lastErr = err
		}
	}

	if successCount > 0 {
		return nil
	}

	return fmt.Errorf("failed to disable GNOME DND: %w", lastErr)
}

func (d *DNDManager) isEnabledGNOME() (bool, error) {
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.notifications", "show-banners")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(string(output)) == "false", nil
}

// KDE Plasma implementation with comprehensive settings
func (d *DNDManager) enableKDE() error {
	commands := [][]string{
		// KDE Plasma 5
		{"kwriteconfig5", "--file", "plasmanotifyrc", "--group", "DoNotDisturb", "--key", "Enabled", "true"},
		{"kwriteconfig5", "--file", "plasmanotifyrc", "--group", "Notifications", "--key", "PopupsEnabled", "false"},
		{"kwriteconfig5", "--file", "plasmanotifyrc", "--group", "Notifications", "--key", "SoundsEnabled", "false"},
		// KDE Plasma 6 (if available)
		{"kwriteconfig6", "--file", "plasmanotifyrc", "--group", "DoNotDisturb", "--key", "Enabled", "true"},
		{"kwriteconfig6", "--file", "plasmanotifyrc", "--group", "Notifications", "--key", "PopupsEnabled", "false"},
	}

	var lastErr error
	successCount := 0

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err == nil {
			successCount++
		} else {
			lastErr = err
		}
	}

	if successCount > 0 {
		// Restart plasmashell to apply changes
		_ = exec.Command("killall", "-SIGUSR1", "plasmashell").Run()
		return nil
	}

	return fmt.Errorf("failed to enable KDE DND: %w", lastErr)
}

func (d *DNDManager) disableKDE() error {
	commands := [][]string{
		// KDE Plasma 5
		{"kwriteconfig5", "--file", "plasmanotifyrc", "--group", "DoNotDisturb", "--key", "Enabled", "false"},
		{"kwriteconfig5", "--file", "plasmanotifyrc", "--group", "Notifications", "--key", "PopupsEnabled", "true"},
		{"kwriteconfig5", "--file", "plasmanotifyrc", "--group", "Notifications", "--key", "SoundsEnabled", "true"},
		// KDE Plasma 6 (if available)
		{"kwriteconfig6", "--file", "plasmanotifyrc", "--group", "DoNotDisturb", "--key", "Enabled", "false"},
		{"kwriteconfig6", "--file", "plasmanotifyrc", "--group", "Notifications", "--key", "PopupsEnabled", "true"},
	}

	var lastErr error
	successCount := 0

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err == nil {
			successCount++
		} else {
			lastErr = err
		}
	}

	if successCount > 0 {
		// Restart plasmashell to apply changes
		_ = exec.Command("killall", "-SIGUSR1", "plasmashell").Run()
		return nil
	}

	return fmt.Errorf("failed to disable KDE DND: %w", lastErr)
}

func (d *DNDManager) isEnabledKDE() (bool, error) {
	// Try KDE Plasma 5 first
	cmd := exec.Command("kreadconfig5", "--file", "plasmanotifyrc", "--group", "DoNotDisturb", "--key", "Enabled")
	output, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output)) == "true", nil
	}

	// Try KDE Plasma 6
	cmd = exec.Command("kreadconfig6", "--file", "plasmanotifyrc", "--group", "DoNotDisturb", "--key", "Enabled")
	output, err = cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output)) == "true", nil
	}

	return false, err
}

// XFCE implementation
func (d *DNDManager) enableXFCE() error {
	commands := [][]string{
		{"xfconf-query", "-c", "xfce4-notifyd", "-p", "/do-not-disturb", "-s", "true"},
		{"xfconf-query", "-c", "xfce4-notifyd", "-p", "/notification-log", "-s", "false"},
	}

	var lastErr error
	successCount := 0

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err == nil {
			successCount++
		} else {
			lastErr = err
		}
	}

	if successCount > 0 {
		return nil
	}

	return fmt.Errorf("failed to enable XFCE DND: %w", lastErr)
}

func (d *DNDManager) disableXFCE() error {
	commands := [][]string{
		{"xfconf-query", "-c", "xfce4-notifyd", "-p", "/do-not-disturb", "-s", "false"},
		{"xfconf-query", "-c", "xfce4-notifyd", "-p", "/notification-log", "-s", "true"},
	}

	var lastErr error
	successCount := 0

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err == nil {
			successCount++
		} else {
			lastErr = err
		}
	}

	if successCount > 0 {
		return nil
	}

	return fmt.Errorf("failed to disable XFCE DND: %w", lastErr)
}

func (d *DNDManager) isEnabledXFCE() (bool, error) {
	cmd := exec.Command("xfconf-query", "-c", "xfce4-notifyd", "-p", "/do-not-disturb")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(string(output)) == "true", nil
}

// MATE implementation
func (d *DNDManager) enableMATE() error {
	commands := [][]string{
		{"gsettings", "set", "org.mate.NotificationDaemon", "popup-location", "none"},
		{"dconf", "write", "/org/mate/notification-daemon/popup-location", "'none'"},
	}

	var lastErr error
	successCount := 0

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err == nil {
			successCount++
		} else {
			lastErr = err
		}
	}

	if successCount > 0 {
		return nil
	}

	return fmt.Errorf("failed to enable MATE DND: %w", lastErr)
}

func (d *DNDManager) disableMATE() error {
	commands := [][]string{
		{"gsettings", "set", "org.mate.NotificationDaemon", "popup-location", "top_right"},
		{"dconf", "write", "/org/mate/notification-daemon/popup-location", "'top_right'"},
	}

	var lastErr error
	successCount := 0

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err == nil {
			successCount++
		} else {
			lastErr = err
		}
	}

	if successCount > 0 {
		return nil
	}

	return fmt.Errorf("failed to disable MATE DND: %w", lastErr)
}

func (d *DNDManager) isEnabledMATE() (bool, error) {
	cmd := exec.Command("gsettings", "get", "org.mate.NotificationDaemon", "popup-location")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return strings.Contains(strings.TrimSpace(string(output)), "none"), nil
}

// Cinnamon implementation
func (d *DNDManager) enableCinnamon() error {
	commands := [][]string{
		{"gsettings", "set", "org.cinnamon.desktop.notifications", "display-notifications", "false"},
		{"dconf", "write", "/org/cinnamon/desktop/notifications/display-notifications", "false"},
	}

	var lastErr error
	successCount := 0

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err == nil {
			successCount++
		} else {
			lastErr = err
		}
	}

	if successCount > 0 {
		return nil
	}

	return fmt.Errorf("failed to enable Cinnamon DND: %w", lastErr)
}

func (d *DNDManager) disableCinnamon() error {
	commands := [][]string{
		{"gsettings", "set", "org.cinnamon.desktop.notifications", "display-notifications", "true"},
		{"dconf", "write", "/org/cinnamon/desktop/notifications/display-notifications", "true"},
	}

	var lastErr error
	successCount := 0

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err == nil {
			successCount++
		} else {
			lastErr = err
		}
	}

	if successCount > 0 {
		return nil
	}

	return fmt.Errorf("failed to disable Cinnamon DND: %w", lastErr)
}

func (d *DNDManager) isEnabledCinnamon() (bool, error) {
	cmd := exec.Command("gsettings", "get", "org.cinnamon.desktop.notifications", "display-notifications")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(string(output)) == "false", nil
}

// Tiling Window Manager implementation (i3, sway, dwm, etc.)
func (d *DNDManager) enableTilingWM() error {
	// For tiling WMs, we'll try to control common notification daemons
	commands := [][]string{
		// dunst (most common)
		{"dunstctl", "set-paused", "true"},
		// mako (Wayland)
		{"makoctl", "set-paused", "true"},
		// notification-daemon
		{"killall", "-STOP", "notification-daemon"},
		// Generic approach - try to pause any notification daemon
		{"pkill", "-STOP", "-f", "notify"},
	}

	var lastErr error
	successCount := 0

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err == nil {
			successCount++
		} else {
			lastErr = err
		}
	}

	if successCount > 0 {
		return nil
	}

	return fmt.Errorf("failed to enable tiling WM DND: %w", lastErr)
}

func (d *DNDManager) disableTilingWM() error {
	commands := [][]string{
		// dunst
		{"dunstctl", "set-paused", "false"},
		// mako
		{"makoctl", "set-paused", "false"},
		// notification-daemon
		{"killall", "-CONT", "notification-daemon"},
		// Generic approach
		{"pkill", "-CONT", "-f", "notify"},
	}

	var lastErr error
	successCount := 0

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err == nil {
			successCount++
		} else {
			lastErr = err
		}
	}

	if successCount > 0 {
		return nil
	}

	return fmt.Errorf("failed to disable tiling WM DND: %w", lastErr)
}

func (d *DNDManager) isEnabledTilingWM() (bool, error) {
	// Check dunst first
	cmd := exec.Command("dunstctl", "is-paused")
	output, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output)) == "true", nil
	}

	// Check mako
	cmd = exec.Command("makoctl", "mode")
	output, err = cmd.Output()
	if err == nil {
		return strings.Contains(strings.TrimSpace(string(output)), "dnd"), nil
	}

	// For other notification daemons, we can't easily check status
	// so we'll return false as the safe default
	return false, nil
}

// Generic Linux implementation (fallback)
func (d *DNDManager) enableGenericLinux() error {
	// Try to control common notification daemons
	commands := [][]string{
		// Send a test notification to see if notifications are working
		{"notify-send", "--urgency=low", "--expire-time=1000", "Rune DND", "Do Not Disturb mode enabled"},
		// Try to pause notification daemons
		{"pkill", "-STOP", "-f", "notification"},
		{"pkill", "-STOP", "-f", "notify"},
	}

	for _, cmd := range commands {
		_ = exec.Command(cmd[0], cmd[1:]...).Run()
	}

	// Always return success for generic method as it's a best-effort approach
	return nil
}
func (d *DNDManager) disableGenericLinux() error {
	commands := [][]string{
		// Resume notification daemons
		{"pkill", "-CONT", "-f", "notification"},
		{"pkill", "-CONT", "-f", "notify"},
		// Send a test notification
		{"notify-send", "--urgency=normal", "--expire-time=1000", "Rune DND", "Do Not Disturb mode disabled"},
	}

	for _, cmd := range commands {
		_ = exec.Command(cmd[0], cmd[1:]...).Run()
	}

	return nil
}

// Windows implementation using Focus Assist
func (d *DNDManager) enableWindows() error {
	// Method 1: Try using Windows 10/11 Focus Assist via registry
	if err := d.enableWindowsFocusAssist(); err == nil {
		return nil
	}

	// Method 2: Try using Windows Action Center API via PowerShell
	if err := d.enableWindowsActionCenter(); err == nil {
		return nil
	}

	// Method 3: Try using Windows Notification Settings
	if err := d.enableWindowsNotificationSettings(); err == nil {
		return nil
	}

	return fmt.Errorf("could not enable Focus Assist - please enable it manually in Windows Settings > System > Focus assist")
}

func (d *DNDManager) disableWindows() error {
	// Method 1: Try using Windows 10/11 Focus Assist via registry
	if err := d.disableWindowsFocusAssist(); err == nil {
		return nil
	}

	// Method 2: Try using Windows Action Center API via PowerShell
	if err := d.disableWindowsActionCenter(); err == nil {
		return nil
	}

	// Method 3: Try using Windows Notification Settings
	if err := d.disableWindowsNotificationSettings(); err == nil {
		return nil
	}

	return fmt.Errorf("could not disable Focus Assist - please disable it manually in Windows Settings > System > Focus assist")
}

func (d *DNDManager) isEnabledWindows() (bool, error) {
	// Method 1: Check Focus Assist registry setting
	if enabled, err := d.isEnabledWindowsFocusAssist(); err == nil {
		return enabled, nil
	}

	// Method 2: Check via PowerShell WinRT API
	if enabled, err := d.isEnabledWindowsWinRT(); err == nil {
		return enabled, nil
	}

	// Method 3: Check notification settings
	if enabled, err := d.isEnabledWindowsNotifications(); err == nil {
		return enabled, nil
	}

	// Default to false if we can't detect
	return false, nil
}

// Windows Focus Assist implementation using registry
func (d *DNDManager) enableWindowsFocusAssist() error {
	// Windows 10/11 Focus Assist registry path
	script := `
try {
    # Set Focus Assist to Priority only (1) or Alarms only (2)
    # 0 = Off, 1 = Priority only, 2 = Alarms only
    $registryPath = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\QuietHours"
    
    # Create the registry path if it doesn't exist
    if (!(Test-Path $registryPath)) {
        New-Item -Path $registryPath -Force | Out-Null
    }
    
    # Enable Focus Assist (Priority only mode)
    Set-ItemProperty -Path $registryPath -Name "Enabled" -Value 1 -Type DWord -Force
    
    # Set to Priority only mode
    Set-ItemProperty -Path $registryPath -Name "Profile" -Value 1 -Type DWord -Force
    
    Write-Output "success"
} catch {
    Write-Error $_.Exception.Message
    exit 1
}
`
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", script)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to enable Focus Assist via registry: %w", err)
	}

	if !strings.Contains(string(output), "success") {
		return fmt.Errorf("Focus Assist registry update failed")
	}

	return nil
}

func (d *DNDManager) disableWindowsFocusAssist() error {
	script := `
try {
    $registryPath = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\QuietHours"
    
    # Create the registry path if it doesn't exist
    if (!(Test-Path $registryPath)) {
        New-Item -Path $registryPath -Force | Out-Null
    }
    
    # Disable Focus Assist
    Set-ItemProperty -Path $registryPath -Name "Enabled" -Value 0 -Type DWord -Force
    
    # Set to Off mode
    Set-ItemProperty -Path $registryPath -Name "Profile" -Value 0 -Type DWord -Force
    
    Write-Output "success"
} catch {
    Write-Error $_.Exception.Message
    exit 1
}
`
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", script)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to disable Focus Assist via registry: %w", err)
	}

	if !strings.Contains(string(output), "success") {
		return fmt.Errorf("Focus Assist registry update failed")
	}

	return nil
}

func (d *DNDManager) isEnabledWindowsFocusAssist() (bool, error) {
	script := `
try {
    $registryPath = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\QuietHours"
    
    if (Test-Path $registryPath) {
        $enabled = Get-ItemProperty -Path $registryPath -Name "Enabled" -ErrorAction SilentlyContinue
        $profile = Get-ItemProperty -Path $registryPath -Name "Profile" -ErrorAction SilentlyContinue
        
        if ($enabled -and $enabled.Enabled -eq 1) {
            Write-Output "true"
        } else {
            Write-Output "false"
        }
    } else {
        Write-Output "false"
    }
} catch {
    Write-Output "false"
}
`
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", script)
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(string(output)) == "true", nil
}

// Windows Action Center implementation
func (d *DNDManager) enableWindowsActionCenter() error {
	script := `
try {
    # Try to use Windows Runtime API to control Focus Assist
    Add-Type -AssemblyName System.Runtime.WindowsRuntime
    
    # This requires Windows 10 version 1903 or later
    $asTaskGeneric = ([System.WindowsRuntimeSystemExtensions].GetMethods() | Where-Object { $_.Name -eq 'AsTask' -and $_.GetParameters().Count -eq 1 -and $_.GetParameters()[0].ParameterType.Name -eq 'IAsyncOperation' })[0]
    
    # Note: This is a simplified approach. Full implementation would require more complex WinRT calls
    Write-Output "winrt_not_available"
} catch {
    Write-Output "winrt_not_available"
}
`
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", script)
	_, _ = cmd.Output()

	// For now, this method is not fully implemented as it requires complex WinRT API calls
	return fmt.Errorf("WinRT API method not yet implemented")
}

func (d *DNDManager) disableWindowsActionCenter() error {
	// Similar to enable, this would require complex WinRT API calls
	return fmt.Errorf("WinRT API method not yet implemented")
}

func (d *DNDManager) isEnabledWindowsWinRT() (bool, error) {
	// This would require complex WinRT API calls to check Focus Assist status
	return false, fmt.Errorf("WinRT API method not yet implemented")
}

// Windows Notification Settings implementation (fallback)
func (d *DNDManager) enableWindowsNotificationSettings() error {
	script := `
try {
    # Try to disable notifications via registry as a fallback
    $registryPath = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\PushNotifications"
    
    if (!(Test-Path $registryPath)) {
        New-Item -Path $registryPath -Force | Out-Null
    }
    
    # Disable toast notifications
    Set-ItemProperty -Path $registryPath -Name "ToastEnabled" -Value 0 -Type DWord -Force
    
    Write-Output "success"
} catch {
    Write-Error $_.Exception.Message
    exit 1
}
`
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", script)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to disable notifications via registry: %w", err)
	}

	if !strings.Contains(string(output), "success") {
		return fmt.Errorf("notification settings update failed")
	}

	return nil
}

func (d *DNDManager) disableWindowsNotificationSettings() error {
	script := `
try {
    $registryPath = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\PushNotifications"
    
    if (!(Test-Path $registryPath)) {
        New-Item -Path $registryPath -Force | Out-Null
    }
    
    # Enable toast notifications
    Set-ItemProperty -Path $registryPath -Name "ToastEnabled" -Value 1 -Type DWord -Force
    
    Write-Output "success"
} catch {
    Write-Error $_.Exception.Message
    exit 1
}
`
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", script)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to enable notifications via registry: %w", err)
	}

	if !strings.Contains(string(output), "success") {
		return fmt.Errorf("notification settings update failed")
	}

	return nil
}

func (d *DNDManager) isEnabledWindowsNotifications() (bool, error) {
	script := `
try {
    $registryPath = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\PushNotifications"
    
    if (Test-Path $registryPath) {
        $toastEnabled = Get-ItemProperty -Path $registryPath -Name "ToastEnabled" -ErrorAction SilentlyContinue
        
        if ($toastEnabled -and $toastEnabled.ToastEnabled -eq 0) {
            Write-Output "true"
        } else {
            Write-Output "false"
        }
    } else {
        Write-Output "false"
    }
} catch {
    Write-Output "false"
}
`
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", script)
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(string(output)) == "true", nil
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
