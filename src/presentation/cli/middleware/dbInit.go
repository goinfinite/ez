package cliMiddleware

import (
	"fmt"
	"os"

	"github.com/speedianet/sfm/src/infra/db"
	"gorm.io/gorm"
)

func DatabaseInit() *gorm.DB {
	dbSvc, err := db.DatabaseService()
	if err != nil {
		fmt.Println("DatabaseConnectionError")
		os.Exit(1)
	}

	return dbSvc
}
