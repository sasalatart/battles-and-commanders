package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Setup initializes viper configs
func Setup() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	if err := viper.ReadInConfig(); err != nil {
		panic(errors.Wrap(err, "Reading viper config"))
	}
}
