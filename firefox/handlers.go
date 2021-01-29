package firefox

import "github.com/andrewarchi/browser/jsonutil"

// Handlers registers handlers for MIME types and URI schemes.
type Handlers struct {
	DefaultHandlersVersion map[string]int      `json:"defaultHandlersVersion"` // key: locale (i.e. "en-US")
	MIMETypes              map[string]MimeType `json:"mimeTypes"`              // key: MIME type (i.e. "image/jpeg")
	Schemes                map[string]Scheme   `json:"schemes"`                // key: URI scheme (i.e. "mailto")
}

// MimeType registers an action to perform for a MIME type and assigns
// file extensions to that MIME type.
type MimeType struct {
	Action     int      `json:"action"` // i.e. 0, 3
	Ask        bool     `json:"ask,omitempty"`
	Extensions []string `json:"extensions,omitempty"` // i.e. "jpg", "jpeg"
}

// Scheme registers an action to perform for a URI scheme and a list of
// handler applications.
type Scheme struct {
	Action    int              `json:"action"` // i.e. 2, 4
	Ask       bool             `json:"ask,omitempty"`
	StubEntry bool             `json:"stubEntry,omitempty"` // true when handler unchanged from default
	Handlers  []*SchemeHandler `json:"handlers,omitempty"`
}

// SchemeHandler is an application that can handle a URI scheme.
type SchemeHandler struct {
	Name        string `json:"name"`
	Path        string `json:"path,omitempty"`        // for local apps
	URITemplate string `json:"uriTemplate,omitempty"` // for web apps
}

// ParseHandlers parses the handlers.json file in a Firefox profile.
func ParseHandlers(filename string) (*Handlers, error) {
	var handlers Handlers
	if err := jsonutil.DecodeFile(filename, &handlers); err != nil {
		return nil, err
	}
	return &handlers, nil
}
