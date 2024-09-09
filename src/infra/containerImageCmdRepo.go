package infra

import (
	"errors"
	"io"
	"os"
	"strings"

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

func (repo *ContainerImageCmdRepo) readArchiveFilesDirectory(
	accountId valueObject.AccountId,
) (archiveFilesDir valueObject.UnixFilePath, err error) {
	accountQueryRepo := NewAccountQueryRepo(repo.persistentDbSvc)
	accountEntity, err := accountQueryRepo.ReadById(accountId)
	if err != nil {
		return archiveFilesDir, err
	}

	archiveDirStr := accountEntity.HomeDirectory.String() + "/archives"
	accountIdStr := accountId.String()
	_, err = infraHelper.RunCmd(
		"install", "-d", "-m", "755", "-o", accountIdStr, "-g", accountIdStr, archiveDirStr,
	)
	if err != nil {
		return archiveFilesDir, errors.New("MakeArchiveDirError: " + err.Error())
	}

	return valueObject.NewUnixFilePath(archiveDirStr)
}

func (repo *ContainerImageCmdRepo) ImportArchiveFile(
	importDto dto.ImportContainerImageArchiveFile,
) (imageId valueObject.ContainerImageId, err error) {
	inputFileHandler, err := importDto.ArchiveFile.Open()
	if err != nil {
		return imageId, errors.New("OpenArchiveFileError: " + err.Error())
	}
	defer inputFileHandler.Close()

	archiveDir, err := repo.readArchiveFilesDirectory(importDto.AccountId)
	if err != nil {
		return imageId, err
	}

	rawOutputFilePath := archiveDir.String() + "/" + importDto.ArchiveFile.Filename
	outputFilePath, err := valueObject.NewUnixFilePath(rawOutputFilePath)
	if err != nil {
		return imageId, errors.New("ArchiveFilePathError: " + err.Error())
	}
	outputFilePathStr := outputFilePath.String()

	outputFileHandler, err := os.Create(outputFilePathStr)
	if err != nil {
		return imageId, errors.New("CreateArchiveFileError: " + err.Error())
	}
	defer outputFileHandler.Close()

	_, err = io.Copy(outputFileHandler, inputFileHandler)
	if err != nil {
		return imageId, errors.New("CopyArchiveFileError: " + err.Error())
	}

	_, err = infraHelper.RunCmdAsUser(
		importDto.AccountId,
		"brotli", "--decompress", "--rm", outputFilePathStr,
	)
	if err != nil {
		return imageId, errors.New("DecompressImageError: " + err.Error())
	}

	outputFilePathStr = strings.TrimSuffix(outputFilePathStr, ".br")

	rawImageId, err := infraHelper.RunCmdAsUser(
		importDto.AccountId,
		"podman", "load", "--quiet", "--input", outputFilePathStr,
	)
	if err != nil {
		return imageId, errors.New("PodmanLoadError: " + err.Error())
	}

	if len(rawImageId) == 0 {
		return imageId, errors.New("EmptyImageId")
	}
	rawImageId = strings.TrimSuffix(rawImageId, "Loaded image: sha256:")
	rawImageId = strings.TrimSpace(rawImageId)
	if len(rawImageId) > 12 {
		rawImageId = rawImageId[:12]
	}

	imageId, err = valueObject.NewContainerImageId(rawImageId)
	if err != nil {
		return imageId, errors.New("NewImageIdError: " + err.Error())
	}

	_, err = infraHelper.RunCmdAsUser(
		importDto.AccountId, "rm", "-f", outputFilePathStr,
	)
	if err != nil {
		return imageId, errors.New("RemoveArchiveFileError: " + err.Error())
	}

	return imageId, nil
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
	archiveDir, err := repo.readArchiveFilesDirectory(createDto.AccountId)
	if err != nil {
		return archiveFile, err
	}

	imageIdStr := createDto.ImageId.String()
	archiveDirStr := archiveDir.String()
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

	accountIdStr := createDto.AccountId.String()
	_, err = infraHelper.RunCmd("chown", "-R", accountIdStr, archiveDirStr)
	if err != nil {
		return archiveFile, errors.New("ChownArchiveDirError: " + err.Error())
	}

	finalFilePath, err := valueObject.NewUnixFilePath(tarFilePath + ".br")
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
	archiveDir, err := repo.readArchiveFilesDirectory(deleteDto.AccountId)
	if err != nil {
		return err
	}

	rawFilePath := archiveDir.String() + "/" + deleteDto.ImageId.String() + ".tar.br"
	filePath, err := valueObject.NewUnixFilePath(rawFilePath)
	if err != nil {
		return errors.New("ArchiveFilePathError: " + err.Error())
	}

	_, err = infraHelper.RunCmd("rm", "-f", filePath.String())
	if err != nil {
		return errors.New("RemoveArchiveFileError: " + err.Error())
	}

	return nil
}
