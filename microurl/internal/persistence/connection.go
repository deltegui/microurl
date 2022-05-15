package persistence

import (
	"log"
	"microurl/internal/config"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Connection struct {
	filename string
	db       *gorm.DB
}

func Connect(config config.Configuration) Connection {
	db, err := gorm.Open(sqlite.Open(config.DB.Conn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error while opening connection: %s.", err)
	}
	return Connection{
		filename: config.DB.Conn,
		db:       db,
	}
}

func (conn Connection) MigrateAll() {
	conn.db.AutoMigrate(User{})
}

func (conn Connection) Destroy() {
	if err := os.Remove(conn.filename); err != nil {
		log.Println("[ERROR] Cannot remove database file: ", conn.filename)
	}
}
