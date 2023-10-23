package repository

import (
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type ContainerProfileCmdRepo interface {
	Add(addDto dto.AddContainerProfile) error
	Update(updateDto dto.UpdateContainerProfile) error
	Delete(profileId valueObject.ContainerProfileId) error
}
