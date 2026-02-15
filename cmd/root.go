package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SKevo18/mediawiki-updater/internal/updater"
	"github.com/spf13/cobra"
)

var (
	configFile string
	targetDir  string
	verbose    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mediawiki-updater",
	Short: "A tool to download and update MediaWiki core, extensions, and skins",
	Long: `MediaWiki Updater is a CLI tool that downloads and installs MediaWiki core 
along with configured extensions and skins based on an INI configuration file.

The tool supports:
- Downloading MediaWiki core from official releases
- Installing extensions and skins from ExtDist or Git repositories
- Configurable version management
- Preserving specified files during updates`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUpdate()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.ini", "path to configuration file")
	rootCmd.PersistentFlags().StringVarP(&targetDir, "target", "t", ".", "target directory for MediaWiki installation")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
}

func runUpdate() error {
	// Validate target directory
	absTargetDir, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("invalid target directory: %w", err)
	}

	// Create target directory if it doesn't exist
	if err := os.MkdirAll(absTargetDir, 0o755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Validate config file
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("configuration file not found: %s", configFile)
	}

	if verbose {
		fmt.Printf("Using configuration file: %s\n", configFile)
		fmt.Printf("Target directory: %s\n", absTargetDir)
	}

	// Create updater instance
	opts := updater.Options{
		ConfigPath: configFile,
		TargetDir:  absTargetDir,
	}

	updaterInstance, err := updater.NewUpdater(opts)
	if err != nil {
		return err
	}

	// Perform update
	fmt.Println("Starting MediaWiki update process...")
	if err := updaterInstance.Update(absTargetDir); err != nil {
		return err
	}

	fmt.Println("MediaWiki update completed successfully!")
	return nil
}
