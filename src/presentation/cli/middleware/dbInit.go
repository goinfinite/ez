package cliMiddleware

import (
	"fmt"
	"os"

	"github.com/goinfinite/fleet/src/infra/db"
)

func DatabaseInit() *db.DatabaseService {
	dbSvc, err := db.NewDatabaseService()
	if err != nil {
		fmt.Println("DatabaseConnectionError")
		os.Exit(1)
	}

	return dbSvc
}
