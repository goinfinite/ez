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

type DatabaseService struct {
	Orm *gorm.DB
}

func NewDatabaseService() (*DatabaseService, error) {
	ormSvc, err := gorm.Open(sqlite.Open(DatabaseFilePath), &gorm.Config{})
	if err != nil {
		return nil, errors.New("DatabaseConnectionError")
	}

	dbSvc := &DatabaseService{Orm: ormSvc}
	err = dbSvc.dbMigrate()
	if err != nil {
		return nil, err
	}

	return dbSvc, nil
}

func (dbSvc DatabaseService) isTableEmpty(model interface{}) (bool, error) {
	var count int64
	err := dbSvc.Orm.Model(&model).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (dbSvc DatabaseService) seedDatabase(seedModels ...interface{}) error {
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

			entryFormatOrmWillAccept := reflect.New(seedModelType)
			entryFormatOrmWillAccept.Elem().Set(entryInnerStructure)
			adjustedEntry := entryFormatOrmWillAccept.Interface()

			err = dbSvc.Orm.Create(adjustedEntry).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (dbSvc DatabaseService) dbMigrate() error {
	err := dbSvc.Orm.AutoMigrate(
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

	err = dbSvc.seedDatabase(modelsWithInitialEntries...)
	if err != nil {
		return errors.New("AddDefaultDatabaseEntriesError")
	}

	return nil
}
