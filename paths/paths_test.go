package paths

import (
	"os"
	"path/filepath"
	"strings"
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
			expected:    homeDir,
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

func TestFileExists_PositiveAndDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "exists.txt")
	if err := os.WriteFile(filePath, []byte("ok"), 0644); err != nil {
		t.Fatalf("write test file failed: %v", err)
	}

	if !FileExists(filePath) {
		t.Fatalf("FileExists(%s) = false; want true", filePath)
	}
	if FileExists(tmpDir) {
		t.Fatalf("FileExists(%s) = true for directory; want false", tmpDir)
	}
}

func TestGetConfigPath_ExplicitPathPassThrough(t *testing.T) {
	got, err := GetConfigPath("/tmp/my-config.toml")
	if err != nil {
		t.Fatalf("GetConfigPath returned unexpected error: %v", err)
	}
	if got != "/tmp/my-config.toml" {
		t.Fatalf("GetConfigPath returned %q; want %q", got, "/tmp/my-config.toml")
	}
}

func TestGetConfigPath_AliasAndDefaultUseConfigPathEnv(t *testing.T) {
	cfgDir := t.TempDir()
	t.Setenv("CONFIG_PATH", cfgDir)
	t.Setenv("XDG_CONFIG_HOME", "")

	aliasPath, err := GetConfigPath("@custom")
	if err != nil {
		t.Fatalf("GetConfigPath(@custom) failed: %v", err)
	}
	wantAlias := filepath.Join(cfgDir, "custom.toml")
	if aliasPath != wantAlias {
		t.Fatalf("alias path = %q; want %q", aliasPath, wantAlias)
	}

	defaultPath, err := GetConfigPath("")
	if err != nil {
		t.Fatalf("GetConfigPath(\"\") failed: %v", err)
	}
	wantDefault := filepath.Join(cfgDir, "config.toml")
	if defaultPath != wantDefault {
		t.Fatalf("default path = %q; want %q", defaultPath, wantDefault)
	}
}

func TestOverrideHistoryPath_Restore(t *testing.T) {
	tmpDir := t.TempDir()
	restore := OverrideHistoryPath(tmpDir)
	t.Cleanup(restore)

	got, err := GetHistoryPath()
	if err != nil {
		t.Fatalf("GetHistoryPath with override failed: %v", err)
	}
	if got != tmpDir {
		t.Fatalf("GetHistoryPath with override = %q; want %q", got, tmpDir)
	}

	restore()

	t.Setenv("CONFIG_PATH", t.TempDir())
	gotAfter, err := GetHistoryPath()
	if err != nil {
		t.Fatalf("GetHistoryPath after restore failed: %v", err)
	}
	if !strings.HasSuffix(gotAfter, filepath.Join("history")) {
		t.Fatalf("GetHistoryPath after restore = %q; expected to end with /history", gotAfter)
	}
}

func TestFallbackToXDGOrHome(t *testing.T) {
	xdgRoot := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdgRoot)

	got, err := fallbackToXDGOrHome()
	if err != nil {
		t.Fatalf("fallbackToXDGOrHome with XDG_CONFIG_HOME failed: %v", err)
	}
	want := filepath.Join(xdgRoot, "picochat")
	if got != want {
		t.Fatalf("fallbackToXDGOrHome = %q; want %q", got, want)
	}

	t.Setenv("XDG_CONFIG_HOME", "")
	gotHome, err := fallbackToXDGOrHome()
	if err != nil {
		t.Fatalf("fallbackToXDGOrHome with HOME fallback failed: %v", err)
	}
	if !strings.HasSuffix(gotHome, filepath.Join(".config", "picochat")) {
		t.Fatalf("fallbackToXDGOrHome HOME result = %q; expected suffix .config/picochat", gotHome)
	}
}

func TestFallbackToExecutableDir(t *testing.T) {
	got, err := fallbackToExecutableDir()
	if err != nil {
		t.Fatalf("fallbackToExecutableDir failed: %v", err)
	}
	if got == "" {
		t.Fatal("fallbackToExecutableDir returned empty path")
	}
}
