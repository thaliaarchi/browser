package chrome

import (
	"os"
	"path/filepath"
	"time"
)

// GetFirstRun retrieves the time that Chrome was first ran from the
// "First Run" file in the Chrome root.
func GetFirstRun(chromeDir string) (time.Time, error) {
	fi, err := os.Stat(filepath.Join(chromeDir, "First Run"))
	if err != nil {
		return time.Time{}, err
	}
	return fi.ModTime(), nil
}
