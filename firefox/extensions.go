package firefox

import (
	"github.com/andrewarchi/browser/jsonutil"
	"github.com/andrewarchi/browser/jsonutil/timefmt"
)

// ExtensionSettings contains preferences and commands set by extensions
// in extension-settings.json.
type ExtensionSettings struct {
	Version              int                 `json:"version"` // i.e. 2
	Commands             map[string]Command  `json:"commands"`
	URLOverrides         jsonutil.UnknownObj `json:"url_overrides"`
	Prefs                map[string]Pref     `json:"prefs"`
	DefaultSearch        jsonutil.UnknownObj `json:"default_search"`
	HomepageNotification jsonutil.UnknownObj `json:"homepageNotification"`
	TabHideNotification  jsonutil.UnknownObj `json:"tabHideNotification"`
	NewTabNotification   jsonutil.UnknownObj `json:"newTabNotification"`
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

// ParseExtensionSettings parses extension-settings.json in a Firefox
// profile.
func ParseExtensionSettings(filename string) (*ExtensionSettings, error) {
	var settings ExtensionSettings
	if err := jsonutil.DecodeFile(filename, &settings); err != nil {
		return nil, err
	}
	return &settings, nil
}

// ExtensionPermissions lists additional permissions granted to an
// extension in extension-preferences.json.
type ExtensionPermissions struct {
	Permissions []string `json:"permissions"` // i.e. "clipboardWrite" or "internal:privateBrowsingAllowed"
	Origins     []string `json:"origins"`     // Origins given access to
}

// ParseExtensionPreferences parses extension-preferences.json in a
// Firefox profile.
func ParseExtensionPreferences(filename string) (map[string]ExtensionPermissions, error) {
	var prefs map[string]ExtensionPermissions
	if err := jsonutil.DecodeFile(filename, &prefs); err != nil {
		return nil, err
	}
	return prefs, nil
}

type Extensions struct {
	SchemaVersion int     `json:"schemaVersion"` // i.e. 33
	Addons        []Addon `json:"addons"`
}

type Addon struct {
	ID                     string                 `json:"id"` // ID or GUID
	SyncGUID               string                 `json:"syncGUID"`
	Version                string                 `json:"version"` // Addon version
	Type                   string                 `json:"type"`    // "extension", "theme", "locale", "dictionary"
	Loader                 jsonutil.UnknownType   `json:"loader"`
	UpdateURL              string                 `json:"updateURL"`
	OptionsURL             string                 `json:"optionsURL"`
	OptionsType            int                    `json:"optionsType"`
	OptionsBrowserStyle    bool                   `json:"optionsBrowserStyle"`
	AboutURL               string                 `json:"aboutURL"`
	DefaultLocale          Locale                 `json:"defaultLocale"`
	Visible                bool                   `json:"visible"`
	Active                 bool                   `json:"active"`
	UserDisabled           bool                   `json:"userDisabled"`
	AppDisabled            bool                   `json:"appDisabled"`
	EmbedderDisabled       bool                   `json:"embedderDisabled"`
	InstallDate            int64                  `json:"installDate"`
	UpdateDate             timefmt.UnixMilli      `json:"updateDate,omitempty"`
	ApplyBackgroundUpdates interface{}            `json:"applyBackgroundUpdates"` // i.e. 1 or "1"
	Path                   string                 `json:"path"`
	Skinnable              bool                   `json:"skinnable"`
	SourceURI              string                 `json:"sourceURI"`
	ReleaseNotesURI        string                 `json:"releaseNotesURI"`
	SoftDisabled           bool                   `json:"softDisabled"`
	ForeignInstall         bool                   `json:"foreignInstall"`
	StrictCompatibility    bool                   `json:"strictCompatibility"`
	Locales                []Locale               `json:"locales"`
	TargetApplications     []TargetApplication    `json:"targetApplications"`
	TargetPlatforms        []jsonutil.UnknownType `json:"targetPlatforms"`
	SignedState            int                    `json:"signedState,omitempty"` // i.e. 2
	SignedDate             timefmt.UnixMilli      `json:"signedDate"`
	Seen                   bool                   `json:"seen"`
	Dependencies           []interface{}          `json:"dependencies"`
	Incognito              string                 `json:"incognito,omitempty"` // i.e. "not_allowed", "spanning"
	UserPermissions        *ExtensionPermissions  `json:"userPermissions"`
	OptionalPermissions    *ExtensionPermissions  `json:"optionalPermissions"`
	Icons                  map[int]string         `json:"icons"` // Sized icon paths
	IconURL                string                 `json:"iconURL"`
	BlocklistState         int                    `json:"blocklistState"` // i.e. 2
	BlocklistURL           string                 `json:"blocklistURL"`
	StartupData            *StartupData           `json:"startupData"`
	Hidden                 bool                   `json:"hidden"`
	InstallTelemetryInfo   *InstallTelemetryInfo  `json:"installTelemetryInfo"`
	RecommendationState    *RecommendationState   `json:"recommendationState"`
	RootURI                string                 `json:"rootURI"`
	Location               string                 `json:"location"` // i.e. "app-builtin"
}

// Locale contains addon information in a locale.
type Locale struct {
	Name         string               `json:"name"` // Addon name
	Description  string               `json:"description,omitempty"`
	Creator      string               `json:"creator,omitempty"`
	HomepageURL  string               `json:"homepageURL,omitempty"`
	Developers   jsonutil.UnknownType `json:"developers"`
	Translators  jsonutil.UnknownType `json:"translators"`
	Contributors jsonutil.UnknownType `json:"contributors"`
	Locales      []string             `json:"locales"`
}

type TargetApplication struct {
	ID         string `json:"id"`
	MinVersion string `json:"minVersion"`
	MaxVersion string `json:"maxVersion"`
}

type StartupData struct {
	Dictionaries        map[string]string    `json:"dictionaries,omitempty"` // key: locale, value: path
	ChromeEntries       [][]string           `json:"chromeEntries"`
	LangpackID          string               `json:"langpackId,omitempty"`
	L10nRegistrySources *L10nRegistrySources `json:"l10nRegistrySources,omitempty"`
	Languages           []string             `json:"languages"`
}

type L10nRegistrySources struct {
	Toolkit string `json:"toolkit"`
	Browser string `json:"browser"`
}

type InstallTelemetryInfo struct {
	Source    string `json:"source"`           // i.e. "app-system-local"
	Method    string `json:"method,omitempty"` // i.e. "sideload"
	SourceURL string `json:"sourceURL,omitempty"`
}

type RecommendationState struct {
	ValidNotAfter  timefmt.UnixMilli `json:"validNotAfter"`
	ValidNotBefore timefmt.UnixMilli `json:"validNotBefore"`
	States         []string          `json:"states"` // i.e. "line"
}

// ParseExtensions parses extensions.json in a Firefox profile.
func ParseExtensions(filename string) (*Extensions, error) {
	var extensions Extensions
	if err := jsonutil.DecodeFile(filename, &extensions); err != nil {
		return nil, err
	}
	return &extensions, nil
}
