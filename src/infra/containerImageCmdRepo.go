package infra

import (
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	infraHelper "github.com/speedianet/control/src/infra/helper"
)

type ContainerImageCmdRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewContainerImageCmdRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerImageCmdRepo {
	return &ContainerImageCmdRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *ContainerImageCmdRepo) CreateSnapshot(
	createDto dto.CreateContainerSnapshotImage,
) (imageId valueObject.ContainerImageId, err error) {
	unixTimeNow := valueObject.NewUnixTimeNow()
	containerIdStr := createDto.ContainerId.String()
	snapshotName := containerIdStr + ":" + unixTimeNow.String()

	rawImageId, err := infraHelper.RunCmdAsUser(
		createDto.AccountId,
		"podman", "commit", "--quiet",
		"--author", "control:"+createDto.OperatorAccountId.String(),
		containerIdStr,
		"localhost/"+createDto.AccountId.String()+"/"+snapshotName,
	)
	if err != nil {
		return imageId, err
	}
	if len(rawImageId) > 12 {
		rawImageId = rawImageId[:12]
	}

	return valueObject.NewContainerImageId(rawImageId)
}

func (repo *ContainerImageCmdRepo) Export(
	exportDto dto.ExportContainerImage,
) (downloadUrl valueObject.Url, err error) {
	return valueObject.NewUrl("http://localhost")
}

func (repo *ContainerImageCmdRepo) Delete(
	deleteDto dto.DeleteContainerImage,
) error {
	_, err := infraHelper.RunCmdAsUser(
		deleteDto.AccountId, "podman", "image", "rm", deleteDto.ImageId.String(),
	)
	return err
}
