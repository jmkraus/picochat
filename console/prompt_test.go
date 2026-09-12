package console

import "testing"

func TestNumberedPrompt(t *testing.T) {
	tests := []struct {
		name string
		slot int
		want string
	}{
		{name: "first slot", slot: 0, want: "[1]\u276F "},
		{name: "third slot", slot: 2, want: "[3]\u276F "},
		{name: "last slot", slot: 7, want: "[8]\u276F "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumberedPrompt(tt.slot); got != tt.want {
				t.Errorf("NumberedPrompt(%d) = %q, want %q", tt.slot, got, tt.want)
			}
		})
	}
}

func TestNumberedPromptPlain(t *testing.T) {
	if got, want := numberedPrompt(0), "[1]\u276F "; got != want {
		t.Errorf("numberedPrompt(0) = %q, want %q", got, want)
	}
}

func TestPromptWidth(t *testing.T) {
	if got, want := PromptWidth(), 5; got != want {
		t.Errorf("PromptWidth() = %d, want %d", got, want)
	}
}
