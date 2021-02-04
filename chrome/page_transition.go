// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package chrome

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/andrewarchi/browser/jsonutil"
)

// Page transition types documentation:
// https://developer.chrome.com/docs/extensions/reference/history/#transition-types
// Chromium source:
// https://source.chromium.org/chromium/chromium/src/+/master:ui/base/page_transition_types.h

// PageTransition is a type of transition between pages. A type is made
// of one core value and a set of zero or more qualifiers.
type PageTransition uint32

// Values for PageTransition:
const (
	TransitionFirst            PageTransition = TransitionLink
	TransitionLink             PageTransition = 0
	TransitionTyped            PageTransition = 1
	TransitionAutoBookmark     PageTransition = 2
	TransitionAutoSubframe     PageTransition = 3
	TransitionManualSubframe   PageTransition = 4
	TransitionGenerated        PageTransition = 5
	TransitionAutoToplevel     PageTransition = 6
	TransitionFormSubmit       PageTransition = 7
	TransitionReload           PageTransition = 8
	TransitionKeyword          PageTransition = 9
	TransitionKeywordGenerated PageTransition = 10
	TransitionLastCore         PageTransition = TransitionKeywordGenerated
	TransitionCoreMask         PageTransition = 0xFF

	TransitionFromAPI3       PageTransition = 0x00200000
	TransitionFromAPI2       PageTransition = 0x00400000
	TransitionBlocked        PageTransition = 0x00800000
	TransitionForwardBack    PageTransition = 0x01000000
	TransitionFromAddressBar PageTransition = 0x02000000
	TransitionHomePage       PageTransition = 0x04000000
	TransitionFromAPI        PageTransition = 0x08000000
	TransitionChainStart     PageTransition = 0x10000000
	TransitionChainEnd       PageTransition = 0x20000000
	TransitionClientRedirect PageTransition = 0x40000000
	TransitionServerRedirect PageTransition = 0x80000000
	TransitionIsRedirectMask PageTransition = 0xC0000000
	TransitionQualifierMask  PageTransition = 0xFFFFFF00
)

// MarshalText implements the encoding.TextMarshaler i
func (typ PageTransition) MarshalText() ([]byte, error) {
	return []byte(typ.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler
func (typ *PageTransition) UnmarshalText(data []byte) error {
	t, err := PageTransitionFromString(string(data))
	if err != nil {
		return err
	}
	*typ = t
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (typ PageTransition) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(typ.String())), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (typ *PageTransition) UnmarshalJSON(data []byte) error {
	return jsonutil.QuotedUnmarshal(data, typ)
}

// PageTransitionFromString returns the page transition core value
// corresponding to the string.
func PageTransitionFromString(typ string) (PageTransition, error) {
	switch strings.ToLower(typ) {
	case "link":
		return TransitionLink, nil
	case "typed":
		return TransitionTyped, nil
	case "auto_bookmark":
		return TransitionAutoBookmark, nil
	case "auto_subframe":
		return TransitionAutoSubframe, nil
	case "manual_subframe":
		return TransitionManualSubframe, nil
	case "generated":
		return TransitionGenerated, nil
	case "auto_toplevel":
		return TransitionAutoToplevel, nil
	case "form_submit":
		return TransitionFormSubmit, nil
	case "reload":
		return TransitionReload, nil
	case "keyword":
		return TransitionKeyword, nil
	case "keyword_generated":
		return TransitionKeywordGenerated, nil
	default:
		return 0, fmt.Errorf("chrome: unrecognized transition type: %q", typ)
	}
}

func (typ PageTransition) String() string {
	switch typ & TransitionCoreMask {
	case TransitionLink:
		return "link"
	case TransitionTyped:
		return "typed"
	case TransitionAutoBookmark:
		return "auto_bookmark"
	case TransitionAutoSubframe:
		return "auto_subframe"
	case TransitionManualSubframe:
		return "manual_subframe"
	case TransitionGenerated:
		return "generated"
	case TransitionAutoToplevel:
		return "auto_toplevel"
	case TransitionFormSubmit:
		return "form_submit"
	case TransitionReload:
		return "reload"
	case TransitionKeyword:
		return "keyword"
	case TransitionKeywordGenerated:
		return "keyword_generated"
	default:
		return fmt.Sprintf("transition_type(%d)", uint8(typ))
	}
}
