package infra

import (
	"encoding/json"
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

	requestContainersDto := dto.ReadContainersRequest{
		Pagination:  useCase.ContainersDefaultPagination,
		ContainerId: []valueObject.ContainerId{createDto.ContainerId},
	}

	containerEntity, err := containerQueryRepo.ReadFirst(requestContainersDto)
	if err != nil {
		return imageId, err
	}

	accountQueryRepo := NewAccountQueryRepo(repo.persistentDbSvc)
	accountEntity, err := accountQueryRepo.ReadById(containerEntity.AccountId)
	if err != nil {
		return imageId, err
	}

	containerEntity.Entrypoint = nil
	containerEntity.PortBindings = []valueObject.PortBinding{}
	containerEntity.Envs = []valueObject.ContainerEnv{}
	containerEntity.StartedAt = nil
	containerEntity.StoppedAt = nil
	containerEntityJsonBytes, err := json.Marshal(containerEntity)
	if err != nil {
		return imageId, errors.New("MarshalContainerEntityError: " + err.Error())
	}
	encodedContainerEntityJson := infraHelper.EncodeStr(string(containerEntityJsonBytes))

	commitParams := []string{
		"commit",
		"--quiet",
		"--author", "ez:" + createDto.OperatorAccountId.String(),
		"--change", "LABEL=ez.originContainerDetails=" + encodedContainerEntityJson,
	}

	mappingQueryRepo := NewMappingQueryRepo(repo.persistentDbSvc)
	containerMappingEntities, err := mappingQueryRepo.ReadByContainerId(createDto.ContainerId)
	if err != nil {
		return imageId, err
	}

	if len(containerMappingEntities) > 0 {
		containerMappingEntitiesJsonBytes, err := json.Marshal(containerMappingEntities)
		if err != nil {
			return imageId, errors.New("MarshalContainerMappingEntitiesError: " + err.Error())
		}
		encodedMappingEntitiesJson := infraHelper.EncodeStr(string(containerMappingEntitiesJsonBytes))
		commitParams = append(
			commitParams,
			"--change", "LABEL=ez.originContainerMappings="+encodedMappingEntitiesJson,
		)
	}

	containerIdStr := createDto.ContainerId.String()
	containerHostnameStrSimplified := strings.ReplaceAll(
		containerEntity.Hostname.String(), ".", "-",
	)
	containerHostnameStrSimplified = strings.ToLower(containerHostnameStrSimplified)
	unixTimeStr := valueObject.NewUnixTimeNow().String()

	snapshotName := containerIdStr + "-" + containerHostnameStrSimplified + ":" + unixTimeStr
	imageAddress := "localhost/" + accountEntity.Username.String() + "/" + snapshotName

	commitParams = append(commitParams, containerIdStr, imageAddress)

	rawImageId, err := infraHelper.RunCmdAsUser(
		containerEntity.AccountId,
		"podman", commitParams...,
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

func (repo *ContainerImageCmdRepo) ImageArchiveFileLocator(
	originalArchiveFilePath valueObject.UnixFilePath,
) (adjustedArchiveFilePath valueObject.UnixFilePath, err error) {
	originalArchiveFilePathStr := originalArchiveFilePath.String()

	possibleCompressionExts := []string{".br", ".gz", ".xz", ".zip"}
	for _, possibleExt := range possibleCompressionExts {
		archiveFileWithPossibleExt := originalArchiveFilePathStr + possibleExt
		_, err = os.Stat(archiveFileWithPossibleExt)
		if err == nil {
			return valueObject.NewUnixFilePath(archiveFileWithPossibleExt)
		}
	}

	for _, possibleExt := range possibleCompressionExts {
		rawUncompressedArchiveFilePath := strings.ReplaceAll(
			originalArchiveFilePathStr, possibleExt, "",
		)
		_, err = os.Stat(rawUncompressedArchiveFilePath)
		if err == nil {
			return valueObject.NewUnixFilePath(rawUncompressedArchiveFilePath)
		}
	}

	return originalArchiveFilePath, errors.New("ArchiveFilePathNotFound")
}

func (repo *ContainerImageCmdRepo) decompressImageArchiveFile(
	compressedArchiveFilePath valueObject.UnixFilePath,
	accountId valueObject.AccountId,
) (uncompressedArchiveFilePath valueObject.UnixFilePath, err error) {
	compressedArchiveFilePathStr := compressedArchiveFilePath.String()
	_, err = os.Stat(compressedArchiveFilePathStr)
	if err != nil {
		compressedArchiveFilePath, err = repo.ImageArchiveFileLocator(compressedArchiveFilePath)
		if err != nil {
			return uncompressedArchiveFilePath, errors.New("ArchiveFilePathNotFound: " + err.Error())
		}
	}

	compressedArchiveExt, err := compressedArchiveFilePath.ReadFileExtension()
	if err != nil {
		return uncompressedArchiveFilePath, errors.New("ReadFileExtensionError: " + err.Error())
	}
	compressionFormat, err := valueObject.NewCompressionFormat(compressedArchiveExt.String())
	if err != nil {
		return uncompressedArchiveFilePath, errors.New("UnsupportedArchiveFileExtension")
	}

	decompressionCmd := ""
	compressionFormatStr := compressionFormat.String()
	switch compressionFormatStr {
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
		return uncompressedArchiveFilePath, errors.New("UnsupportedArchiveFileExtension")
	}

	if len(decompressionCmd) == 0 {
		return compressedArchiveFilePath, nil
	}

	rawUncompressedArchiveFilePath := strings.TrimSuffix(
		compressedArchiveFilePathStr, "."+compressionFormatStr,
	)
	uncompressedArchiveFilePath, err = valueObject.NewUnixFilePath(rawUncompressedArchiveFilePath)
	if err != nil {
		return uncompressedArchiveFilePath, errors.New("DefineUncompressedArchiveFilePathError")
	}
	uncompressedLocalArchiveFileStr := uncompressedArchiveFilePath.String()

	_ = os.Remove(uncompressedLocalArchiveFileStr)
	_, err = infraHelper.RunCmdAsUserWithSubShell(
		accountId,
		decompressionCmd+" "+compressedArchiveFilePathStr,
	)
	if err != nil {
		return uncompressedArchiveFilePath, errors.New("DecompressCmdError: " + err.Error())
	}

	return uncompressedArchiveFilePath, nil
}

func (repo *ContainerImageCmdRepo) ImportArchive(
	importDto dto.ImportContainerImageArchive,
) (imageId valueObject.ContainerImageId, err error) {
	wasArchiveFilePathProvided := importDto.ArchiveFilePath != nil
	if !wasArchiveFilePathProvided {
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
		localArchiveFilePath, err := valueObject.NewUnixFilePath(rawOutputFilePath)
		if err != nil {
			return imageId, errors.New("ArchiveFilePathError: " + err.Error())
		}

		localArchiveFileHandler, err := os.Create(localArchiveFilePath.String())
		if err != nil {
			return imageId, errors.New("CreateArchiveError: " + err.Error())
		}
		defer localArchiveFileHandler.Close()

		_, err = io.Copy(localArchiveFileHandler, inputFileHandler)
		if err != nil {
			return imageId, errors.New("CopyArchiveFileError: " + err.Error())
		}

		importDto.ArchiveFilePath = &localArchiveFilePath
	}
	localArchiveFilePathStr := importDto.ArchiveFilePath.String()

	accountIdInt := int(importDto.AccountId)
	_ = os.Chown(localArchiveFilePathStr, accountIdInt, accountIdInt)

	uncompressedArchiveFilePath, err := repo.decompressImageArchiveFile(
		*importDto.ArchiveFilePath, importDto.AccountId,
	)
	if err != nil {
		return imageId, errors.New("DecompressImageArchiveError: " + err.Error())
	}
	importDto.ArchiveFilePath = &uncompressedArchiveFilePath
	localArchiveFilePathStr = uncompressedArchiveFilePath.String()

	err = os.Chown(localArchiveFilePathStr, accountIdInt, accountIdInt)
	if err != nil {
		return imageId, errors.New("ChownUncompressedArchiveFileError: " + err.Error())
	}

	rawImageId, err := infraHelper.RunCmdAsUser(
		importDto.AccountId,
		"podman", "load", "--quiet", "--input", localArchiveFilePathStr,
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

	if wasArchiveFilePathProvided {
		return imageId, nil
	}

	return imageId, os.Remove(localArchiveFilePathStr)
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
	rawArchiveFileName := imageEntity.AccountId.String() + "-" + imageEntity.Id.String()

	imageOrgNameStr := ""
	imageOrgName, err := imageEntity.ImageAddress.ReadOrgName()
	if err == nil {
		imageOrgNameStr = imageOrgName.String()
	}
	if imageOrgNameStr != "" {
		rawArchiveFileName += "-" + imageOrgNameStr
	}

	imageNameStr := ""
	imageName, err := imageEntity.ImageAddress.ReadImageName()
	if err == nil {
		imageNameStr = imageName.String()
	}
	if imageNameStr != "" {
		rawArchiveFileName += "-" + imageNameStr
	}

	imageTagStr := ""
	imageTag, err := imageEntity.ImageAddress.ReadImageTag()
	if err == nil {
		imageTagStr = imageTag.String()
	}
	if imageTagStr != "" {
		rawArchiveFileName += "-" + imageTagStr
	}

	return valueObject.NewUnixFileName(rawArchiveFileName + ".tar")
}

func (repo *ContainerImageCmdRepo) compressImageArchiveFile(
	uncompressedArchiveFilePath valueObject.UnixFilePath,
	accountId valueObject.AccountId,
	compressionFormat *valueObject.CompressionFormat,
) (compressedArchiveFilePath valueObject.UnixFilePath, err error) {
	uncompressedArchiveFilePathStr := uncompressedArchiveFilePath.String()

	compressionCmd := "brotli --quality=4 --rm"
	compressionSuffix := ".br"
	if compressionFormat != nil {
		switch compressionFormatStr := compressionFormat.String(); compressionFormatStr {
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
			compressionCmd = "zip -q -m -6 " + uncompressedArchiveFilePathStr + ".zip"
			compressionSuffix = ".zip"
		case "xz":
			compressionCmd = "xz -1 --memlimit=10%"
			compressionSuffix = ".xz"
		default:
			return compressedArchiveFilePath, errors.New("UnsupportedCompressionFormat")
		}
	}

	compressedArchiveFilePath, err = valueObject.NewUnixFilePath(
		uncompressedArchiveFilePathStr + compressionSuffix,
	)
	if err != nil {
		return compressedArchiveFilePath, errors.New("NewCompressedArchiveFilePathError")
	}

	if compressionCmd == "" {
		return uncompressedArchiveFilePath, nil
	}

	_, err = infraHelper.RunCmdAsUserWithSubShell(
		accountId, compressionCmd+" "+uncompressedArchiveFilePathStr,
	)
	if err != nil {
		return compressedArchiveFilePath, errors.New("CompressImageArchiveError: " + err.Error())
	}

	return compressedArchiveFilePath, nil
}

func (repo *ContainerImageCmdRepo) CreateArchive(
	createDto dto.CreateContainerImageArchive,
) (archiveFile entity.ContainerImageArchive, err error) {
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
	if createDto.DestinationPath != nil {
		archiveDir = *createDto.DestinationPath
	}
	archiveDirStr := archiveDir.String()

	rawArchiveFilePath := archiveDirStr + "/" + archiveFileName.String()
	archiveFilePath, err := valueObject.NewUnixFilePath(rawArchiveFilePath)
	if err != nil {
		return archiveFile, errors.New("DefineNewArchiveFilePathError")
	}
	archiveFilePathStr := archiveFilePath.String()

	accountIdStr := imageEntity.AccountId.String()
	_, err = infraHelper.RunCmd(
		"chown", "-R", accountIdStr+":"+accountIdStr, archiveDirStr,
	)
	if err != nil {
		return archiveFile, errors.New("ChownArchiveDirError: " + err.Error())
	}

	imageIdStr := imageEntity.Id.String()
	_ = os.Remove(archiveFilePathStr)
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

	compressedImageArchiveFilePath, err := repo.compressImageArchiveFile(
		archiveFilePath, imageEntity.AccountId, createDto.CompressionFormat,
	)
	if err != nil {
		return archiveFile, errors.New("CompressImageArchiveError: " + err.Error())
	}

	archiveFileInfo, err := os.Stat(compressedImageArchiveFilePath.String())
	if err != nil {
		return archiveFile, errors.New("StatFinalFileError: " + err.Error())
	}

	sizeBytes, err := valueObject.NewByte(archiveFileInfo.Size())
	if err != nil {
		return archiveFile, errors.New("CalculateArchiveFileSizeError")
	}

	serverHostname, err := infraHelper.ReadServerHostname()
	if err != nil {
		return archiveFile, errors.New("InvalidServerHostname: " + err.Error())
	}

	downloadUrl, _ := valueObject.NewUrl(
		"https://" + serverHostname.String() +
			"/v1/container/image/archive/" + accountIdStr + "/" + imageIdStr + "/",
	)

	return entity.NewContainerImageArchive(
		createDto.ImageId, imageEntity.AccountId, compressedImageArchiveFilePath, sizeBytes,
		&downloadUrl, nil, valueObject.NewUnixTimeNow(),
	), nil
}

func (repo *ContainerImageCmdRepo) DeleteArchive(
	archiveEntity entity.ContainerImageArchive,
) error {
	_, err := infraHelper.RunCmd("rm", "-f", archiveEntity.UnixFilePath.String())
	return err
}
