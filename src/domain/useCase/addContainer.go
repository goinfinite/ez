package useCase

import (
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/repository"
)

func AddContainer(
	containerCmdRepo repository.ContainerCmdRepo,
	addContainer dto.AddContainer,
) error {
	return containerCmdRepo.Add(addContainer)
}
