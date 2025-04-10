package files

import (
	"os"
	"path/filepath"
	"strings"
)

func GetFilesWithNameContaining(dir string, substr string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var matchingFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.Contains(entry.Name(), substr) {
			fullPath := filepath.Join(dir, entry.Name())
			matchingFiles = append(matchingFiles, fullPath)
		}
	}

	return matchingFiles, nil
}
