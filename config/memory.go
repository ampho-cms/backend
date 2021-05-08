// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package config

import (
	"reflect"
	"time"

	"github.com/spf13/cast"
)

// Memory is the in memory config backend.
type Memory struct {
	defaults map[string]interface{}
	storage  map[string]interface{}
}

// Backend returns backend configuration engine.
func (m *Memory) Backend() interface{} {
	return m.storage
}

// Get can retrieve any value given the key to use.
func (m *Memory) Get(key string) interface{} {
	if v := m.storage[key]; v != nil {
		return v
	}

	return m.defaults[key]
}

// Sub returns new Config instance representing a sub tree of this instance.
func (m *Memory) Sub(key string) Config {
	sub := NewMemory()
	data := m.Get(key)
	if data == nil {
		return nil
	}

	if reflect.TypeOf(data).Kind() == reflect.Map {
		sub.storage = cast.ToStringMap(data)
		return sub
	}

	return nil
}

// GetString returns the value associated with the key as a string.
func (m *Memory) GetString(key string) string {
	return cast.ToString(m.Get(key))
}

// GetBool returns the value associated with the key as a boolean.
func (m *Memory) GetBool(key string) bool {
	return cast.ToBool(m.Get(key))
}

// GetInt returns the value associated with the key as an integer.
func (m *Memory) GetInt(key string) int {
	return cast.ToInt(m.Get(key))
}

// GetFloat returns the value associated with the key as a float64.
func (m *Memory) GetFloat(key string) float64 {
	return cast.ToFloat64(m.Get(key))
}

// GetTime returns the value associated with the key as time.
func (m *Memory) GetTime(key string) time.Time {
	return cast.ToTime(m.Get(key))
}

// GetDuration returns the value associated with the key as a duration.
func (m *Memory) GetDuration(key string) time.Duration {
	return cast.ToDuration(m.Get(key))
}

// GetIntSlice returns the value associated with the key as a slice of int values.
func (m *Memory) GetIntSlice(key string) []int {
	return cast.ToIntSlice(m.Get(key))
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (m *Memory) GetStringSlice(key string) []string {
	return cast.ToStringSlice(m.Get(key))
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (m *Memory) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(m.Get(key))
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (m *Memory) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(m.Get(key))
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (m *Memory) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(m.Get(key))
}

// IsSet checks to see if the key has been set in any of the data locations.
func (m *Memory) IsSet(key string) bool {
	return m.Get(key) != nil
}

// InConfig checks to see if the given key (or an alias) is in the config file.
func (m *Memory) InConfig(key string) bool {
	if _, ok := m.storage[key]; ok {
		return true
	}

	if _, ok := m.defaults[key]; ok {
		return true
	}

	return false
}

// SetDefault sets the default value for this key.
func (m *Memory) SetDefault(key string, value interface{}) {
	m.defaults[key] = value
}

// Set sets the value for the key in the override register.
func (m *Memory) Set(key string, value interface{}) {
	m.storage[key] = value
}

// NewMemory creates a new Memory configurator backend.
func NewMemory() *Memory {
	return &Memory{}
}
