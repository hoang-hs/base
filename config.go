package base

import (
	"github.com/spf13/viper"
)

type Config struct {
	Mode string `mapstructure:"mode"`
	Name string `mapstructure:"name"`

	Tracer struct {
		Enabled bool `mapstructure:"enabled"`
		Jaeger  struct {
			Endpoint string `mapstructure:"endpoint"`
			Active   bool   `mapstructure:"active"`
		} `mapstructure:"jaeger"`
	} `mapstructure:"tracer"`
}

var common *Config

func Get() *Config {
	return common
}

func LoadConfig(pathConfig string) {
	viper.SetConfigFile(pathConfig)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&common)
	if err != nil {
		panic(err)
	}
}
