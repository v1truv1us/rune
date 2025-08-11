# Notifications

Rune includes a comprehensive notification system that helps you stay aware of important work-life balance events without being intrusive.

## Features

The notification system supports several types of notifications:

- **Break Reminders**: Gentle reminders to take breaks during long work sessions
- **End-of-Day Reminders**: Notifications when you've reached your target work hours or when it's time to wrap up
- **Session Complete**: Confirmations when you finish a work session
- **Idle Detection**: Alerts when you've been idle for an extended period
- **Custom Notifications**: Support for custom notification types

## Platform Support

Notifications are supported on:

- **macOS**: Uses `terminal-notifier` (preferred) with fallback to `osascript` for native notification center integration
- **Linux**: Uses `notify-send` for desktop environment notifications
- **Windows**: Uses PowerShell and Windows Toast notifications

## Configuration

Notifications can be configured in your `config.yaml` file:

```yaml
settings:
  notifications:
    enabled: true                    # Enable/disable all notifications
    break_reminders: true           # Break reminder notifications
    end_of_day_reminders: true      # End-of-day notifications
    session_complete: true          # Session completion notifications
    idle_detection: true            # Idle detection notifications
    sound: true                     # Enable notification sounds
```

## Testing Notifications

You can test your notification setup using the built-in test commands:

```bash
# Diagnose notification tooling and setup on your OS
rune debug notifications

# Test all notification types
rune test notifications

# Test DND integration
rune test dnd
```

## Notification Types

### Break Reminders

Sent when you've been working continuously for your configured break interval:

- **Title**: "🧘 Time for a Break"
- **Message**: Shows how long you've been working and suggests taking a break
- **Priority**: Normal
- **Sound**: Yes (if enabled)

### End-of-Day Reminders

Sent when you approach or exceed your target work hours:

- **Title**: "🌅 End of Workday"
- **Message**: Shows total work time and remaining time to target (if applicable)
- **Priority**: High
- **Sound**: Yes (if enabled)

### Session Complete

Sent when you stop a work session:

- **Title**: "✅ Session Complete"
- **Message**: Shows session duration and project name
- **Priority**: Normal
- **Sound**: No

### Idle Detection

Sent when you've been idle for an extended period:

- **Title**: "💤 Idle Time Detected"
- **Message**: Shows idle duration and asks if session should be paused
- **Priority**: Normal
- **Sound**: No

## Integration with Do Not Disturb

The notification system is integrated with Rune's DND functionality:

- Notifications respect your system's Do Not Disturb settings
- Rune can automatically enable DND when starting work sessions
- Break and end-of-day notifications can still appear even when DND is enabled (depending on your system settings)

## Troubleshooting

### Notifications Not Appearing

1. **Check system permissions**: Ensure Rune has permission to send notifications
2. **Verify configuration**: Make sure notifications are enabled in your config
3. **Test the system**: Run `rune test notifications` to verify functionality
4. **Check DND settings**: Ensure your system's Do Not Disturb isn't blocking notifications

### macOS Specific Issues

- Grant notification permissions in System Settings > Notifications
- If using Terminal/Ghostty, ensure your terminal has notification permissions
- **Focus Mode Configuration**: If using Focus/Do Not Disturb mode, add your terminal application (Terminal, iTerm2, Ghostty, etc.) to the allowed apps list:
  1. Go to System Settings > Focus
  2. Select your Focus mode (Do Not Disturb, Work, etc.)
  3. Under "Allowed Notifications" click "Apps"
  4. Add your terminal application to the list
  5. Restart your terminal application for changes to take effect
- For optimal notification delivery, install `terminal-notifier`: `brew install terminal-notifier`

### Linux Specific Issues

- Ensure `notify-send` is installed: `sudo apt install libnotify-bin` (Ubuntu/Debian)
- Check your desktop environment's notification settings

### Windows Specific Issues

- Ensure Windows notifications are enabled in Settings > System > Notifications & actions
- PowerShell execution policy may need to be adjusted

## Privacy

Rune's notification system:

- Only sends notifications locally to your system
- Does not transmit any data over the network
- Respects your system's notification and privacy settings
- Can be completely disabled if not needed

## Examples

### Basic Configuration

```yaml
settings:
  notifications:
    enabled: true
    break_reminders: true
    end_of_day_reminders: true
    session_complete: false  # Disable if too noisy
    idle_detection: true
    sound: false  # Quiet notifications
```

### Developer-Focused Configuration

```yaml
settings:
  notifications:
    enabled: true
    break_reminders: true      # Important for health
    end_of_day_reminders: true # Work-life balance
    session_complete: false    # Can be distracting
    idle_detection: true       # Helpful for time tracking
    sound: false              # Maintain focus
```

### Minimal Configuration

```yaml
settings:
  notifications:
    enabled: true
    break_reminders: false
    end_of_day_reminders: true  # Only end-of-day reminders
    session_complete: false
    idle_detection: false
    sound: false
```