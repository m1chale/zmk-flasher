package files

import (
	"fmt"
	"os"
)


func EnsureDeleted(path string) error {
	err := os.RemoveAll(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete %s: %w", path, err)
	}
	return nil
}
