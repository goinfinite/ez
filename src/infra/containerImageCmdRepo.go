package infra

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	infraHelper "github.com/goinfinite/ez/src/infra/helper"
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
	containerQueryRepo := NewContainerQueryRepo(repo.persistentDbSvc)

	readContainersRequestDto := dto.ReadContainersRequest{
		Pagination:  useCase.ContainersDefaultPagination,
		ContainerId: &createDto.ContainerId,
	}

	readContainersResponseDto, err := containerQueryRepo.Read(readContainersRequestDto)
	if err != nil || len(readContainersResponseDto.Containers) == 0 {
		return imageId, errors.New("ContainerNotFound")
	}
	containerEntity := readContainersResponseDto.Containers[0]
	containerIdStr := createDto.ContainerId.String()
	containerHostnameStrSimplified := strings.ReplaceAll(
		containerEntity.Hostname.String(), ".", "-",
	)
	containerHostnameStrSimplified = strings.ToLower(containerHostnameStrSimplified)
	snapshotName := containerIdStr + "-" +
		containerHostnameStrSimplified +
		":" + valueObject.NewUnixTimeNow().String()

	accountQueryRepo := NewAccountQueryRepo(repo.persistentDbSvc)
	accountEntity, err := accountQueryRepo.ReadById(containerEntity.AccountId)
	if err != nil {
		return imageId, err
	}

	rawImageId, err := infraHelper.RunCmdAsUser(
		containerEntity.AccountId,
		"podman", "commit", "--quiet",
		"--author", "ez:"+createDto.OperatorAccountId.String(),
		containerIdStr,
		"localhost/"+accountEntity.Username.String()+"/"+snapshotName,
	)
	if err != nil {
		return imageId, err
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

	outputFileExtension, err := outputFilePath.ReadFileExtension()
	if err != nil {
		if !strings.HasSuffix(outputFilePathStr, ".br") {
			return imageId, errors.New("ReadFileExtensionError: " + err.Error())
		}
	}

	decompressionCmd := ""
	outputFileExtStr := outputFileExtension.String()
	switch outputFileExtStr {
	case "tar":
		decompressionCmd = ""
	case "br":
		decompressionCmd = "brotli --decompress --rm"
	case "gz":
		decompressionCmd = "gunzip"
	case "zip":
		decompressionCmd = "unzip -q"
	case "xz":
		decompressionCmd = "unxz"
	default:
		return imageId, errors.New("UnsupportedArchiveFileExtension")
	}

	_, err = io.Copy(outputFileHandler, inputFileHandler)
	if err != nil {
		return imageId, errors.New("CopyArchiveFileError: " + err.Error())
	}

	accountIdStr := importDto.AccountId.String()
	_, err = infraHelper.RunCmd(
		"chown", accountIdStr+":"+accountIdStr, outputFilePathStr,
	)
	if err != nil {
		return imageId, errors.New("ChownArchiveError: " + err.Error())
	}

	archiveFilePathStr := outputFilePathStr
	shouldDecompress := len(decompressionCmd) > 0
	if shouldDecompress {
		archiveFilePathStr = strings.TrimSuffix(outputFilePathStr, "."+outputFileExtStr)
		_, err = infraHelper.RunCmdAsUserWithSubShell(
			importDto.AccountId, "rm -f "+archiveFilePathStr,
		)
		if err != nil {
			return imageId, errors.New("RemoveExistingTarFileError: " + err.Error())
		}

		_, err = infraHelper.RunCmdAsUserWithSubShell(
			importDto.AccountId, decompressionCmd+" "+outputFilePathStr,
		)
		if err != nil {
			return imageId, errors.New("DecompressImageError: " + err.Error())
		}
	}

	rawImageId, err := infraHelper.RunCmdAsUser(
		importDto.AccountId,
		"podman", "load", "--quiet", "--input", archiveFilePathStr,
	)
	if err != nil {
		return imageId, errors.New("PodmanLoadError: " + err.Error())
	}

	if len(rawImageId) == 0 {
		return imageId, errors.New("PodmanLoadEmptyImageId")
	}
	rawImageId = strings.TrimPrefix(rawImageId, "Loaded image: sha256:")
	rawImageId = strings.TrimSpace(rawImageId)

	imageId, err = valueObject.NewContainerImageId(rawImageId)
	if err != nil {
		return imageId, err
	}

	_, err = infraHelper.RunCmdAsUser(
		importDto.AccountId, "rm", "-f", archiveFilePathStr,
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

func (repo *ContainerImageCmdRepo) archiveFileNameFactory(
	imageEntity entity.ContainerImage,
) (archiveFileName valueObject.UnixFileName, err error) {
	imageOrgNameStr := ""
	imageOrgName, err := imageEntity.ImageAddress.ReadOrgName()
	if err == nil {
		imageOrgNameStr = imageOrgName.String()
	}

	imageNameStr := ""
	imageName, err := imageEntity.ImageAddress.ReadImageName()
	if err == nil {
		imageNameStr = imageName.String()
	}

	imageTagStr := ""
	imageTag, err := imageEntity.ImageAddress.ReadImageTag()
	if err == nil {
		imageTagStr = imageTag.String()
	}

	rawArchiveFileName := imageEntity.Id.String()
	if imageOrgNameStr != "" {
		rawArchiveFileName += "-" + imageOrgNameStr
	}
	if imageNameStr != "" {
		rawArchiveFileName += "-" + imageNameStr
	}
	if imageTagStr != "" {
		rawArchiveFileName += "-" + imageTagStr
	}
	rawArchiveFileName += ".tar"

	return valueObject.NewUnixFileName(rawArchiveFileName)
}

func (repo *ContainerImageCmdRepo) CreateArchiveFile(
	createDto dto.CreateContainerImageArchiveFile,
) (archiveFile entity.ContainerImageArchiveFile, err error) {
	containerImageQueryRepo := NewContainerImageQueryRepo(repo.persistentDbSvc)
	imageEntity, err := containerImageQueryRepo.ReadById(
		createDto.AccountId, createDto.ImageId,
	)
	if err != nil {
		return archiveFile, err
	}

	archiveFileName, err := repo.archiveFileNameFactory(imageEntity)
	if err != nil {
		return archiveFile, errors.New("NewArchiveFileNameError: " + err.Error())
	}

	archiveDir, err := repo.readArchiveFilesDirectory(imageEntity.AccountId)
	if err != nil {
		return archiveFile, err
	}
	archiveDirStr := archiveDir.String()

	archiveFilePathStr := archiveDirStr + "/" + archiveFileName.String()
	_, err = infraHelper.RunCmdAsUser(
		imageEntity.AccountId, "rm", "-f", archiveFilePathStr,
	)
	if err != nil {
		return archiveFile, errors.New("RemoveExistingTarFileError: " + err.Error())
	}

	imageIdStr := imageEntity.Id.String()
	_, err = infraHelper.RunCmdAsUser(
		imageEntity.AccountId,
		"podman", "save",
		"--format", "docker-archive",
		"--output", archiveFilePathStr,
		imageIdStr,
	)
	if err != nil {
		return archiveFile, errors.New("PodmanSaveError: " + err.Error())
	}

	compressionCmd := "brotli --quality=4 --rm"
	compressionSuffix := ".br"
	if createDto.CompressionFormat != nil {
		compressionFormatStr := createDto.CompressionFormat.String()
		switch compressionFormatStr {
		case "tar":
			compressionCmd = ""
			compressionSuffix = ""
		case "br":
			compressionCmd = "brotli --quality=4 --rm"
			compressionSuffix = ".br"
		case "gzip":
			compressionCmd = "gzip -6"
			compressionSuffix = ".gz"
		case "zip":
			compressionCmd = "zip -q -m -6 " + archiveFilePathStr + ".zip"
			compressionSuffix = ".zip"
		case "xz":
			compressionCmd = "xz -1 --memlimit=10%"
			compressionSuffix = ".xz"
		default:
			return archiveFile, errors.New("UnsupportedCompressionFormat")
		}
	}

	if compressionCmd != "" {
		_, err = infraHelper.RunCmdAsUserWithSubShell(
			imageEntity.AccountId, compressionCmd+" "+archiveFilePathStr,
		)
		if err != nil {
			return archiveFile, errors.New("CompressImageError: " + err.Error())
		}
	}

	accountIdStr := imageEntity.AccountId.String()
	_, err = infraHelper.RunCmd(
		"chown", "-R", accountIdStr+":"+accountIdStr, archiveDirStr,
	)
	if err != nil {
		return archiveFile, errors.New("ChownArchiveDirError: " + err.Error())
	}

	finalFilePath, err := valueObject.NewUnixFilePath(
		archiveFilePathStr + compressionSuffix,
	)
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
		createDto.ImageId, imageEntity.AccountId, finalFilePath, downloadUrl,
		sizeBytes, valueObject.NewUnixTimeNow(),
	), nil
}

func (repo *ContainerImageCmdRepo) DeleteArchiveFile(
	archiveFileEntity entity.ContainerImageArchiveFile,
) error {
	_, err := infraHelper.RunCmd("rm", "-f", archiveFileEntity.UnixFilePath.String())
	return err
}
