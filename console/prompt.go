package console

import (
	"fmt"
	"strconv"

	"github.com/mattn/go-runewidth"
)

const (
	numPrompt string = "[%s]\u276F " // requires a parameter
)

// numberedPrompt returns the plain prompt without formatting.
func numberedPrompt(slot int) string {
	num := strconv.Itoa(slot + 1)
	return fmt.Sprintf(numPrompt, num)
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
	num := Colorize(Blue, strconv.Itoa(slot+1))
	return fmt.Sprintf(numPrompt, num)
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
