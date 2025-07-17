package commands

import (
	"fmt"
	"time"
)

// formatDuration formats a duration as "Xh Ym"
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

// formatRelativeTime formats a time relative to now (e.g., "2h 30m ago", "just now", "yesterday")
func formatRelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	// Handle future times (shouldn't happen in normal usage, but just in case)
	if diff < 0 {
		return "in the future"
	}

	// Less than 1 minute ago
	if diff < time.Minute {
		return "just now"
	}

	// Less than 1 hour ago - show minutes
	if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1m ago"
		}
		return fmt.Sprintf("%dm ago", minutes)
	}

	// Less than 24 hours ago - show hours and minutes
	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		minutes := int(diff.Minutes()) % 60
		if hours == 1 && minutes == 0 {
			return "1h ago"
		} else if minutes == 0 {
			return fmt.Sprintf("%dh ago", hours)
		} else if hours == 1 {
			return fmt.Sprintf("1h %dm ago", minutes)
		}
		return fmt.Sprintf("%dh %dm ago", hours, minutes)
	}

	// Check if it's yesterday
	yesterday := now.AddDate(0, 0, -1)
	if t.Year() == yesterday.Year() && t.Month() == yesterday.Month() && t.Day() == yesterday.Day() {
		return "yesterday"
	}

	// Less than 7 days ago - show days
	if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}

	// Less than 30 days ago - show weeks
	if diff < 30*24*time.Hour {
		weeks := int(diff.Hours() / (24 * 7))
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	}

	// Less than 365 days ago - show months
	if diff < 365*24*time.Hour {
		months := int(diff.Hours() / (24 * 30))
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}

	// More than a year ago - show years
	years := int(diff.Hours() / (24 * 365))
	if years == 1 {
		return "1 year ago"
	}
	return fmt.Sprintf("%d years ago", years)
}
