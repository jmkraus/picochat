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
	if err.Error() != "no history files found" {
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

	if got["type"] != "object" {
		t.Fatalf("type = %v, want object", got["type"])
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

func TestGetMimeType_Image(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "image.jpg")
	if err := os.WriteFile(file, []byte("dummy"), 0644); err != nil {
		t.Fatalf("write image.jpg: %v", err)
	}

	mt, err := GetMimeType(file)
	if err != nil {
		t.Fatalf("GetMimeType returned error: %v", err)
	}
	if mt != "image/jpeg" {
		t.Fatalf("mime = %q, want %q", mt, "image/jpeg")
	}
}

func TestGetMimeType_NonImageReturnsError(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "note.txt")
	if err := os.WriteFile(file, []byte("hello"), 0644); err != nil {
		t.Fatalf("write note.txt: %v", err)
	}

	_, err := GetMimeType(file)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestStripDataURLPrefix(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "data url with base64 marker",
			in:   "data:image/png;base64,QUJDRA==",
			want: "QUJDRA==",
		},
		{
			name: "plain base64 unchanged",
			in:   "QUJDRA==",
			want: "QUJDRA==",
		},
		{
			name: "data url without base64 marker unchanged",
			in:   "data:image/png,QUJDRA==",
			want: "data:image/png,QUJDRA==",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StripDataURLPrefix(tt.in)
			if got != tt.want {
				t.Fatalf("StripDataURLPrefix(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
