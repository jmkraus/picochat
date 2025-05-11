package utils

import (
	"os"
	"path/filepath"
	"picochat/paths"
	"strings"
)

func ListHistoryFiles() ([]string, error) {
	basePath := filepath.Join(filepath.Dir(paths.GetConfigPath()), "history")
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
