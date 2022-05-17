package config

import "github.com/deltegui/configloader/v2"

type DBConfig struct {
	Driver string `configName:"driver"`
	Conn   string `configName:"conn"`
}

type TLSConfig struct {
	CRT     string `configName:"crt"`
	KEY     string `configName:"key"`
	Enabled bool   `configName:"enabled"`
}

type Configuration struct {
	ListenURL  string    `configName:"url"`
	JWTKey     string    `configName:"jwt"`
	SessionKey string    `configName:"session"`
	DB         DBConfig  `configPrefix:"db"`
	TLS        TLSConfig `configPrefix:"tls"`
}

func Load() Configuration {
	return *configloader.NewConfigLoaderFor(&Configuration{}).
		AddHook(configloader.CreateFileHook("./config.json")).
		AddHook(configloader.CreateParamsHook()).
		AddHook(configloader.CreateEnvHook()).
		Retrieve().(*Configuration)
}
