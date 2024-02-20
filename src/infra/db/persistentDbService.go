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
	Orm *gorm.DB
}

func NewPersistentDatabaseService() (*PersistentDatabaseService, error) {
	ormSvc, err := gorm.Open(sqlite.Open(DatabaseFilePath), &gorm.Config{})
	if err != nil {
		return nil, errors.New("DatabaseConnectionError")
	}

	persistDbSvc := &PersistentDatabaseService{Orm: ormSvc}
	err = persistDbSvc.dbMigrate()
	if err != nil {
		return nil, err
	}

	return persistDbSvc, nil
}

func (persistDbSvc PersistentDatabaseService) isTableEmpty(model interface{}) (bool, error) {
	var count int64
	err := persistDbSvc.Orm.Model(&model).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (persistDbSvc PersistentDatabaseService) seedDatabase(seedModels ...interface{}) error {
	for _, seedModel := range seedModels {
		isTableEmpty, err := persistDbSvc.isTableEmpty(seedModel)
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

			entryFormatOrmWillAccept := reflect.New(seedModelType)
			entryFormatOrmWillAccept.Elem().Set(entryInnerStructure)
			adjustedEntry := entryFormatOrmWillAccept.Interface()

			err = persistDbSvc.Orm.Create(adjustedEntry).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (persistDbSvc PersistentDatabaseService) dbMigrate() error {
	err := persistDbSvc.Orm.AutoMigrate(
		&dbModel.Account{},
		&dbModel.AccountQuota{},
		&dbModel.AccountQuotaUsage{},
		&dbModel.Container{},
		&dbModel.ContainerPortBinding{},
		&dbModel.ContainerProfile{},
		&dbModel.Mapping{},
		&dbModel.MappingTarget{},
		&dbModel.LicenseInfo{},
	)
	if err != nil {
		return errors.New("DatabaseMigrationError")
	}

	modelsWithInitialEntries := []interface{}{
		&dbModel.ContainerProfile{},
	}

	err = persistDbSvc.seedDatabase(modelsWithInitialEntries...)
	if err != nil {
		return errors.New("AddDefaultDatabaseEntriesError")
	}

	return nil
}
