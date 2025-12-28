package paths

import (
	"os"
	"testing"
)

func TestEnsureSuffix(t *testing.T) {
	if EnsureSuffix("foo", ".chat") != "foo.chat" {
		t.Errorf("Expected suffix to be appended")
	}
	if EnsureSuffix("bar.chat", ".chat") != "bar.chat" {
		t.Errorf("Suffix should not be duplicated")
	}
}

func TestExpandHomeDir(t *testing.T) {
	// Get the home directory to use in tests
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Could not get user home directory: %v", err)
	}

	// Test cases
	testCases := []struct {
		path        string
		expected    string
		description string
	}{
		{
			path:        "~/Documents/file.txt",
			expected:    homeDir + "/Documents/file.txt",
			description: "Expand home directory with file path",
		},
		{
			path:        "~/file.txt",
			expected:    homeDir + "/file.txt",
			description: "Expand home directory with file in root",
		},
		{
			path:        "/absolute/path/file.txt",
			expected:    "/absolute/path/file.txt",
			description: "Absolute path should remain unchanged",
		},
		{
			path:        "relative/path/file.txt",
			expected:    "relative/path/file.txt",
			description: "Relative path should remain unchanged",
		},
		{
			path:        "~/",
			expected:    homeDir + "/",
			description: "Expand home directory only",
		},
		{
			path:        "",
			expected:    "",
			description: "Empty string should remain unchanged",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result, err := ExpandHomeDir(tc.path)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestFileExists_Negative(t *testing.T) {
	// Define a path that does not exist
	nonExistentFilePath := "non_existent_file.txt"

	// Call the FileExists function
	result := FileExists(nonExistentFilePath)

	// Assert that the result is false
	if result {
		t.Errorf("FileExists(%s) = true; want false", nonExistentFilePath)
	}
}
