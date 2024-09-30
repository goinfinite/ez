package cliMiddleware

import (
	"fmt"
	"os"

	"github.com/goinfinite/ez/src/infra/db"
)

func PersistentDatabaseInit() *db.PersistentDatabaseService {
	persistentDbSvc, err := db.NewPersistentDatabaseService()
	if err != nil {
		fmt.Println("PersistentDatabaseConnectionError")
		os.Exit(1)
	}

	return persistentDbSvc
}

func TransientDatabaseInit() *db.TransientDatabaseService {
	transientDbSvc, err := db.NewTransientDatabaseService()
	if err != nil {
		fmt.Println("TransientDatabaseConnectionError")
		os.Exit(1)
	}

	return transientDbSvc
}
