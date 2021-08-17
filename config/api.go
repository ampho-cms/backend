// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// DefaultSearchPaths returns default search directory paths for using in New.
//
// Default paths slice include:
//     - $HOME;
//     - directory where executable is located;
//     - current working directory.
func DefaultSearchPaths() []string {
	r := []string{"$HOME"}

	if execDir, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		r = append(r, execDir)
	}

	r = append(r, ".")

	return r
}

// New creates a new configuration instance using default parameters.
//
// Returned configuration uses files and environment variables for getting config data.
func New(name, typ string, searchPath ...string) (Config, error) {
	vp := viper.New()

	vp.SetEnvPrefix(name)

	vpName := name
	if typ != "" {
		vp.SetConfigType(typ)
		vpName += "." + typ
	}
	vp.SetConfigName(vpName)

	for _, p := range searchPath {
		vp.AddConfigPath(p)
	}

	if err := vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	return &Viper{name, vp}, nil
}

// NewTesting creates a new configuration instance using default parameters suitable for using in tests.
//
// This configuration supports only direct and ENV settings.
func NewTesting(name string) Config {
	vp := viper.New()
	vp.SetEnvPrefix(name)

	return &Viper{name, vp}
}
