package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// SimpleINI represents a simple INI file structure that supports duplicate keys
type SimpleINI struct {
	sections map[string]map[string][]string
}

// NewSimpleINI creates a new SimpleINI instance
func NewSimpleINI() *SimpleINI {
	return &SimpleINI{
		sections: make(map[string]map[string][]string),
	}
}

// LoadINIFile loads an INI file and returns a SimpleINI instance
func LoadINIFile(filename string) (*SimpleINI, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	ini := NewSimpleINI()
	scanner := bufio.NewScanner(file)

	currentSection := ""
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for section headers
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.TrimSpace(line[1 : len(line)-1])
			if ini.sections[currentSection] == nil {
				ini.sections[currentSection] = make(map[string][]string)
			}
			continue
		}

		// Parse key-value pairs
		if currentSection == "" {
			return nil, fmt.Errorf("key-value pair found outside of section at line %d", lineNumber)
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line format at line %d: %s", lineNumber, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Add value to the slice for this key
		ini.sections[currentSection][key] = append(ini.sections[currentSection][key], value)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return ini, nil
}

// GetSection returns all key-value pairs for a given section
func (ini *SimpleINI) GetSection(sectionName string) map[string][]string {
	if section, exists := ini.sections[sectionName]; exists {
		return section
	}
	return make(map[string][]string)
}

// GetSectionKeys returns all keys in a section
func (ini *SimpleINI) GetSectionKeys(sectionName string) []string {
	section := ini.GetSection(sectionName)
	keys := make([]string, 0, len(section))
	for key := range section {
		keys = append(keys, key)
	}
	return keys
}

// GetValues returns all values for a specific key in a section
func (ini *SimpleINI) GetValues(sectionName, key string) []string {
	section := ini.GetSection(sectionName)
	if values, exists := section[key]; exists {
		return values
	}
	return []string{}
}

// GetFirstValue returns the first value for a specific key in a section
func (ini *SimpleINI) GetFirstValue(sectionName, key string) string {
	values := ini.GetValues(sectionName, key)
	if len(values) > 0 {
		return values[0]
	}
	return ""
}
