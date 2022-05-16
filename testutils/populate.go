package testutils

import (
	"microurl/internal"
	"microurl/internal/config"
	"microurl/internal/persistence"
)

func loadTestConfig() config.Configuration {
	return config.Configuration{
		ListenURL:  "localhost:3000",
		JWTKey:     "blablamykeyblabla",
		SessionKey: "blablamykeyblabla",
		DB: config.DBConfig{
			Driver: "sqlite",
			Conn:   ":memory:",
		},
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
}

var users = []internal.User{
	{
		Name:     "manolo",
		Password: "$2a$10$CgHCJHQNLHNlqFD5zy0dJOH1XTMLLmi4DPB6rd1vnEwwFGrcH/1QO",
	},
	{
		Name:     "ambrosio",
		Password: "$2a$10$.k95rpYLCwSI3gqJzgHgG.xuuyzyBQHfgQfqjTfB/JdzuhRTTb.S2",
	},
	{
		Name:     "paola",
		Password: "$2a$10$Ep1jgVKZE7GB.xJzdryxXOOEuI0lM8W2A18cxxmkhYIS.hfCAdL4G",
	},
}

var urls = []internal.URL{
	{
		Name:     "My youtube",
		Original: "http://youtube.com/hola",
		Owner:    "manolo",
	},
	{
		Name:     "Hola!",
		Original: "https://hello.com/hola",
		Owner:    "manolo",
	},
	{
		Name:     "A video that I liked",
		Original: "https://youtube.com/xasfhasd",
		Owner:    "paola",
	},
	{
		Name:     "My web",
		Original: "https://manolo.com/manolo",
		Owner:    "ambrosio",
	},
}

type Populator struct {
	UserRepo   persistence.GormUserRepository
	URLRepo    persistence.GormURLRepository
	populators []func()
}

func newPopulator(conn persistence.Connection) Populator {
	urlRepo := persistence.NewGormURLRepository(conn)
	userRepo := persistence.NewGormUserRepository(conn)
	p := Populator{
		UserRepo: userRepo,
		URLRepo:  urlRepo,
	}
	p.populators = []func(){
		p.PopulateUsers,
		p.PopulateURLs,
	}
	return p
}

func (populator Populator) PopulateAll() {
	for _, populate := range populator.populators {
		populate()
	}
}

func (p Populator) PopulateUsers() {
	populate(users, p.UserRepo.Save)
}

func (p Populator) PopulateURLs() {
	populate(urls, func(url internal.URL) error { return p.URLRepo.Save(&url) })
}

func populate[T any](elements []T, save func(T) error) {
	for _, e := range elements {
		if err := save(e); err != nil {
			panic(err)
		}
	}
}

func DBTransaction(test func(Populator)) {
	conf := loadTestConfig()
	conn := persistence.Connect(conf)
	conn.MigrateAll()
	test(newPopulator(conn))
}
