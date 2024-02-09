package cliInit

import "github.com/speedianet/control/src/infra/db"

func DatabaseService() *db.DatabaseService {
	dbSvc, err := db.NewDatabaseService()
	if err != nil {
		panic("DatabaseConnectionError:" + err.Error())
	}

	return dbSvc
}
