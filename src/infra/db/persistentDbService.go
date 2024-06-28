package db

import (
	"errors"
	"reflect"

	"github.com/glebarez/sqlite"
	dbModel "github.com/speedianet/control/src/infra/db/model"
	"gorm.io/gorm"
)

const (
	DatabaseFilePath = "/var/speedia/control.db"
)

type PersistentDatabaseService struct {
	Handler *gorm.DB
}

func NewPersistentDatabaseService() (*PersistentDatabaseService, error) {
	ormSvc, err := gorm.Open(
		sqlite.Open(DatabaseFilePath),
		&gorm.Config{},
	)
	if err != nil {
		return nil, errors.New("DatabaseConnectionError")
	}

	dbSvc := &PersistentDatabaseService{Handler: ormSvc}
	err = dbSvc.dbMigrate()
	if err != nil {
		return nil, err
	}

	return dbSvc, nil
}

func (dbSvc PersistentDatabaseService) isTableEmpty(model interface{}) (bool, error) {
	var count int64
	err := dbSvc.Handler.Model(&model).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (dbSvc PersistentDatabaseService) seedDatabase(seedModels ...interface{}) error {
	for _, seedModel := range seedModels {
		isTableEmpty, err := dbSvc.isTableEmpty(seedModel)
		if err != nil {
			return err
		}

		if !isTableEmpty {
			continue
		}

		seedModelType := reflect.TypeOf(seedModel).Elem()
		seedModelFieldsAndMethods := reflect.ValueOf(seedModel)

		seedModelInitialEntriesMethod := seedModelFieldsAndMethods.MethodByName(
			"InitialEntries",
		)
		seedModelInitialEntriesMethodResults := seedModelInitialEntriesMethod.Call(
			[]reflect.Value{},
		)
		initialEntries := seedModelInitialEntriesMethodResults[0].Interface()

		for _, entry := range initialEntries.([]interface{}) {
			entryInnerStructure := reflect.ValueOf(entry)

			entryFormatHandlerWillAccept := reflect.New(seedModelType)
			entryFormatHandlerWillAccept.Elem().Set(entryInnerStructure)
			adjustedEntry := entryFormatHandlerWillAccept.Interface()

			err = dbSvc.Handler.Create(adjustedEntry).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (dbSvc PersistentDatabaseService) dbMigrate() error {
	err := dbSvc.Handler.AutoMigrate(
		&dbModel.Account{},
		&dbModel.AccountQuota{},
		&dbModel.AccountQuotaUsage{},
		&dbModel.Container{},
		&dbModel.ContainerPortBinding{},
		&dbModel.ContainerProfile{},
		&dbModel.ContainerProxy{},
		&dbModel.Mapping{},
		&dbModel.MappingTarget{},
		&dbModel.LicenseInfo{},
		&dbModel.ScheduledTask{},
		&dbModel.SecurityEvent{},
	)
	if err != nil {
		return errors.New("DatabaseMigrationError")
	}

	modelsWithInitialEntries := []interface{}{
		&dbModel.ContainerProfile{},
	}

	err = dbSvc.seedDatabase(modelsWithInitialEntries...)
	if err != nil {
		return errors.New("AddDefaultDatabaseEntriesError")
	}

	return nil
}
