package config

import (
	"fmt"
)

type AppConfig struct {
	Logger LoggerConfig `yaml:"logger"`
	Device DeviceConfig `yaml:"device"`
}

type LoggerConfig struct {
	Level      string `yaml:"level"`       // One of: debug, info, warn, error
	MaxSize    int    `yaml:"max_size"`    // Megabytes before rotating
	MaxBackups int    `yaml:"max_backups"` // Max number of old log files to retain
	MaxAge     int    `yaml:"max_age"`     // Days to retain old log files
	Compress   bool   `yaml:"compress"`    // Compress rotated files using gzip
}

func (cfg *AppConfig) String() string {
	str := fmt.Sprintf("Client app config: %s", cfg.Device)
	return str
}

func GetDefaultAppConfig(devCfg DeviceConfig) AppConfig {
	appCfg := AppConfig{
		Device: devCfg,
	}
	return appCfg
}

func GetAppConfig(appPar *AppParams, devCfg DeviceConfig) (error, AppConfig) {
	appCfg := GetDefaultAppConfig(devCfg)
	return nil, appCfg
}
