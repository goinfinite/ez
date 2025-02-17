package infra

import (
	"testing"

	testHelpers "github.com/goinfinite/ez/src/devUtils"
	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

func TestMappingQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	mappingQueryRepo := NewMappingQueryRepo(persistentDbSvc)

	t.Run("Read", func(t *testing.T) {
		_, err := mappingQueryRepo.Read()
		if err != nil {
			t.Errorf("ReadError: %v", err)
			return
		}
	})

	t.Run("ReadById", func(t *testing.T) {
		_, err := mappingQueryRepo.ReadById(1)
		if err != nil {
			t.Errorf("ReadByIdError: %v", err)
			return
		}
	})

	t.Run("ReadByContainerId", func(t *testing.T) {
		containerQueryRepo := NewContainerQueryRepo(persistentDbSvc)
		containerEntity, err := containerQueryRepo.ReadFirst(dto.ReadContainersRequest{})
		if err != nil {
			t.Errorf("ReadContainersError: %v", err)
			return
		}

		_, err = mappingQueryRepo.ReadByContainerId(containerEntity.Id)
		if err != nil {
			t.Errorf("ReadByContainerIdError: %v", err)
			return
		}
	})

	t.Run("ReadTargetById", func(t *testing.T) {
		_, err := mappingQueryRepo.ReadTargetById(1)
		if err != nil {
			t.Errorf("ReadTargetByIdError: %v", err)
			return
		}
	})

	t.Run("ReadByProtocol", func(t *testing.T) {
		protocol, _ := valueObject.NewNetworkProtocol("http")
		_, err := mappingQueryRepo.GetByProtocol(protocol)
		if err != nil {
			t.Errorf("ReadByProtocolError: %v", err)
			return
		}
	})
}
