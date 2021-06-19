// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package config

import (
	"github.com/spf13/viper"
)

// NewDefault creates a default config.
func NewDefault(name, typ string, searchPaths []string, defaults map[string]interface{}) (Config, error) {

	vp := viper.New()
	vp.SetConfigName(name)

	if typ != "" {
		vp.SetConfigType(typ)
	}

	if len(searchPaths) > 0 {
		for _, p := range searchPaths {
			vp.AddConfigPath(p)
		}
	}

	if defaults != nil {
		for k, v := range defaults {
			vp.SetDefault(k, v)
		}
	}

	if err := vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}


	return NewViper(vp), nil
}
