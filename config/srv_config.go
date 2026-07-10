package config

type SrvConfig struct {
	//Logger   *core.LogConfig		`yaml:"logger"`
	Handlers HandlerList `yaml:"handlers"`
}

func GetDefaultSrvConfig() *SrvConfig {
	appCfg := &SrvConfig{}
	return appCfg
}
