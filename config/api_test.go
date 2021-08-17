// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

// Package config_test provides tests of config package.
//
// TODO: add tests for the following methods:
// 		- GetTime
//		- GetDuration
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
		require.Equal(t, fmt.Sprintf("%v", v), cfg.GetString(k), msg)
	}
}

func testGetBool(t *testing.T, cfg config.Config) {
	values := append(randValues(), []interface{}{"true", "TRUE", "True", "TrUe", "1"})

	for _, v := range values {
		k := util.RandAsciiAlphaNum(8)
		cfg.Set(k, v)
		msg := fmt.Sprintf("value was: %v", v)

		switch v.(type) {
		case bool:
			require.Equal(t, v, cfg.GetBool(k))
		case int:
			if v == 0 {
				require.Equal(t, false, cfg.GetBool(k), msg)
			} else {
				require.Equal(t, true, cfg.GetBool(k), msg)
			}
		case string:
			switch v {
			case "1", "true", "TRUE":
				require.Equal(t, true, cfg.GetBool(k), msg)
			default:
				require.Equal(t, false, cfg.GetBool(k), msg)
			}
		default:
			require.Equal(t, false, cfg.GetBool(k), msg)
		}
	}
}

func testGetInt(t *testing.T, cfg config.Config) {
	for _, v := range randValues() {
		k := util.RandAsciiAlphaNum(8)
		cfg.Set(k, v)
		msg := fmt.Sprintf("value was: %v", v)

		switch v.(type) {
		case bool:
			if v.(bool) == true {
				require.Equal(t, 1, cfg.GetInt(k), msg)
			} else {
				require.Equal(t, 0, cfg.GetInt(k), msg)
			}
		case int:
			require.Equal(t, v, cfg.GetInt(k), msg)
		case float64:
			require.Equal(t, int(v.(float64)), cfg.GetInt(k), msg)
		case string:
			n, err := strconv.Atoi(v.(string))
			if err != nil {
				require.Equal(t, 0, cfg.GetInt(k), msg)
			} else {
				require.Equal(t, n, cfg.GetInt(k), msg)
			}
		default:
			t.Errorf("unexpected type")
		}
	}
}

func testGetFloat(t *testing.T, cfg config.Config) {
	for _, v := range randValues() {
		k := util.RandAsciiAlphaNum(8)
		cfg.Set(k, v)
		msg := fmt.Sprintf("value was: %v", v)

		switch v.(type) {
		case bool:
			if v.(bool) == true {
				require.Equal(t, 1.0, cfg.GetFloat(k), msg)
			} else {
				require.Equal(t, 0.0, cfg.GetFloat(k), msg)
			}
		case int:
			require.Equal(t, float64(v.(int)), cfg.GetFloat(k), msg)
		case float64:
			require.Equal(t, v, cfg.GetFloat(k), msg)
		case string:
			n, err := strconv.ParseFloat(v.(string), 64)
			if err != nil {
				require.Equal(t, 0.0, cfg.GetFloat(k), msg)
			} else {
				require.Equal(t, n, cfg.GetFloat(k), msg)
			}
		default:
			t.Errorf("unexpected type")
		}
	}
}

func TestConfig(t *testing.T) {
	rand.Seed(time.Now().Unix())
	cfg := config.NewTesting(util.RandAsciiAlphaNum(8))

	tests := map[string]func(*testing.T, config.Config){
		"Name":      testName,
		"Get":       testGet,
		"GetString": testGetString,
		"GetBool":   testGetBool,
		"GetInt":    testGetInt,
		"GetFloat":  testGetFloat,
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			fn(t, cfg)
		})
	}
}
