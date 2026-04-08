package config

import (
	"fmt"
)

type SrvConfig struct {
	//Logger   *core.LogConfig		`yaml:"logger"`
	Handlers HandlerList `yaml:"handlers"`
}

func (cfg *SrvConfig) String() string {
	str := fmt.Sprintf("Server app config: %s",
		cfg.Handlers)
	return str
}

func GetDefaultSrvConfig() *SrvConfig {
	appCfg := &SrvConfig{}
	return appCfg
}

func GetSrvConfig(appPar *AppParams) (error, *SrvConfig) {
	appCfg := GetDefaultSrvConfig()
	return nil, appCfg
}
