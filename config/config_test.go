// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

// Package config_test provides tests of config package.
//
// TODO: add tests for the following methods:
//		- GetIntSlice
//		- GetStringSlice
//		- GetStringMap
//		- GetStringMapString
//		- GetStringMapStringSlice
//		- GetStringMapStringSlice
// 		- IsSet
// 		- InConfig
// 		- SetDefault
// 		- All
// 		- AllKeys

package config_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"

	"ampho.xyz/core/config"
	"ampho.xyz/core/util"
)

func randValues() []interface{} {
	return []interface{}{
		// bool
		true,
		false,

		// int
		0,
		rand.Intn(1000),
		-rand.Intn(1000),
		0.0,

		// float64
		rand.Float64() * float64(rand.Intn(1000)),
		-rand.Float64() * float64(rand.Intn(1000)),

		// string
		util.RandAscii(8),

		// number as a string
		"0",
		strconv.Itoa(rand.Intn(1000)),
		strconv.Itoa(-rand.Intn(1000)),
		"0.0",
		fmt.Sprintf("%f", rand.Float64()*float64(rand.Intn(1000))),
		fmt.Sprintf("%f", -rand.Float64()*float64(rand.Intn(1000))),
	}
}

func testName(t *testing.T, cfg config.Config) {
	require.NotEmpty(t, cfg.Name())
}

func testGet(t *testing.T, cfg config.Config) {
	for _, v := range randValues() {
		k := util.RandAsciiAlphaNum(8)
		cfg.Set(k, v)
		msg := fmt.Sprintf("value was: %v", v)
		require.Equal(t, v, cfg.Get(k), msg)
	}
}

func testGetString(t *testing.T, cfg config.Config) {
	for _, v := range randValues() {
		k := util.RandAsciiAlphaNum(8)
		cfg.Set(k, v)
		msg := fmt.Sprintf("value was: %v", v)
		require.Equal(t, cast.ToString(v), cfg.GetString(k), msg)
	}
}

func testGetBool(t *testing.T, cfg config.Config) {
	values := append(randValues(), []interface{}{"true", "TRUE", "True", "TrUe", "1"})

	for _, v := range values {
		k := util.RandAsciiAlphaNum(8)
		cfg.Set(k, v)
		msg := fmt.Sprintf("value was: %v", v)
		require.Equal(t, cast.ToBool(v), cfg.GetBool(k), msg)
	}
}

func testGetInt(t *testing.T, cfg config.Config) {
	for _, v := range randValues() {
		k := util.RandAsciiAlphaNum(8)
		cfg.Set(k, v)
		msg := fmt.Sprintf("value was: %v", v)
		require.Equal(t, cast.ToInt(v), cfg.GetInt(k), msg)
	}
}

func testGetFloat(t *testing.T, cfg config.Config) {
	for _, v := range randValues() {
		k := util.RandAsciiAlphaNum(8)
		cfg.Set(k, v)
		msg := fmt.Sprintf("value was: %v", v)
		require.Equal(t, cast.ToFloat64(v), cfg.GetFloat(k), msg)
	}
}

func testGetTime(t *testing.T, cfg config.Config) {
	tests := []struct {
		input  interface{}
		expect time.Time
		isErr  bool
	}{
		{"2009-11-10 23:00:00 +0000 UTC", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},  // Time.String()
		{"Tue Nov 10 23:00:00 2009", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},       // ANSI
		{"Tue Nov 10 23:00:00 UTC 2009", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},   // UnixDate
		{"Tue Nov 10 23:00:00 +0000 2009", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false}, // RubyDate
		{"10 Nov 09 23:00 UTC", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},            // RFC822
		//{"10 Nov 09 23:00 +0000", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},           // RFC822Z
		{"Tuesday, 10-Nov-09 23:00:00 UTC", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false}, // RFC850
		{"Tue, 10 Nov 2009 23:00:00 UTC", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},   // RFC1123
		//{"Tue, 10 Nov 2009 23:00:00 +0000", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false}, // RFC1123Z
		{"2009-11-10T23:00:00Z", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false}, // RFC3339
		//{"2018-10-21T23:21:29+0200", time.Date(2018, 10, 21, 21, 21, 29, 0, time.UTC), false},      // RFC3339 without timezone hh:mm colon
		{"2009-11-10T23:00:00Z", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},   // RFC3339Nano
		{"11:00PM", time.Date(0, 1, 1, 23, 0, 0, 0, time.UTC), false},                     // Kitchen
		{"Nov 10 23:00:00", time.Date(0, 11, 10, 23, 0, 0, 0, time.UTC), false},           // Stamp
		{"Nov 10 23:00:00.000", time.Date(0, 11, 10, 23, 0, 0, 0, time.UTC), false},       // StampMilli
		{"Nov 10 23:00:00.000000", time.Date(0, 11, 10, 23, 0, 0, 0, time.UTC), false},    // StampMicro
		{"Nov 10 23:00:00.000000000", time.Date(0, 11, 10, 23, 0, 0, 0, time.UTC), false}, // StampNano
		//{"2016-03-06 15:28:01-00:00", time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC), false},        // RFC3339 without T
		//{"2016-03-06 15:28:01-0000", time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC), false},         // RFC3339 without T or timezone hh:mm colon
		{"2016-03-06 15:28:01", time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC), false},
		//{"2016-03-06 15:28:01 -0000", time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC), false},
		//{"2016-03-06 15:28:01 -00:00", time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC), false},
		{"2006-01-02", time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC), false},
		{"02 Jan 2006", time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC), false},
		{1472574600, time.Unix(1472574600, 0), false},
		{int64(1234567890), time.Unix(1234567890, 0), false},
		{int32(1234567890), time.Unix(1234567890, 0), false},
		{uint(1482597504), time.Unix(1482597504, 0), false},
		{uint64(1234567890), time.Unix(1234567890, 0), false},
		{uint32(1234567890), time.Unix(1234567890, 0), false},
		{time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), false},

		// errors
		{"2006", time.Time{}, true},
		{testing.T{}, time.Time{}, true},
	}

	for _, v := range tests {
		k := util.RandAsciiAlphaNum(8)
		cfg.Set(k, v.input)
		msg := fmt.Sprintf("input was: %#v", v.input)

		if v.isErr {
			require.Equal(t, time.Time{}, cfg.GetTime(k), msg)
		} else {
			require.Equal(t, v.expect, cfg.GetTime(k), msg)
		}
	}
}

func testGetDuration(t *testing.T, cfg config.Config) {
	var td time.Duration = 5

	tests := []struct {
		input  interface{}
		expect time.Duration
		isErr  bool
	}{
		{time.Duration(5), td, false},
		{5, td, false},
		{int64(5), td, false},
		{int32(5), td, false},
		{int16(5), td, false},
		{int8(5), td, false},
		{uint(5), td, false},
		{uint64(5), td, false},
		{uint32(5), td, false},
		{uint16(5), td, false},
		{uint8(5), td, false},
		{float64(5), td, false},
		{float32(5), td, false},
		{"5", td, false},
		{"5ns", td, false},
		{"5us", time.Microsecond * td, false},
		{"5µs", time.Microsecond * td, false},
		{"5ms", time.Millisecond * td, false},
		{"5s", time.Second * td, false},
		{"5m", time.Minute * td, false},
		{"5h", time.Hour * td, false},

		// errors
		{"test", 0, true},
		{testing.T{}, 0, true},
	}

	for _, v := range tests {
		k := util.RandAsciiAlphaNum(8)
		cfg.Set(k, v.input)
		msg := fmt.Sprintf("input was: %#v", v.input)

		if v.isErr {
			require.Equal(t, time.Duration(0), cfg.GetDuration(k), msg)
		} else {
			require.Equal(t, v.expect, cfg.GetDuration(k), msg)
		}
	}
}

func testGetIntSlice(t *testing.T, cfg config.Config) {
	tests := []struct {
		input  interface{}
		expect []int
		isErr  bool
	}{
		{[]int{1, 3}, []int{1, 3}, false},
		{[]interface{}{1.2, 3.2}, []int{1, 3}, false},
		{[]string{"2", "3"}, []int{2, 3}, false},
		{[2]string{"2", "3"}, []int{2, 3}, false},

		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
	}

	for _, v := range tests {
		k := util.RandAsciiAlphaNum(8)
		cfg.Set(k, v.input)
		msg := fmt.Sprintf("input was: %#v", v.input)

		if v.isErr {
			require.Equal(t, []int{}, cfg.GetIntSlice(k), msg)
		} else {
			require.Equal(t, v.expect, cfg.GetIntSlice(k), msg)
		}
	}
}

func TestConfig(t *testing.T) {
	rand.Seed(time.Now().Unix())
	cfg := config.NewTesting(util.RandAsciiAlphaNum(8))

	tests := map[string]func(*testing.T, config.Config){
		"Name":        testName,
		"Get":         testGet,
		"GetString":   testGetString,
		"GetBool":     testGetBool,
		"GetInt":      testGetInt,
		"GetFloat":    testGetFloat,
		"GetTime":     testGetTime,
		"GetDuration": testGetDuration,
		"GetIntSlice": testGetIntSlice,
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			fn(t, cfg)
		})
	}
}
