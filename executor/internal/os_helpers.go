package internal

import (
	"fmt"
	"os"
)

func IsFileExists(filename string) error {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return fmt.Errorf("file not exists: %s", filename)
	}
	if info.IsDir() {
		return fmt.Errorf("given directory path: %s", filename)
	}
	return nil
}
