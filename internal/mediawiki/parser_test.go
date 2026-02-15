package mediawiki

import (
	"testing"
)

func TestExtractMajorMinor(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"1.43.1", "1.43", false},
		{"1.42.0", "1.42", false},
		{"2.0.1", "2.0", false},
		{"1.43", "1.43", false},
		{"invalid", "", true},
		{"1", "", true},
		{"", "", true},
	}

	for _, test := range tests {
		result, err := parser.extractMajorMinor(test.input)

		if test.hasError {
			if err == nil {
				t.Errorf("Expected error for input %s, but got none", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input %s: %v", test.input, err)
			}
			if result != test.expected {
				t.Errorf("Expected %s for input %s, got %s", test.expected, test.input, result)
			}
		}
	}
}

func TestNewParser(t *testing.T) {
	parser := NewParser()
	if parser == nil {
		t.Error("NewParser returned nil")
	}
}

// Note: GetDownloadURL requires network access, so we don't test it in unit tests
// It would require mocking the HTTP client or using integration tests
