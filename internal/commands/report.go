package commands

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ferg-cod3s/rune/internal/tracking"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate time reports",
	Long: `Generate time tracking reports for different periods.

This command can show:
- Daily work summaries
- Weekly productivity reports
- Project-based time allocation
- Monthly trends and insights`,
	RunE: runReport,
}

var (
	today   bool
	week    bool
	month   bool
	project string
	format  string
	output  string
)

func init() {
	rootCmd.AddCommand(reportCmd)

	reportCmd.Flags().BoolVar(&today, "today", false, "Show today's report")
	reportCmd.Flags().BoolVar(&week, "week", false, "Show this week's report")
	reportCmd.Flags().BoolVar(&month, "month", false, "Show this month's report")
	reportCmd.Flags().StringVar(&project, "project", "", "Filter by project name")
	reportCmd.Flags().StringVar(&format, "format", "text", "Output format: text, csv, json")
	reportCmd.Flags().StringVar(&output, "output", "", "Output file (default: stdout)")
}

func runReport(cmd *cobra.Command, args []string) error {
	var sessions []*tracking.Session
	var totalDuration time.Duration
	var err error

	// Initialize tracker
	tracker, err := tracking.NewTracker()
	if err != nil {
		return fmt.Errorf("failed to initialize tracker: %w", err)
	}
	defer tracker.Close()

	// Get data based on period
	if today {
		sessions, totalDuration, err = getTodayData(tracker)
	} else if week {
		sessions, totalDuration, err = getWeekData(tracker)
	} else if month {
		sessions, totalDuration, err = getMonthData(tracker)
	} else {
		// Default to today's report
		sessions, totalDuration, err = getTodayData(tracker)
	}

	if err != nil {
		return err
	}

	// Filter by project if specified
	if project != "" {
		var filteredSessions []*tracking.Session
		var filteredDuration time.Duration
		for _, session := range sessions {
			if session.Project == project {
				filteredSessions = append(filteredSessions, session)
				filteredDuration += session.Duration
			}
		}
		sessions = filteredSessions
		totalDuration = filteredDuration
	}

	// Output based on format
	switch format {
	case "csv":
		return exportCSV(sessions, totalDuration)
	case "json":
		return exportJSON(sessions, totalDuration)
	default:
		// Show text report
		if today {
			return showTodayReport(tracker)
		} else if week {
			return showWeekReport(tracker)
		} else if month {
			return showMonthReport(tracker)
		} else {
			return showTodayReport(tracker)
		}
	}
}

func showTodayReport(tracker *tracking.Tracker) error {
	fmt.Println("📈 Today's Report")
	fmt.Println("=================")
	fmt.Println()

	// Get daily total
	dailyTotal, err := tracker.GetDailyTotal()
	if err != nil {
		return fmt.Errorf("failed to get daily total: %w", err)
	}

	// Get project stats for today
	projectStats, err := tracker.GetProjectStats()
	if err != nil {
		return fmt.Errorf("failed to get project stats: %w", err)
	}

	// Get today's sessions
	sessions, err := tracker.GetSessionHistory(50) // Get more to filter by today
	if err != nil {
		return fmt.Errorf("failed to get session history: %w", err)
	}

	// Filter sessions for today
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	var todaySessions []*tracking.Session
	for _, session := range sessions {
		if session.StartTime.After(today) && session.StartTime.Before(tomorrow) {
			if project == "" || session.Project == project {
				todaySessions = append(todaySessions, session)
			}
		}
	}

	fmt.Printf("Total Time:    %s\n", formatDuration(dailyTotal))
	fmt.Printf("Sessions:      %d\n", len(todaySessions))

	if len(projectStats) > 0 {
		fmt.Println("\nProject Breakdown:")
		for proj, duration := range projectStats {
			if project == "" || proj == project {
				fmt.Printf("  %-15s %s\n", proj, formatDuration(duration))
			}
		}
	}

	if len(todaySessions) > 0 {
		fmt.Println("\nToday's Sessions:")
		for _, session := range todaySessions {
			fmt.Printf("  %s  %-15s  %s\n",
				session.StartTime.Format("15:04"),
				session.Project,
				formatDuration(session.Duration))
		}
	}

	return nil
}

func showWeekReport(tracker *tracking.Tracker) error {
	fmt.Println("📈 This Week's Report")
	fmt.Println("=====================")
	fmt.Println()

	// Get weekly total
	weeklyTotal, err := tracker.GetWeeklyTotal()
	if err != nil {
		return fmt.Errorf("failed to get weekly total: %w", err)
	}

	// Calculate daily average (divide by 7 days)
	dailyAverage := weeklyTotal / 7

	// Get project stats
	projectStats, err := tracker.GetProjectStats()
	if err != nil {
		return fmt.Errorf("failed to get project stats: %w", err)
	}

	fmt.Printf("Total Time:    %s\n", formatDuration(weeklyTotal))
	fmt.Printf("Daily Average: %s\n", formatDuration(dailyAverage))

	if len(projectStats) > 0 {
		fmt.Println("\nProject Breakdown:")
		for proj, duration := range projectStats {
			if project == "" || proj == project {
				fmt.Printf("  %-15s %s\n", proj, formatDuration(duration))
			}
		}
	}

	return nil
}

func showMonthReport(tracker *tracking.Tracker) error {
	fmt.Println("📈 This Month's Report")
	fmt.Println("======================")
	fmt.Println()

	// Get monthly total (approximate - get all sessions and filter)
	sessions, err := tracker.GetSessionHistory(1000) // Get many sessions
	if err != nil {
		return fmt.Errorf("failed to get session history: %w", err)
	}

	// Filter for this month
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)

	var monthlyTotal time.Duration
	var monthlySessions []*tracking.Session
	for _, session := range sessions {
		if session.StartTime.After(monthStart) && session.StartTime.Before(monthEnd) {
			if project == "" || session.Project == project {
				monthlyTotal += session.Duration
				monthlySessions = append(monthlySessions, session)
			}
		}
	}

	// Calculate daily average
	daysInMonth := monthEnd.Sub(monthStart).Hours() / 24
	dailyAverage := time.Duration(float64(monthlyTotal) / daysInMonth)

	// Get project stats
	projectStats, err := tracker.GetProjectStats()
	if err != nil {
		return fmt.Errorf("failed to get project stats: %w", err)
	}

	fmt.Printf("Total Time:    %s\n", formatDuration(monthlyTotal))
	fmt.Printf("Daily Average: %s\n", formatDuration(dailyAverage))
	fmt.Printf("Sessions:      %d\n", len(monthlySessions))

	if len(projectStats) > 0 {
		fmt.Println("\nProject Breakdown:")
		for proj, duration := range projectStats {
			if project == "" || proj == project {
				fmt.Printf("  %-15s %s\n", proj, formatDuration(duration))
			}
		}
	}

	return nil
}

// getTodayData returns today's sessions and total duration
func getTodayData(tracker *tracking.Tracker) ([]*tracking.Session, time.Duration, error) {
	sessions, err := tracker.GetSessionHistory(50)
	if err != nil {
		return nil, 0, err
	}

	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	var todaySessions []*tracking.Session
	var totalDuration time.Duration

	for _, session := range sessions {
		if session.StartTime.After(today) && session.StartTime.Before(tomorrow) {
			todaySessions = append(todaySessions, session)
			totalDuration += session.Duration
		}
	}

	return todaySessions, totalDuration, nil
}

// getWeekData returns this week's sessions and total duration
func getWeekData(tracker *tracking.Tracker) ([]*tracking.Session, time.Duration, error) {
	sessions, err := tracker.GetSessionHistory(500)
	if err != nil {
		return nil, 0, err
	}

	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	weekStart = weekStart.Truncate(24 * time.Hour)
	weekEnd := weekStart.Add(7 * 24 * time.Hour)
	var weekSessions []*tracking.Session
	var totalDuration time.Duration

	for _, session := range sessions {
		if session.StartTime.After(weekStart) && session.StartTime.Before(weekEnd) {
			weekSessions = append(weekSessions, session)
			totalDuration += session.Duration
		}
	}

	return weekSessions, totalDuration, nil
}

// getMonthData returns this month's sessions and total duration
func getMonthData(tracker *tracking.Tracker) ([]*tracking.Session, time.Duration, error) {
	sessions, err := tracker.GetSessionHistory(1000)
	if err != nil {
		return nil, 0, err
	}

	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)
	var monthSessions []*tracking.Session
	var totalDuration time.Duration

	for _, session := range sessions {
		if session.StartTime.After(monthStart) && session.StartTime.Before(monthEnd) {
			monthSessions = append(monthSessions, session)
			totalDuration += session.Duration
		}
	}

	return monthSessions, totalDuration, nil
}

// exportCSV exports sessions to CSV format
func exportCSV(sessions []*tracking.Session, totalDuration time.Duration) error {
	var writer *csv.Writer
	if output != "" {
		file, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer file.Close()
		writer = csv.NewWriter(file)
	} else {
		writer = csv.NewWriter(os.Stdout)
	}
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"Date", "Start Time", "End Time", "Project", "Duration (minutes)", "State"}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write sessions
	for _, session := range sessions {
		endTime := ""
		if session.EndTime != nil {
			endTime = session.EndTime.Format("15:04:05")
		}

		durationMinutes := strconv.FormatFloat(session.Duration.Minutes(), 'f', 2, 64)

		record := []string{
			session.StartTime.Format("2006-01-02"),
			session.StartTime.Format("15:04:05"),
			endTime,
			session.Project,
			durationMinutes,
			session.State.String(),
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	// Write summary row
	totalMinutes := strconv.FormatFloat(totalDuration.Minutes(), 'f', 2, 64)
	summaryRecord := []string{"TOTAL", "", "", "", totalMinutes, ""}
	if err := writer.Write(summaryRecord); err != nil {
		return fmt.Errorf("failed to write CSV summary: %w", err)
	}

	return nil
}

// ReportData represents the structure for JSON export
type ReportData struct {
	GeneratedAt   time.Time              `json:"generated_at"`
	TotalDuration string                 `json:"total_duration"`
	Sessions      []*tracking.Session    `json:"sessions"`
	Summary       map[string]interface{} `json:"summary"`
}

// exportJSON exports sessions to JSON format
func exportJSON(sessions []*tracking.Session, totalDuration time.Duration) error {
	// Calculate project breakdown
	projectStats := make(map[string]time.Duration)
	for _, session := range sessions {
		projectStats[session.Project] += session.Duration
	}

	// Convert project stats to string format for JSON
	projectStatsStr := make(map[string]string)
	for project, duration := range projectStats {
		projectStatsStr[project] = formatDuration(duration)
	}

	data := ReportData{
		GeneratedAt:   time.Now(),
		TotalDuration: formatDuration(totalDuration),
		Sessions:      sessions,
		Summary: map[string]interface{}{
			"total_sessions":    len(sessions),
			"project_breakdown": projectStatsStr,
		},
	}

	var jsonData []byte
	var err error
	jsonData, err = json.MarshalIndent(data, "", "  ")

	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if output != "" {
		if err := os.WriteFile(output, jsonData, 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	} else {
		fmt.Println(string(jsonData))
	}

	return nil
}
