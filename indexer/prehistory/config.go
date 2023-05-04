package prehistory

const (
	defaultIndexFromPrehistory = false
	defaultFromBlock           = 0
	defaultBatchSize           = 10000
	defaultPrehistoryChainId   = 1313161554
)

type Config struct {
	IndexFromPrehistory bool   `mapstructure:"indexFromPrehistory"`
	PrehistoryHeight    uint64 `mapstructure:"prehistoryHeight"`
	From                uint64 `mapstructure:"from"`
	To                  uint64 `mapstructure:"to"`
	ArchiveURL          string `mapstructure:"archiveURL"`
	BatchSize           uint64 `mapstructure:"batchSize"`
	PrehistoryChainId   uint64 `mapstructure:"prehistoryChainId"`
}

func DefaultConfig() *Config {
	return &Config{
		IndexFromPrehistory: defaultIndexFromPrehistory,
		From:                defaultFromBlock,
		BatchSize:           defaultBatchSize,
		PrehistoryChainId:   defaultPrehistoryChainId,
	}
}
