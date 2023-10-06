package testHelpers

import (
	"fmt"
	"path"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/speedianet/sfm/src/infra/db"
	"gorm.io/gorm"
)

func LoadEnvVars() {
	_, fileDirectory, _, _ := runtime.Caller(0)
	rootDir := path.Dir(path.Dir(path.Dir(fileDirectory)))
	testEnvPath := rootDir + "/.env"

	loadEnvErr := godotenv.Load(testEnvPath)
	if loadEnvErr != nil {
		panic(fmt.Errorf("LoadEnvFileError: %s", loadEnvErr))
	}
}

func GetDbSvc() *gorm.DB {
	dbSvc, err := db.DatabaseService()
	if err != nil {
		panic(fmt.Errorf("GetDbSvcError: %s", err))
	}
	return dbSvc
}
