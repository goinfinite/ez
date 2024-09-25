package db

import (
	"errors"

	"github.com/glebarez/sqlite"
	dbModel "github.com/speedianet/control/src/infra/db/model"
	"gorm.io/gorm"
)

type TrailDatabaseService struct {
	Handler *gorm.DB
}

func NewTrailDatabaseService() (*TrailDatabaseService, error) {
	ormSvc, err := gorm.Open(
		sqlite.Open("/var/speedia/trail.db"),
		&gorm.Config{},
	)
	if err != nil {
		return nil, errors.New("DatabaseConnectionError")
	}

	dbSvc := &TrailDatabaseService{Handler: ormSvc}
	err = dbSvc.RunMigrations()
	if err != nil {
		return nil, err
	}

	return dbSvc, nil
}

func (service *TrailDatabaseService) RunMigrations() error {
	err := service.Handler.AutoMigrate(
		&dbModel.ActivityRecord{},
		&dbModel.ActivityRecordAffectedResource{},
	)
	if err != nil {
		return errors.New("TrailDatabaseMigrationError: " + err.Error())
	}

	return nil
}
