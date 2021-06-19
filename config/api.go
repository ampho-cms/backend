// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package config

import "time"

var defaultConfig Config

// SetConfig sets default config.
func SetConfig(cfg Config) {
	defaultConfig = cfg
}

// GetConfig returns default config.
func GetConfig() Config {
	if defaultConfig == nil {
		panic("Default config is not set. Did you forget to call SetConfig()?")
	}

	return defaultConfig
}

// Get can retrieve any value given the key to use.
func Get(key string) interface{} {
	return GetConfig().Get(key)
}

// Sub returns new Config instance representing a sub tree of this instance.
func Sub(key string) Config {
	return GetConfig().Sub(key)
}

// GetString returns the value associated with the key as a string.
func GetString(key string) string {
	return GetConfig().GetString(key)
}

// GetBool returns the value associated with the key as a boolean.
func GetBool(key string) bool {
	return GetConfig().GetBool(key)
}

// GetInt returns the value associated with the key as an integer.
func GetInt(key string) int {
	return GetConfig().GetInt(key)
}

// GetFloat returns the value associated with the key as a float64.
func GetFloat(key string) float64 {
	return GetConfig().GetFloat(key)
}

// GetTime returns the value associated with the key as time.
func GetTime(key string) time.Time {
	return GetConfig().GetTime(key)
}

// GetDuration returns the value associated with the key as a duration.
func GetDuration(key string) time.Duration {
	return GetConfig().GetDuration(key)
}

// GetIntSlice returns the value associated with the key as a slice of int values.
func GetIntSlice(key string) []int {
	return GetConfig().GetIntSlice(key)
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func GetStringSlice(key string) []string {
	return GetConfig().GetStringSlice(key)
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func GetStringMap(key string) map[string]interface{} {
	return GetConfig().GetStringMap(key)
}

// GetStringMapString returns the value associated with the key as a map of strings.
func GetStringMapString(key string) map[string]string {
	return GetConfig().GetStringMapString(key)
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func GetStringMapStringSlice(key string) map[string][]string {
	return GetConfig().GetStringMapStringSlice(key)
}

// IsSet checks to see if the key has been set in any of the data locations.
func IsSet(key string) bool {
	return GetConfig().IsSet(key)
}

// InConfig checks to see if the given key (or an alias) is in the config file.
func InConfig(key string) bool {
	return GetConfig().InConfig(key)
}

// SetDefault sets the default value for this key.
func SetDefault(key string, value interface{}) {
	GetConfig().SetDefault(key, value)
}

// Set sets the value for the key in the override register.
func Set(key string, value interface{}) {
	GetConfig().Set(key, value)
}

// All returns all values.
func All() map[string]interface{} {
	return GetConfig().All()
}

// AllKeys returns all keys.
func AllKeys() []string {
	return GetConfig().AllKeys()
}
