package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	Mode string `mapstructure:"mode"`

	Server struct {
		Name string `mapstructure:"name"`
		Http struct {
			Address string `mapstructure:"address"`
			Prefix  string `mapstructure:"prefix"`
		} `mapstructure:"http"`
	} `mapstructure:"server"`
	Swagger struct {
		Enabled bool `mapstructure:"enabled"`
	} `mapstructure:"swagger"`

	Postgresql struct {
		Host        string `mapstructure:"host"`
		Port        string `mapstructure:"port"`
		User        string `mapstructure:"user"`
		DbName      string `mapstructure:"db_name"`
		SslMode     string `mapstructure:"ssl_mode"`
		Password    string `mapstructure:"password"`
		AutoMigrate bool   `mapstructure:"auto_migrate"`
		MaxLifeTime int    `mapstructure:"max_life_time"`
	}
	Mongo struct {
		Uri string `mapstructure:"uri"`
		DB  string `mapstructure:"db"`
	}

	Redis struct {
		Hosts    []string `mapstructure:"hosts"`
		Username string   `mapstructure:"username"`
		Password string   `mapstructure:"password"`
	} `mapstructure:"redis"`

	Kafka   *Kafka   `mapstructure:"kafka"`
	Observe *Observe `mapstructure:"observe"`
}

type Kafka struct {
	Host     string `mapstructure:"host"`
	Consumer struct {
		GroupID string `mapstructure:"group_id"`
		Topic   string `mapstructure:"topic"`
	} `mapstructure:"consumer"`
}

type Observe struct {
	Trace struct {
		Enabled      bool `mapstructure:"enabled"`
		OtlpExporter struct {
			Endpoint string `mapstructure:"endpoint"`
			Active   bool   `mapstructure:"active"`
		} `mapstructure:"otlp_exporter"`
		SampleRate float64 `mapstructure:"sample_rate"`
	} `mapstructure:"trace"`
	Metric struct {
		Enabled bool `mapstructure:"enabled"`
	} `mapstructure:"metric"`
}

var common *Config

func Get() *Config {
	return common
}

func LoadConfig(pathConfig string) error {
	viper.SetConfigFile(pathConfig)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&common)

	return nil
}
