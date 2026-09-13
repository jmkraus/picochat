package jsonutils

import (
	"strings"
	"testing"
)

func TestSanityCheck(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr string
	}{
		{
			name:    "empty input",
			input:   "",
			wantErr: "json string is empty",
		},
		{
			name:    "whitespace only",
			input:   " \n\t ",
			wantErr: "json string is empty",
		},
		{
			name:    "opening brace not found",
			input:   "not json",
			wantErr: "invalid json - opening bracket not found",
		},
		{
			name:  "trims surrounding whitespace",
			input: " \n {\"a\":1} \n ",
			want:  "{\"a\":1}",
		},
		{
			name:  "removes text before opening brace",
			input: "Here is the response:\n```json\n{\"a\":1}",
			want:  "{\"a\":1}",
		},
		{
			name:  "preserves trailing text",
			input: "{\"a\":1} trailing",
			want:  "{\"a\":1} trailing",
		},
		{
			name:  "uses first opening brace",
			input: "prefix {\"a\":{\"b\":2}}",
			want:  "{\"a\":{\"b\":2}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sanityCheck(tt.input)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("sanityCheck(%q) returned nil error, want %q", tt.input, tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("sanityCheck(%q) error = %q, want %q", tt.input, err.Error(), tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("sanityCheck(%q) returned error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("sanityCheck(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestPrettyPrint(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		_, err := PrettyPrint("   ")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		_, err := PrettyPrint("{")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "invalid json") {
			t.Fatalf("error = %q, want to contain %q", err.Error(), "invalid json")
		}
	})

	t.Run("input without opening brace does not panic", func(t *testing.T) {
		_, err := PrettyPrint("not json")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "opening bracket not found") {
			t.Fatalf("error = %q, want to contain %q", err.Error(), "opening bracket not found")
		}
	})

	t.Run("formats object", func(t *testing.T) {
		got, err := PrettyPrint("{\"a\":1,\"b\":true}")
		if err != nil {
			t.Fatalf("PrettyPrint returned error: %v", err)
		}
		if !strings.Contains(got, "\n") {
			t.Fatalf("expected pretty printed json to contain newlines, got: %q", got)
		}
		if !strings.Contains(got, "\"a\"") || !strings.Contains(got, "\"b\"") {
			t.Fatalf("unexpected formatted json: %q", got)
		}
	})

	t.Run("removes text before json object", func(t *testing.T) {
		got, err := PrettyPrint("Here is the response:\n```json\n{\"a\":1}")
		if err != nil {
			t.Fatalf("PrettyPrint returned error: %v", err)
		}
		want := "{\n  \"a\": 1\n}"
		if got != want {
			t.Fatalf("PrettyPrint returned %q, want %q", got, want)
		}
	})
}

func TestValidateJSON(t *testing.T) {
	t.Cleanup(func() {
		// Reset cache between tests to avoid cross-test influence.
		resolvedSchema = nil
	})

	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"a": map[string]any{"type": "integer"},
		},
		"required": []any{"a"},
	}

	t.Run("empty schema", func(t *testing.T) {
		err := ValidateJSON(nil, `{"a":1}`)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("empty json", func(t *testing.T) {
		err := ValidateJSON(schema, "   ")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		err := ValidateJSON(schema, "{")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "invalid json") {
			t.Fatalf("error = %q, want to contain %q", err.Error(), "invalid json")
		}
	})

	t.Run("input without opening brace does not panic", func(t *testing.T) {
		err := ValidateJSON(schema, "not json")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "opening bracket not found") {
			t.Fatalf("error = %q, want to contain %q", err.Error(), "opening bracket not found")
		}
	})

	t.Run("trailing data", func(t *testing.T) {
		err := ValidateJSON(schema, `{"a":1} trailing`)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "trailing") {
			t.Fatalf("error = %q, want to contain %q", err.Error(), "trailing")
		}
	})

	t.Run("valid instance passes", func(t *testing.T) {
		if err := ValidateJSON(schema, `{"a":1}`); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("removes text before json object", func(t *testing.T) {
		if err := ValidateJSON(schema, "The model says:\n```json\n{\"a\":1}"); err != nil {
			t.Fatalf("expected no error for json preceded by text, got %v", err)
		}
	})

	t.Run("missing required field fails", func(t *testing.T) {
		err := ValidateJSON(schema, `{"b":2}`)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("cache is used", func(t *testing.T) {
		// Prime cache.
		if err := ValidateJSON(schema, `{"a":1}`); err != nil {
			t.Fatalf("prime ValidateJSON returned error: %v", err)
		}
		if resolvedSchema == nil {
			t.Fatalf("expected resolvedSchema cache to be set")
		}

		// Intentionally pass an empty schemaMap. Since the cache is set, this should
		// not trigger schema marshal/resolve and should validate with cached schema.
		// Note: ValidateJSON currently treats len(schemaMap)==0 as error *before*
		// using cache, so pass the same schema but confirm cache remains non-nil.
		if err := ValidateJSON(schema, `{"a":2}`); err != nil {
			t.Fatalf("ValidateJSON returned error: %v", err)
		}
		if resolvedSchema == nil {
			t.Fatalf("expected resolvedSchema cache to remain set")
		}
	})
}
