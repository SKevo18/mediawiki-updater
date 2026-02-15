package updater

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/SKevo18/mediawiki-updater/internal/config"
	"github.com/SKevo18/mediawiki-updater/internal/downloader"
	"github.com/SKevo18/mediawiki-updater/internal/extractor"
	"github.com/SKevo18/mediawiki-updater/internal/mediawiki"
)

// Updater manages the MediaWiki update process
type Updater struct {
	config      *config.Config
	downloader  *downloader.Downloader
	extractor   *extractor.Extractor
	mwParser    *mediawiki.Parser
	ignorePaths []string
}

// Options contains configuration options for the updater
type Options struct {
	ConfigPath  string
	TargetDir   string
	IgnorePaths []string
}

// NewUpdater creates a new Updater instance
func NewUpdater(opts Options) (*Updater, error) {
	cfg, err := config.LoadConfig(opts.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	ignorePaths := opts.IgnorePaths
	if len(ignorePaths) == 0 {
		ignorePaths = []string{"LocalSettings.php", ".htaccess", "images"}
	}

	return &Updater{
		config:      cfg,
		downloader:  downloader.NewDownloader(),
		extractor:   extractor.NewExtractor(),
		mwParser:    mediawiki.NewParser(),
		ignorePaths: ignorePaths,
	}, nil
}

// Update performs the complete MediaWiki update process
func (u *Updater) Update(targetDir string) error {
	tempDir, err := os.MkdirTemp("", "mediawiki-temp-")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Download MediaWiki core
	if err := u.downloadMediaWikiCore(tempDir); err != nil {
		return fmt.Errorf("failed to download MediaWiki core: %w", err)
	}

	// Download extensions
	if err := u.downloadExtensions(tempDir); err != nil {
		return fmt.Errorf("failed to download extensions: %w", err)
	}

	// Download skins
	if err := u.downloadSkins(tempDir); err != nil {
		return fmt.Errorf("failed to download skins: %w", err)
	}

	// Copy contents to target directory
	if err := u.extractor.CopyContents(tempDir, targetDir, u.ignorePaths); err != nil {
		return fmt.Errorf("failed to copy contents: %w", err)
	}

	return nil
}

// downloadMediaWikiCore downloads the MediaWiki core
func (u *Updater) downloadMediaWikiCore(tempDir string) error {
	version := u.config.MediaWiki.Version
	if version == "" {
		return fmt.Errorf("MediaWiki version not specified in config")
	}

	fmt.Printf("Downloading MediaWiki core version %s...\n", version)

	downloadURL, err := u.mwParser.GetDownloadURL(version)
	if err != nil {
		return err
	}

	return u.downloader.DownloadAndExtract(downloadURL, tempDir, true)
}

// downloadExtensions downloads all configured extensions
func (u *Updater) downloadExtensions(tempDir string) error {
	if len(u.config.Extensions) == 0 {
		fmt.Println("No extensions configured, skipping...")
		return nil
	}

	extensionsDir := filepath.Join(tempDir, "extensions")
	if err := os.MkdirAll(extensionsDir, 0o755); err != nil {
		return err
	}

	versionTag, err := u.getVersionTag()
	if err != nil {
		return err
	}

	fmt.Printf("Downloading %d extensions...\n", len(u.config.Extensions))

	for _, ext := range u.config.Extensions {
		fmt.Printf("  - %s (from %s)\n", ext.Name, ext.Distributor)
		if err := u.downloader.DownloadComponent(ext, extensionsDir, versionTag); err != nil {
			fmt.Printf("    WARNING: Failed to download extension %s: %v\n", ext.Name, err)
			// Continue with other extensions instead of failing completely
		}
	}

	return nil
}

// downloadSkins downloads all configured skins
func (u *Updater) downloadSkins(tempDir string) error {
	if len(u.config.Skins) == 0 {
		fmt.Println("No skins configured, skipping...")
		return nil
	}

	skinsDir := filepath.Join(tempDir, "skins")
	if err := os.MkdirAll(skinsDir, 0o755); err != nil {
		return err
	}

	versionTag, err := u.getVersionTag()
	if err != nil {
		return err
	}

	fmt.Printf("Downloading %d skins...\n", len(u.config.Skins))

	for _, skin := range u.config.Skins {
		fmt.Printf("  - %s (from %s)\n", skin.Name, skin.Distributor)
		if err := u.downloader.DownloadComponent(skin, skinsDir, versionTag); err != nil {
			fmt.Printf("    WARNING: Failed to download skin %s: %v\n", skin.Name, err)
			// Continue with other skins instead of failing completely
		}
	}

	return nil
}

// getVersionTag converts the MediaWiki version to the format used by ExtDist (e.g., "1.43.1" -> "REL1_43")
func (u *Updater) getVersionTag() (string, error) {
	version := u.config.MediaWiki.Version
	if version == "" {
		return "", fmt.Errorf("MediaWiki version not specified in config")
	}

	// Extract major.minor version and convert to REL format
	re := regexp.MustCompile(`^(\d+)\.(\d+)(\.\d+.*)?$`)
	matches := re.FindStringSubmatch(version)
	if len(matches) < 3 {
		return "", fmt.Errorf("invalid version format: %s", version)
	}

	return fmt.Sprintf("REL%s_%s", matches[1], matches[2]), nil
}
