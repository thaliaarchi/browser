// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package firefox parses many formats of Firefox data.
package firefox

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/andrewarchi/browser/iniutil"
	"gopkg.in/ini.v1"
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
	Version              int       // e.g. 2
	Profiles             []Profile `ini:"-"`
	Installs             []Install `ini:"-"`
}

// Profile is a Firefox profile.
type Profile struct {
	ID         int    `ini:"-"` // sequential (e.g. 0, 1, 2)
	Name       string // e.g. "default", "default-release", "dev-edition-default"
	IsRelative bool
	Path       string
	Default    bool
}

// Install is a Firefox installation.
type Install struct {
	ID      uint64 `ini:"-"` // displayed in uppercase hex
	Default string // default profile path
	Locked  bool
}

// AbsPath returns the absolute to the profile, relative to the firefox
// root.
func (p *Profile) AbsPath(firefoxDir string) string {
	if p.IsRelative {
		return filepath.Clean(filepath.Join(firefoxDir, p.Path))
	}
	return p.Path
}

// ParseProfiles parses profiles.ini in the Firefox root.
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
		switch {
		case name == "DEFAULT":
			if len(section.KeyStrings()) != 0 {
				return nil, errors.New("firefox: root section has bare keys")
			}
		case name == "General":
			if err := iniutil.Decode(section, &info); err != nil {
				return nil, err
			}
		case strings.HasPrefix(name, "Profile"):
			profile, err := parseProfile(section, strings.TrimPrefix(name, "Profile"))
			if err != nil {
				return nil, err
			}
			info.Profiles = append(info.Profiles, *profile)
		case strings.HasPrefix(name, "Install"):
			install, err := parseInstall(section, strings.TrimPrefix(name, "Install"))
			if err != nil {
				return nil, err
			}
			info.Installs = append(info.Installs, *install)
		default:
			return nil, fmt.Errorf("firefox: unknown section: %q", name)
		}
	}
	return &info, nil
}

// ParseInstalls parses installs.ini in the Firefox root.
func ParseInstalls(firefoxDir string) ([]Install, error) {
	filename := filepath.Join(firefoxDir, "installs.ini")
	f, err := ini.Load(filename)
	if err != nil {
		return nil, err
	}
	var installs []Install

	// Sections are listed in reverse order
	sections := f.Sections()
	for i := len(sections) - 1; i >= 0; i-- {
		section := sections[i]
		name := section.Name()
		if name == "DEFAULT" {
			if len(section.KeyStrings()) != 0 {
				return nil, errors.New("firefox: root section has bare keys")
			}
		} else {
			install, err := parseInstall(section, name)
			if err != nil {
				return nil, err
			}
			installs = append(installs, *install)
		}
	}
	return installs, nil
}

func parseProfile(section *ini.Section, id string) (*Profile, error) {
	var profile Profile
	n, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("firefox: invalid profile ID: %w", err)
	}
	profile.ID = n

	if err := iniutil.Decode(section, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

func parseInstall(section *ini.Section, id string) (*Install, error) {
	var install Install
	if len(id) != 16 {
		return nil, fmt.Errorf("firefox: install ID not 8 bytes: %q", section.Name())
	}
	b, err := hex.DecodeString(id)
	if err != nil {
		return nil, err
	}
	install.ID = binary.BigEndian.Uint64(b)

	if err := iniutil.Decode(section, &install); err != nil {
		return nil, err
	}
	return &install, nil
}
