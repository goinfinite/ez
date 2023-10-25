package repository

import (
	"github.com/goinfinite/fleet/src/domain/dto"
	"github.com/goinfinite/fleet/src/domain/valueObject"
)

type ContainerProfileCmdRepo interface {
	Add(addDto dto.AddContainerProfile) error
	Update(updateDto dto.UpdateContainerProfile) error
	Delete(profileId valueObject.ContainerProfileId) error
}
