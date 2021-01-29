package firefox

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	profilesDir, err := ProfilesDir()
	if err != nil {
		t.Fatal(err)
	}
	profiles, err := filepath.Glob(filepath.Join(profilesDir, "*"))
	if err != nil {
		t.Fatal(err)
	}
	for _, profile := range profiles {
		if fi, err := os.Stat(profile); err == nil && !fi.IsDir() {
			continue
		}

		addons := filepath.Join(profile, "addons.json")
		_, err = ParseAddons(addons)
		checkError(t, addons, err)

		containers := filepath.Join(profile, "containers.json")
		_, err = ParseContainers(containers)
		checkError(t, containers, err)

		extensions := filepath.Join(profile, "extensions.json")
		_, err = ParseExtensions(extensions)
		checkError(t, extensions, err)

		extensionPreferences := filepath.Join(profile, "extension-preferences.json")
		_, err = ParseExtensionPreferences(extensionPreferences)
		checkError(t, extensionPreferences, err)

		extensionSettings := filepath.Join(profile, "extension-settings.json")
		_, err = ParseExtensionSettings(extensionSettings)
		checkError(t, extensionSettings, err)

		handlers := filepath.Join(profile, "handlers.json")
		_, err = ParseHandlers(handlers)
		checkError(t, handlers, err)

		times := filepath.Join(profile, "times.json")
		_, err = ParseTimes(times)
		checkError(t, times, err)

		bookmarkBackups, err := filepath.Glob(filepath.Join(profile, "bookmarkbackups", "bookmarks-*.jsonlz4"))
		if err != nil {
			t.Error(err)
			continue
		}
		for _, bookmarkBackup := range bookmarkBackups {
			_, err = ParseBookmarkBackup(bookmarkBackup)
			checkError(t, bookmarkBackup, err)
		}
	}
}

func checkError(t *testing.T, filename string, err error) {
	t.Helper()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s: %s", filename, err)
	}
}
