# Linux Do Not Disturb Integration

Rune provides comprehensive Do Not Disturb (DND) functionality across all major Linux desktop environments. This document explains how the Linux DND integration works and how to troubleshoot common issues.

## Supported Desktop Environments

### GNOME (including Ubuntu, Pop!_OS)
- **Method**: Uses `gsettings` to control notification banners and sounds
- **Requirements**: `gsettings` command available
- **Settings Modified**:
  - `org.gnome.desktop.notifications.show-banners`
  - `org.gnome.desktop.notifications.show-in-lock-screen`
  - `org.gnome.desktop.notifications.application-children`

### KDE Plasma (5 & 6)
- **Method**: Uses `kwriteconfig5`/`kwriteconfig6` and `kreadconfig5`/`kreadconfig6`
- **Requirements**: KDE configuration tools available
- **Settings Modified**:
  - `plasmanotifyrc` → `DoNotDisturb.Enabled`
  - `plasmanotifyrc` → `Notifications.PopupsEnabled`
  - `plasmanotifyrc` → `Notifications.SoundsEnabled`

### XFCE
- **Method**: Uses `xfconf-query` to control XFCE notification daemon
- **Requirements**: `xfconf-query` command available
- **Settings Modified**:
  - `xfce4-notifyd` → `/do-not-disturb`
  - `xfce4-notifyd` → `/notification-log`

### MATE
- **Method**: Uses `gsettings` and `dconf` for MATE notification daemon
- **Requirements**: `gsettings` command available
- **Settings Modified**:
  - `org.mate.NotificationDaemon.popup-location`

### Cinnamon
- **Method**: Uses `gsettings` and `dconf` for Cinnamon notifications
- **Requirements**: `gsettings` command available
- **Settings Modified**:
  - `org.cinnamon.desktop.notifications.display-notifications`

### Tiling Window Managers (i3, sway, dwm, awesome, bspwm)
- **Method**: Controls notification daemons directly
- **Supported Daemons**:
  - **dunst**: Uses `dunstctl set-paused true/false`
  - **mako**: Uses `makoctl set-paused true/false`
  - **notification-daemon**: Uses process signals (STOP/CONT)
- **Requirements**: One of the supported notification daemons

### Generic/Fallback
- **Method**: Best-effort approach using process signals
- **Behavior**: Attempts to pause notification processes
- **Always succeeds**: Returns success even if methods fail

## How It Works

### Desktop Environment Detection

Rune automatically detects your desktop environment using:

1. **Environment Variables** (checked in order):
   - `XDG_CURRENT_DESKTOP`
   - `DESKTOP_SESSION`
   - `XDG_SESSION_DESKTOP`
   - `GDMSESSION`

2. **Running Processes** (if environment variables fail):
   - `gnome-shell` → GNOME
   - `plasmashell` → KDE
   - `xfce4-panel` → XFCE
   - `mate-panel` → MATE
   - `cinnamon` → Cinnamon
   - `i3`, `sway`, `dwm`, `awesome`, `bspwm` → Tiling WMs

### Fallback Strategy

If the detected desktop environment method fails, Rune tries methods in this order:
1. GNOME (`gsettings`)
2. KDE (`kwriteconfig5`/`kwriteconfig6`)
3. XFCE (`xfconf-query`)
4. MATE (`gsettings`)
5. Cinnamon (`gsettings`)
6. Tiling WM notification daemons
7. Generic process control

## Usage

### Basic Commands

```bash
# Enable Do Not Disturb
rune start

# Disable Do Not Disturb
rune stop

# Check DND status
rune status
```

### Manual DND Control

```bash
# Enable DND manually (if auto-detection fails)
rune dnd enable

# Disable DND manually
rune dnd disable

# Check if DND is enabled
rune dnd status
```

## Troubleshooting

### Common Issues

#### 1. "DND not supported" Error
**Cause**: Desktop environment not detected or unsupported
**Solutions**:
- Check if you're running a supported desktop environment
- Verify required tools are installed (see requirements below)
- Try manual commands for your specific DE

#### 2. Commands Fail Silently
**Cause**: Missing required tools or permissions
**Solutions**:
- Install missing packages (see requirements below)
- Check if you have permission to modify notification settings
- Try running commands manually to see detailed error messages

#### 3. DND Doesn't Take Effect
**Cause**: Desktop environment requires restart or manual refresh
**Solutions**:
- Log out and log back in
- Restart your notification daemon
- For KDE: Rune automatically restarts `plasmashell`

#### 4. Status Detection Fails
**Cause**: Unable to read configuration files or settings
**Solutions**:
- This is often normal - DND functionality may still work
- Check manually in your desktop environment's settings
- Some DEs don't provide reliable status checking

### Requirements by Desktop Environment

#### GNOME/Ubuntu/Pop!_OS
```bash
# Usually pre-installed, but if missing:
sudo apt install glib2.0-bin  # For gsettings
```

#### KDE Plasma
```bash
# Usually pre-installed with KDE, but if missing:
sudo apt install kde-cli-tools  # For kwriteconfig5/kreadconfig5
# For Plasma 6:
sudo apt install kde-cli-tools-6  # For kwriteconfig6/kreadconfig6
```

#### XFCE
```bash
# Usually pre-installed with XFCE, but if missing:
sudo apt install xfconf  # For xfconf-query
```

#### MATE
```bash
# Usually pre-installed with MATE, but if missing:
sudo apt install mate-desktop-common  # For gsettings schemas
```

#### Cinnamon
```bash
# Usually pre-installed with Cinnamon, but if missing:
sudo apt install cinnamon-desktop-data  # For gsettings schemas
```

#### Tiling Window Managers

**For dunst** (most common):
```bash
sudo apt install dunst
```

**For mako** (Wayland/sway):
```bash
sudo apt install mako-notifier
```

**For notification-daemon**:
```bash
sudo apt install notification-daemon
```

### Manual Testing

You can test DND functionality manually for your desktop environment:

#### GNOME
```bash
# Enable DND
gsettings set org.gnome.desktop.notifications show-banners false

# Disable DND
gsettings set org.gnome.desktop.notifications show-banners true

# Check status
gsettings get org.gnome.desktop.notifications show-banners
```

#### KDE
```bash
# Enable DND
kwriteconfig5 --file plasmanotifyrc --group DoNotDisturb --key Enabled true

# Disable DND
kwriteconfig5 --file plasmanotifyrc --group DoNotDisturb --key Enabled false

# Check status
kreadconfig5 --file plasmanotifyrc --group DoNotDisturb --key Enabled
```

#### XFCE
```bash
# Enable DND
xfconf-query -c xfce4-notifyd -p /do-not-disturb -s true

# Disable DND
xfconf-query -c xfce4-notifyd -p /do-not-disturb -s false

# Check status
xfconf-query -c xfce4-notifyd -p /do-not-disturb
```

#### Tiling WMs (dunst)
```bash
# Enable DND
dunstctl set-paused true

# Disable DND
dunstctl set-paused false

# Check status
dunstctl is-paused
```

## Advanced Configuration

### Custom Notification Daemon

If you use a custom notification daemon not supported by Rune, you can:

1. Create a wrapper script that controls your daemon
2. Place it in your PATH with the name `rune-dnd-custom`
3. Make it executable and accept `enable`/`disable`/`status` arguments

Example wrapper script:
```bash
#!/bin/bash
case "$1" in
    enable)
        # Your custom enable command
        my-notification-daemon --pause
        ;;
    disable)
        # Your custom disable command
        my-notification-daemon --resume
        ;;
    status)
        # Your custom status check
        my-notification-daemon --is-paused && echo "true" || echo "false"
        ;;
esac
```

### Environment-Specific Overrides

You can force Rune to use a specific desktop environment method by setting:
```bash
export RUNE_FORCE_DE="gnome"  # or kde, xfce, mate, cinnamon, i3, etc.
```

## Debugging

### Enable Debug Logging

```bash
export RUNE_DEBUG=true
rune start  # Will show detailed DND operation logs
```

### Check Desktop Environment Detection

```bash
# Check what Rune detects
rune debug desktop-environment

# Check environment variables
echo $XDG_CURRENT_DESKTOP
echo $DESKTOP_SESSION
```

### Test Notification System

```bash
# Send a test notification
notify-send "Test" "This is a test notification"

# Test Rune's notification system
rune test notifications
```

## Limitations

1. **Wayland Restrictions**: Some desktop environments on Wayland may have limited notification control
2. **Custom Setups**: Heavily customized desktop environments may not be detected correctly
3. **Permission Requirements**: Some operations may require specific user permissions
4. **Session Restart**: Some changes may require logging out and back in to take effect

## Contributing

If you encounter issues with a specific desktop environment or have suggestions for improvements:

1. Check the [GitHub Issues](https://github.com/ferg-cod3s/rune/issues)
2. Provide details about your desktop environment and distribution
3. Include output from `rune debug desktop-environment`
4. Test manual commands and include results

The Linux DND integration is actively maintained and we welcome contributions for additional desktop environments or notification daemons.