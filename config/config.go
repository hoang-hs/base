package config

type Config struct {
	Mode string `mapstructure:"mode"`

	Server struct {
		Name string `mapstructure:"name"`
		Http struct {
			Address string `mapstructure:"address"`
			Prefix  string `mapstructure:"prefix"`
		} `mapstructure:"http"`
	} `mapstructure:"server"`

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

	Tracer struct {
		Enabled bool `mapstructure:"enabled"`
		Jaeger  struct {
			Endpoint string `mapstructure:"endpoint"`
			Active   bool   `mapstructure:"active"`
		} `mapstructure:"jaeger"`
	} `mapstructure:"tracer"`
}

var config *Config

func Get() *Config {
	return config
}
