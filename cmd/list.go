package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/SKevo18/mediawiki-updater/internal/httputil"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available MediaWiki versions, extensions, or skins",
	Long: `List available components from various sources.

Available subcommands:
- versions: List available MediaWiki versions
- extensions: List available extensions from ExtDist  
- skins: List available skins from ExtDist`,
}

// versionsCmd lists available MediaWiki versions
var versionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "List available MediaWiki versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listMediaWikiVersions()
	},
}

// extensionsCmd lists available extensions
var extensionsCmd = &cobra.Command{
	Use:   "extensions",
	Short: "List available extensions from ExtDist",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listComponents("https://extdist.wmflabs.org/dist/extensions/", "Extensions")
	},
}

// skinsCmd lists available skins
var skinsCmd = &cobra.Command{
	Use:   "skins",
	Short: "List available skins from ExtDist",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listComponents("https://extdist.wmflabs.org/dist/skins/", "Skins")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(versionsCmd)
	listCmd.AddCommand(extensionsCmd)
	listCmd.AddCommand(skinsCmd)
}

func listMediaWikiVersions() error {
	fmt.Println("Fetching available MediaWiki versions...")

	// Get the main releases page
	resp, err := httputil.Get("https://releases.wikimedia.org/mediawiki/")
	if err != nil {
		return fmt.Errorf("failed to fetch releases page: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse releases page: %w", err)
	}

	var versions []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href := s.AttrOr("href", "")
		// Look for version directories (e.g., "1.43/")
		if strings.Contains(href, ".") && strings.HasSuffix(href, "/") {
			version := strings.TrimSuffix(href, "/")
			if strings.Count(version, ".") == 1 { // Only major.minor versions
				versions = append(versions, version)
			}
		}
	})

	sort.Strings(versions)

	fmt.Printf("\nAvailable MediaWiki version series:\n")
	for _, version := range versions {
		fmt.Printf("- %s\n", version)
	}

	return nil
}

func listComponents(url, componentType string) error {
	fmt.Printf("Fetching available %s from ExtDist...\n", strings.ToLower(componentType))

	resp, err := httputil.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch %s page: %w", strings.ToLower(componentType), err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse %s page: %w", strings.ToLower(componentType), err)
	}

	componentMap := make(map[string][]string)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href := s.AttrOr("href", "")
		if strings.HasSuffix(href, ".tar.gz") {
			// Extract component name and version
			filename := strings.TrimSuffix(href, ".tar.gz")
			parts := strings.Split(filename, "-")
			if len(parts) >= 2 {
				// First part is component name, rest is version
				componentName := parts[0]
				version := strings.Join(parts[1:], "-")
				componentMap[componentName] = append(componentMap[componentName], version)
			}
		}
	})

	if len(componentMap) == 0 {
		fmt.Printf("No %s found.\n", strings.ToLower(componentType))
		return nil
	}

	// Sort component names
	var componentNames []string
	for name := range componentMap {
		componentNames = append(componentNames, name)
	}
	sort.Strings(componentNames)

	fmt.Printf("\nAvailable %s:\n", componentType)
	for _, name := range componentNames {
		versions := componentMap[name]
		sort.Strings(versions)
		fmt.Printf("- %s (versions: %s)\n", name, strings.Join(versions, ", "))
	}

	return nil
}
