package config

import (
	"github.com/spf13/viper"
)

type DbSettings struct {
	MigrationConnectionString string `mapstructure:"MigrationConnectionString"`
	ConnectionString          string `mapstructure:"ConnectionString"`
}

type Config struct {
	DbSettings DbSettings `mapstructure:"DbSettings"`
}

func LoadConfig(environment string) (*Config, error) {
	viper.SetConfigName("appsettings." + environment)
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
