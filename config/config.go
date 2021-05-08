// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

// Package config contains configuration related things.
package config

import "time"

// Config is the configuration interface.
type Config interface {
	// Backend returns backend configuration engine
	Backend() interface{}

	// Get can retrieve any value given the key to use.
	Get(key string) interface{}

	// Sub returns new Config instance representing a sub tree of this instance.
	Sub(key string) Config

	// GetString returns the value associated with the key as a string.
	GetString(key string) string

	// GetBool returns the value associated with the key as a boolean.
	GetBool(key string) bool

	// GetInt returns the value associated with the key as an integer.
	GetInt(key string) int

	// GetFloat returns the value associated with the key as a float64.
	GetFloat(key string) float64

	// GetTime returns the value associated with the key as time.
	GetTime(key string) time.Time

	// GetDuration returns the value associated with the key as a duration.
	GetDuration(key string) time.Duration

	// GetIntSlice returns the value associated with the key as a slice of int values.
	GetIntSlice(key string) []int

	// GetStringSlice returns the value associated with the key as a slice of strings.
	GetStringSlice(key string) []string

	// GetStringMap returns the value associated with the key as a map of interfaces.
	GetStringMap(key string) map[string]interface{}

	// GetStringMapString returns the value associated with the key as a map of strings.
	GetStringMapString(key string) map[string]string

	// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
	GetStringMapStringSlice(key string) map[string][]string

	// IsSet checks to see if the key has been set in any of the data locations.
	IsSet(key string) bool

	// InConfig checks to see if the given key (or an alias) is in the config file.
	InConfig(key string) bool

	// SetDefault sets the default value for this key.
	SetDefault(key string, value interface{})

	// Set sets the value for the key in the override register.
	Set(key string, value interface{})
}
