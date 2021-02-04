// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package tabcloud parses browser windows saved by the TabCloud Chrome
// extension.
package tabcloud

import (
	"fmt"
	"io"

	"github.com/andrewarchi/browser/jsonutil"
)

/*
	TabCloud

	     https://chrometabcloud.appspot.com/login     Sign in with Google account
	     https://chrometabcloud.appspot.com/logout    Sign out
	GET  https://chrometabcloud.appspot.com/tabcloud  Retrieve list of windows
	POST https://chrometabcloud.appspot.com/add	      Add a window
	POST https://chrometabcloud.appspot.com/update    Update a window
	POST https://chrometabcloud.appspot.com/move      Rearrange windows
	POST https://chrometabcloud.appspot.com/remove    Delete a window

	To retrieve your windows, open /login in a browser, signin,
	then open /tabcloud.

	Source:
	chrome-extension://npecfdijgoblfcgagoijgmgejmcpnhof/popup.js
*/

type tabCloudResponse struct {
	Status  string   `json:"status"`
	Windows []Window `json:"windows"`
}

// Window is a named set of tabs.
type Window struct {
	Name string `json:"name"`
	Tabs []Tab  `json:"tabs"`
}

// Tab is a tab within a window.
type Tab struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Favicon string `json:"favicon"`
	Pinned  bool   `json:"pinned"`
}

// Parse parses a list of saved windows from the TabCloud Chrome
// extension.
func Parse(r io.Reader) ([]Window, error) {
	var resp tabCloudResponse
	if err := jsonutil.Decode(r, &resp); err != nil {
		return nil, err
	}
	if resp.Status != "loggedin" {
		return nil, fmt.Errorf("tabcloud: status not loggedin: %s", resp.Status)
	}
	return resp.Windows, nil
}
