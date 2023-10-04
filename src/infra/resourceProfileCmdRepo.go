package infra

import (
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/valueObject"
)

type ResourceProfileCmdRepo struct {
}

func (repo ResourceProfileCmdRepo) Add(
	addDto dto.AddResourceProfile,
) error {
	return nil
}

func (repo ResourceProfileCmdRepo) Update(
	updateDto dto.UpdateResourceProfile,
) error {
	return nil
}

func (repo ResourceProfileCmdRepo) Delete(
	profileId valueObject.ResourceProfileId,
) error {
	return nil
}
