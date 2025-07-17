package commands

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{0, "0h 0m"},
		{30 * time.Second, "0h 0m"},
		{1 * time.Minute, "0h 1m"},
		{30 * time.Minute, "0h 30m"},
		{1 * time.Hour, "1h 0m"},
		{1*time.Hour + 30*time.Minute, "1h 30m"},
		{2*time.Hour + 45*time.Minute, "2h 45m"},
		{25 * time.Hour, "25h 0m"},
	}

	for _, test := range tests {
		result := formatDuration(test.duration)
		if result != test.expected {
			t.Errorf("formatDuration(%v) = %q, expected %q", test.duration, result, test.expected)
		}
	}
}

func TestFormatRelativeTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "just now",
			time:     now.Add(-30 * time.Second),
			expected: "just now",
		},
		{
			name:     "1 minute ago",
			time:     now.Add(-1 * time.Minute),
			expected: "1m ago",
		},
		{
			name:     "5 minutes ago",
			time:     now.Add(-5 * time.Minute),
			expected: "5m ago",
		},
		{
			name:     "1 hour ago",
			time:     now.Add(-1 * time.Hour),
			expected: "1h ago",
		},
		{
			name:     "1 hour 30 minutes ago",
			time:     now.Add(-1*time.Hour - 30*time.Minute),
			expected: "1h 30m ago",
		},
		{
			name:     "2 hours ago",
			time:     now.Add(-2 * time.Hour),
			expected: "2h ago",
		},
		{
			name:     "2 hours 15 minutes ago",
			time:     now.Add(-2*time.Hour - 15*time.Minute),
			expected: "2h 15m ago",
		},
		{
			name:     "yesterday",
			time:     now.AddDate(0, 0, -1),
			expected: "yesterday",
		},
		{
			name:     "2 days ago",
			time:     now.AddDate(0, 0, -2),
			expected: "2 days ago",
		},
		{
			name:     "1 week ago",
			time:     now.AddDate(0, 0, -7),
			expected: "1 week ago",
		},
		{
			name:     "2 weeks ago",
			time:     now.AddDate(0, 0, -14),
			expected: "2 weeks ago",
		},
		{
			name:     "1 month ago",
			time:     now.AddDate(0, -1, 0),
			expected: "1 month ago",
		},
		{
			name:     "3 months ago",
			time:     now.AddDate(0, -3, 0),
			expected: "3 months ago",
		},
		{
			name:     "1 year ago",
			time:     now.AddDate(-1, 0, 0),
			expected: "1 year ago",
		},
		{
			name:     "2 years ago",
			time:     now.AddDate(-2, 0, 0),
			expected: "2 years ago",
		},
		{
			name:     "future time",
			time:     now.Add(1 * time.Hour),
			expected: "in the future",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := formatRelativeTime(test.time)
			if result != test.expected {
				t.Errorf("formatRelativeTime(%v) = %q, expected %q", test.time, result, test.expected)
			}
		})
	}
}
