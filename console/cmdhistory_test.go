package console

import "testing"

func TestCommandHistory_Internal(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		wantPrev []string
		wantNext []string
	}{
		{
			name:     "basic history navigation",
			commands: []string{"/help", "/models", "/info"},
			wantPrev: []string{"/info", "/models", "/help", "/help"}, // Prev stops at oldest
			wantNext: []string{"/models", "/info", ""},               // Next stops at newest
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &commandHistory{}
			for _, c := range tt.commands {
				h.add(c)
			}

			for i, expected := range tt.wantPrev {
				got := h.prev()
				if got != expected {
					t.Errorf("Prev[%d]: expected %q, got %q", i, expected, got)
				}
			}

			for i, expected := range tt.wantNext {
				got := h.next()
				if got != expected {
					t.Errorf("Next[%d]: expected %q, got %q", i, expected, got)
				}
			}
		})
	}
}

func TestCommandHistory_Global(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		wantPrev []string
		wantNext []string
	}{
		{
			name:     "global history behaves like instance history",
			commands: []string{"/start", "/status", "/stop"},
			wantPrev: []string{"/stop", "/status", "/start", "/start"},
			wantNext: []string{"/status", "/stop", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global history between test cases
			cmdHistory.entries = nil
			cmdHistory.index = 0

			for _, c := range tt.commands {
				AddCommand(c)
			}

			for i, expected := range tt.wantPrev {
				got := PrevCommand()
				if got != expected {
					t.Errorf("Prev[%d]: expected %q, got %q", i, expected, got)
				}
			}

			for i, expected := range tt.wantNext {
				got := NextCommand()
				if got != expected {
					t.Errorf("Next[%d]: expected %q, got %q", i, expected, got)
				}
			}
		})
	}
}
