package firefox

import (
	"fmt"
	"os"
	"runtime"
)

// ProfilesDir returns the path for the Firefox Profiles directory.
func ProfilesDir() (string, error) {
	// MozillaZine lists paths for Windows 2000 and XP and an alternate
	// macOS path, but Go is unlikely to run on those systems.
	//
	// http://kb.mozillazine.org/Profile_folder_-_Firefox#Navigating_to_the_profile_folder

	if runtime.GOOS == "windows" {
		if appdata := os.Getenv("AppData"); appdata != "" {
			return appdata + `\Mozilla\Firefox\Profiles`, nil
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	switch runtime.GOOS {
	case "windows":
		return home + `\AppData\Roaming\Mozilla\Firefox\Profiles`, nil
	case "darwin":
		return home + "/Library/Application Support/Firefox/Profiles", nil
	case "linux":
		return home + "/.mozilla/firefox", nil
	default:
		return "", fmt.Errorf("firefox: unsupported GOOS: %s", runtime.GOOS)
	}
}
