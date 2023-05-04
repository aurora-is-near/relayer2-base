package log

const (
	defaultLogFilePath = "/tmp/relayer/log/relayer.log"
	defaultLogLevel    = "info"
	defaultLogToFile   = true
	defaultLogToStdOut = true
)

type Config struct {
	LogToFile    bool   `mapstructure:"logToFile"`
	LogToConsole bool   `mapstructure:"logToConsole"`
	Level        string `mapstructure:"level"`
	FilePath     string `mapstructure:"filePath"`
}

func DefaultConfig() *Config {
	return &Config{
		LogToFile:    defaultLogToFile,
		LogToConsole: defaultLogToStdOut,
		Level:        defaultLogLevel,
		FilePath:     defaultLogFilePath,
	}
}
