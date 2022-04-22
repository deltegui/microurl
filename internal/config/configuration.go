package config

import "github.com/deltegui/configloader"

type DBConfig struct {
	Driver string `configName:"dbdriver"`
	Conn   string `configName:"dbconn"`
}

type TLSConfig struct {
	CRT     string `configName:"tlscrt"`
	KEY     string `configName:"tlskey"`
	Enabled bool   `configName:"tlsenabled"`
}

type Configuration struct {
	ListenURL string `paramName:"url"`
	DB        DBConfig
	TLS       TLSConfig
}

func Load() Configuration {
	return *configloader.NewConfigLoaderFor(&Configuration{}).
		AddHook(configloader.CreateFileHook("./config.json")).
		AddHook(configloader.CreateParamsHook()).
		AddHook(configloader.CreateEnvHook()).
		Retrieve().(*Configuration)
}
