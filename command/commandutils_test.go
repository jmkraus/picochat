package command

import (
	"testing"
)

func TestExtractCodeBlock(t *testing.T) {
	text := "Some explanation.\n```go\nfmt.Println(\"hi\")\n```"
	code, found := extractCodeBlock(text)
	if code != "fmt.Println(\"hi\")\n" {
		t.Errorf("Unexpected extracted code block: %q", code)
	}
	if code == "fmt.Println(\"hi\")\n" && found == false {
		t.Errorf("'Found' flag for ExtractCode reported False, but should be True")
	}
}

func TestExtractCodeBlock_Empty(t *testing.T) {
	text := "No code block here"
	code, found := extractCodeBlock(text)
	if code != "" {
		t.Errorf("Expected empty string, got %q", code)
	}
	if code == "" && found == true {
		t.Errorf("'Found' flag for ExtractCode reported True, but should be False")
	}
}
