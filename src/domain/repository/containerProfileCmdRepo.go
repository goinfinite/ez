package repository

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ContainerProfileCmdRepo interface {
	Add(addDto dto.AddContainerProfile) error
	Update(updateDto dto.UpdateContainerProfile) error
	Delete(profileId valueObject.ContainerProfileId) error
}
