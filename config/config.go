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
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config") // For running Go Delve from cmd/x/main.go

	if err := viper.ReadInConfig(); err != nil {
		panic(errors.Wrap(err, "Reading viper config"))
	}
	mustBindEnv("POSTGRES_HOST")
	mustBindEnv("POSTGRES_PORT")
	mustBindEnv("POSTGRES_PASS")
}

func mustBindEnv(key string) {
	if err := viper.BindEnv(key); err != nil {
		panic(errors.Wrapf(err, "Failed binding env var %q", key))
	}
}
