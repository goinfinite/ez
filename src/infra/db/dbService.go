package db

import (
	"errors"
	"reflect"

	"github.com/glebarez/sqlite"
	dbModel "github.com/goinfinite/fleet/src/infra/db/model"
	"gorm.io/gorm"
)

type DatabaseService struct {
	Orm *gorm.DB
}

func NewDatabaseService() (*DatabaseService, error) {
	ormSvc, err := gorm.Open(sqlite.Open("/var/infinite/fleet.db"), &gorm.Config{})
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

		seedModelDefaultEntryMethod := seedModelFieldsAndMethods.MethodByName(
			"DefaultEntry",
		)
		seedModelDefaultEntryMethodResults := seedModelDefaultEntryMethod.Call(
			[]reflect.Value{},
		)
		firstAndOnlyResult := seedModelDefaultEntryMethodResults[0].Interface()
		defaultEntryInnerStructure := reflect.ValueOf(firstAndOnlyResult)

		defaultEntryFormatOrmWillAccept := reflect.New(seedModelType)
		defaultEntryFormatOrmWillAccept.Elem().Set(defaultEntryInnerStructure)
		adjustedDefaultEntry := defaultEntryFormatOrmWillAccept.Interface()

		err = dbSvc.Orm.Create(adjustedDefaultEntry).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (dbSvc DatabaseService) dbMigrate() error {
	err := dbSvc.Orm.AutoMigrate(
		&dbModel.Account{},
		&dbModel.AccountQuota{},
		&dbModel.AccountQuotaUsage{},
		&dbModel.ContainerProfile{},
	)
	if err != nil {
		return errors.New("DatabaseMigrationError")
	}

	modelsWithDefaultEntries := []interface{}{
		&dbModel.ContainerProfile{},
	}

	err = dbSvc.seedDatabase(modelsWithDefaultEntries...)
	if err != nil {
		return errors.New("AddDefaultDatabaseEntriesError")
	}

	return nil
}
