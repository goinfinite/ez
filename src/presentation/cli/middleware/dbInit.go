package cliMiddleware

import (
	"fmt"
	"os"

	"github.com/speedianet/control/src/infra/db"
)

func DatabaseInit() *db.PersistentDatabaseService {
	persistentDbSvc, err := db.NewPersistentDatabaseService()
	if err != nil {
		fmt.Println("DatabaseConnectionError")
		os.Exit(1)
	}

	return persistentDbSvc
}
