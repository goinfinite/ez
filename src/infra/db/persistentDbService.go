package db

import (
	"errors"
	"reflect"

	"github.com/glebarez/sqlite"
	dbModel "github.com/goinfinite/ez/src/infra/db/model"
	"gorm.io/gorm"
)

const (
	PersistentDatabaseFilePath = "/var/infinite/ez.db"
)

type PersistentDatabaseService struct {
	Handler *gorm.DB
}

func NewPersistentDatabaseService() (*PersistentDatabaseService, error) {
	ormSvc, err := gorm.Open(
		sqlite.Open(PersistentDatabaseFilePath),
		&gorm.Config{},
	)
	if err != nil {
		return nil, errors.New("PersistentDatabaseConnectionError")
	}

	service := &PersistentDatabaseService{Handler: ormSvc}
	err = service.RunMigrations()
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (service *PersistentDatabaseService) isTableEmpty(model interface{}) (bool, error) {
	var count int64
	err := service.Handler.Model(&model).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (service *PersistentDatabaseService) seedDatabase(seedModels ...interface{}) error {
	for _, seedModel := range seedModels {
		isTableEmpty, err := service.isTableEmpty(seedModel)
		if err != nil {
			return err
		}

		if !isTableEmpty {
			continue
		}

		seedModelType := reflect.TypeOf(seedModel).Elem()
		seedModelFieldsAndMethods := reflect.ValueOf(seedModel)

		initialEntriesMethod := seedModelFieldsAndMethods.MethodByName(
			"InitialEntries",
		)
		initialEntriesMethodResults := initialEntriesMethod.Call(
			[]reflect.Value{},
		)
		initialEntries := initialEntriesMethodResults[0].Interface()

		for _, entry := range initialEntries.([]interface{}) {
			entryInnerStructure := reflect.ValueOf(entry)

			entryFormatHandlerWillAccept := reflect.New(seedModelType)
			entryFormatHandlerWillAccept.Elem().Set(entryInnerStructure)
			adjustedEntry := entryFormatHandlerWillAccept.Interface()

			err = service.Handler.Create(adjustedEntry).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (service *PersistentDatabaseService) RunMigrations() error {
	err := service.Handler.AutoMigrate(
		&dbModel.Account{},
		&dbModel.AccountQuota{},
		&dbModel.AccountQuotaUsage{},
		&dbModel.BackupDestination{},
		&dbModel.BackupJob{},
		&dbModel.BackupTask{},
		&dbModel.Container{},
		&dbModel.ContainerPortBinding{},
		&dbModel.ContainerProfile{},
		&dbModel.ContainerProxy{},
		&dbModel.Mapping{},
		&dbModel.MappingTarget{},
		&dbModel.ScheduledTask{},
		&dbModel.ScheduledTaskTag{},
	)
	if err != nil {
		return errors.New("PersistentDatabaseMigrationError")
	}

	modelsWithInitialEntries := []interface{}{
		&dbModel.ContainerProfile{},
	}

	err = service.seedDatabase(modelsWithInitialEntries...)
	if err != nil {
		return errors.New("AddDefaultPersistentDatabaseEntriesError")
	}

	return nil
}
