package execpath

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// findExecutable searches for the executable in the system's PATH.
func FindExecutable(cmd string, path string) (string, error) { //Capitalized function name.
	paths := strings.Split(path, string(os.PathListSeparator))

	for _, dir := range paths {
		exePath := filepath.Join(dir, cmd)

		// Check if the file exists and is executable
		info, err := os.Stat(exePath)
		if err == nil && !info.IsDir() && (info.Mode()&0111 != 0) { // Check if executable bit is set
			return exePath, nil
		}
	}

	return "", fmt.Errorf("command not found: %s", cmd)
}
