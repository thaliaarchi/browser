package firefox

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-ini/ini"
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

// ProfileInfo contains Firefox profiles and installs.
type ProfileInfo struct {
	StartWithLastProfile bool
	Version              int // i.e. 2
	Profiles             []Profile
	Installs             []Install
}

// Profile is a Firefox profile.
type Profile struct {
	ID         int    // sequential (i.e. 0, 1, 2)
	Name       string // i.e. "default", "default-release", "dev-edition-default"
	IsRelative bool
	Path       string
	Default    bool
}

// Install is a Firefox installation.
type Install struct {
	ID      uint64 // displayed in uppercase hex
	Default string // default profile path
	Locked  bool
}

// ParseProfiles parses the profiles.ini in the Firefox root.
func ParseProfiles(firefoxDir string) (*ProfileInfo, error) {
	filename := filepath.Join(firefoxDir, "profiles.ini")
	f, err := ini.Load(filename)
	if err != nil {
		return nil, err
	}
	var info ProfileInfo

	// Sections are listed in reverse order
	sections := f.Sections()
	for i := len(sections) - 1; i >= 0; i-- {
		section := sections[i]
		name := section.Name()
		keys := section.Keys()

		switch {
		case name == "DEFAULT":
			if len(keys) != 0 {
				return nil, fmt.Errorf("firefox: root section has bare keys: %s", filename)
			}
			continue
		case name == "General":
			if err := parseGeneralSection(keys, &info); err != nil {
				return nil, err
			}
		case strings.HasPrefix(name, "Profile"):
			profile, err := parseProfileSection(name, keys)
			if err != nil {
				return nil, err
			}
			info.Profiles = append(info.Profiles, *profile)
		case strings.HasPrefix(name, "Install"):
			install, err := parseInstallSection(name, keys)
			if err != nil {
				return nil, err
			}
			info.Installs = append(info.Installs, *install)
		default:
			return nil, fmt.Errorf("firefox: unknown section: %s", name)
		}
	}
	return &info, nil
}

func parseGeneralSection(keys []*ini.Key, info *ProfileInfo) error {
	for _, key := range keys {
		switch keyName := key.Name(); keyName {
		case "StartWithLastProfile":
			s, err := key.Bool()
			if err != nil {
				return err
			}
			info.StartWithLastProfile = s
		case "Version":
			v, err := key.Int()
			if err != nil {
				return err
			}
			info.Version = v
		default:
			return fmt.Errorf("firefox: unknown key in General: %s", keyName)
		}
	}
	return nil
}

func parseProfileSection(name string, keys []*ini.Key) (*Profile, error) {
	var profile Profile
	id, err := strconv.Atoi(strings.TrimPrefix(name, "Profile"))
	if err != nil {
		return nil, fmt.Errorf("firefox: profile ID: %w", err)
	}
	profile.ID = id

	for _, key := range keys {
		switch keyName := key.Name(); keyName {
		case "Name":
			profile.Name = key.String()
		case "IsRelative":
			r, err := key.Bool()
			if err != nil {
				return nil, err
			}
			profile.IsRelative = r
		case "Path":
			profile.Path = key.String()
		case "Default":
			d, err := key.Bool()
			if err != nil {
				return nil, err
			}
			profile.Default = d
		default:
			return nil, fmt.Errorf("firefox: unknown key in %s: %s", name, keyName)
		}
	}
	return &profile, nil
}

func parseInstallSection(name string, keys []*ini.Key) (*Install, error) {
	var install Install
	id := strings.TrimPrefix(name, "Install")
	if len(id) != 16 {
		return nil, fmt.Errorf("firefox: install ID not 8 bytes: %s", name)
	}
	b, err := hex.DecodeString(id)
	if err != nil {
		return nil, err
	}
	install.ID = binary.BigEndian.Uint64(b)

	for _, key := range keys {
		switch keyName := key.Name(); keyName {
		case "Default":
			install.Default = key.String()
		case "Locked":
			l, err := key.Bool()
			if err != nil {
				return nil, err
			}
			install.Locked = l
		}
	}
	return &install, nil
}
