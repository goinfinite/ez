package testHelpers

import (
	"path"
	"runtime"

	"github.com/goinfinite/ez/src/infra/db"
	"github.com/joho/godotenv"
)

func LoadEnvVars() {
	_, fileDirectory, _, _ := runtime.Caller(0)
	rootDir := path.Dir(path.Dir(path.Dir(fileDirectory)))
	testEnvPath := rootDir + "/.env"

	err := godotenv.Load(testEnvPath)
	if err != nil {
		panic("LoadEnvFileError: " + err.Error())
	}
}

func GetPersistentDbSvc() *db.PersistentDatabaseService {
	persistentDbSvc, err := db.NewPersistentDatabaseService()
	if err != nil {
		panic("GetPersistentDbSvcError: " + err.Error())
	}
	return persistentDbSvc
}

func GetTransientDbSvc() *db.TransientDatabaseService {
	transientDbSvc, err := db.NewTransientDatabaseService()
	if err != nil {
		panic("GetTransientDbSvcError: " + err.Error())
	}
	return transientDbSvc
}

func GetTrailDbSvc() *db.TrailDatabaseService {
	trailDbSvc, err := db.NewTrailDatabaseService()
	if err != nil {
		panic("GetTrailDbSvcError: " + err.Error())
	}
	return trailDbSvc
}
