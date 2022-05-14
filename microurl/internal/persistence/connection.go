package persistence

import (
	"log"
	"microurl/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Connection struct {
	db *gorm.DB
}

func Connect(config config.Configuration) Connection {
	db, err := gorm.Open(sqlite.Open(config.DB.Conn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error while opening connection: %s.", err)
	}
	return Connection{db}
}

func (conn Connection) MigrateAll() {
	conn.db.AutoMigrate(User{})
}
