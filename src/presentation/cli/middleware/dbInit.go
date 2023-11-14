package cliMiddleware

import (
	"fmt"
	"os"

	"github.com/speedianet/control/src/infra/db"
)

func DatabaseInit() *db.DatabaseService {
	dbSvc, err := db.NewDatabaseService()
	if err != nil {
		fmt.Println("DatabaseConnectionError")
		os.Exit(1)
	}

	return dbSvc
}
