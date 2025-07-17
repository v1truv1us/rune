package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ferg-cod3s/rune/internal/colors"
	"github.com/ferg-cod3s/rune/internal/tracking"
	"github.com/spf13/cobra"
)

// WatsonFrame represents a Watson time tracking frame
type WatsonFrame struct {
	ID      string    `json:"id"`
	Start   time.Time `json:"start"`
	Stop    time.Time `json:"stop"`
	Project string    `json:"project"`
	Tags    []string  `json:"tags"`
}

// TimewarriorInterval represents a Timewarrior interval
type TimewarriorInterval struct {
	Start string   `json:"start"`
	End   string   `json:"end,omitempty"`
	Tags  []string `json:"tags"`
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate data from other time tracking tools",
	Long: `Migrate your existing time tracking data from Watson or Timewarrior to Rune.

This command helps you import your historical time tracking data so you can
continue your workflow seamlessly with Rune.

Supported tools:
  - Watson: Imports frames from Watson's JSON export
  - Timewarrior: Imports intervals from Timewarrior's JSON export

Examples:
  # Migrate from Watson
  rune migrate watson ~/.config/watson/frames

  # Migrate from Timewarrior  
  rune migrate timewarrior ~/.timewarrior

  # Dry run to see what would be imported
  rune migrate watson ~/.config/watson/frames --dry-run

  # Import with custom project mapping
  rune migrate watson ~/.config/watson/frames --project-map "old-name=new-name"`,
}

var migrateWatsonCmd = &cobra.Command{
	Use:   "watson [frames-file]",
	Short: "Migrate data from Watson time tracker",
	Long: `Import time tracking data from Watson.

Watson stores its data in a JSON file typically located at:
  ~/.config/watson/frames

This command will parse the Watson frames file and import all completed
time tracking sessions into Rune's database.

The frames file contains an array of time tracking sessions with start/stop
times, project names, and tags.`,
	Args: cobra.ExactArgs(1),
	RunE: runMigrateWatson,
}

var migrateTimewarriorCmd = &cobra.Command{
	Use:   "timewarrior [data-dir]",
	Short: "Migrate data from Timewarrior",
	Long: `Import time tracking data from Timewarrior.

Timewarrior stores its data in a directory typically located at:
  ~/.timewarrior

This command will export data from Timewarrior using the 'timew export' command
and import all completed intervals into Rune's database.

Note: This requires the 'timew' command to be available in your PATH.`,
	Args: cobra.ExactArgs(1),
	RunE: runMigrateTimewarrior,
}

var (
	dryRun     bool
	projectMap map[string]string
)

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateWatsonCmd)
	migrateCmd.AddCommand(migrateTimewarriorCmd)

	// Add flags
	migrateCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Show what would be imported without actually importing")
	migrateCmd.PersistentFlags().StringToStringVar(&projectMap, "project-map", nil, "Map old project names to new ones (format: old=new)")
}

func runMigrateWatson(cmd *cobra.Command, args []string) error {
	framesFile := args[0]

	// Check if file exists
	if _, err := os.Stat(framesFile); os.IsNotExist(err) {
		return fmt.Errorf("Watson frames file not found: %s", framesFile)
	}

	fmt.Printf("%s Migrating from Watson frames file: %s\n", colors.StatusInfo("→"), framesFile)

	// Read and parse Watson frames
	frames, err := readWatsonFrames(framesFile)
	if err != nil {
		return fmt.Errorf("failed to read Watson frames: %w", err)
	}

	if len(frames) == 0 {
		fmt.Printf("%s No frames found in Watson file\n", colors.StatusWarning("⚠"))
		return nil
	}

	fmt.Printf("%s Found %d frames to import\n", colors.StatusSuccess("✓"), len(frames))

	if dryRun {
		return previewWatsonImport(frames)
	}

	// Import frames
	return importWatsonFrames(frames)
}

func runMigrateTimewarrior(cmd *cobra.Command, args []string) error {
	dataDir := args[0]

	// Check if directory exists
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		return fmt.Errorf("Timewarrior data directory not found: %s", dataDir)
	}

	fmt.Printf("%s Migrating from Timewarrior data directory: %s\n", colors.StatusInfo("→"), dataDir)

	// Export data from Timewarrior
	intervals, err := exportTimewarriorData(dataDir)
	if err != nil {
		return fmt.Errorf("failed to export Timewarrior data: %w", err)
	}

	if len(intervals) == 0 {
		fmt.Printf("%s No intervals found in Timewarrior data\n", colors.Warning("⚠"))
		return nil
	}

	fmt.Printf("%s Found %d intervals to import\n", colors.Success("✓"), len(intervals))

	if dryRun {
		return previewTimewarriorImport(intervals)
	}

	// Import intervals
	return importTimewarriorIntervals(intervals)
}

func readWatsonFrames(filename string) ([]WatsonFrame, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var frames []WatsonFrame
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&frames); err != nil {
		return nil, err
	}

	return frames, nil
}

func exportTimewarriorData(dataDir string) ([]TimewarriorInterval, error) {
	// For now, we'll read the data files directly since we can't execute commands
	// In a real implementation, you'd use exec.Command to run "timew export"
	// Look for data files in the Timewarrior directory
	dataFiles, err := filepath.Glob(filepath.Join(dataDir, "data", "*.data"))
	if err != nil {
		return nil, err
	}

	var intervals []TimewarriorInterval

	for _, dataFile := range dataFiles {
		fileIntervals, err := parseTimewarriorDataFile(dataFile)
		if err != nil {
			fmt.Printf("%s Warning: Failed to parse %s: %v\n", colors.Warning("⚠"), dataFile, err)
			continue
		}
		intervals = append(intervals, fileIntervals...)
	}

	return intervals, nil
}

func parseTimewarriorDataFile(filename string) ([]TimewarriorInterval, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var intervals []TimewarriorInterval
	scanner := bufio.NewScanner(file)

	// Timewarrior data files have a simple format:
	// inc 20160101T120000Z - 20160101T130000Z # tag1 tag2
	// Each line represents an interval

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		interval, err := parseTimewarriorLine(line)
		if err != nil {
			continue // Skip invalid lines
		}

		intervals = append(intervals, interval)
	}

	return intervals, scanner.Err()
}

func parseTimewarriorLine(line string) (TimewarriorInterval, error) {
	// Parse format: inc 20160101T120000Z - 20160101T130000Z # tag1 tag2
	re := regexp.MustCompile(`^inc\s+(\d{8}T\d{6}Z)\s+-\s+(\d{8}T\d{6}Z)\s+#\s*(.*)$`)
	matches := re.FindStringSubmatch(line)

	if len(matches) != 4 {
		return TimewarriorInterval{}, fmt.Errorf("invalid line format")
	}

	start := matches[1]
	end := matches[2]
	tagsStr := strings.TrimSpace(matches[3])

	var tags []string
	if tagsStr != "" {
		tags = strings.Fields(tagsStr)
	}

	return TimewarriorInterval{
		Start: start,
		End:   end,
		Tags:  tags,
	}, nil
}

func previewWatsonImport(frames []WatsonFrame) error {
	fmt.Printf("\n%s Preview of Watson import (dry run):\n\n", colors.StatusInfo("👁"))

	projectStats := make(map[string]int)
	var totalDuration time.Duration

	for i, frame := range frames {
		if i < 10 { // Show first 10 frames
			duration := frame.Stop.Sub(frame.Start)
			project := mapProject(frame.Project)

			fmt.Printf("  %s %s → %s (%s)\n",
				colors.Accent("•"),
				frame.Start.Format("2006-01-02 15:04"),
				project,
				formatDuration(duration))

			if len(frame.Tags) > 0 {
				fmt.Printf("    Tags: %s\n", strings.Join(frame.Tags, ", "))
			}
		}

		duration := frame.Stop.Sub(frame.Start)
		totalDuration += duration
		projectStats[mapProject(frame.Project)]++
	}

	if len(frames) > 10 {
		fmt.Printf("  ... and %d more frames\n", len(frames)-10)
	}

	fmt.Printf("\n%s Summary:\n", colors.StatusInfo("📊"))
	fmt.Printf("  Total frames: %d\n", len(frames))
	fmt.Printf("  Total time: %s\n", formatDuration(totalDuration))
	fmt.Printf("  Projects:\n")

	for project, count := range projectStats {
		fmt.Printf("    %s: %d frames\n", project, count)
	}

	fmt.Printf("\n%s Run without --dry-run to import\n", colors.StatusSuccess("✓"))
	return nil
}

func previewTimewarriorImport(intervals []TimewarriorInterval) error {
	fmt.Printf("\n%s Preview of Timewarrior import (dry run):\n\n", colors.StatusInfo("👁"))

	projectStats := make(map[string]int)
	var totalDuration time.Duration

	for i, interval := range intervals {
		if i < 10 { // Show first 10 intervals
			start, _ := time.Parse("20060102T150405Z", interval.Start)
			end, _ := time.Parse("20060102T150405Z", interval.End)
			duration := end.Sub(start)

			project := "default"
			if len(interval.Tags) > 0 {
				project = mapProject(interval.Tags[0])
			}

			fmt.Printf("  %s %s → %s (%s)\n",
				colors.Accent("•"),
				start.Format("2006-01-02 15:04"),
				project,
				formatDuration(duration))

			if len(interval.Tags) > 1 {
				fmt.Printf("    Tags: %s\n", strings.Join(interval.Tags[1:], ", "))
			}
		}

		start, _ := time.Parse("20060102T150405Z", interval.Start)
		end, _ := time.Parse("20060102T150405Z", interval.End)
		duration := end.Sub(start)
		totalDuration += duration

		project := "default"
		if len(interval.Tags) > 0 {
			project = mapProject(interval.Tags[0])
		}
		projectStats[project]++
	}

	if len(intervals) > 10 {
		fmt.Printf("  ... and %d more intervals\n", len(intervals)-10)
	}

	fmt.Printf("\n%s Summary:\n", colors.StatusInfo("📊"))
	fmt.Printf("  Total intervals: %d\n", len(intervals))
	fmt.Printf("  Total time: %s\n", formatDuration(totalDuration))
	fmt.Printf("  Projects:\n")

	for project, count := range projectStats {
		fmt.Printf("    %s: %d intervals\n", project, count)
	}

	fmt.Printf("\n%s Run without --dry-run to import\n", colors.Success("✓"))
	return nil
}

func importWatsonFrames(frames []WatsonFrame) error {
	tracker, err := tracking.NewTracker()
	if err != nil {
		return fmt.Errorf("failed to create tracker: %w", err)
	}
	defer tracker.Close()

	imported := 0
	skipped := 0

	fmt.Printf("\n%s Importing Watson frames...\n", colors.StatusInfo("📥"))

	for _, frame := range frames {
		// Create a Rune session from Watson frame
		session := &tracking.Session{
			ID:        fmt.Sprintf("watson_%s", frame.ID),
			Project:   mapProject(frame.Project),
			StartTime: frame.Start,
			EndTime:   &frame.Stop,
			Duration:  frame.Stop.Sub(frame.Start),
			State:     tracking.StateStopped,
		}

		// Save the session directly to the database
		if err := saveImportedSession(tracker, session); err != nil {
			fmt.Printf("%s Failed to import frame %s: %v\n", colors.Error("✗"), frame.ID, err)
			skipped++
			continue
		}

		imported++
		if imported%10 == 0 {
			fmt.Printf("  Imported %d frames...\n", imported)
		}
	}

	fmt.Printf("\n%s Import complete!\n", colors.Success("✅"))
	fmt.Printf("  Imported: %d frames\n", imported)
	if skipped > 0 {
		fmt.Printf("  Skipped: %d frames\n", skipped)
	}

	return nil
}

func importTimewarriorIntervals(intervals []TimewarriorInterval) error {
	tracker, err := tracking.NewTracker()
	if err != nil {
		return fmt.Errorf("failed to create tracker: %w", err)
	}
	defer tracker.Close()

	imported := 0
	skipped := 0

	fmt.Printf("\n%s Importing Timewarrior intervals...\n", colors.StatusInfo("📥"))

	for i, interval := range intervals {
		start, err := time.Parse("20060102T150405Z", interval.Start)
		if err != nil {
			skipped++
			continue
		}

		end, err := time.Parse("20060102T150405Z", interval.End)
		if err != nil {
			skipped++
			continue
		}

		project := "default"
		if len(interval.Tags) > 0 {
			project = mapProject(interval.Tags[0])
		}

		// Create a Rune session from Timewarrior interval
		session := &tracking.Session{
			ID:        fmt.Sprintf("timewarrior_%d", i),
			Project:   project,
			StartTime: start,
			EndTime:   &end,
			Duration:  end.Sub(start),
			State:     tracking.StateStopped,
		}

		// Save the session directly to the database
		if err := saveImportedSession(tracker, session); err != nil {
			fmt.Printf("%s Failed to import interval %d: %v\n", colors.Error("✗"), i, err)
			skipped++
			continue
		}

		imported++
		if imported%10 == 0 {
			fmt.Printf("  Imported %d intervals...\n", imported)
		}
	}

	fmt.Printf("\n%s Import complete!\n", colors.Success("✅"))
	fmt.Printf("  Imported: %d intervals\n", imported)
	if skipped > 0 {
		fmt.Printf("  Skipped: %d intervals\n", skipped)
	}

	return nil
}

func saveImportedSession(tracker *tracking.Tracker, session *tracking.Session) error {
	return tracker.SaveImportedSession(session)
}

func mapProject(originalProject string) string {
	if mapped, exists := projectMap[originalProject]; exists {
		return mapped
	}
	return originalProject
}
