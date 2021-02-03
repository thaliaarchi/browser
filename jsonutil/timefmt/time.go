// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package timefmt provides types for representing time formats in JSON.
package timefmt

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Epoch uint8

const (
	Unix    Epoch = iota // 1970-01-01 00:00:00 UTC
	Windows              // 1601-01-01 00:00:00 UTC
)

type Unit uint8

const (
	Sec   Unit = 0
	Milli Unit = 3
	Micro Unit = 6
	Nano  Unit = 9
)

var exp = [10]int32{1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9}

func FromInt(n, nsec int64, unit Unit, epoch Epoch) time.Time {
	if n == 0 && nsec == 0 {
		return time.Time{}
	}
	if n < 0 {
		// TODO handle time before epoch
		panic(fmt.Sprintf("negative time: %d", n))
	}
	e := int64(exp[unit])
	e0 := int64(exp[Nano-unit])
	sec := n / e
	nsec += (n % e) * e0
	switch epoch {
	case Unix:
		return time.Unix(sec, nsec).UTC()
	case Windows:
		return time.Date(1601, 1, 1, 0, 0, int(sec), int(nsec), time.UTC)
	default:
		panic(fmt.Sprintf("illegal epoch: %d", epoch))
	}
}

var windowsEpoch = time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC)

func ToInt(t time.Time, unit Unit, epoch Epoch) (n, nsec int64) {
	if t.IsZero() {
		return 0, 0
	}
	e := int64(exp[unit])
	e0 := int64(exp[Nano-unit])
	switch epoch {
	case Unix:
		sec, nsec := t.Unix(), t.UnixNano()/1e9
		return sec*e + nsec/e0, nsec % e0
	case Windows:
		nsec := int64(t.Sub(windowsEpoch))
		// TODO sub loses precision on large times
		return nsec / e0, nsec % e0
	default:
		panic(fmt.Sprintf("illegal epoch: %d", epoch))
	}
}

func Parse(s string, unit Unit, epoch Epoch) (time.Time, error) {
	n, nsec, err := splitFrac(s, unit)
	if err != nil {
		return time.Time{}, err
	}
	return FromInt(n, nsec, unit, epoch), nil
}

func splitFrac(num string, unit Unit) (n, nsec int64, err error) {
	if i := strings.IndexByte(num, '.'); i != -1 {
		frac := num[i+1:]
		nsec, err = strconv.ParseInt(frac, 10, 64)
		if err != nil {
			return
		}
		nsec *= int64(exp[int(Nano-unit)-len(frac)])
		num = num[:i]
	}
	n, err = strconv.ParseInt(num, 10, 64)
	return
}

func Format(t time.Time, unit Unit, epoch Epoch) string {
	return string(FormatBytes(t, unit, epoch))
}

func FormatBytes(t time.Time, unit Unit, epoch Epoch) []byte {
	return Append(nil, t, unit, epoch)
}

func Append(b []byte, t time.Time, unit Unit, epoch Epoch) []byte {
	n, nsec := ToInt(t, unit, epoch)
	b = strconv.AppendInt(b, n, 10)
	if nsec != 0 {
		b = append(b, '.')
		for nsec%10 == 0 {
			nsec /= 10
		}
		b = strconv.AppendInt(b, nsec, 10)
	}
	return b
}

func (e Epoch) String() string {
	switch e {
	case Unix:
		return "unix"
	case Windows:
		return "windows"
	default:
		return fmt.Sprintf("epoch(%d)", e)
	}
}

func (u Unit) String() string {
	switch u {
	case Sec:
		return "sec"
	case Milli:
		return "milli"
	case Micro:
		return "micro"
	default:
		return fmt.Sprintf("unit(%d)", u)
	}
}
