package repository

import "github.com/speedianet/control/src/domain/dto"

type ContainerImageCmdRepo interface {
	Delete(deleteDto dto.DeleteContainerImage) error
}
