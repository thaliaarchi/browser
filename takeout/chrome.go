package takeout

import (
	"time"

	"github.com/andrewarchi/browser/bookmark"
	"github.com/andrewarchi/browser/timefmt"
)

type Chrome struct {
	ExportTime time.Time `json:"-"`
	// Autofill.json
	Autofill        []AutofillProfile `json:"Autofill"` // Appears in older exports
	AutofillProfile []AutofillProfile `json:"Autofill Profile"`
	// Bookmarks.html
	Bookmarks []bookmark.BookmarkEntry
	// BrowserHistory.json
	BrowserHistory []BrowserHistory `json:"Browser History"`
	// Dictionary.csv - TODO unknown structure
	// Extensions.json
	Extensions        []Extension        `json:"Extensions"`
	ExtensionSettings []ExtensionSetting `json:"Extension Settings"`
	// SearchEngines.json
	SearchEngines []SearchEngine `json:"Search Engines"`
	// SyncSettings.json
	Apps         []App         `json:"Apps"`
	AppSettings  []AppSetting  `json:"App Settings"`
	Preferences  []Preference  `json:"Preferences"`
	Themes       []Theme       `json:"Themes"`
	ManagedUsers []interface{} `json:"Managed Users"` // TODO unknown structure
}

type AutofillProfile struct {
	GUID                          string          `json:"guid"`
	NameFull                      []string        `json:"name_full"`
	NameFirst                     []string        `json:"name_first"`
	NameMiddle                    []string        `json:"name_middle"`
	NameLast                      []string        `json:"name_last"`
	AddressHomeStreetAddress      string          `json:"address_home_street_address"`
	AddressHomeLine1              string          `json:"address_home_line1"`
	AddressHomeLine2              string          `json:"address_home_line2"`
	AddressHomeCity               string          `json:"address_home_city"`
	AddressHomeState              string          `json:"address_home_state"`
	AddressHomeZip                string          `json:"address_home_zip"`
	AddressHomeCountry            string          `json:"address_home_country"`
	AddressHomeSortingCode        string          `json:"address_home_sorting_code"`
	AddressHomeLanguageCode       string          `json:"address_home_language_code"`
	AddressHomeDependentLocality  string          `json:"address_home_dependent_locality"`
	EmailAddress                  []string        `json:"email_address"`
	PhoneHomeWholeNumber          []string        `json:"phone_home_whole_number"`
	Origin                        string          `json:"origin"`
	IsClientValidityStatesUpdated bool            `json:"is_client_validity_states_updated"`
	UseCount                      int             `json:"use_count"`
	ValidityStateBitfield         uint64          `json:"validity_state_bitfield"` // TODO unknown states
	CompanyName                   string          `json:"company_name"`
	UseDate                       timefmt.UnixSec `json:"use_date"`
}

type BrowserHistory struct {
	FaviconURL     string            `json:"favicon_url,omitempty"`
	PageTransition PageTransition    `json:"page_transition"`
	Title          string            `json:"title"`
	URL            string            `json:"url"`
	ClientID       string            `json:"client_id"` // base64-encoded
	Time           timefmt.UnixMicro `json:"time_usec"`
}

type PageTransition string

// Known values for page_transition.
const (
	TransitionAutoBookmark PageTransition = "AUTO_BOOKMARK"
	TransitionAutoToplevel PageTransition = "AUTO_TOPLEVEL"
	TransitionFormSubmit   PageTransition = "FORM_SUBMIT"
	TransitionGenerated    PageTransition = "GENERATED"
	TransitionKeyword      PageTransition = "KEYWORD"
	TransitionLink         PageTransition = "LINK"
	TransitionReload       PageTransition = "RELOAD"
	TransitionTyped        PageTransition = "TYPED"
)

type Extension struct {
	IncognitoEnabled     bool   `json:"incognito_enabled"`
	RemoteInstall        bool   `json:"remote_install"`
	DisableReasons       *int64 `json:"disable_reasons,omitempty"`
	InstalledByCustodian *bool  `json:"installed_by_custodian,omitempty"`
	UpdateURL            string `json:"update_url"`
	Name                 string `json:"name"`
	ID                   string `json:"id"`
	Version              string `json:"version"`
	Enabled              bool   `json:"enabled"`
}

type ExtensionSetting struct {
	ExtensionID string `json:"extension_id"`
	Value       string `json:"value"`
	Key         string `json:"key"`
}

type SearchEngine struct {
	ShortName                   string              `json:"short_name"`
	Keyword                     string              `json:"keyword"`
	URL                         string              `json:"url"`
	SuggestionsURL              string              `json:"suggestions_url"`
	FaviconURL                  string              `json:"favicon_url"`
	ImageURL                    *string             `json:"image_url,omitempty"`
	NewTabURL                   string              `json:"new_tab_url"`
	InstantURL                  *string             `json:"instant_url,omitempty"`
	OriginatingURL              string              `json:"originating_url"`
	ImageURLPostParams          *string             `json:"image_url_post_params,omitempty"`
	SafeForAutoreplace          bool                `json:"safe_for_autoreplace"`
	DateCreated                 timefmt.ChromeMicro `json:"date_created"`
	LastModified                timefmt.ChromeMicro `json:"last_modified"`
	SearchTermsReplacementKey   *string             `json:"search_terms_replacement_key,omitempty"`
	DeprecatedShowInDefaultList *bool               `json:"deprecated_show_in_default_list,omitempty"`
	SyncGUID                    string              `json:"sync_guid"`
	InputEncodings              InputEncodings      `json:"input_encodings"`
	AlternateUrls               []string            `json:"alternate_urls,omitempty"`
	PrepopulateID               int64               `json:"prepopulate_id"`
}

type InputEncodings string

// Known values for input_encodings.
const (
	EncodingEmpty       InputEncodings = ""
	EncodingUTF8        InputEncodings = "UTF-8"
	EncodingISO88591    InputEncodings = "ISO-8859-1"
	EncodingWindows1252 InputEncodings = "windows-1252"
	EncodingInput       InputEncodings = "inputEncoding"
)

type App struct {
	AppLaunchOrdinal string      `json:"app_launch_ordinal"`
	Extension        Extension   `json:"extension"`
	PageOrdinal      PageOrdinal `json:"page_ordinal"`
}

type AppSetting struct {
	ExtensionSetting ExtensionSetting `json:"extension_setting"`
}

type Preference struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Theme struct {
	UseSystemThemeByDefault bool `json:"use_system_theme_by_default"`
	UseCustomTheme          bool `json:"use_custom_theme"`
}

type PageOrdinal string

// Known values for page_ordinal.
const (
	PageOrdinalN PageOrdinal = "n"
	PageOrdinalT PageOrdinal = "t"
)
