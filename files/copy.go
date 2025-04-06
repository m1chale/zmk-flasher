package files

import (
	"io"
	"os"
	"strings"
)

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	if syncErr := destinationFile.Sync(); syncErr != nil {
		// ignore "device not configured" errors (MCU likely rebooted)
		if strings.Contains(syncErr.Error(), "device not configured") ||
			strings.Contains(syncErr.Error(), "input/output error") {
			return nil
		}
		return syncErr
	}

	return nil
}
