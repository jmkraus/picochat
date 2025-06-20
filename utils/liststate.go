package utils

import (
	"fmt"
	"strings"
	"sync"
)

// State for selectable lists like models or history files.
type ListState struct {
	mu    sync.RWMutex
	items []string
}

// NewListState returns a new list state.
func NewListState() *ListState {
	return &ListState{
		items: []string{},
	}
}

// Set sets a new list of items and resets internal state.
func (ls *ListState) Set(list []string) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	ls.items = make([]string, len(list))
	copy(ls.items, list)
}

// GetItem returns the item by 1-based index (e.g. 1 for first item).
func (ls *ListState) GetItem(index int) (string, bool) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	if index <= 0 || index > len(ls.items) {
		return "", false
	}
	return ls.items[index-1], true
}

// Size returns number of stored items.
func (ls *ListState) Size() int {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	return len(ls.items)
}

// FormatForDisplay returns a numbered list with a heading, suitable for printing.
func (ls *ListState) FormatForDisplay(heading string) string {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	if len(ls.items) == 0 {
		return fmt.Sprintf("No %s found.", strings.ToLower(heading))
	}

	var lines []string
	for i, item := range ls.items {
		lines = append(lines, fmt.Sprintf("(%02d) %s", i+1, item))
	}
	return fmt.Sprintf("Available %s:\n%s", strings.ToLower(heading), strings.Join(lines, "\n"))
}
