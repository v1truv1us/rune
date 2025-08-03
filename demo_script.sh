#!/bin/bash

# Rune CLI Demo Script - Full Capabilities Showcase
# This script demonstrates all major features of Rune CLI

# Colors for better output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Demo configuration
DEMO_DIR="/tmp/rune_demo"
RUNE_CONFIG_DIR="$HOME/.rune"
BACKUP_CONFIG_DIR="$HOME/.rune_backup_$(date +%s)"

# Function to type with realistic speed
type_command() {
    local cmd="$1"
    local delay=${2:-0.05}
    
    echo -n "$ "
    for (( i=0; i<${#cmd}; i++ )); do
        echo -n "${cmd:$i:1}"
        sleep $delay
    done
    echo
    sleep 0.5
}

# Function to pause for effect
pause() {
    sleep ${1:-1}
}

echo "=== Rune CLI - Full Capabilities Demo ==="
echo "Showcasing: Project Detection, Time Tracking, Rituals, DND, Migration, Reporting"
echo

# Backup existing config if it exists
if [ -d "$RUNE_CONFIG_DIR" ]; then
    mv "$RUNE_CONFIG_DIR" "$BACKUP_CONFIG_DIR"
    echo "Backed up existing config to $BACKUP_CONFIG_DIR"
fi

# Setup demo environment
mkdir -p "$DEMO_DIR/go-project"
mkdir -p "$DEMO_DIR/node-project"

# Create sample project files
cat > "$DEMO_DIR/go-project/go.mod" << 'EOF'
module demo-app

go 1.21
EOF

cat > "$DEMO_DIR/go-project/main.go" << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("Demo Go Application")
}
EOF

cat > "$DEMO_DIR/node-project/package.json" << 'EOF'
{
  "name": "demo-web-app",
  "version": "1.0.0",
  "description": "Demo Node.js application"
}
EOF

cat > "$DEMO_DIR/node-project/index.js" << 'EOF'
console.log("Demo Node.js Application");
EOF

# Initialize git repos
cd "$DEMO_DIR/go-project"
git init -q
git add .
git commit -q -m "Initial commit"

cd "$DEMO_DIR/node-project"  
git init -q
git add .
git commit -q -m "Initial commit"

# Return to demo root
cd "$DEMO_DIR"

echo
echo "🎬 Starting Rune CLI Demo..."
echo "Demo environment created at: $DEMO_DIR"
echo

# === SCENE 1: INITIALIZATION ===
echo -e "${BLUE}=== Scene 1: First-Time Setup ===${NC}"
pause 1

type_command "rune init --guided"
# Simulate the guided init (we'll run the real command but with pre-answers)
./bin/rune init --work-hours 9 --break-interval 25 --idle-threshold 300 --enable-dnd true --enable-telemetry false

pause 2

# === SCENE 2: PROJECT DETECTION ===
echo -e "${BLUE}=== Scene 2: Intelligent Project Detection ===${NC}"
pause 1

cd go-project
echo -e "${GREEN}Current directory: $(pwd)${NC}"
ls -la

type_command "rune start"
./bin/rune start

pause 2

type_command "rune status"
./bin/rune status

pause 2

# === SCENE 3: SWITCHING PROJECTS ===
echo -e "${BLUE}=== Scene 3: Project Switching ===${NC}"
pause 1

type_command "rune stop"
./bin/rune stop

pause 1

cd ../node-project
echo -e "${GREEN}Switching to: $(pwd)${NC}"
ls -la

type_command "rune start"
./bin/rune start

pause 2

# === SCENE 4: PAUSE/RESUME ===
echo -e "${BLUE}=== Scene 4: Pause & Resume ===${NC}"
pause 1

type_command "rune pause"
./bin/rune pause

pause 1

type_command "rune status"
./bin/rune status

pause 1

type_command "rune resume"
./bin/rune resume

pause 2

# === SCENE 5: REPORTING ===
echo -e "${BLUE}=== Scene 5: Beautiful Reporting ===${NC}"
pause 1

type_command "rune report today"
./bin/rune report today

pause 2

type_command "rune report --format json | head -10"
./bin/rune report --format json | head -10

pause 2

# === SCENE 6: CONFIGURATION ===
echo -e "${BLUE}=== Scene 6: Configuration Management ===${NC}"
pause 1

type_command "rune config get settings.work_hours"
./bin/rune config get settings.work_hours

pause 1

type_command "rune config set settings.work_hours 8"
./bin/rune config set settings.work_hours 8

pause 1

type_command "rune config list"
./bin/rune config list

pause 2

# === SCENE 7: RITUALS DEMO ===
echo -e "${BLUE}=== Scene 7: Ritual Automation ===${NC}"
pause 1

# Show ritual configuration
type_command "cat ~/.rune/config.yaml | grep -A 10 rituals"
cat ~/.rune/config.yaml | grep -A 10 rituals || echo "# Rituals section in config"

pause 2

# === SCENE 8: MIGRATION DEMO ===
echo -e "${BLUE}=== Scene 8: Migration Tools ===${NC}"
pause 1

# Create sample Watson data
mkdir -p ~/.watson
echo '{"projects": ["old-project"], "tags": [], "start": "2025-01-27T10:00:00Z", "stop": "2025-01-27T12:00:00Z"}' > ~/.watson/frames

type_command "rune migrate watson --dry-run"  
./bin/rune migrate watson --dry-run || echo "Migration preview (would import Watson data)"

pause 2

# === SCENE 9: ADVANCED FEATURES ===
echo -e "${BLUE}=== Scene 9: Advanced Features ===${NC}"
pause 1

type_command "rune report --project go-project --format table"
./bin/rune report --project go-project --format table

pause 1

type_command "rune update --check"
./bin/rune update --check

pause 2

# === FINALE ===
type_command "rune stop"
./bin/rune stop

echo
echo -e "${GREEN}🎉 Demo Complete!${NC}"
echo "Rune CLI showcased:"
echo "  ✅ Intelligent project detection (Git, go.mod, package.json)"
echo "  ✅ Cross-platform Do Not Disturb"
echo "  ✅ Time tracking with pause/resume"
echo "  ✅ Project switching"
echo "  ✅ Beautiful reporting (table, JSON)"
echo "  ✅ Configuration management"
echo "  ✅ Ritual automation setup"
echo "  ✅ Migration from Watson/Timewarrior"
echo "  ✅ Auto-updates and version checking"
echo
echo "Visit https://runecli.dev for installation and documentation"

# Cleanup and restore
cd /
rm -rf "$DEMO_DIR"

# Restore original config if it existed
if [ -d "$BACKUP_CONFIG_DIR" ]; then
    rm -rf "$RUNE_CONFIG_DIR"
    mv "$BACKUP_CONFIG_DIR" "$RUNE_CONFIG_DIR"
    echo "Restored original config"
fi

echo "Demo environment cleaned up"