package config

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

func GetDefaultAppConfig(devCfg DeviceConfig) AppConfig {
	appCfg := AppConfig{
		Device: devCfg,
	}
	return appCfg
}
