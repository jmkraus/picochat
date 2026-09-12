package console

import (
	"fmt"

	"github.com/mattn/go-runewidth"
)

const (
	NumPrompt string = "[%d]\u276F " // requires a parameter
)

// numberedPrompt returns the plain prompt without formatting.
// (Kept for compatibility if esc sequences are introduced again.)
//
// Parameters:
//
//	slot (int) - the slot number
//
// Returns:
//
//	string - the unformatted prompt string
func numberedPrompt(slot int) string {
	return NumberedPrompt(slot)
}

// NumberedPrompt returns the prompt with formatting.
//
// Parameters:
//
//	slot (int) - the slot number
//
// Returns:
//
//	string - the formatted prompt string for CLI output
func NumberedPrompt(slot int) string {
	return fmt.Sprintf(NumPrompt, slot+1)
}

// PromptWidth returns the visible width of the numbered prompt.
//
// Parameters:
//
//	none
//
// Returns:
//
//	int - the width of the symbols
func PromptWidth() int {
	return runewidth.StringWidth(numberedPrompt(0))
}
