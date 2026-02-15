package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/SKevo18/mediawiki-updater/internal/config"
	"github.com/SKevo18/mediawiki-updater/internal/extractor"
	"github.com/SKevo18/mediawiki-updater/internal/httputil"
)

const (
	ExtDistURL = "https://extdist.wmflabs.org/dist/extensions/"
	SkinsURL   = "https://extdist.wmflabs.org/dist/skins/"
)

// Downloader handles downloading files from various sources
type Downloader struct {
	extractor *extractor.Extractor
}

// NewDownloader creates a new Downloader instance
func NewDownloader() *Downloader {
	return &Downloader{
		extractor: extractor.NewExtractor(),
	}
}

// DownloadFile downloads a file from URL to the specified path
func (d *Downloader) DownloadFile(url, targetPath string) error {
	resp, err := httputil.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	file, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// DownloadAndExtract downloads a file and extracts it to the target directory
func (d *Downloader) DownloadAndExtract(url, targetDir string, isMediaWiki bool) error {
	tempFile, err := os.CreateTemp("", "mw-temp-*.tar.gz")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if err := d.DownloadFile(url, tempFile.Name()); err != nil {
		return err
	}

	// Reopen for reading
	tempFile.Close()
	file, err := os.Open(tempFile.Name())
	if err != nil {
		return err
	}
	defer file.Close()

	if isMediaWiki {
		return d.extractor.ExtractMediaWikiCore(file, targetDir)
	}
	return d.extractor.ExtractArchive(file, targetDir, nil)
}

// DownloadComponent downloads a component (extension or skin) based on its configuration
func (d *Downloader) DownloadComponent(component config.ComponentConfig, targetDir, versionTag string) error {
	switch component.Distributor {
	case "extdist":
		return d.downloadFromExtDist(component, targetDir, versionTag)
	case "git":
		return d.downloadFromGit(component, targetDir)
	default:
		return fmt.Errorf("unknown distributor: %s", component.Distributor)
	}
}

// downloadFromExtDist downloads a component from the ExtDist service
func (d *Downloader) downloadFromExtDist(component config.ComponentConfig, targetDir, versionTag string) error {
	// Determine base URL based on target directory
	baseURL := ExtDistURL
	if strings.Contains(targetDir, "skins") {
		baseURL = SkinsURL
	}

	// Use component version if specified, otherwise use the global version tag
	version := versionTag
	if component.Version != "" {
		version = component.Version
	}

	downloadURL, err := d.getExtDistDownloadURL(baseURL, component.Name, version)
	if err != nil {
		return err
	}

	if downloadURL == "" {
		return fmt.Errorf("component not found: %s", component.Name)
	}

	return d.DownloadAndExtract(downloadURL, targetDir, false)
}

// downloadFromGit clones a Git repository
func (d *Downloader) downloadFromGit(component config.ComponentConfig, targetDir string) error {
	// component.Name should be the git repository URL for git distributor
	repoURL := component.Name
	version := component.Version
	if version == "" {
		version = "master"
	}

	// Create target directory for this component
	componentDir := filepath.Join(targetDir, extractRepoName(repoURL))

	// Clone the repository
	cmd := exec.Command("git", "clone", "--branch", version, "--depth", "1", repoURL, componentDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone git repository %s: %w", repoURL, err)
	}

	// Remove .git directory to clean up
	gitDir := filepath.Join(componentDir, ".git")
	if err := os.RemoveAll(gitDir); err != nil {
		return fmt.Errorf("failed to remove .git directory: %w", err)
	}

	return nil
}

// getExtDistDownloadURL finds the download URL for a component from ExtDist
func (d *Downloader) getExtDistDownloadURL(baseURL, componentName, version string) (string, error) {
	fmt.Printf("    Checking %s for %s-%s\n", baseURL, componentName, version)

	resp, err := httputil.Get(baseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	var downloadURL string
	searchPattern := fmt.Sprintf("%s-%s", componentName, version)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href := s.AttrOr("href", "")
		if strings.Contains(href, searchPattern) && strings.HasSuffix(href, ".tar.gz") {
			downloadURL = href
		}
	})

	if downloadURL == "" {
		fmt.Printf("    No match found for pattern: %s\n", searchPattern)
		return "", nil
	}

	fullURL := baseURL + downloadURL
	fmt.Printf("    Found: %s\n", fullURL)
	return fullURL, nil
}

// extractRepoName extracts repository name from a Git URL
func extractRepoName(repoURL string) string {
	// Remove .git suffix and extract last part of URL
	repoURL = strings.TrimSuffix(repoURL, ".git")
	parts := strings.Split(repoURL, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "unknown"
}
