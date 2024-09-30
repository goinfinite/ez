package infra

import (
	"testing"

	testHelpers "github.com/goinfinite/ez/src/devUtils"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func TestMappingQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	mappingQueryRepo := NewMappingQueryRepo(persistentDbSvc)

	t.Run("ReadMappings", func(t *testing.T) {
		mappingList, err := mappingQueryRepo.Read()
		if err != nil {
			t.Error(err)
			return
		}

		if len(mappingList) == 0 {
			t.Error("NoMappingsFound")
			return
		}
	})

	t.Run("GetMappingById", func(t *testing.T) {
		mapping, err := mappingQueryRepo.ReadById(1)
		if err != nil {
			t.Error(err)
			return
		}

		if mapping.Id.Uint64() != 1 {
			t.Error("MappingNotFound")
			return
		}
	})

	t.Run("ReadTargetById", func(t *testing.T) {
		_, err := mappingQueryRepo.ReadTargetById(1)
		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("GetByProtocol", func(t *testing.T) {
		protocol, _ := valueObject.NewNetworkProtocol("http")
		_, err := mappingQueryRepo.GetByProtocol(protocol)
		if err != nil {
			t.Error(err)
			return
		}
	})
}
