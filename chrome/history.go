// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package chrome

import (
	"fmt"
	"strings"
)

// Transition types documentation:
// https://developer.chrome.com/docs/extensions/reference/history/#transition-types
//
// See page_transition_types.h in Chromium.

type TransitionType uint8

// TODO make sure TransitionType values have the same int vales as Chrome.

const (
	_                TransitionType = iota
	Link                            // link
	Typed                           // typed
	AutoBookmark                    // auto_bookmark
	AutoSubframe                    // auto_subframe
	ManualSubframe                  // manual_subframe
	Generated                       // generated
	AutoToplevel                    // auto_toplevel
	FormSubmit                      // form_submit
	Reload                          // reload
	Keyword                         // keyword
	KeywordGenerated                // keyword_generated
)

func TransitionTypeFromString(typ string) (TransitionType, error) {
	switch strings.ToLower(typ) {
	case "link":
		return Link, nil
	case "typed":
		return Typed, nil
	case "auto_bookmark":
		return AutoBookmark, nil
	case "auto_subframe":
		return AutoSubframe, nil
	case "manual_subframe":
		return ManualSubframe, nil
	case "generated":
		return Generated, nil
	case "auto_toplevel":
		return AutoToplevel, nil
	case "form_submit":
		return FormSubmit, nil
	case "reload":
		return Reload, nil
	case "keyword":
		return Keyword, nil
	case "keyword_generated":
		return KeywordGenerated, nil
	default:
		return 0, fmt.Errorf("chrome: illegal transition type: %q", typ)
	}
}

func (typ TransitionType) String() string {
	switch typ {
	case Link:
		return "link"
	case Typed:
		return "typed"
	case AutoBookmark:
		return "auto_bookmark"
	case AutoSubframe:
		return "auto_subframe"
	case ManualSubframe:
		return "manual_subframe"
	case Generated:
		return "generated"
	case AutoToplevel:
		return "auto_toplevel"
	case FormSubmit:
		return "form_submit"
	case Reload:
		return "reload"
	case Keyword:
		return "keyword"
	case KeywordGenerated:
		return "keyword_generated"
	default:
		return fmt.Sprintf("transition_type(%d)", uint8(typ))
	}
}
