package config

import (
	"fmt"
)

type AppConfig struct {
	Device *DeviceConfig `yaml:"device"`
}

func (cfg *AppConfig) String() string {
	str := fmt.Sprintf("Client app config: %s",
		cfg.Device)
	return str
}

func GetDefaultAppConfig(devCfg *DeviceConfig) *AppConfig {
	appCfg := &AppConfig{
		Device: devCfg,
	}
	return appCfg
}

func GetAppConfig(appPar *AppParams, devCfg *DeviceConfig) (error, *AppConfig) {
	appCfg := GetDefaultAppConfig(devCfg)
	return nil, appCfg
}
