// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package config

import (
	"time"

	"github.com/spf13/viper"
)

// Viper is the Viper config backend.
type Viper struct {
	backend *viper.Viper
}

// Backend returns backend configuration engine.
func (v *Viper) Backend() interface{} {
	return v.backend
}

// Get can retrieve any value given the key to use.
func (v *Viper) Get(key string) interface{} {
	return v.backend.Get(key)
}

// Sub returns new Config instance representing a sub tree of this instance.
func (v *Viper) Sub(key string) Config {
	return &Viper{v.backend.Sub(key)}
}

// GetString returns the value associated with the key as a string.
func (v *Viper) GetString(key string) string {
	return v.backend.GetString(key)
}

// GetBool returns the value associated with the key as a boolean.
func (v *Viper) GetBool(key string) bool {
	return v.backend.GetBool(key)
}

// GetInt returns the value associated with the key as an integer.
func (v *Viper) GetInt(key string) int {
	return v.backend.GetInt(key)
}

// GetFloat returns the value associated with the key as a float64.
func (v *Viper) GetFloat(key string) float64 {
	return v.backend.GetFloat64(key)
}

// GetTime returns the value associated with the key as time.
func (v *Viper) GetTime(key string) time.Time {
	return v.backend.GetTime(key)
}

// GetDuration returns the value associated with the key as a duration.
func (v *Viper) GetDuration(key string) time.Duration {
	return v.backend.GetDuration(key)
}

// GetIntSlice returns the value associated with the key as a slice of int values.
func (v *Viper) GetIntSlice(key string) []int {
	return v.backend.GetIntSlice(key)
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (v *Viper) GetStringSlice(key string) []string {
	return v.backend.GetStringSlice(key)
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (v *Viper) GetStringMap(key string) map[string]interface{} {
	return v.backend.GetStringMap(key)
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (v *Viper) GetStringMapString(key string) map[string]string {
	return v.backend.GetStringMapString(key)
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (v *Viper) GetStringMapStringSlice(key string) map[string][]string {
	return v.backend.GetStringMapStringSlice(key)
}

// IsSet checks to see if the key has been set in any of the data locations.
func (v *Viper) IsSet(key string) bool {
	return v.backend.IsSet(key)
}

// InConfig checks to see if the given key (or an alias) is in the config file.
func (v *Viper) InConfig(key string) bool {
	return v.backend.InConfig(key)
}

// SetDefault sets the default value for this key.
func (v *Viper) SetDefault(key string, value interface{}) {
	v.backend.SetDefault(key, value)
}

// Set sets the value for the key in the override register.
func (v *Viper) Set(key string, value interface{}) {
	v.backend.Set(key, value)
}

// NewViper creates a new Viper backend.
func NewViper(backend *viper.Viper) Config {
	return &Viper{backend}
}
