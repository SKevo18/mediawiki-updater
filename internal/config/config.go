package config

import (
	"fmt"
	"strings"
)

// Config represents the main configuration structure
type Config struct {
	MediaWiki  MediaWikiConfig
	Extensions []ComponentConfig
	Skins      []ComponentConfig
}

// MediaWikiConfig holds MediaWiki core configuration
type MediaWikiConfig struct {
	Version string `ini:"version"`
}

// ComponentConfig represents an extension or skin configuration
type ComponentConfig struct {
	Distributor string
	Name        string
	Version     string
}

// LoadConfig loads configuration from an INI file
func LoadConfig(configPath string) (*Config, error) {
	ini, err := LoadINIFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	config := &Config{}

	// Load MediaWiki section
	config.MediaWiki.Version = ini.GetFirstValue("mediawiki", "version")

	// Load Extensions section
	config.Extensions = parseComponentsFromINI(ini, "extensions")

	// Load Skins section
	config.Skins = parseComponentsFromINI(ini, "skins")

	return config, nil
}

// parseComponentsFromINI parses component configurations from a SimpleINI section
func parseComponentsFromINI(ini *SimpleINI, sectionName string) []ComponentConfig {
	var components []ComponentConfig

	section := ini.GetSection(sectionName)
	for distributor, values := range section {
		for _, value := range values {
			// Parse format: <extension>|<optional version>
			parts := strings.SplitN(value, "|", 2)
			name := parts[0]
			version := ""
			if len(parts) > 1 {
				version = parts[1]
			}

			components = append(components, ComponentConfig{
				Distributor: distributor,
				Name:        name,
				Version:     version,
			})
		}
	}

	return components
}
