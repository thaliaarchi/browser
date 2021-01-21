package takeout

import "github.com/andrewarchi/archive/bookmark"

type Chrome struct {
	// Autofill.json
	AutofillProfile []AutofillProfile `json:"Autofill Profile"`
	// Bookmarks.html
	Bookmarks []bookmark.BookmarkEntry
	// BrowserHistory.json
	BrowserHistory []BrowserHistory `json:"Browser History"`
	// Dictionary.csv TODO
	// Extensions.json
	Extensions        []Extension        `json:"Extensions"`
	ExtensionSettings []ExtensionSetting `json:"Extension Settings"`
	// SearchEngines.json
	SearchEngines []SearchEngine `json:"Search Engines"`
	// SyncSettings.json
	Apps        []App        `json:"Apps"`
	AppSettings []AppSetting `json:"App Settings"`
	Preferences []Preference `json:"Preferences"`
	Themes      []Theme      `json:"Themes"`
}

type AutofillProfile struct {
	NameFirst                     []string `json:"name_first"`
	AddressHomeCountry            string   `json:"address_home_country"`
	AddressHomeSortingCode        string   `json:"address_home_sorting_code"`
	AddressHomeState              string   `json:"address_home_state"`
	AddressHomeDependentLocality  string   `json:"address_home_dependent_locality"`
	AddressHomeCity               string   `json:"address_home_city"`
	AddressHomeLanguageCode       string   `json:"address_home_language_code"`
	NameFull                      []string `json:"name_full"`
	Origin                        string   `json:"origin"`
	NameLast                      []string `json:"name_last"`
	IsClientValidityStatesUpdated bool     `json:"is_client_validity_states_updated"`
	NameMiddle                    []string `json:"name_middle"`
	UseCount                      int64    `json:"use_count"`
	EmailAddress                  []string `json:"email_address"`
	ValidityStateBitfield         int64    `json:"validity_state_bitfield"`
	CompanyName                   string   `json:"company_name"`
	AddressHomeLine1              string   `json:"address_home_line1"`
	AddressHomeLine2              string   `json:"address_home_line2"`
	GUID                          string   `json:"guid"`
	AddressHomeZip                string   `json:"address_home_zip"`
	AddressHomeStreetAddress      string   `json:"address_home_street_address"`
	PhoneHomeWholeNumber          []string `json:"phone_home_whole_number"`
	UseDate                       int64    `json:"use_date"`
}

type BrowserHistory struct {
	FaviconURL     string         `json:"favicon_url,omitempty"`
	PageTransition PageTransition `json:"page_transition"`
	Title          string         `json:"title"`
	URL            string         `json:"url"`
	ClientID       string         `json:"client_id"` // base64-encoded
	TimeUsec       int64          `json:"time_usec"`
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
	SuggestionsURL              string         `json:"suggestions_url"`
	ImageURLPostParams          *string        `json:"image_url_post_params,omitempty"`
	FaviconURL                  string         `json:"favicon_url"`
	SafeForAutoreplace          bool           `json:"safe_for_autoreplace"`
	DateCreated                 int64          `json:"date_created"`
	ImageURL                    *string        `json:"image_url,omitempty"`
	URL                         string         `json:"url"`
	NewTabURL                   string         `json:"new_tab_url"`
	InstantURL                  *string        `json:"instant_url,omitempty"`
	OriginatingURL              string         `json:"originating_url"`
	SearchTermsReplacementKey   *string        `json:"search_terms_replacement_key,omitempty"`
	DeprecatedShowInDefaultList *bool          `json:"deprecated_show_in_default_list,omitempty"`
	SyncGUID                    string         `json:"sync_guid"`
	ShortName                   string         `json:"short_name"`
	Keyword                     string         `json:"keyword"`
	InputEncodings              InputEncodings `json:"input_encodings"`
	AlternateUrls               []string       `json:"alternate_urls"`
	PrepopulateID               int64          `json:"prepopulate_id"`
	LastModified                int64          `json:"last_modified"`
}

type InputEncodings string

// Known values for input_encodings.
const (
	EncodingEmpty         InputEncodings = ""
	EncodingISO88591      InputEncodings = "ISO-8859-1"
	EncodingInputEncoding InputEncodings = "inputEncoding"
	EncodingUTF8          InputEncodings = "UTF-8"
	EncodingWindows1252   InputEncodings = "windows-1252"
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
	N PageOrdinal = "n"
)
