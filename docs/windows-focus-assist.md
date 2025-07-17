# Windows Focus Assist Integration

Rune provides comprehensive Windows Focus Assist integration to help you maintain focus during work sessions by automatically managing notifications and distractions.

## Overview

Windows Focus Assist is Microsoft's built-in notification management system that helps reduce distractions by filtering notifications. Rune integrates with Focus Assist to automatically enable it when you start a work session and disable it when you finish.

## Features

### Automatic Focus Assist Control
- **Enable on Start**: Automatically enables Focus Assist when starting a work session
- **Disable on Stop**: Automatically disables Focus Assist when ending a work session
- **Status Detection**: Checks current Focus Assist status
- **Multiple Methods**: Uses multiple approaches for maximum compatibility

### Focus Assist Modes
Rune enables Focus Assist in "Priority only" mode, which:
- Blocks most notifications except priority ones
- Allows alarms and reminders to come through
- Maintains important system notifications
- Provides a good balance between focus and functionality

## Implementation Methods

Rune uses a multi-layered approach to ensure compatibility across different Windows versions:

### 1. Registry-Based Control (Primary)
- **Path**: `HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\QuietHours`
- **Method**: Direct registry manipulation for reliable control
- **Compatibility**: Windows 10 version 1903+ and Windows 11
- **Advantages**: Fast, reliable, works without user interaction

### 2. Notification Settings (Fallback)
- **Path**: `HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\PushNotifications`
- **Method**: Controls toast notifications directly
- **Compatibility**: Windows 10 and Windows 11
- **Advantages**: Works when Focus Assist isn't available

### 3. WinRT API (Future)
- **Method**: Windows Runtime API calls
- **Status**: Planned for future implementation
- **Advantages**: Native Windows integration, more granular control

## Configuration

### Basic Setup
Focus Assist integration is enabled by default in Rune configurations. No additional setup is required.

```yaml
# In your ~/.rune/config.yaml
integrations:
  # Focus Assist is automatically managed
  # No additional configuration needed
```

### Advanced Configuration
You can customize Focus Assist behavior in your ritual configurations:

```yaml
rituals:
  start:
    global:
      - name: "Enable Focus Mode"
        command: "rune dnd enable"
        optional: true
  stop:
    global:
      - name: "Disable Focus Mode"
        command: "rune dnd disable"
        optional: true
```

## Usage

### Automatic Usage
Focus Assist is automatically managed when you use Rune's session commands:

```bash
# Automatically enables Focus Assist
rune start my-project

# Automatically disables Focus Assist
rune stop
```

### Manual Control
You can also manually control Focus Assist:

```bash
# Enable Focus Assist
rune dnd enable

# Disable Focus Assist
rune dnd disable

# Check current status
rune dnd status
```

### Testing
Test your Focus Assist integration:

```bash
# Test DND functionality
rune test dnd
```

## Troubleshooting

### Common Issues

#### "PowerShell not found" Error
**Problem**: PowerShell is not available in the system PATH.
**Solution**: 
1. Ensure PowerShell is installed (comes with Windows 10/11)
2. Add PowerShell to your system PATH
3. Try running `powershell` from Command Prompt to verify

#### "Access Denied" Error
**Problem**: Insufficient permissions to modify registry settings.
**Solution**:
1. Run Command Prompt as Administrator
2. Execute Rune commands from the elevated prompt
3. Alternatively, manually enable Focus Assist in Windows Settings

#### Focus Assist Not Working
**Problem**: Focus Assist appears to enable but notifications still come through.
**Solution**:
1. Check Windows Settings > System > Focus assist
2. Verify that "Priority only" mode is active
3. Review your priority notification settings
4. Ensure apps aren't bypassing Focus Assist

### Manual Fallback
If automatic Focus Assist control isn't working, you can manually manage it:

1. **Windows Settings Method**:
   - Open Settings > System > Focus assist
   - Select "Priority only" or "Alarms only"
   - Configure priority senders as needed

2. **Action Center Method**:
   - Click the Action Center icon in the system tray
   - Click the "Focus assist" tile to toggle

3. **Keyboard Shortcut**:
   - Windows key + U to open Ease of Access settings
   - Navigate to Focus assist settings

## Compatibility

### Windows Versions
- **Windows 11**: Full support with all methods
- **Windows 10 (1903+)**: Full support with registry method
- **Windows 10 (older)**: Limited support via notification settings
- **Windows 8.1 and earlier**: Not supported

### PowerShell Versions
- **PowerShell 5.1**: Fully supported (default on Windows 10/11)
- **PowerShell Core 6+**: Supported
- **Windows PowerShell ISE**: Supported

### Permissions
- **Standard User**: Basic functionality available
- **Administrator**: Full functionality with all methods
- **Restricted User**: Limited functionality

## Security Considerations

### Registry Access
- Rune only modifies user-specific registry keys (`HKCU`)
- No system-wide changes are made
- Changes are reversible and don't affect system stability

### PowerShell Execution
- Scripts use `-ExecutionPolicy Bypass` for reliability
- No persistent execution policy changes
- Scripts are embedded and not stored on disk

### Privacy
- No data is collected about your Focus Assist usage
- All operations are performed locally
- No network requests are made for DND functionality

## Advanced Usage

### Custom Focus Assist Scripts
You can create custom PowerShell scripts for advanced Focus Assist control:

```powershell
# Custom script: focus-assist-priority.ps1
$registryPath = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\QuietHours"
Set-ItemProperty -Path $registryPath -Name "Profile" -Value 1 -Type DWord
```

### Integration with Other Tools
Combine Rune's Focus Assist with other productivity tools:

```yaml
rituals:
  start:
    global:
      - name: "Enable Focus Mode"
        command: "rune dnd enable"
      - name: "Close Distracting Apps"
        command: "taskkill /f /im slack.exe /im discord.exe"
        optional: true
      - name: "Start Focus Music"
        command: "start spotify:playlist:focus"
        optional: true
```

### Monitoring and Logging
Track Focus Assist usage for productivity insights:

```bash
# Check DND status and log it
rune dnd status >> focus-log.txt

# Create a daily focus report
echo "$(date): Focus session started" >> daily-focus.log
```

## API Reference

### DND Manager Methods

#### `Enable() error`
Enables Focus Assist using the best available method.

#### `Disable() error`
Disables Focus Assist and restores normal notifications.

#### `IsEnabled() (bool, error)`
Returns the current Focus Assist status.

#### `enableWindowsFocusAssist() error`
Directly enables Focus Assist via registry (Windows-specific).

#### `disableWindowsFocusAssist() error`
Directly disables Focus Assist via registry (Windows-specific).

#### `isEnabledWindowsFocusAssist() (bool, error)`
Checks Focus Assist status via registry (Windows-specific).

### Error Handling
All methods return descriptive errors for troubleshooting:
- Registry access errors
- PowerShell execution errors
- Permission-related errors
- Compatibility issues

## Contributing

### Testing on Windows
When contributing Windows Focus Assist improvements:

1. Test on multiple Windows versions (10, 11)
2. Test with different user permission levels
3. Verify registry changes don't persist after disable
4. Test PowerShell compatibility

### Adding New Methods
To add new Focus Assist control methods:

1. Implement the method in `internal/dnd/dnd.go`
2. Add it to the fallback chain in `enableWindows()`
3. Add corresponding tests in `internal/dnd/dnd_windows_test.go`
4. Update this documentation

## Related Documentation

- [Do Not Disturb Overview](notifications.md)
- [Configuration Guide](getting-started/configuration.md)
- [Ritual Automation](rituals.md)
- [Cross-Platform Support](platform-support.md)