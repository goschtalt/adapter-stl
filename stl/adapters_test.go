// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: Apache-2.0

package stl

import (
	"errors"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/goschtalt/goschtalt"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalAdapterInternals(t *testing.T) {
	unknownErr := errors.New("unknownErr")
	tests := []struct {
		description string
		from        any
		to          any
		fn          func(reflect.Value, reflect.Value) (any, error)
		expect      any
		expectErr   error
	}{
		{
			description: "stringToDuration",
			from:        "1s",
			to:          time.Duration(1),
			fn:          stringToDuration,
			expect:      time.Second,
		}, {
			description: "stringToDuration - fail",
			from:        "dogs",
			to:          time.Duration(1),
			fn:          stringToDuration,
			expectErr:   unknownErr,
		}, {
			description: "stringToDuration - didn't match",
			from:        "dogs",
			to:          time.Time{},
			fn:          stringToDuration,
			expectErr:   goschtalt.ErrNotApplicable,
		}, {
			description: "stringToIP",
			from:        "127.0.0.1",
			to:          net.IP{},
			fn:          stringToIP,
			expect:      net.ParseIP("127.0.0.1"),
		}, {
			description: "stringToIP - fail",
			from:        "dogs",
			to:          net.IP{},
			fn:          stringToIP,
			expectErr:   unknownErr,
		}, {
			description: "stringToIP - didn't match",
			from:        "dogs",
			to:          time.Time{},
			fn:          stringToIP,
			expectErr:   goschtalt.ErrNotApplicable,
		}, {
			description: "stringToTime",
			from:        "2022-01-30",
			to:          time.Time{},
			fn:          stringToTime("2006-01-02"),
			expect:      time.Date(2022, time.January, 30, 0, 0, 0, 0, time.UTC),
		}, {
			description: "stringToTime - fail",
			from:        "dogs",
			to:          time.Time{},
			fn:          stringToTime("2006-01-02"),
			expectErr:   unknownErr,
		}, {
			description: "stringToTime - didn't match",
			from:        "dogs",
			to:          net.IP{},
			fn:          stringToTime("2006-01-02"),
			expectErr:   goschtalt.ErrNotApplicable,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)

			got, err := tc.fn(reflect.ValueOf(tc.from), reflect.ValueOf(tc.to))

			if errors.Is(unknownErr, tc.expectErr) {
				// Accept nil, or the zero value of the type
				if nil != got {
					assert.Equal(reflect.Zero(reflect.TypeOf(got)).Interface(), got)
				}
				assert.Error(err)
				return
			}

			if tc.expectErr != nil {
				// Accept nil, or the zero value of the type
				if nil != got {
					assert.Equal(reflect.Zero(reflect.TypeOf(got)).Interface(), got)
				}
				assert.ErrorIs(err, tc.expectErr)
				return
			}

			assert.Equal(tc.expect, got)
			assert.NoError(err)
		})
	}
}

func TestValueAdapterInternals(t *testing.T) {
	unknownErr := errors.New("unknownErr")
	tests := []struct {
		description string
		from        any
		fn          func(reflect.Value) (any, error)
		expect      any
		expectErr   error
	}{
		{
			description: "durationToCfg",
			from:        time.Second,
			fn:          durationToCfg,
			expect:      "1s",
		}, {
			description: "durationToCfg - didn't match",
			from:        time.Time{},
			fn:          durationToCfg,
			expectErr:   goschtalt.ErrNotApplicable,
		}, {
			description: "ipToCfg",
			from:        net.ParseIP("127.0.0.1"),
			fn:          ipToCfg,
			expect:      "127.0.0.1",
		}, {
			description: "ipToCfg - didn't match",
			from:        "dogs",
			fn:          ipToCfg,
			expectErr:   goschtalt.ErrNotApplicable,
		}, {
			description: "timeToCfg",
			from:        time.Date(2022, time.January, 30, 0, 0, 0, 0, time.UTC),
			fn:          timeToCfg("2006-01-02"),
			expect:      "2022-01-30",
		}, {
			description: "timeToCfg - didn't match",
			from:        net.IP{127, 0, 0, 1},
			fn:          timeToCfg("2006-01-02"),
			expectErr:   goschtalt.ErrNotApplicable,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)

			got, err := tc.fn(reflect.ValueOf(tc.from))

			if errors.Is(unknownErr, tc.expectErr) {
				// Accept nil, or the zero value of the type
				if nil != got {
					assert.Equal(reflect.Zero(reflect.TypeOf(got)).Interface(), got)
				}
				assert.Error(err)
				return
			}

			if tc.expectErr != nil {
				// Accept nil, or the zero value of the type
				if nil != got {
					assert.Equal(reflect.Zero(reflect.TypeOf(got)).Interface(), got)
				}
				assert.ErrorIs(err, tc.expectErr)
				return
			}

			assert.Equal(tc.expect, got)
			assert.NoError(err)
		})
	}
}

func TestEndToEnd(t *testing.T) {
	type all struct {
		D  time.Duration
		T  time.Time
		IP net.IP
	}
	tests := []struct {
		description string
		from        all
		unmarshal   []goschtalt.UnmarshalOption
		value       []goschtalt.ValueOption
		expectErr   error
	}{
		{
			description: "String <-> time.Duration",
			from:        all{D: time.Second},
			unmarshal:   []goschtalt.UnmarshalOption{AdaptStringToDuration()},
			value:       []goschtalt.ValueOption{AdaptDurationToCfg()},
		}, {
			description: "String <-> time.Time",
			from:        all{T: time.Date(2022, time.August, 15, 11, 10, 9, 0, time.UTC)},
			unmarshal:   []goschtalt.UnmarshalOption{AdaptStringToTime(time.RFC3339)},
			value:       []goschtalt.ValueOption{AdaptTimeToCfg(time.RFC3339)},
		}, {
			description: "String <-> net.IP",
			from:        all{IP: net.ParseIP("127.0.0.1")},
			unmarshal:   []goschtalt.UnmarshalOption{AdaptStringToIP()},
			value:       []goschtalt.ValueOption{AdaptIPToCfg()},
		}, {
			description: "String <-> all",
			from: all{
				D:  time.Hour,
				T:  time.Date(2002, time.August, 15, 0, 0, 0, 0, time.UTC),
				IP: net.ParseIP("192.168.1.1"),
			},
			unmarshal: []goschtalt.UnmarshalOption{
				AdaptStringToDuration(),
				AdaptStringToTime("2006-01-02"),
				AdaptStringToIP(),
			},
			value: []goschtalt.ValueOption{
				AdaptDurationToCfg(),
				AdaptTimeToCfg("2006-01-02"),
				AdaptIPToCfg(),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)

			cfg, err := goschtalt.New(
				goschtalt.AutoCompile(),
				goschtalt.DefaultUnmarshalOptions(tc.unmarshal...),
				goschtalt.DefaultValueOptions(tc.value...),
				goschtalt.AddValue("rec", goschtalt.Root, tc.from),
			)

			assert.NoError(err)

			var got all

			err = cfg.Unmarshal(goschtalt.Root, &got)

			if tc.expectErr != nil {
				assert.Error(err)
				return
			}

			assert.Empty(cmp.Diff(tc.from, got))
			assert.NoError(err)
		})
	}
}
