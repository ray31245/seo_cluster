package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

// note: dsn is "publish_manager.db"
func NewDB(dsn string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("NewDB: %w", err)
	}

	return &DB{db: db}, nil
}

func (d *DB) Close() error {
	db, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("Close: %w", err)
	}

	err = db.Close()
	if err != nil {
		return fmt.Errorf("Close: %w", err)
	}

	return nil
}
