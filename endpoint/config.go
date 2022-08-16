package endpoint

type Config struct {
	DisabledEndpoints []string
}

func DefaultConfig() Config {
	return Config{DisabledEndpoints: []string{}}
}
