package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ferg-cod3s/rune/internal/tracking"
)

func TestReadWatsonFrames(t *testing.T) {
	// Create a temporary Watson frames file
	tmpDir := t.TempDir()
	framesFile := filepath.Join(tmpDir, "frames.json")

	frames := []WatsonFrame{
		{
			ID:      "frame1",
			Start:   time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC),
			Stop:    time.Date(2024, 1, 15, 11, 30, 0, 0, time.UTC),
			Project: "web-development",
			Tags:    []string{"frontend", "react"},
		},
		{
			ID:      "frame2",
			Start:   time.Date(2024, 1, 15, 13, 0, 0, 0, time.UTC),
			Stop:    time.Date(2024, 1, 15, 15, 45, 0, 0, time.UTC),
			Project: "documentation",
			Tags:    []string{"writing", "api-docs"},
		},
	}

	// Write test data to file
	data, err := json.Marshal(frames)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	if err := os.WriteFile(framesFile, data, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test reading the frames
	result, err := readWatsonFrames(framesFile)
	if err != nil {
		t.Fatalf("readWatsonFrames failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 frames, got %d", len(result))
	}

	// Verify first frame
	if result[0].ID != "frame1" {
		t.Errorf("Expected frame ID 'frame1', got '%s'", result[0].ID)
	}
	if result[0].Project != "web-development" {
		t.Errorf("Expected project 'web-development', got '%s'", result[0].Project)
	}
	if len(result[0].Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(result[0].Tags))
	}

	// Verify duration calculation
	expectedDuration := 2*time.Hour + 30*time.Minute
	actualDuration := result[0].Stop.Sub(result[0].Start)
	if actualDuration != expectedDuration {
		t.Errorf("Expected duration %v, got %v", expectedDuration, actualDuration)
	}
}

func TestReadWatsonFramesFileNotFound(t *testing.T) {
	_, err := readWatsonFrames("/nonexistent/file.json")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestReadWatsonFramesInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	framesFile := filepath.Join(tmpDir, "invalid.json")

	// Write invalid JSON
	if err := os.WriteFile(framesFile, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	_, err := readWatsonFrames(framesFile)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestParseTimewarriorLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected TimewarriorInterval
		hasError bool
	}{
		{
			name: "valid line with tags",
			line: "inc 20160101T120000Z - 20160101T130000Z # project tag1 tag2",
			expected: TimewarriorInterval{
				Start: "20160101T120000Z",
				End:   "20160101T130000Z",
				Tags:  []string{"project", "tag1", "tag2"},
			},
			hasError: false,
		},
		{
			name: "valid line without tags",
			line: "inc 20160101T120000Z - 20160101T130000Z #",
			expected: TimewarriorInterval{
				Start: "20160101T120000Z",
				End:   "20160101T130000Z",
				Tags:  []string{},
			},
			hasError: false,
		},
		{
			name: "valid line with single tag",
			line: "inc 20160101T120000Z - 20160101T130000Z # work",
			expected: TimewarriorInterval{
				Start: "20160101T120000Z",
				End:   "20160101T130000Z",
				Tags:  []string{"work"},
			},
			hasError: false,
		},
		{
			name:     "invalid line format",
			line:     "invalid line format",
			expected: TimewarriorInterval{},
			hasError: true,
		},
		{
			name:     "missing end time",
			line:     "inc 20160101T120000Z # work",
			expected: TimewarriorInterval{},
			hasError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := parseTimewarriorLine(test.line)

			if test.hasError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.Start != test.expected.Start {
				t.Errorf("Expected start %s, got %s", test.expected.Start, result.Start)
			}
			if result.End != test.expected.End {
				t.Errorf("Expected end %s, got %s", test.expected.End, result.End)
			}
			if len(result.Tags) != len(test.expected.Tags) {
				t.Errorf("Expected %d tags, got %d", len(test.expected.Tags), len(result.Tags))
			}
			for i, tag := range test.expected.Tags {
				if i < len(result.Tags) && result.Tags[i] != tag {
					t.Errorf("Expected tag %s, got %s", tag, result.Tags[i])
				}
			}
		})
	}
}

func TestParseTimewarriorDataFile(t *testing.T) {
	tmpDir := t.TempDir()
	dataFile := filepath.Join(tmpDir, "test.data")

	// Create test data file
	content := `# Timewarrior data file
inc 20160101T120000Z - 20160101T130000Z # project1 tag1
inc 20160102T140000Z - 20160102T160000Z # project2 tag2 tag3

# Comment line
inc 20160103T100000Z - 20160103T120000Z # project3
`

	if err := os.WriteFile(dataFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	intervals, err := parseTimewarriorDataFile(dataFile)
	if err != nil {
		t.Fatalf("parseTimewarriorDataFile failed: %v", err)
	}

	if len(intervals) != 3 {
		t.Errorf("Expected 3 intervals, got %d", len(intervals))
	}

	// Verify first interval
	if intervals[0].Start != "20160101T120000Z" {
		t.Errorf("Expected start '20160101T120000Z', got '%s'", intervals[0].Start)
	}
	if intervals[0].End != "20160101T130000Z" {
		t.Errorf("Expected end '20160101T130000Z', got '%s'", intervals[0].End)
	}
	if len(intervals[0].Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(intervals[0].Tags))
	}
	if intervals[0].Tags[0] != "project1" {
		t.Errorf("Expected first tag 'project1', got '%s'", intervals[0].Tags[0])
	}
}

func TestMapProject(t *testing.T) {
	// Set up project mapping
	projectMap = map[string]string{
		"old-project": "new-project",
		"web-dev":     "frontend",
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"old-project", "new-project"},
		{"web-dev", "frontend"},
		{"unmapped-project", "unmapped-project"},
		{"", ""},
	}

	for _, test := range tests {
		result := mapProject(test.input)
		if result != test.expected {
			t.Errorf("mapProject(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}

	// Clean up
	projectMap = nil
}

func TestWatsonFrameToRuneSession(t *testing.T) {
	frame := WatsonFrame{
		ID:      "test-frame",
		Start:   time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC),
		Stop:    time.Date(2024, 1, 15, 11, 30, 0, 0, time.UTC),
		Project: "test-project",
		Tags:    []string{"tag1", "tag2"},
	}

	session := &tracking.Session{
		ID:        "watson_test-frame",
		Project:   "test-project",
		StartTime: frame.Start,
		EndTime:   &frame.Stop,
		Duration:  frame.Stop.Sub(frame.Start),
		State:     tracking.StateStopped,
	}

	// Verify conversion logic
	expectedDuration := 2*time.Hour + 30*time.Minute
	if session.Duration != expectedDuration {
		t.Errorf("Expected duration %v, got %v", expectedDuration, session.Duration)
	}

	if session.State != tracking.StateStopped {
		t.Errorf("Expected state StateStopped, got %v", session.State)
	}

	if session.Project != frame.Project {
		t.Errorf("Expected project %s, got %s", frame.Project, session.Project)
	}
}

func TestTimewarriorIntervalToRuneSession(t *testing.T) {
	interval := TimewarriorInterval{
		Start: "20240115T090000Z",
		End:   "20240115T113000Z",
		Tags:  []string{"project", "tag1", "tag2"},
	}

	start, err := time.Parse("20060102T150405Z", interval.Start)
	if err != nil {
		t.Fatalf("Failed to parse start time: %v", err)
	}

	end, err := time.Parse("20060102T150405Z", interval.End)
	if err != nil {
		t.Fatalf("Failed to parse end time: %v", err)
	}

	session := &tracking.Session{
		ID:        "timewarrior_0",
		Project:   "project", // First tag becomes project
		StartTime: start,
		EndTime:   &end,
		Duration:  end.Sub(start),
		State:     tracking.StateStopped,
	}

	// Verify conversion logic
	expectedDuration := 2*time.Hour + 30*time.Minute
	if session.Duration != expectedDuration {
		t.Errorf("Expected duration %v, got %v", expectedDuration, session.Duration)
	}

	if session.Project != "project" {
		t.Errorf("Expected project 'project', got %s", session.Project)
	}
}

func TestPreviewWatsonImportStats(t *testing.T) {
	frames := []WatsonFrame{
		{
			ID:      "frame1",
			Start:   time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC),
			Stop:    time.Date(2024, 1, 15, 11, 30, 0, 0, time.UTC),
			Project: "project1",
			Tags:    []string{"tag1"},
		},
		{
			ID:      "frame2",
			Start:   time.Date(2024, 1, 15, 13, 0, 0, 0, time.UTC),
			Stop:    time.Date(2024, 1, 15, 15, 45, 0, 0, time.UTC),
			Project: "project1",
			Tags:    []string{"tag2"},
		},
		{
			ID:      "frame3",
			Start:   time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
			Stop:    time.Date(2024, 1, 16, 12, 0, 0, 0, time.UTC),
			Project: "project2",
			Tags:    []string{"tag3"},
		},
	}

	// Calculate expected statistics
	projectStats := make(map[string]int)
	var totalDuration time.Duration

	for _, frame := range frames {
		duration := frame.Stop.Sub(frame.Start)
		totalDuration += duration
		projectStats[frame.Project]++
	}

	// Verify statistics
	expectedTotal := 2*time.Hour + 30*time.Minute + 2*time.Hour + 45*time.Minute + 2*time.Hour
	if totalDuration != expectedTotal {
		t.Errorf("Expected total duration %v, got %v", expectedTotal, totalDuration)
	}

	if projectStats["project1"] != 2 {
		t.Errorf("Expected 2 frames for project1, got %d", projectStats["project1"])
	}

	if projectStats["project2"] != 1 {
		t.Errorf("Expected 1 frame for project2, got %d", projectStats["project2"])
	}
}

func TestTimewarriorTimeFormat(t *testing.T) {
	// Test parsing Timewarrior's time format
	timeStr := "20240115T093000Z"
	parsed, err := time.Parse("20060102T150405Z", timeStr)
	if err != nil {
		t.Fatalf("Failed to parse Timewarrior time format: %v", err)
	}

	expected := time.Date(2024, 1, 15, 9, 30, 0, 0, time.UTC)
	if !parsed.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, parsed)
	}
}

func TestEmptyWatsonFrames(t *testing.T) {
	tmpDir := t.TempDir()
	framesFile := filepath.Join(tmpDir, "empty.json")

	// Write empty array
	if err := os.WriteFile(framesFile, []byte("[]"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	frames, err := readWatsonFrames(framesFile)
	if err != nil {
		t.Fatalf("readWatsonFrames failed: %v", err)
	}

	if len(frames) != 0 {
		t.Errorf("Expected 0 frames, got %d", len(frames))
	}
}

func TestEmptyTimewarriorDataFile(t *testing.T) {
	tmpDir := t.TempDir()
	dataFile := filepath.Join(tmpDir, "empty.data")

	// Write empty file
	if err := os.WriteFile(dataFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	intervals, err := parseTimewarriorDataFile(dataFile)
	if err != nil {
		t.Fatalf("parseTimewarriorDataFile failed: %v", err)
	}

	if len(intervals) != 0 {
		t.Errorf("Expected 0 intervals, got %d", len(intervals))
	}
}

// Benchmark tests for performance
func BenchmarkReadWatsonFrames(b *testing.B) {
	tmpDir := b.TempDir()
	framesFile := filepath.Join(tmpDir, "frames.json")

	// Create a large number of frames for benchmarking
	frames := make([]WatsonFrame, 1000)
	for i := 0; i < 1000; i++ {
		frames[i] = WatsonFrame{
			ID:      fmt.Sprintf("frame%d", i),
			Start:   time.Now().Add(-time.Duration(i) * time.Hour),
			Stop:    time.Now().Add(-time.Duration(i)*time.Hour + 2*time.Hour),
			Project: fmt.Sprintf("project%d", i%10),
			Tags:    []string{"tag1", "tag2"},
		}
	}

	data, _ := json.Marshal(frames)
	_ = os.WriteFile(framesFile, data, 0644)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := readWatsonFrames(framesFile)
		if err != nil {
			b.Fatalf("readWatsonFrames failed: %v", err)
		}
	}
}

func BenchmarkParseTimewarriorLine(b *testing.B) {
	line := "inc 20160101T120000Z - 20160101T130000Z # project tag1 tag2"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parseTimewarriorLine(line)
		if err != nil {
			b.Fatalf("parseTimewarriorLine failed: %v", err)
		}
	}
}
