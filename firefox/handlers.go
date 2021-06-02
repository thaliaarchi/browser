// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package firefox

import "github.com/andrewarchi/browser/jsonutil"

// Handlers registers handlers for MIME types and URI schemes.
type Handlers struct {
	DefaultHandlersVersion map[string]int      `json:"defaultHandlersVersion"` // key: locale (e.g. "en-US")
	MIMETypes              map[string]MimeType `json:"mimeTypes"`              // key: MIME type (e.g. "image/jpeg")
	Schemes                map[string]Scheme   `json:"schemes"`                // key: URI scheme (e.g. "mailto")
}

// MimeType registers an action to perform for a MIME type and assigns
// file extensions to that MIME type.
type MimeType struct {
	Action     int      `json:"action"` // e.g. 0, 3
	Ask        bool     `json:"ask,omitempty"`
	Extensions []string `json:"extensions,omitempty"` // e.g. "jpg", "jpeg"
}

// Scheme registers an action to perform for a URI scheme and a list of
// handler applications.
type Scheme struct {
	Action    int              `json:"action"` // e.g. 2, 4
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

// ParseHandlers parses handlers.json in a Firefox profile.
func ParseHandlers(filename string) (*Handlers, error) {
	var handlers Handlers
	if err := jsonutil.DecodeFile(filename, &handlers); err != nil {
		return nil, err
	}
	return &handlers, nil
}
