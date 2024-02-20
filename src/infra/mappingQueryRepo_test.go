package infra

import (
	"testing"

	testHelpers "github.com/speedianet/control/src/devUtils"
	"github.com/speedianet/control/src/domain/valueObject"
)

func TestMappingQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	mappingQueryRepo := NewMappingQueryRepo(persistentDbSvc)

	t.Run("GetMappings", func(t *testing.T) {
		mappingList, err := mappingQueryRepo.Get()
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
		mapping, err := mappingQueryRepo.GetById(1)
		if err != nil {
			t.Error(err)
			return
		}

		if mapping.Id.Get() != 1 {
			t.Error("MappingNotFound")
			return
		}
	})

	t.Run("GetTargetById", func(t *testing.T) {
		_, err := mappingQueryRepo.GetTargetById(1)
		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("GetByProtocol", func(t *testing.T) {
		protocol := valueObject.NewNetworkProtocolPanic("http")
		_, err := mappingQueryRepo.GetByProtocol(protocol)
		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("FindOneMappingWithoutHostname", func(t *testing.T) {
		_, err := mappingQueryRepo.FindOne(
			nil,
			valueObject.NewNetworkPortPanic(80),
			valueObject.NewNetworkProtocolPanic("http"),
		)
		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("FindOneMappingWithHostname", func(t *testing.T) {
		hostname := valueObject.NewFqdnPanic("speedia.net")

		_, err := mappingQueryRepo.FindOne(
			&hostname,
			valueObject.NewNetworkPortPanic(80),
			valueObject.NewNetworkProtocolPanic("http"),
		)
		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("FindOneMappingWithUnpublishedPort", func(t *testing.T) {
		_, err := mappingQueryRepo.FindOne(
			nil,
			valueObject.NewNetworkPortPanic(8080),
			valueObject.NewNetworkProtocolPanic("http"),
		)
		if err == nil {
			t.Error("ShouldNotFindAnyMapping")
			return
		}
	})
}
