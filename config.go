package base

type Config struct {
	Mode string `mapstructure:"mode"`

	Server struct {
		Name string `mapstructure:"name"`
		Http struct {
			Address string `mapstructure:"address"`
			Prefix  string `mapstructure:"prefix"`
		} `mapstructure:"http"`
	} `mapstructure:"server"`

	Tracer struct {
		Enabled bool `mapstructure:"enabled"`
		Jaeger  struct {
			Endpoint string `mapstructure:"endpoint"`
			Active   bool   `mapstructure:"active"`
		} `mapstructure:"jaeger"`
	} `mapstructure:"tracer"`
}
