package colors

import (
	"os"
	"testing"
)

func TestColorModes(t *testing.T) {
	// Save original state
	originalMode := colorMode
	originalTheme := currentTheme
	defer func() {
		colorMode = originalMode
		currentTheme = originalTheme
	}()

	tests := []struct {
		name          string
		mode          ColorMode
		expectedTheme Theme
	}{
		{
			name:          "ColorNever",
			mode:          ColorNever,
			expectedTheme: NoColorTheme,
		},
		{
			name:          "ColorAlways",
			mode:          ColorAlways,
			expectedTheme: DefaultTheme,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			SetColorMode(test.mode)
			if colorMode != test.mode {
				t.Errorf("Expected color mode %v, got %v", test.mode, colorMode)
			}
			if currentTheme != test.expectedTheme {
				t.Errorf("Expected theme to be set correctly for mode %v", test.mode)
			}
		})
	}
}

func TestColorFunctions(t *testing.T) {
	// Test with colors enabled
	SetColorMode(ColorAlways)

	tests := []struct {
		name     string
		function func(string) string
		input    string
		contains string
	}{
		{"Primary", Primary, "test", "\033[38;2;107;70;193m"},
		{"Secondary", Secondary, "test", "\033[38;2;46;52;64m"},
		{"Accent", Accent, "test", "\033[38;2;255;183;0m"},
		{"Glow", Glow, "test", "\033[38;2;136;192;208m"},
		{"Success", Success, "test", "\033[38;2;0;255;0m"},
		{"Error", Error, "test", "\033[38;2;255;85;85m"},
		{"Warning", Warning, "test", "\033[38;2;255;183;0m"},
		{"Muted", Muted, "test", "\033[38;2;128;128;128m"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.function(test.input)
			if result == test.input {
				t.Errorf("Expected colored output, got plain text")
			}
			if len(result) <= len(test.input) {
				t.Errorf("Expected result to be longer than input due to color codes")
			}
		})
	}
}

func TestColorFunctionsDisabled(t *testing.T) {
	// Test with colors disabled
	SetColorMode(ColorNever)

	tests := []struct {
		name     string
		function func(string) string
		input    string
	}{
		{"Primary", Primary, "test"},
		{"Secondary", Secondary, "test"},
		{"Accent", Accent, "test"},
		{"Glow", Glow, "test"},
		{"Success", Success, "test"},
		{"Error", Error, "test"},
		{"Warning", Warning, "test"},
		{"Muted", Muted, "test"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.function(test.input)
			if result != test.input {
				t.Errorf("Expected plain text %q, got %q", test.input, result)
			}
		})
	}
}

func TestStatusFunctions(t *testing.T) {
	SetColorMode(ColorAlways)

	tests := []struct {
		name     string
		function func(string) string
		input    string
		symbol   string
	}{
		{"StatusRunning", StatusRunning, "test", "▶"},
		{"StatusStopped", StatusStopped, "test", "⏹"},
		{"StatusPaused", StatusPaused, "test", "⏸"},
		{"StatusSuccess", StatusSuccess, "test", "✓"},
		{"StatusError", StatusError, "test", "✗"},
		{"StatusInfo", StatusInfo, "test", "ℹ"},
		{"StatusWarning", StatusWarning, "test", "⚠"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.function(test.input)
			if result == test.input {
				t.Errorf("Expected status symbol to be added")
			}
		})
	}
}

func TestFormattingFunctions(t *testing.T) {
	SetColorMode(ColorAlways)

	input := "test"

	boldResult := Bold(input)
	if boldResult == input {
		t.Errorf("Expected bold formatting to be applied")
	}

	dimResult := Dim(input)
	if dimResult == input {
		t.Errorf("Expected dim formatting to be applied")
	}

	// Test with colors disabled
	SetColorMode(ColorNever)

	boldResultDisabled := Bold(input)
	if boldResultDisabled != input {
		t.Errorf("Expected no formatting when colors disabled")
	}

	dimResultDisabled := Dim(input)
	if dimResultDisabled != input {
		t.Errorf("Expected no formatting when colors disabled")
	}
}

func TestEnvironmentVariables(t *testing.T) {
	// Save original environment
	originalNoColor := os.Getenv("NO_COLOR")
	originalForceColor := os.Getenv("FORCE_COLOR")
	originalTerm := os.Getenv("TERM")

	defer func() {
		os.Setenv("NO_COLOR", originalNoColor)
		os.Setenv("FORCE_COLOR", originalForceColor)
		os.Setenv("TERM", originalTerm)
	}()

	// Test NO_COLOR
	os.Setenv("NO_COLOR", "1")
	os.Unsetenv("FORCE_COLOR")
	Initialize()
	if colorMode != ColorNever {
		t.Errorf("Expected ColorNever when NO_COLOR is set")
	}

	// Test FORCE_COLOR
	os.Unsetenv("NO_COLOR")
	os.Setenv("FORCE_COLOR", "1")
	Initialize()
	if colorMode != ColorAlways {
		t.Errorf("Expected ColorAlways when FORCE_COLOR is set")
	}

	// Test dumb terminal
	os.Unsetenv("NO_COLOR")
	os.Unsetenv("FORCE_COLOR")
	os.Setenv("TERM", "dumb")
	Initialize()
	if colorMode != ColorNever {
		t.Errorf("Expected ColorNever for dumb terminal")
	}
}

func TestHexColor(t *testing.T) {
	SetColorMode(ColorAlways)

	tests := []struct {
		name     string
		hex      string
		expected bool // whether it should return a non-empty string
	}{
		{"Valid hex with #", "#FF0000", true},
		{"Valid hex without #", "FF0000", true},
		{"Invalid hex - too short", "#FF", false},
		{"Invalid hex - too long", "#FF00000", false},
		{"Invalid hex - non-hex chars", "#GGGGGG", false},
		{"Empty string", "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Hex(test.hex)
			isEmpty := result == ""
			if test.expected && isEmpty {
				t.Errorf("Expected non-empty result for valid hex %s", test.hex)
			}
			if !test.expected && !isEmpty {
				t.Errorf("Expected empty result for invalid hex %s", test.hex)
			}
		})
	}

	// Test with colors disabled
	SetColorMode(ColorNever)
	result := Hex("#FF0000")
	if result != "" {
		t.Errorf("Expected empty result when colors disabled")
	}
}

func TestRGBColor(t *testing.T) {
	SetColorMode(ColorAlways)

	result := RGB(255, 0, 0)
	expected := "\033[38;2;255;0;0m"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test with colors disabled
	SetColorMode(ColorNever)
	result = RGB(255, 0, 0)
	if result != "" {
		t.Errorf("Expected empty result when colors disabled")
	}
}

func TestUtilityFunctions(t *testing.T) {
	SetColorMode(ColorAlways)

	tests := []struct {
		name     string
		function func(string) string
		input    string
	}{
		{"Header", Header, "Test Header"},
		{"Subheader", Subheader, "Test Subheader"},
		{"Logo", Logo, "Test Logo"},
		{"Duration", Duration, "1h 30m"},
		{"Project", Project, "my-project"},
		{"Time", Time, "14:30"},
		{"RelativeTime", RelativeTime, "2h ago"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.function(test.input)
			if result == test.input {
				t.Errorf("Expected styling to be applied to %s", test.name)
			}
		})
	}
}

func TestGetters(t *testing.T) {
	SetColorMode(ColorAlways)

	if GetColorMode() != ColorAlways {
		t.Errorf("Expected GetColorMode to return ColorAlways")
	}

	theme := GetTheme()
	if theme != DefaultTheme {
		t.Errorf("Expected GetTheme to return DefaultTheme")
	}
}

func TestCustomTheme(t *testing.T) {
	// Save original theme
	originalTheme := currentTheme
	defer func() {
		currentTheme = originalTheme
	}()

	customTheme := Theme{
		Primary:   "\033[31m", // Red
		Secondary: "\033[32m", // Green
		Accent:    "\033[33m", // Yellow
		Glow:      "\033[34m", // Blue
		Success:   "\033[35m", // Magenta
		Error:     "\033[36m", // Cyan
		Warning:   "\033[37m", // White
		Muted:     "\033[90m", // Bright Black
	}

	SetColorMode(ColorAlways)
	SetTheme(customTheme)

	if GetTheme() != customTheme {
		t.Errorf("Expected custom theme to be set")
	}

	// Test that colors are not set when ColorNever
	SetColorMode(ColorNever)
	SetTheme(customTheme)

	if GetTheme() == customTheme {
		t.Errorf("Expected theme not to be set when ColorNever")
	}
}
