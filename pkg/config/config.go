// Package config provides a configuration loading and validation functionality.
package config

import (
	"bytes"
	"errors"
	"strings"

	"github.com/creasty/defaults"
	validator "github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// LoadConfig loads the configuration from a os environment variables.
func LoadConfig[T any](c T) (T, error) {
	err := SetDefault(c)
	if err != nil {
		return c, err
	}

	bs, e := yaml.Marshal(c)
	if e != nil {
		return c, errors.New("error marshalling config")
	}

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if e = viper.ReadConfig(bytes.NewBuffer(bs)); e != nil {
		return c, errors.New("error reading config")
	}

	bindKeys()

	e = viper.Unmarshal(c)
	if e != nil {
		return c, errors.New("error unmarshalling config")
	}

	validator := validator.New()
	if e = validator.Struct(c); e != nil {
		return c, e
	}

	return c, nil
}

// Viper doesn't marshall environment variables automatically. It requires `bind` to be called for every param.
// We are binding the keys explicitly here.
// Issue - https://github.com/spf13/viper/issues/522
func bindKeys() {
	for _, key := range viper.AllKeys() {
		_ = viper.BindEnv(key, strings.ToUpper(strings.ReplaceAll(key, ".", "_")))
	}
}

// SetDefault sets default values for the struct fields.
func SetDefault[T any](c T) error {
	if err := defaults.Set(c); err != nil {
		return err
	}

	return nil
}
