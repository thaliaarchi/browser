package firefox

import (
	"github.com/andrewarchi/browser/jsonutil"
	"github.com/andrewarchi/browser/jsonutil/timefmt"
)

// Addons contains information from addons.mozilla.org for installed
// addons.
type Addons struct {
	Schema int         `json:"schema"` // i.e. 6
	Addons []AddonInfo `json:"addons"`
}

// AddonInfo contains addon information from addons.mozilla.org.
type AddonInfo struct {
	ID              string            `json:"id"`    // ID "addon@example.com" or GUID "{01234567-89ab-cdef-0123-456789abcdef}"
	Icons           map[int]string    `json:"icons"` // key: icon size, value: path
	Type            string            `json:"type"`  // i.e. "extension", "locale", "dictionary"
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Creator         Person            `json:"creator"`
	Developers      []Person          `json:"developers"`
	Description     string            `json:"description"`
	FullDescription string            `json:"fullDescription"`
	Screenshots     []Screenshot      `json:"screenshots"`
	HomepageURL     string            `json:"homepageURL"`
	SupportURL      string            `json:"supportURL"`
	ContributionURL string            `json:"contributionURL"`
	AverageRating   float64           `json:"averageRating"` // out of 5
	ReviewCount     int               `json:"reviewCount"`
	ReviewURL       string            `json:"reviewURL"`
	WeeklyDownloads int               `json:"weeklyDownloads"`
	SourceURI       string            `json:"sourceURI"` // URI to .xpi
	UpdateDate      timefmt.UnixMilli `json:"updateDate"`
}

// Person is an addon creator or developer.
type Person struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Screenshot is an image of an addon displayed on addons.mozilla.org.
type Screenshot struct {
	URL             string `json:"url"`
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	ThumbnailURL    string `json:"thumbnailURL"`
	ThumbnailWidth  int    `json:"thumbnailWidth"`
	ThumbnailHeight int    `json:"thumbnailHeight"`
	Caption         string `json:"caption,omitempty"`
}

// ParseAddons parses addons.json in a Firefox profile.
func ParseAddons(filename string) (*Addons, error) {
	var addons Addons
	if err := jsonutil.DecodeFile(filename, &addons); err != nil {
		return nil, err
	}
	return &addons, nil
}
