package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test.ini")

	configContent := `[mediawiki]
version=1.43.1

[extensions]
extdist=VisualEditor
extdist=WikiEditor|REL1_42

[skins]
extdist=Vector
git=https://github.com/example/skin.git|main
`

	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Load and test the configuration
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test MediaWiki section
	if config.MediaWiki.Version != "1.43.1" {
		t.Errorf("Expected MediaWiki version 1.43.1, got %s", config.MediaWiki.Version)
	}

	// Test Extensions section
	if len(config.Extensions) != 2 {
		t.Errorf("Expected 2 extensions, got %d", len(config.Extensions))
	}

	// Test first extension
	if config.Extensions[0].Distributor != "extdist" || config.Extensions[0].Name != "VisualEditor" || config.Extensions[0].Version != "" {
		t.Errorf("First extension not parsed correctly: %+v", config.Extensions[0])
	}

	// Test second extension with version
	if config.Extensions[1].Distributor != "extdist" || config.Extensions[1].Name != "WikiEditor" || config.Extensions[1].Version != "REL1_42" {
		t.Errorf("Second extension not parsed correctly: %+v", config.Extensions[1])
	}

	// Test Skins section
	if len(config.Skins) != 2 {
		t.Errorf("Expected 2 skins, got %d", len(config.Skins))
	}

	// Test extdist skin
	if config.Skins[0].Distributor != "extdist" || config.Skins[0].Name != "Vector" {
		t.Errorf("ExtDist skin not parsed correctly: %+v", config.Skins[0])
	}

	// Test git skin
	if config.Skins[1].Distributor != "git" || config.Skins[1].Name != "https://github.com/example/skin.git" || config.Skins[1].Version != "main" {
		t.Errorf("Git skin not parsed correctly: %+v", config.Skins[1])
	}
}

func TestLoadConfigFileNotExists(t *testing.T) {
	_, err := LoadConfig("nonexistent.ini")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}
