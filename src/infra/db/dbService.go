package db

import (
	"errors"
	"reflect"

	"github.com/glebarez/sqlite"
	dbModel "github.com/speedianet/sfm/src/infra/db/model"
	"gorm.io/gorm"
)

func isTableEmpty(dbSvc *gorm.DB, model interface{}) (bool, error) {
	var count int64
	err := dbSvc.Model(&model).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func seedDatabase(dbSvc *gorm.DB, seedModels ...interface{}) error {
	for _, seedModel := range seedModels {
		isTableEmpty, err := isTableEmpty(dbSvc, seedModel)
		if err != nil {
			return err
		}

		if !isTableEmpty {
			continue
		}

		seedModelType := reflect.TypeOf(seedModel).Elem()
		seedModelFieldsAndMethods := reflect.ValueOf(seedModel)

		defaultEntryMethod := seedModelFieldsAndMethods.MethodByName("DefaultEntry")
		defaultEntryMethodResults := defaultEntryMethod.Call([]reflect.Value{})
		actualDefaultEntry := defaultEntryMethodResults[0].Interface()
		defaultEntryInnerStructure := reflect.ValueOf(actualDefaultEntry)

		defaultEntryFormatOrmWillAccept := reflect.New(seedModelType)
		defaultEntryFormatOrmWillAccept.Elem().Set(defaultEntryInnerStructure)
		newDefaultEntry := defaultEntryFormatOrmWillAccept.Interface()

		err = dbSvc.Create(newDefaultEntry).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func DatabaseService() (*gorm.DB, error) {
	dbSvc, err := gorm.Open(sqlite.Open("/var/speedia/sfm.db"), &gorm.Config{})
	if err != nil {
		return nil, errors.New("DatabaseConnectionError")
	}

	err = dbSvc.AutoMigrate(
		&dbModel.Account{},
		&dbModel.AccountQuota{},
		&dbModel.AccountQuotaUsage{},
		&dbModel.ContainerProfile{},
	)
	if err != nil {
		return nil, errors.New("DatabaseMigrationError")
	}

	modelsWithDefaultEntries := []interface{}{
		&dbModel.ContainerProfile{},
	}

	err = seedDatabase(dbSvc, modelsWithDefaultEntries...)
	if err != nil {
		return nil, errors.New("AddDefaultDatabaseEntriesError")
	}

	return dbSvc, nil
}
