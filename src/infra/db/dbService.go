package db

import (
	"errors"

	"github.com/glebarez/sqlite"
	dbModel "github.com/speedianet/sfm/src/infra/db/model"
	"gorm.io/gorm"
)

func DatabaseService() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("/var/speedia/sfm.db"), &gorm.Config{})
	if err != nil {
		return nil, errors.New("DatabaseConnectionError")
	}

	err = db.AutoMigrate(
		&dbModel.Account{},
		&dbModel.AccountQuota{},
		&dbModel.AccountQuotaUsage{},
		&dbModel.ResourceProfile{},
	)
	if err != nil {
		return nil, errors.New("DatabaseMigrationError")
	}

	return db, nil
}
