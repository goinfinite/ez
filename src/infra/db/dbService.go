package db

import (
	"errors"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func DatabaseService() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("/var/speedia/sfm.db"), &gorm.Config{})
	if err != nil {
		return nil, errors.New("DatabaseConnectionError")
	}

	return db, nil
}
