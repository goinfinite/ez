package cliInit

import "github.com/speedianet/control/src/infra/db"

func PersistentDatabaseService() *db.PersistentDatabaseService {
	persistDbSvc, err := db.NewPersistentDatabaseService()
	if err != nil {
		panic("DatabaseConnectionError:" + err.Error())
	}

	return persistDbSvc
}
