package commands

import (
	"fmt"
	"os"

	"github.com/ferg-cod3s/rune/internal/colors"
	"github.com/ferg-cod3s/rune/internal/config"
	"github.com/ferg-cod3s/rune/internal/telemetry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	version = "dev"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rune",
	Short: "Ancient wisdom for modern workflows",
	Long: `Rune is a developer-first CLI productivity platform that automates daily work 
rituals, enforces healthy work-life boundaries, and integrates seamlessly 
with existing developer workflows.

Cast your daily runes and master your workflow.`,
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	// Ensure telemetry is properly closed on exit
	defer telemetry.Close()
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig, initTelemetry, initColors)

	// Version template will be set dynamically in initColors

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rune/config.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable colored output")

	// Bind flags to viper
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("no-color", rootCmd.PersistentFlags().Lookup("no-color"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".rune" (without extension).
		viper.AddConfigPath(home + "/.rune")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}

// initTelemetry initializes telemetry tracking
func initTelemetry() {
	// Get telemetry configuration from environment variables or config
	segmentWriteKey := os.Getenv("RUNE_SEGMENT_WRITE_KEY")
	sentryDSN := os.Getenv("RUNE_SENTRY_DSN")

	if os.Getenv("RUNE_DEBUG") == "true" {
		fmt.Printf("DEBUG: initTelemetry called\n")
		fmt.Printf("DEBUG: Env Segment Key: %s\n", segmentWriteKey)
		fmt.Printf("DEBUG: Env Sentry DSN: %s\n", sentryDSN)
	}

	// Try to load from config if environment variables are not set
	if cfg, err := config.Load(); err == nil {
		if segmentWriteKey == "" {
			segmentWriteKey = cfg.Integrations.Telemetry.SegmentWriteKey
		}
		if sentryDSN == "" {
			sentryDSN = cfg.Integrations.Telemetry.SentryDSN
		}
		if os.Getenv("RUNE_DEBUG") == "true" {
			fmt.Printf("DEBUG: Config loaded - Segment: %s, Sentry: %s\n", cfg.Integrations.Telemetry.SegmentWriteKey, cfg.Integrations.Telemetry.SentryDSN)
		}
	} else if os.Getenv("RUNE_DEBUG") == "true" {
		fmt.Printf("DEBUG: Config load failed: %v\n", err)
	}

	if os.Getenv("RUNE_DEBUG") == "true" {
		fmt.Printf("DEBUG: Final keys - Segment: %s, Sentry: %s\n", segmentWriteKey, sentryDSN)
	}

	telemetry.Initialize(segmentWriteKey, sentryDSN)
}

// initColors initializes the color system
func initColors() {
	// Check for --no-color flag
	if viper.GetBool("no-color") {
		colors.SetColorMode(colors.ColorNever)
	} else {
		colors.Initialize()
	}

	// Set version template with appropriate coloring
	logoText := ` ______     __  __     __   __     ______   
/\  == \   /\ \/\ \   /\ "-.\ \   /\  ___\  
\ \  __<   \ \ \_\ \  \ \ \-.  \  \ \  __\  
 \ \_\ \_\  \ \_____\  \ \_\\"\_\  \ \_____\
  \/_/ /_/   \/_____/   \/_/ \/_/   \/_____/ 
`
	rootCmd.SetVersionTemplate(colors.Logo(logoText) + `
version {{.Version}}

`)

	// Set long description with appropriate coloring
	rootCmd.Long = colors.Logo(logoText) + `
` + colors.Glow("Rune is a developer-first CLI productivity platform that automates daily work") + `
` + colors.Glow("rituals, enforces healthy work-life boundaries, and integrates seamlessly") + `
` + colors.Glow("with existing developer workflows.") + `

` + colors.Accent("Cast your daily runes and master your workflow.")
}
