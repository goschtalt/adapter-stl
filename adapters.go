// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package stl

import (
	"fmt"
	"net"
	"reflect"
	"time"

	"github.com/goschtalt/goschtalt"
)

// AdaptStringToDuration converts a string to a time.Duration if possible, or
// returns an error indicating the failure.
func AdaptStringToDuration() goschtalt.UnmarshalOption {
	return goschtalt.AdaptFromCfg(stringToDuration, "AdaptStringToDuration")
}

func stringToDuration(from, to reflect.Value) (any, error) {
	if from.Kind() == reflect.String && to.Type() == reflect.TypeOf(time.Duration(1)) {
		return time.ParseDuration(from.Interface().(string))
	}

	return nil, goschtalt.ErrNotApplicable
}

// AdaptDurationToCfg converts a time.Duration into its configuration form.  The
// configuration form is a string.
func AdaptDurationToCfg() goschtalt.ValueOption {
	return goschtalt.AdaptToCfg(durationToCfg, "AdaptDurationToCfg")
}

func durationToCfg(from reflect.Value) (any, error) {
	if from.Type() == reflect.TypeOf(time.Duration(1)) {
		return from.Interface().(time.Duration).String(), nil
	}

	return nil, goschtalt.ErrNotApplicable
}

// AdaptStringToIP converts a string to a net.IP if possible, or returns an
// error indicating the failure.
func AdaptStringToIP() goschtalt.UnmarshalOption {
	return goschtalt.AdaptFromCfg(stringToIP, "AdaptStringToIP")
}

func stringToIP(from, to reflect.Value) (any, error) {
	if from.Kind() == reflect.String && to.Type() == reflect.TypeOf(net.IP{}) {
		ip := net.ParseIP(from.Interface().(string))
		if ip == nil {
			return nil, fmt.Errorf("failed parsing ip %v", from)
		}
		return ip, nil
	}

	return nil, goschtalt.ErrNotApplicable
}

// AdaptIPToCfg converts a net.IP into its configuration form.  The
// configuration form is a string.
func AdaptIPToCfg() goschtalt.ValueOption {
	return goschtalt.AdaptToCfg(ipToCfg, "AdaptIPToCfg")
}

func ipToCfg(from reflect.Value) (any, error) {
	if from.Type() == reflect.TypeOf(net.IP{}) {
		return from.Interface().(net.IP).String(), nil
	}

	return nil, goschtalt.ErrNotApplicable
}

// AdaptStringToTime converts a string to a time.Time if possible, or returns an
// error indicating the failure.  The specified layout is used as the string
// form.
func AdaptStringToTime(layout string) goschtalt.UnmarshalOption {
	return goschtalt.AdaptFromCfg(stringToTime(layout), "AdaptStringToTime")
}

func stringToTime(layout string) func(reflect.Value, reflect.Value) (any, error) {
	return func(from, to reflect.Value) (any, error) {
		if from.Kind() == reflect.String && to.Type() == reflect.TypeOf(time.Time{}) {
			a, e := time.Parse(layout, from.Interface().(string))
			return a, e
		}

		return nil, goschtalt.ErrNotApplicable
	}
}

// AdaptTimeToCfg converts a time.Time into its configuration form. The
// configuration form is a string matching the specified layout.
func AdaptTimeToCfg(layout string) goschtalt.ValueOption {
	return goschtalt.AdaptToCfg(timeToCfg(layout), "AdaptTimeToCfg")
}

func timeToCfg(layout string) func(reflect.Value) (any, error) {
	return func(from reflect.Value) (any, error) {
		if from.Type() == reflect.TypeOf(time.Time{}) {
			a := from.Interface().(time.Time).Format(layout)
			return a, nil
		}

		return nil, goschtalt.ErrNotApplicable
	}
}
