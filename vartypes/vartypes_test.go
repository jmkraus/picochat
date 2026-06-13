package vartypes

import "testing"

func TestVarTypeString(t *testing.T) {
	tests := []struct {
		name string
		in   VarType
		want string
	}{
		{name: "float", in: VarFloat, want: "float"},
		{name: "int", in: VarInt, want: "int"},
		{name: "string", in: VarString, want: "string"},
		{name: "bool", in: VarBool, want: "bool"},
		{name: "unknown", in: VarType(99), want: "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.String(); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestConvertBoolAliases(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{in: "true", want: true},
		{in: "YES", want: true},
		{in: "1", want: true},
		{in: "false", want: false},
		{in: "No", want: false},
		{in: "0", want: false},
	}

	for _, tt := range tests {
		got, err := Convert(VarBool, tt.in)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tt.in, err)
		}
		val, ok := got.(bool)
		if !ok {
			t.Fatalf("expected bool type for %q, got %T", tt.in, got)
		}
		if val != tt.want {
			t.Fatalf("for %q expected %v, got %v", tt.in, tt.want, val)
		}
	}
}

func TestConvertUnsupportedTypeAndInvalidBool(t *testing.T) {
	if _, err := Convert(VarUnknown, "x"); err == nil {
		t.Fatal("expected unsupported type error, got nil")
	}

	if _, err := Convert(VarBool, "maybe"); err == nil {
		t.Fatal("expected invalid bool error, got nil")
	}
}

func TestConvertFloat(t *testing.T) {
	t.Run("valid float", func(t *testing.T) {
		got, err := Convert(VarFloat, "3.14")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		v, ok := got.(float64)
		if !ok {
			t.Fatalf("expected float64, got %T", got)
		}
		if v != 3.14 {
			t.Fatalf("expected 3.14, got %v", v)
		}
	})

	t.Run("invalid float", func(t *testing.T) {
		if _, err := Convert(VarFloat, "abc"); err == nil {
			t.Fatal("expected error for invalid float, got nil")
		}
	})
}

func TestConvertInt(t *testing.T) {
	t.Run("valid int", func(t *testing.T) {
		got, err := Convert(VarInt, "42")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		v, ok := got.(int)
		if !ok {
			t.Fatalf("expected int, got %T", got)
		}
		if v != 42 {
			t.Fatalf("expected 42, got %v", v)
		}
	})

	t.Run("invalid int", func(t *testing.T) {
		if _, err := Convert(VarInt, "4.2"); err == nil {
			t.Fatal("expected error for invalid int, got nil")
		}
	})
}

func TestConvertString(t *testing.T) {
	got, err := Convert(VarString, "HeLLo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := got.(string)
	if !ok {
		t.Fatalf("expected string, got %T", got)
	}
	if v != "HeLLo" {
		t.Fatalf("expected %q, got %q", "HeLLo", v)
	}
}
