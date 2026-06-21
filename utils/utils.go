package utils

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// ListAvailableModels lists all models via /tags API call
//
// Parameters:
//
//	baseUrl (string) - Base URL for the API
//
// Returns:
//
//	string - formatted list of available models
//	error
func ListAvailableModels(models []string) (string, error) {
	if len(models) == 0 {
		return "", fmt.Errorf("no models available")
	}

	// Sort case-insensitively
	sort.Slice(models, func(i, j int) bool {
		return strings.ToLower(models[i]) < strings.ToLower(models[j])
	})

	setModelsList(models)
	return FormatList(models, "Language models", true), nil
}

// FormatList formats a list of items with an optional heading and numbering
//
// Parameters:
//
//	content ([]string) - Items to format
//	heading (string)   - Heading for the list
//	numbered (bool)    - Whether to number the items
//
// Returns:
//
//	string - formatted list
func FormatList(content []string, heading string, numbered bool) string {
	if len(content) == 0 {
		return fmt.Sprintf("No %s found.", strings.ToLower(heading))
	}

	var lines []string
	for i, item := range content {
		if numbered {
			lines = append(lines, fmt.Sprintf("(%02d) %s", i+1, item))
		} else {
			lines = append(lines, " - "+item)
		}
	}

	return fmt.Sprintf("%s:\n%s", capitalize(heading), strings.Join(lines, "\n"))
}

// MarkdownTable creates a markdown table from a 2-dim array and adjusts
// the column widths according to content.
//
// Parameters:
//
//	tableData ([][]string) - rows and columns of the table
//
// Returns:
//
//	string - markdown table
func MarkdownTable(tableData [][]string) string {
	if len(tableData) == 0 || len(tableData[0]) == 0 {
		return ""
	}

	numColumns := len(tableData[0])
	separator := make([]string, numColumns)
	for i := range separator {
		separator[i] = "---"
	}

	rows := make([][]string, 0, len(tableData)+1)
	rows = append(rows, tableData[0], separator)
	rows = append(rows, tableData[1:]...)

	maxWidths := make([]int, numColumns)
	for _, row := range rows {
		for colIdx := range numColumns {
			col := ""
			if colIdx < len(row) {
				col = row[colIdx]
			}
			if len(col) > maxWidths[colIdx] {
				maxWidths[colIdx] = len(col)
			}
		}
	}

	pad := func(s string, width int, fill byte) string {
		if len(s) >= width {
			return s
		}
		return s + strings.Repeat(string(fill), width-len(s))
	}

	var builder strings.Builder
	for rowIdx, row := range rows {
		fill := byte(' ')
		if rowIdx == 1 {
			fill = '-'
		}

		builder.WriteByte('|')
		for colIdx := range numColumns {
			col := ""
			if colIdx < len(row) {
				col = row[colIdx]
			}
			builder.WriteByte(' ')
			builder.WriteString(pad(col, maxWidths[colIdx], fill))
			builder.WriteString(" |")
		}
		if rowIdx < len(rows)-1 {
			builder.WriteByte('\n')
		}
	}

	return builder.String()
}

// capitalize capitalizes the first letter of a string
//
// Parameters:
//
//	s (string) - Input string
//
// Returns:
//
//	string - String with first letter capitalized
func capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(strings.ToLower(s))
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
