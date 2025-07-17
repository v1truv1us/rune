package colors

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ANSI color codes
const (
	Reset    = "\033[0m"
	BoldCode = "\033[1m"
	DimCode  = "\033[2m"
)

// Theme represents a color theme for the CLI
type Theme struct {
	Primary   string // Deep Purple - main brand color
	Secondary string // Norse Blue - secondary elements
	Accent    string // Mystic Gold - highlights and warnings
	Glow      string // Rune Cyan - success and info
	Success   string // Terminal Green - success messages
	Error     string // Red - error messages
	Warning   string // Yellow - warning messages
	Muted     string // Gray - secondary text
}

// Predefined themes
var (
	// DefaultTheme - full color theme
	DefaultTheme = Theme{
		Primary:   "\033[38;2;107;70;193m",  // #6B46C1
		Secondary: "\033[38;2;46;52;64m",    // #2E3440
		Accent:    "\033[38;2;255;183;0m",   // #FFB700
		Glow:      "\033[38;2;136;192;208m", // #88C0D0
		Success:   "\033[38;2;0;255;0m",     // #00FF00
		Error:     "\033[38;2;255;85;85m",   // #FF5555
		Warning:   "\033[38;2;255;183;0m",   // #FFB700 (same as accent)
		Muted:     "\033[38;2;128;128;128m", // #808080
	}

	// MonochromeTheme - no colors, just formatting
	MonochromeTheme = Theme{
		Primary:   BoldCode,
		Secondary: "",
		Accent:    BoldCode,
		Glow:      "",
		Success:   BoldCode,
		Error:     BoldCode,
		Warning:   BoldCode,
		Muted:     DimCode,
	}

	// NoColorTheme - completely plain text
	NoColorTheme = Theme{
		Primary:   "",
		Secondary: "",
		Accent:    "",
		Glow:      "",
		Success:   "",
		Error:     "",
		Warning:   "",
		Muted:     "",
	}
)

// ColorMode represents the color output mode
type ColorMode int

const (
	ColorAuto ColorMode = iota
	ColorAlways
	ColorNever
)

// Global color settings
var (
	currentTheme = DefaultTheme
	colorMode    = ColorAuto
	forceColor   = false
)

// Initialize sets up the color system based on environment
func Initialize() {
	// Check environment variables
	if os.Getenv("NO_COLOR") != "" {
		SetColorMode(ColorNever)
		return
	}

	if os.Getenv("FORCE_COLOR") != "" {
		forceColor = true
		SetColorMode(ColorAlways)
		return
	}

	// Check TERM environment
	term := os.Getenv("TERM")
	if term == "dumb" || term == "" {
		SetColorMode(ColorNever)
		return
	}

	// Auto-detect based on terminal capabilities
	if !isTerminalCapable() {
		SetColorMode(ColorNever)
		return
	}

	SetColorMode(ColorAuto)
}

// SetColorMode sets the color output mode
func SetColorMode(mode ColorMode) {
	colorMode = mode
	switch mode {
	case ColorNever:
		currentTheme = NoColorTheme
	case ColorAlways:
		currentTheme = DefaultTheme
	case ColorAuto:
		if isTerminalCapable() {
			currentTheme = DefaultTheme
		} else {
			currentTheme = NoColorTheme
		}
	}
}

// SetTheme sets a custom theme
func SetTheme(theme Theme) {
	if colorMode != ColorNever {
		currentTheme = theme
	}
}

// isTerminalCapable checks if the terminal supports colors
func isTerminalCapable() bool {
	if forceColor {
		return true
	}

	// Check if stdout is a terminal
	if !isatty(os.Stdout) {
		return false
	}

	// Check COLORTERM
	if os.Getenv("COLORTERM") != "" {
		return true
	}

	// Check TERM for color support
	term := os.Getenv("TERM")
	colorTerms := []string{
		"xterm", "xterm-color", "xterm-256color",
		"screen", "screen-256color",
		"tmux", "tmux-256color",
		"rxvt", "ansi", "color",
	}

	for _, colorTerm := range colorTerms {
		if strings.Contains(term, colorTerm) {
			return true
		}
	}

	return false
}

// isatty checks if the file descriptor is a terminal
func isatty(f *os.File) bool {
	// Simple check - in a real implementation you might want to use
	// a more robust method or a library like github.com/mattn/go-isatty
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// Color functions - these apply colors and automatically reset
func Primary(text string) string {
	if currentTheme.Primary == "" {
		return text
	}
	return currentTheme.Primary + text + Reset
}

func Secondary(text string) string {
	if currentTheme.Secondary == "" {
		return text
	}
	return currentTheme.Secondary + text + Reset
}

func Accent(text string) string {
	if currentTheme.Accent == "" {
		return text
	}
	return currentTheme.Accent + text + Reset
}

func Glow(text string) string {
	if currentTheme.Glow == "" {
		return text
	}
	return currentTheme.Glow + text + Reset
}

func Success(text string) string {
	if currentTheme.Success == "" {
		return text
	}
	return currentTheme.Success + text + Reset
}

func Error(text string) string {
	if currentTheme.Error == "" {
		return text
	}
	return currentTheme.Error + text + Reset
}

func Warning(text string) string {
	if currentTheme.Warning == "" {
		return text
	}
	return currentTheme.Warning + text + Reset
}

func Muted(text string) string {
	if currentTheme.Muted == "" {
		return text
	}
	return currentTheme.Muted + text + Reset
}

// Formatting functions
func Bold(text string) string {
	if colorMode == ColorNever {
		return text
	}
	return BoldCode + text + Reset
}

func Dim(text string) string {
	if colorMode == ColorNever {
		return text
	}
	return DimCode + text + Reset
}

// Utility functions for common patterns
func StatusRunning(text string) string {
	return Glow("▶ ") + text
}

func StatusStopped(text string) string {
	return Muted("⏹ ") + text
}

func StatusPaused(text string) string {
	return Warning("⏸ ") + text
}

func StatusSuccess(text string) string {
	return Success("✓ ") + text
}

func StatusError(text string) string {
	return Error("✗ ") + text
}

func StatusInfo(text string) string {
	return Glow("ℹ ") + text
}

func StatusWarning(text string) string {
	return Warning("⚠ ") + text
}

// Header creates a styled header with the Rune brand colors
func Header(text string) string {
	return Primary(Bold(text))
}

// Subheader creates a styled subheader
func Subheader(text string) string {
	return Glow(text)
}

// Logo applies the primary color to ASCII art
func Logo(text string) string {
	return Primary(text)
}

// Duration formats duration with accent color
func Duration(text string) string {
	return Accent(text)
}

// Project formats project names with glow color
func Project(text string) string {
	return Glow(text)
}

// Time formats time strings with muted color
func Time(text string) string {
	return Muted(text)
}

// RelativeTime formats relative time with muted color
func RelativeTime(text string) string {
	return Muted("(" + text + ")")
}

// Sprintf applies colors to formatted strings
func Sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// Printf prints colored formatted strings
func Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Println prints colored text with newline
func Println(text string) {
	fmt.Println(text)
}

// GetColorMode returns the current color mode
func GetColorMode() ColorMode {
	return colorMode
}

// GetTheme returns the current theme
func GetTheme() Theme {
	return currentTheme
}

// RGB creates a 24-bit color code from RGB values
func RGB(r, g, b int) string {
	if colorMode == ColorNever {
		return ""
	}
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

// Hex creates a color code from hex string (e.g., "#FF0000")
func Hex(hex string) string {
	if colorMode == ColorNever {
		return ""
	}

	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return ""
	}

	r, err1 := strconv.ParseInt(hex[0:2], 16, 0)
	g, err2 := strconv.ParseInt(hex[2:4], 16, 0)
	b, err3 := strconv.ParseInt(hex[4:6], 16, 0)

	if err1 != nil || err2 != nil || err3 != nil {
		return ""
	}

	return RGB(int(r), int(g), int(b))
}
