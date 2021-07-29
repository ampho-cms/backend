// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// DefaultSearchPaths return default search paths.
func DefaultSearchPaths(name string) []string {
	r := []string{"$HOME/." + name}

	if execDir, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		r = append(r, execDir)
	}

	r = append(r, ".")

	return r
}

// New creates a new "default" config instance.
func New(name, typ string, searchPath ...string) (Config, error) {
	vp := viper.New()
	vpName := name

	if typ != "" {
		vp.SetConfigType(typ)
		vpName += "." + typ
	}

	vp.SetConfigName(vpName)

	if len(searchPath) > 0 {
		for _, p := range searchPath {
			vp.AddConfigPath(p)
		}
	}

	if err := vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	return NewViper(name, vp), nil
}

// NewTesting creates a new "default" config instance suitable for usage in unit tests.
func NewTesting(name string) Config {
	return NewViper(name, viper.New())
}
