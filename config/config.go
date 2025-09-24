package config

import (
	"log/slog"

	"github.com/spf13/viper"
)

type Config struct {
	Outlets         []string `mapstructure:"outlets"`
	PollingInterval int64    `mapstructure:"polling_interval"`
	HTTPAddress     string   `mapstructure:"http_address"`
}

func ConfigFromFile() (*Config, error) {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}
	slog.Info("Outlets loaded", "count", len(conf.Outlets))
	return &conf, nil
}
