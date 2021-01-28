package firefox

import (
	"encoding/json"
	"os"

	"github.com/andrewarchi/archive/timefmt"
)

// ExtensionSettings contains preferences and commands set by extensions
// in extension-settings.json.
type ExtensionSettings struct {
	Version              int                `json:"version"` // i.e. 2
	Commands             map[string]Command `json:"commands"`
	URLOverrides         interface{}        `json:"url_overrides"` // TODO unknown structure
	Prefs                map[string]Pref    `json:"prefs"`
	DefaultSearch        interface{}        `json:"default_search"`       // TODO unknown structure
	HomepageNotification interface{}        `json:"homepageNotification"` // TODO unknown structure
	TabHideNotification  interface{}        `json:"tabHideNotification"`  // TODO unknown structure
	NewTabNotification   interface{}        `json:"newTabNotification"`   // TODO unknown structure
}

// Command is a command with values set by extensions.
type Command struct {
	PrecedenceList []ExtensionSetting `json:"precedenceList"`
}

// Pref is a preference with values set by extensions and an initial
// value.
type Pref struct {
	InitialValue   interface{}        `json:"initialValue"`
	PrecedenceList []ExtensionSetting `json:"precedenceList"`
}

// ExtensionSetting is a value set by an extension.
type ExtensionSetting struct {
	ID          string            `json:"id"`
	InstallDate timefmt.UnixMilli `json:"installDate"`
	Value       interface{}       `json:"value"`
	Enabled     bool              `json:"enabled"`
}

// ExtensionPreferences lists additional permissions granted to an
// extension in extension-preferences.json.
type ExtensionPreferences struct {
	Permissions []string `json:"permissions"` // i.e. "clipboardWrite" or "internal:privateBrowsingAllowed"
	Origins     []string `json:"origins"`     // Origins given access to
}

// ParseExtensionSettings parses an extension-settings.json file.
func ParseExtensionSettings(filename string) (*ExtensionSettings, error) {
	var settings ExtensionSettings
	if err := parseJSON(filename, &settings); err != nil {
		return nil, err
	}
	return &settings, nil
}

// ParseExtensionPreferences parses an extension-preferences.json file.
func ParseExtensionPreferences(filename string) (map[string]ExtensionPreferences, error) {
	var prefs map[string]ExtensionPreferences
	if err := parseJSON(filename, &prefs); err != nil {
		return nil, err
	}
	return prefs, nil
}

func parseJSON(filename string, data interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	d := json.NewDecoder(f)
	d.DisallowUnknownFields()
	return d.Decode(data)
}
