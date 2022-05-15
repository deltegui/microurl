package testutils

import (
	"microurl/internal/config"
	"microurl/internal/persistence"

	"github.com/deltegui/configloader"
)

func loadTestConfig() config.Configuration {
	return *configloader.NewConfigLoaderFor(&config.Configuration{}).
		AddHook(configloader.CreateFileHook("./test.json")).
		Retrieve().(*config.Configuration)
}

func DBTransaction(test func(conn persistence.Connection)) {
	conf := loadTestConfig()
	conn := persistence.Connect(conf)
	conn.MigrateAll()
	test(conn)
	conn.Destroy()
}
