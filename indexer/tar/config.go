package tar

const (
	DefaultIndexFromBackup = false
)

type Config struct {
	IndexFromBackup bool   `mapstructure:"indexFromBackup"`
	Dir             string `mapstructure:"backupDir"`
	NamePrefix      string `mapstructure:"backupNamePrefix"`
	From            uint64 `mapstructure:"from"`
	To              uint64 `mapstructure:"to"`
}

func DefaultConfig() *Config {
	return &Config{
		IndexFromBackup: DefaultIndexFromBackup,
	}
}
