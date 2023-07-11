package log

import (
	"github.com/spf13/viper"
)

const (
	defaultLogFilePath = "/tmp/relayer/log/relayer.log"
	defaultLogLevel    = "info"
	defaultLogToFile   = true
	defaultLogToStdOut = true

	configPath = "logger"
)

type Config struct {
	LogToFile    bool   `mapstructure:"logToFile"`
	LogToConsole bool   `mapstructure:"logToConsole"`
	Level        string `mapstructure:"level"`
	FilePath     string `mapstructure:"filePath"`
}

func defaultConfig() *Config {
	return &Config{
		LogToFile:    defaultLogToFile,
		LogToConsole: defaultLogToStdOut,
		Level:        defaultLogLevel,
		FilePath:     defaultLogFilePath,
	}
}

func GetConfig() *Config {
	config := defaultConfig()
	sub := viper.Sub(configPath)
	if sub != nil {
		// cmd.BindSubViper(sub, configPath)
		_ = sub.Unmarshal(&config)
	}
	return config
}
