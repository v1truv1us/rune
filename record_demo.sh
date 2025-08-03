#!/bin/bash

# Rune CLI Demo Recording Script for GIF
# Optimized for terminal recording and conversion to GIF

set -e

DEMO_DIR="/tmp/rune_gif_demo"
RUNE_BIN="./bin/rune"

# Clean up any existing demo
rm -rf "$DEMO_DIR"
mkdir -p "$DEMO_DIR/my-go-project"

# Setup demo project
cd "$DEMO_DIR/my-go-project"
cat > go.mod << 'EOF'
module my-go-project

go 1.21
EOF

cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("Hello, Rune!")
}
EOF

# Initialize git
git init -q
git add .
git commit -q -m "Initial commit"

echo "🎬 Demo environment ready at: $DEMO_DIR/my-go-project"
echo "Starting asciinema recording in 3 seconds..."
sleep 3

# Record the demo
asciinema rec rune-demo.cast --overwrite --command "bash -c '
cd /tmp/rune_gif_demo/my-go-project

echo "# Rune CLI - Full Workflow Demo"
echo "# Intelligent time tracking, focus mode, and automation"
echo

sleep 1

echo "$ pwd"
pwd
sleep 0.5

echo "$ ls -la"
ls -la
sleep 1

echo "$ /home/f3rg/Documents/git/rune/bin/rune start"
/home/f3rg/Documents/git/rune/bin/rune start
sleep 2

echo "$ /home/f3rg/Documents/git/rune/bin/rune status"
/home/f3rg/Documents/git/rune/bin/rune status  
sleep 2

echo "# Simulating work... (25 minutes later)"
sleep 1

echo "$ /home/f3rg/Documents/git/rune/bin/rune status"
/home/f3rg/Documents/git/rune/bin/rune status
sleep 2

echo "$ /home/f3rg/Documents/git/rune/bin/rune pause"
/home/f3rg/Documents/git/rune/bin/rune pause
sleep 1

echo "# Taking a break..."
sleep 1

echo "$ /home/f3rg/Documents/git/rune/bin/rune resume"
/home/f3rg/Documents/git/rune/bin/rune resume
sleep 1

echo "$ /home/f3rg/Documents/git/rune/bin/rune stop"
/home/f3rg/Documents/git/rune/bin/rune stop
sleep 2

echo "$ /home/f3rg/Documents/git/rune/bin/rune report today"
/home/f3rg/Documents/git/rune/bin/rune report today
sleep 3

echo
echo "🎉 Visit https://runecli.dev to get started!"
sleep 2
'"

echo
echo "Recording complete! Converting to GIF..."

# Convert to GIF
/tmp/agg --theme dracula --cols 100 --rows 30 --speed 1.5 rune-demo.cast rune-demo.gif

echo "✅ GIF created: rune-demo.gif"
echo "📏 Optimized for web: 100x30 chars, 1.5x speed"
echo

# Cleanup
rm -rf "$DEMO_DIR"
echo "Demo environment cleaned up"