package useCase

import (
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/repository"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

func AddContainer(
	containerCmdRepo repository.ContainerCmdRepo,
	addContainer dto.AddContainer,
) error {
	defaultSpecs := valueObject.NewContainerSpecs(
		0.5,
		valueObject.Byte(512000000),
	)

	if addContainer.BaseSpecs == nil {
		addContainer.BaseSpecs = &defaultSpecs
	}

	return containerCmdRepo.Add(addContainer)
}
