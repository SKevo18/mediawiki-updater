package mediawiki

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/SKevo18/mediawiki-updater/internal/httputil"
)

const (
	BaseDownloadURL = "https://releases.wikimedia.org/mediawiki/"
)

// Parser handles MediaWiki release page parsing
type Parser struct{}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{}
}

// GetDownloadURL parses the MediaWiki release page to find the download URL for a specific version
func (p *Parser) GetDownloadURL(version string) (string, error) {
	// Extract major.minor version for URL construction (e.g., "1.43.1" -> "1.43")
	majorMinor, err := p.extractMajorMinor(version)
	if err != nil {
		return "", err
	}

	releasePageURL := fmt.Sprintf("%s%s/", BaseDownloadURL, majorMinor)

	resp, err := httputil.Get(releasePageURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch release page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("release page returned status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse release page: %w", err)
	}

	// Look for the exact version tar.gz file
	targetFilename := fmt.Sprintf("mediawiki-%s.tar.gz", version)

	var downloadURL string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href := s.AttrOr("href", "")
		if strings.Contains(href, targetFilename) {
			downloadURL = releasePageURL + href
		}
	})

	if downloadURL == "" {
		return "", fmt.Errorf("no download URL found for MediaWiki version %s", version)
	}

	return downloadURL, nil
}

// extractMajorMinor extracts the major.minor version from a full version string
func (p *Parser) extractMajorMinor(version string) (string, error) {
	// Match pattern like "1.43.1" and extract "1.43"
	re := regexp.MustCompile(`^(\d+\.\d+)(\.\d+.*)?$`)
	matches := re.FindStringSubmatch(version)
	if len(matches) < 2 {
		return "", fmt.Errorf("invalid version format: %s", version)
	}
	return matches[1], nil
}
