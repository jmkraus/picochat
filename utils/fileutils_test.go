package utils

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"picochat/paths"
	"testing"
)

func TestListHistoryFiles(t *testing.T) {
	tmpDir := t.TempDir()
	restore := paths.OverrideHistoryPath(tmpDir)
	defer restore()

	if err := os.WriteFile(filepath.Join(tmpDir, "a.chat"), []byte("a"), 0644); err != nil {
		t.Fatalf("write a.chat: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "b.txt"), []byte("b"), 0644); err != nil {
		t.Fatalf("write b.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "z.chat"), []byte("z"), 0644); err != nil {
		t.Fatalf("write z.chat: %v", err)
	}
	if err := os.Mkdir(filepath.Join(tmpDir, "dir.chat"), 0755); err != nil {
		t.Fatalf("mkdir dir.chat: %v", err)
	}

	got, err := ListHistoryFiles()
	if err != nil {
		t.Fatalf("ListHistoryFiles returned error: %v", err)
	}

	want := "History files:\n(01) a.chat\n(02) z.chat"
	if got != want {
		t.Fatalf("unexpected output\nwant: %q\ngot:  %q", want, got)
	}

	name, ok := GetHistoryByIndex(1)
	if !ok || name != "a.chat" {
		t.Fatalf("GetHistoryByIndex(1) = (%q, %v), want (a.chat, true)", name, ok)
	}
}

func TestListHistoryFiles_NoHistoryFiles(t *testing.T) {
	tmpDir := t.TempDir()
	restore := paths.OverrideHistoryPath(tmpDir)
	defer restore()

	_, err := ListHistoryFiles()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "no history files found." {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestImageToBase64(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "image.bin")
	content := []byte{0, 1, 2, 3, 255}
	if err := os.WriteFile(file, content, 0644); err != nil {
		t.Fatalf("write image.bin: %v", err)
	}

	got, err := ImageToBase64(file)
	if err != nil {
		t.Fatalf("ImageToBase64 returned error: %v", err)
	}

	want := base64.StdEncoding.EncodeToString(content)
	if got != want {
		t.Fatalf("unexpected base64\nwant: %q\ngot:  %q", want, got)
	}
}

func TestLoadSchemaFromFile_JSONLiteral(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "schema.txt")
	if err := os.WriteFile(file, []byte("json\n"), 0644); err != nil {
		t.Fatalf("write schema.txt: %v", err)
	}

	got, err := LoadSchemaFromFile(file)
	if err != nil {
		t.Fatalf("LoadSchemaFromFile returned error: %v", err)
	}

	s, ok := got.(string)
	if !ok || s != "json" {
		t.Fatalf("unexpected value: %#v", got)
	}
}

func TestLoadSchemaFromFile_Object(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "schema.json")
	if err := os.WriteFile(file, []byte(`{"type":"object","required":["name"]}`), 0644); err != nil {
		t.Fatalf("write schema.json: %v", err)
	}

	got, err := LoadSchemaFromFile(file)
	if err != nil {
		t.Fatalf("LoadSchemaFromFile returned error: %v", err)
	}

	obj, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", got)
	}
	if obj["type"] != "object" {
		t.Fatalf("type = %v, want object", obj["type"])
	}
}

func TestLoadSchemaFromFile_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "schema.json")
	if err := os.WriteFile(file, []byte(`{"type":"object"`), 0644); err != nil {
		t.Fatalf("write schema.json: %v", err)
	}

	_, err := LoadSchemaFromFile(file)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
