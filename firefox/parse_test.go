package firefox

import (
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
		name := filepath.Base(profile)
		if _, err := ParseContainers(profile + "/containers.json"); err != nil {
			t.Errorf("%s/containers.json: %s", name, err)
		}
		if _, err := ParseExtensions(profile + "/extensions.json"); err != nil {
			t.Errorf("%s/extensions.json: %s", name, err)
		}
		if _, err := ParseExtensionPreferences(profile + "/extension-preferences.json"); err != nil {
			t.Errorf("%s/extension-preferences.json: %s", name, err)
		}
		if _, err := ParseExtensionSettings(profile + "/extension-settings.json"); err != nil {
			t.Errorf("%s/extension-settings.json: %s", name, err)
		}
		if _, err := ParseHandlers(profile + "/handlers.json"); err != nil {
			t.Errorf("%s/handlers.json: %s", name, err)
		}
		if _, err := ParseTimes(profile + "/times.json"); err != nil {
			t.Errorf("%s/times.json: %s", name, err)
		}
	}
}
