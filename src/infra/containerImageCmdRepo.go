package infra

import (
	"errors"
	"os"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
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

func (repo *ContainerImageCmdRepo) getAccountHomeDir(
	accountId valueObject.AccountId,
) (string, error) {
	// @see https://github.com/speedianet/control-issues-tracker/issues/92
	return infraHelper.RunCmdWithSubShell(
		"awk -F: '$3 == " + accountId.String() + " {print $6}' /etc/passwd",
	)
}

func (repo *ContainerImageCmdRepo) ImportArchiveFile(
	importDto dto.ImportContainerImageArchiveFile,
) (imageId valueObject.ContainerImageId, err error) {
	imageId, _ = valueObject.NewContainerImageId("1234567890abcdef1234567890")
	return imageId, errors.New("NotImplemented")
}

func (repo *ContainerImageCmdRepo) Delete(
	deleteDto dto.DeleteContainerImage,
) error {
	_, err := infraHelper.RunCmdAsUser(
		deleteDto.AccountId, "podman", "image", "rm", deleteDto.ImageId.String(),
	)
	return err
}

func (repo *ContainerImageCmdRepo) CreateArchiveFile(
	createDto dto.CreateContainerImageArchiveFile,
) (archiveFile entity.ContainerImageArchiveFile, err error) {
	accountHomeDir, err := repo.getAccountHomeDir(createDto.AccountId)
	if err != nil {
		return archiveFile, err
	}

	archiveDirStr := accountHomeDir + "/archive"
	accountIdStr := createDto.AccountId.String()
	_, err = infraHelper.RunCmd(
		"install", "-d", "-m", "755", "-o", accountIdStr, "-g", accountIdStr, archiveDirStr,
	)
	if err != nil {
		return archiveFile, errors.New("MakeArchiveDirError: " + err.Error())
	}

	imageIdStr := createDto.ImageId.String()
	tarFilePath := archiveDirStr + "/" + imageIdStr + ".tar"
	_, err = infraHelper.RunCmdAsUser(
		createDto.AccountId, "podman", "save", imageIdStr, "--output", tarFilePath,
	)
	if err != nil {
		return archiveFile, errors.New("PodmanSaveError: " + err.Error())
	}

	_, err = infraHelper.RunCmdAsUser(
		createDto.AccountId, "brotli", "--quality=4", "--rm", tarFilePath,
	)
	if err != nil {
		return archiveFile, errors.New("CompressImageError: " + err.Error())
	}

	_, err = infraHelper.RunCmd("chown", "-R", accountIdStr, archiveDirStr)
	if err != nil {
		return archiveFile, errors.New("ChownArchiveDirError: " + err.Error())
	}

	finalFilePath, err := valueObject.NewUnixFilePath(tarFilePath + ".bry")
	if err != nil {
		return archiveFile, errors.New("NewFinalFilePathError: " + err.Error())
	}

	fileInfo, err := os.Stat(finalFilePath.String())
	if err != nil {
		return archiveFile, errors.New("StatFinalFileError: " + err.Error())
	}

	sizeBytes, err := valueObject.NewByte(fileInfo.Size())
	if err != nil {
		return archiveFile, errors.New("NewSizeBytesError: " + err.Error())
	}

	serverHostname, err := infraHelper.ReadServerHostname()
	if err != nil {
		return archiveFile, errors.New("InvalidServerHostname: " + err.Error())
	}

	downloadUrl, _ := valueObject.NewUrl(
		"https://" + serverHostname.String() +
			"/v1/container/image/archive/" + accountIdStr + "/" + imageIdStr + "/",
	)

	return entity.NewContainerImageArchiveFile(
		createDto.ImageId, createDto.AccountId, finalFilePath, downloadUrl,
		sizeBytes, valueObject.NewUnixTimeNow(),
	), nil
}

func (repo *ContainerImageCmdRepo) DeleteArchiveFile(
	deleteDto dto.DeleteContainerImageArchiveFile,
) error {
	return errors.New("NotImplemented")
}
