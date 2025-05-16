package utils

import (
	"os"
	"picochat/paths"
	"strings"
)

func ListHistoryFiles() ([]string, error) {
	basePath, err := paths.GetHistoryDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entry := range entries {
		if !entry.IsDir() && HasSuffix(entry.Name()) {
			result = append(result, entry.Name())
		}
	}
	return result, nil
}

const fileSuffix = ".chat"

func HasSuffix(filename string) bool {
	return strings.HasSuffix(filename, fileSuffix)
}

func EnsureSuffix(filename string) string {
	if !HasSuffix(filename) {
		return filename + fileSuffix
	}
	return filename
}
