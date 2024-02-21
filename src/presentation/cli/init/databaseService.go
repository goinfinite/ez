package cliInit

import "github.com/speedianet/control/src/infra/db"

func PersistentDatabaseService() *db.PersistentDatabaseService {
	persistentDbSvc, err := db.NewPersistentDatabaseService()
	if err != nil {
		panic("PersistentDatabaseConnectionError:" + err.Error())
	}

	return persistentDbSvc
}

func TransientDatabaseService() *db.TransientDatabaseService {
	transientDbSvc, err := db.NewTransientDatabaseService()
	if err != nil {
		panic("PersistentDatabaseConnectionError:" + err.Error())
	}

	return transientDbSvc
}
