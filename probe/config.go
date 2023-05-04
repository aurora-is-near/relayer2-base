package probe

const (
	defaultAddress   = "127.0.0.1:28080"
	defaultNamespace = "Aurora"
	defaultSubsystem = "Relayer2"
	defaultEnabled   = false
)

type ServerConfig struct {
	Address   string `mapstructure:"address"`
	Namespace string `mapstructure:"namespace"`
	Subsystem string `mapstructure:"subsystem"`
}

type MetricConfig struct {
	Id          string    `mapstructure:"id"`
	Type        string    `mapstructure:"type"`
	Name        string    `mapstructure:"name"`
	Help        string    `mapstructure:"help"`
	Buckets     []float64 `mapstructure:"buckets"`
	LabelNames  []string  `mapstructure:"labelNames"`
	LabelValues []string  `mapstructure:"labelValues"`
}

type Config struct {
	Enabled       bool            `mapstructure:"enable"`
	ServerConfig  *ServerConfig   `mapstructure:"server"`
	MetricConfigs *[]MetricConfig `mapstructure:"metrics"`
}

func DefaultConfig() *Config {
	return &Config{
		Enabled: defaultEnabled,
		ServerConfig: &ServerConfig{
			Address:   defaultAddress,
			Namespace: defaultNamespace,
			Subsystem: defaultSubsystem,
		},
		MetricConfigs: &[]MetricConfig{},
	}
}
