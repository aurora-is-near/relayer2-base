package endpoint

type Config struct {
	DisabledEndpoints []string `mapstructure:"disabledEndpoints"`
}

func DefaultConfig() Config {
	return Config{DisabledEndpoints: []string{}}
}
