// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package config

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// NewDefault creates a default configurator.
func NewDefault(name string, defaults map[string]interface{}) (Config, error) {
	// Executable directory path
	execDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}

	// Configuration engine
	vp := viper.New()
	vp.SetConfigName(name)
	vp.SetConfigType("yaml")
	vp.AddConfigPath("$HOME/." + name)
	vp.AddConfigPath(execDir)
	vp.AddConfigPath(".")
	cfg := NewViper(vp)

	if defaults != nil {
		for k, v := range defaults {
			cfg.SetDefault(k, v)
		}
	}

	// Load configuration
	if err := vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	return cfg, nil
}
