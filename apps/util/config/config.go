package config

import (
	"strings"

	"github.com/spf13/viper"

	"ddv_loc/pkg/translator/config"
)

type (
	AppConfig struct {
		Translator config.Config `mapstructure:"translator"`
	}
)

func New() AppConfig {
	var appConfig AppConfig

	viper.SetConfigFile("etc/config.yml")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		panic("Error reading config file")
	}

	if err := viper.Unmarshal(&appConfig); err != nil {
		panic("Error unmarshal config file")
	}

	return appConfig
}
