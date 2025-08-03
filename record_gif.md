# Manual Recording Instructions for Rune Demo GIF

## Setup
1. Terminal: Set to 100 columns x 30 rows
2. Theme: Use a dark theme (dracula/nord recommended)
3. Directory: `/home/f3rg/Documents/git/rune/demo_project`

## Recording Command
```bash
asciinema rec rune-demo.cast --cols 100 --rows 30
```

## Demo Script (Type these commands):

```bash
# === SCENE 1: Project Detection ===
pwd
ls -la

# === SCENE 2: Start Tracking ===
../bin/rune start

# === SCENE 3: Check Status ===
../bin/rune status

# === SCENE 4: Pause & Resume ===
../bin/rune pause
../bin/rune status
../bin/rune resume

# === SCENE 5: Generate Report ===
../bin/rune stop
../bin/rune report today

# === SCENE 6: Configuration ===
../bin/rune config list
```

## After Recording
```bash
# Exit recording with Ctrl+D

# Convert to GIF
/tmp/agg --theme dracula --cols 100 --rows 30 --speed 1.2 rune-demo.cast rune-demo.gif
```

## Timing Tips
- Pause 1-2 seconds between commands
- Let output display for 2 seconds before next command
- Total target time: 45-60 seconds
- Final GIF should be under 5MB for web optimization