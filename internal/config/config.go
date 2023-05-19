package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	ShutdownTimeout time.Duration  `mapstructure:"shutdown_timeout"`
	Targets         []TargetConfig `mapstructure:"targets" validate:"required"`
}

type TargetConfig struct {
	Path          string   `mapstructure:"path" validate:"required"`
	IncludeRegexp []string `mapstructure:"include_regexp"`
	ExcludeRegexp []string `mapstructure:"exclude_regexp"`
	Commands      []string `mapstructure:"commands"`
	LogFile       string   `mapstructure:"log_file" validate:"filepath"`
}

func NewConfig(configFilePath, configFileName string) (Config, error) {
	myViper := viper.New()
	myViper.AddConfigPath(configFilePath)
	myViper.SetConfigName(configFileName)

	if err := myViper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var result Config
	if err := myViper.Unmarshal(&result); err != nil {
		return Config{}, err
	}

	validate := validator.New()
	if err := validate.Struct(&result); err != nil {
		return Config{}, err
	}

	return result, nil
}
