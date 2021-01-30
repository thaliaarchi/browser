// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bookmark

import (
	"encoding/json"
	"os"
	"testing"
)

func TestBookmarks(t *testing.T) {
	f, err := os.Open("Bookmarks.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	b, err := ParseHTML(f)
	if err != nil {
		t.Fatal(err)
	}
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	if err := e.Encode(b); err != nil {
		t.Error(err)
	}
}
