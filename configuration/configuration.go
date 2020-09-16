package configuration

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	Host string `yaml:"host"`
}

func New(f string) (*Configuration, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetConfigFile(f)

	err := viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	var cfg Configuration

	err = viper.Unmarshal(&cfg)

	return &cfg, err
}
