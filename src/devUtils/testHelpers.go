package testHelpers

import (
	"fmt"
	"path"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/speedianet/control/src/infra/db"
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

func GetPersistentDbSvc() *db.PersistentDatabaseService {
	persistDbSvc, err := db.NewPersistentDatabaseService()
	if err != nil {
		panic(fmt.Errorf("GetPersistentDbSvcError: %s", err))
	}
	return persistDbSvc
}
