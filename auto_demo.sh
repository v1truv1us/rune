#!/bin/bash

# Auto demo script for Rune CLI
cd /home/f3rg/Documents/git/rune/demo_project

echo "$ pwd"
pwd
sleep 1

echo "$ ls -la"
ls -la  
sleep 2

echo "$ ../bin/rune start"
../bin/rune start
sleep 3

echo "$ ../bin/rune status"
../bin/rune status
sleep 3

echo "$ ../bin/rune pause"
../bin/rune pause
sleep 2

echo "$ ../bin/rune resume"  
../bin/rune resume
sleep 2

echo "$ ../bin/rune stop"
../bin/rune stop
sleep 3

echo "$ ../bin/rune report today"
../bin/rune report today
sleep 4

echo ""
echo "🎉 Visit https://runecli.dev to get started!"
sleep 2