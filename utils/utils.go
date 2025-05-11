package utils

import (
	"os"
	"picochat/paths"
	"strings"
)

func ListHistoryFiles() ([]string, error) {
	basePath := paths.GetHistoryDir()
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".chat") {
			result = append(result, entry.Name())
		}
	}
	return result, nil
}

const fileSuffix = ".chat"

func EnsureSuffix(filename string) string {
	if !strings.HasSuffix(filename, fileSuffix) {
		return filename + fileSuffix
	}
	return filename
}
