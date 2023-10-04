package repository

import (
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type ResourceProfileCmdRepo interface {
	Add(addDto dto.AddResourceProfile) error
	Update(updateDto dto.UpdateResourceProfile) error
	Delete(profileId valueObject.ResourceProfileId) error
}
