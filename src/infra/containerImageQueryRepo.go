package infra

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	infraEnvs "github.com/goinfinite/ez/src/infra/envs"
	infraHelper "github.com/goinfinite/ez/src/infra/helper"
)

type ContainerImageQueryRepo struct {
	persistentDbSvc  *db.PersistentDatabaseService
	accountQueryRepo *AccountQueryRepo
}

func NewContainerImageQueryRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerImageQueryRepo {
	return &ContainerImageQueryRepo{
		persistentDbSvc:  persistentDbSvc,
		accountQueryRepo: NewAccountQueryRepo(persistentDbSvc),
	}
}

func (repo *ContainerImageQueryRepo) originContainerDetailsFactory(
	rawOriginContainerDetails string,
) (containerEntity entity.Container, err error) {
	decodedDetails, err := infraHelper.DecodeStr(rawOriginContainerDetails)
	if err != nil {
		return containerEntity, err
	}

	var originContainerDetailsMap map[string]interface{}
	err = json.Unmarshal([]byte(decodedDetails), &originContainerDetailsMap)
	if err != nil {
		return containerEntity, err
	}

	rawContainerId, exists := originContainerDetailsMap["id"]
	if !exists {
		return containerEntity, errors.New("ContainerIdNotFound")
	}
	containerId, err := valueObject.NewContainerId(rawContainerId)
	if err != nil {
		return containerEntity, err
	}

	rawAccountId, exists := originContainerDetailsMap["accountId"]
	if !exists {
		return containerEntity, errors.New("AccountIdNotFound")
	}
	accountId, err := valueObject.NewAccountId(rawAccountId)
	if err != nil {
		return containerEntity, err
	}

	rawHostname, exists := originContainerDetailsMap["hostname"]
	if !exists {
		return containerEntity, errors.New("HostnameNotFound")
	}
	containerHostname, err := valueObject.NewFqdn(rawHostname)
	if err != nil {
		return containerEntity, err
	}

	rawRestartPolicy, exists := originContainerDetailsMap["restartPolicy"]
	if !exists {
		return containerEntity, errors.New("RestartPolicyNotFound")
	}
	restartPolicy, err := valueObject.NewContainerRestartPolicy(rawRestartPolicy)
	if err != nil {
		return containerEntity, err
	}

	rawProfileId, exists := originContainerDetailsMap["profileId"]
	if !exists {
		return containerEntity, errors.New("ProfileIdNotFound")
	}
	profileId, err := valueObject.NewContainerProfileId(rawProfileId)
	if err != nil {
		return containerEntity, err
	}

	return entity.Container{
		Id:            containerId,
		AccountId:     accountId,
		Hostname:      containerHostname,
		RestartPolicy: restartPolicy,
		ProfileId:     profileId,
	}, nil
}

func (repo *ContainerImageQueryRepo) originContainerMappingTargetFactory(
	rawTarget interface{},
) (targetEntity entity.MappingTarget, err error) {
	rawTargetMap, assertOk := rawTarget.(map[string]interface{})
	if !assertOk {
		return targetEntity, errors.New("InvalidContainerMappingTarget")
	}

	targetId, err := valueObject.NewMappingTargetId(rawTargetMap["id"])
	if err != nil {
		return targetEntity, nil
	}

	mappingId, err := valueObject.NewMappingId(rawTargetMap["mappingId"])
	if err != nil {
		return targetEntity, nil
	}

	containerId, err := valueObject.NewContainerId(rawTargetMap["containerId"])
	if err != nil {
		return targetEntity, nil
	}

	containerHostname, err := valueObject.NewFqdn(rawTargetMap["containerHostname"])
	if err != nil {
		return targetEntity, nil
	}

	containerPrivatePort, err := valueObject.NewNetworkPort(rawTargetMap["containerPrivatePort"])
	if err != nil {
		return targetEntity, nil
	}

	return entity.NewMappingTarget(
		targetId, mappingId, containerId, containerHostname, containerPrivatePort,
	), nil
}

func (repo *ContainerImageQueryRepo) originContainerMappingFactory(
	rawMapping interface{},
) (mappingEntity entity.Mapping, err error) {
	rawMappingMap, assertOk := rawMapping.(map[string]interface{})
	if !assertOk {
		return mappingEntity, errors.New("InvalidContainerMappingStructure")
	}

	mappingId, err := valueObject.NewMappingId(rawMappingMap["id"])
	if err != nil {
		return mappingEntity, err
	}

	accountId, err := valueObject.NewAccountId(rawMappingMap["accountId"])
	if err != nil {
		return mappingEntity, err
	}

	var hostnamePtr *valueObject.Fqdn
	if rawMappingMap["hostname"] != nil {
		hostname, err := valueObject.NewFqdn(rawMappingMap["hostname"])
		if err != nil {
			return mappingEntity, err
		}
		hostnamePtr = &hostname
	}

	publicPort, err := valueObject.NewNetworkPort(rawMappingMap["publicPort"])
	if err != nil {
		return mappingEntity, err
	}

	protocol, err := valueObject.NewNetworkProtocol(rawMappingMap["protocol"])
	if err != nil {
		return mappingEntity, err
	}

	targetEntities := []entity.MappingTarget{}
	rawTargets, assertOk := rawMappingMap["targets"].([]interface{})
	if !assertOk {
		return mappingEntity, errors.New("InvalidContainerMappingTargets")
	}
	for _, rawTarget := range rawTargets {
		targetEntity, err := repo.originContainerMappingTargetFactory(rawTarget)
		if err != nil {
			slog.Debug(err.Error(), slog.Any("rawTarget", rawTarget))
			return mappingEntity, err
		}

		targetEntities = append(targetEntities, targetEntity)
	}

	createdAt, err := valueObject.NewUnixTime(rawMappingMap["createdAt"])
	if err != nil {
		return mappingEntity, err
	}

	updatedAt, err := valueObject.NewUnixTime(rawMappingMap["updatedAt"])
	if err != nil {
		return mappingEntity, err
	}

	return entity.NewMapping(
		mappingId, accountId, hostnamePtr, publicPort, protocol,
		targetEntities, createdAt, updatedAt,
	), nil
}

func (repo *ContainerImageQueryRepo) originContainerMappingsFactory(
	rawOriginContainerMappings string,
) (mappingEntities []entity.Mapping, err error) {
	decodedMappings, err := infraHelper.DecodeStr(rawOriginContainerMappings)
	if err != nil {
		return mappingEntities, err
	}

	var originContainerMappingsRawList []interface{}
	err = json.Unmarshal([]byte(decodedMappings), &originContainerMappingsRawList)
	if err != nil {
		return mappingEntities, err
	}

	for _, rawMapping := range originContainerMappingsRawList {
		mappingEntity, err := repo.originContainerMappingFactory(rawMapping)
		if err != nil {
			slog.Debug(err.Error(), slog.Any("rawMapping", rawMapping))
			return mappingEntities, err
		}

		mappingEntities = append(mappingEntities, mappingEntity)
	}

	return mappingEntities, nil
}

func (repo *ContainerImageQueryRepo) containerImageFactory(
	accountId valueObject.AccountId,
	rawContainerImage map[string]interface{},
) (containerImage entity.ContainerImage, err error) {
	rawImageId, assertOk := rawContainerImage["Id"].(string)
	if !assertOk {
		return containerImage, errors.New("InvalidContainerImageId")
	}
	imageId, err := valueObject.NewContainerImageId(rawImageId)
	if err != nil {
		return containerImage, err
	}

	rawConfig, assertOk := rawContainerImage["Config"].(map[string]interface{})
	if !assertOk {
		return containerImage, errors.New("InvalidContainerImageConfig")
	}

	rawLabels, assertOk := rawConfig["Labels"].(map[string]interface{})
	if !assertOk {
		return containerImage, errors.New("InvalidContainerImageLabels")
	}

	rawImageAddress, exists := rawLabels["ez.imageAddress"]
	if !exists {
		rawImageNames, assertOk := rawContainerImage["NamesHistory"].([]interface{})
		if !assertOk {
			rawAccountUsername := "unknown"
			idOutput, err := infraHelper.RunCmd("id", "-nu", accountId.String())
			if err == nil {
				rawAccountUsername = idOutput
			}

			accountUsername, err := valueObject.NewUnixUsername(rawAccountUsername)
			if err != nil {
				return containerImage, errors.New("ParseContainerImageAccountUsernameError")
			}

			rawImageNames = []interface{}{
				"localhost/" + accountUsername.String() + "/" + imageId.String(),
			}
		}
		if len(rawImageNames) == 0 {
			return containerImage, errors.New("ReadContainerImageNamesError")
		}
		rawImageAddress = rawImageNames[0]
	}

	imageAddress, err := valueObject.NewContainerImageAddress(rawImageAddress)
	if err != nil {
		return containerImage, err
	}

	rawImageDigest, assertOk := rawContainerImage["Digest"].(string)
	if !assertOk {
		return containerImage, errors.New("InvalidContainerImageDigest")
	}
	rawImageDigest = strings.TrimPrefix(rawImageDigest, "sha256:")
	if len(rawImageDigest) > 12 {
		rawImageDigest = rawImageDigest[:12]
	}
	imageHash, err := valueObject.NewHash(rawImageDigest)
	if err != nil {
		return containerImage, err
	}

	rawIsa, assertOk := rawContainerImage["Architecture"].(string)
	if !assertOk {
		return containerImage, errors.New("InvalidContainerImageIsa")
	}
	// TODO: support arm, armv7 and arm64 in the future.
	switch rawIsa {
	case "amd64", "x86-64":
		rawIsa = "amd64"
	default:
		return containerImage, errors.New("UnsupportedContainerImageIsa")
	}
	isa, err := valueObject.NewInstructionSetArchitecture(rawIsa)
	if err != nil {
		return containerImage, err
	}

	rawImageSize, assertOk := rawContainerImage["Size"].(float64)
	if !assertOk {
		return containerImage, errors.New("InvalidContainerImageSize")
	}
	sizeBytes, err := valueObject.NewByte(rawImageSize)
	if err != nil {
		return containerImage, err
	}

	portBindings := []valueObject.PortBinding{}
	rawPortBindings, assertOk := rawConfig["ExposedPorts"].(map[string]interface{})
	if assertOk {
		for rawPortBinding := range rawPortBindings {
			rawPortBinding = strings.ReplaceAll(rawPortBinding, "/tcp", "")
			parsedPortBindings, err := valueObject.NewPortBindingFromString(rawPortBinding)
			if err != nil {
				return containerImage, err
			}

			portBindings = append(portBindings, parsedPortBindings...)
		}
		sort.SliceStable(portBindings, func(i, j int) bool {
			return portBindings[i].PublicPort < portBindings[j].PublicPort
		})
	}

	rawEnvs, assertOk := rawConfig["Env"].([]interface{})
	if !assertOk {
		return containerImage, errors.New("InvalidContainerImageEnv")
	}
	envs := []valueObject.ContainerEnv{}
	for _, rawEnv := range rawEnvs {
		parsedEnv, err := valueObject.NewContainerEnv(rawEnv)
		if err != nil {
			return containerImage, err
		}

		envs = append(envs, parsedEnv)
	}

	rawEntrypoint := ""
	rawEntrypointSlice, assertOk := rawConfig["Entrypoint"].([]interface{})
	if assertOk {
		for _, rawEntrypointItem := range rawEntrypointSlice {
			rawEntrypointPart, assertOk := rawEntrypointItem.(string)
			if !assertOk {
				continue
			}
			rawEntrypoint += rawEntrypointPart + " "
		}
	}
	var entrypointPtr *valueObject.ContainerEntrypoint
	if rawEntrypoint != "" {
		entrypoint, err := valueObject.NewContainerEntrypoint(rawEntrypoint)
		if err != nil {
			return containerImage, err
		}
		entrypointPtr = &entrypoint
	}

	originContainerDetails := entity.Container{}
	rawEncodedOriginContainerDetails, assertOk := rawLabels["ez.originContainerDetails"].(string)
	if assertOk {
		originContainerDetails, err = repo.originContainerDetailsFactory(
			rawEncodedOriginContainerDetails,
		)
		if err != nil {
			return containerImage, errors.New("ParseContainerImageDetailsError: " + err.Error())
		}
	}

	originContainerMappings := []entity.Mapping{}
	rawEncodedOriginContainerMappings, assertOk := rawLabels["ez.originContainerMappings"].(string)
	if assertOk {
		originContainerMappings, err = repo.originContainerMappingsFactory(
			rawEncodedOriginContainerMappings,
		)
		if err != nil {
			return containerImage, errors.New("ParseContainerImageMappingsError: " + err.Error())
		}
	}

	rawCreated, assertOk := rawContainerImage["Created"].(string)
	if !assertOk {
		return containerImage, errors.New("InvalidContainerImageCreated")
	}
	createdTime, err := time.Parse(time.RFC3339Nano, rawCreated)
	if err != nil {
		return containerImage, errors.New("ParseContainerImageCreatedError")
	}
	createdAt := valueObject.NewUnixTimeWithGoTime(createdTime)

	return entity.NewContainerImage(
		imageId, accountId, imageAddress, imageHash, isa, sizeBytes, portBindings, envs,
		entrypointPtr, &originContainerDetails, originContainerMappings, createdAt,
	), nil
}

func (repo *ContainerImageQueryRepo) Read() ([]entity.ContainerImage, error) {
	containerImages := []entity.ContainerImage{}

	readAccountsResponseDto, err := repo.accountQueryRepo.Read(
		dto.ReadAccountsRequest{
			Pagination: dto.PaginationUnpaginated,
		},
	)
	if err != nil {
		return containerImages, err
	}

	for _, accountEntity := range readAccountsResponseDto.Accounts {
		accountIdStr := accountEntity.Id.String()

		rawContainerImagesIdsStr, err := infraHelper.RunCmdAsUser(
			accountEntity.Id, "podman", "images", "--format", "{{.Id}}",
		)
		if err != nil {
			slog.Debug(
				"PodmanListImagesIdError",
				slog.String("accountId", accountIdStr),
				slog.Any("error", err),
			)
			continue
		}

		rawContainerImagesIds := strings.Split(rawContainerImagesIdsStr, "\n")
		if len(rawContainerImagesIds) == 0 {
			continue
		}

		for _, rawContainerImageId := range rawContainerImagesIds {
			if rawContainerImageId == "" {
				continue
			}

			imageId, err := valueObject.NewContainerImageId(rawContainerImageId)
			if err != nil {
				slog.Debug(
					"ContainerImageIdParseError",
					slog.String("accountId", accountIdStr),
					slog.String("rawImageId", rawContainerImageId),
					slog.Any("error", err),
				)
				continue
			}

			containerImage, err := repo.ReadById(accountEntity.Id, imageId)
			if err != nil {
				slog.Debug(
					"ContainerImageReadError",
					slog.String("accountId", accountIdStr),
					slog.String("imageId", imageId.String()),
					slog.Any("error", err),
				)
				continue
			}

			containerImages = append(containerImages, containerImage)
		}
	}

	return containerImages, nil
}

func (repo *ContainerImageQueryRepo) ReadById(
	accountId valueObject.AccountId,
	imageId valueObject.ContainerImageId,
) (containerImage entity.ContainerImage, err error) {
	rawContainerImageAttributesStr, err := infraHelper.RunCmdAsUser(
		accountId, "podman", "inspect", imageId.String(), "--format", "{{json .}}",
	)
	if err != nil {
		return containerImage, err
	}

	rawContainerImageAttributes := map[string]interface{}{}
	err = json.Unmarshal([]byte(rawContainerImageAttributesStr), &rawContainerImageAttributes)
	if err != nil {
		return containerImage, err
	}

	return repo.containerImageFactory(accountId, rawContainerImageAttributes)
}

func (repo *ContainerImageQueryRepo) archiveFileFactory(
	archiveFilePath valueObject.UnixFilePath,
	serverHostname valueObject.Fqdn,
) (archiveFile entity.ContainerImageArchive, err error) {
	archiveFileName := archiveFilePath.ReadFileName()
	archiveFileNameParts := strings.Split(archiveFileName.String(), "-")
	if len(archiveFileNameParts) == 0 {
		return archiveFile, errors.New("SplitArchiveFilePartsError")
	}

	imageId, err := valueObject.NewContainerImageId(archiveFileNameParts[1])
	if err != nil {
		imageId, err = valueObject.NewContainerImageId(archiveFileNameParts[0])
		if err != nil {
			return archiveFile, errors.New("ArchiveFileImageIdParseError")
		}
	}

	fileInfo, err := os.Stat(archiveFilePath.String())
	if err != nil {
		return archiveFile, errors.New("ArchiveFileStatError")
	}

	rawOwnerAccountId := fileInfo.Sys().(*syscall.Stat_t).Uid
	accountId, err := valueObject.NewAccountId(rawOwnerAccountId)
	if err != nil {
		return archiveFile, errors.New("ArchiveFileOwnerAccountIdParseError")
	}

	sizeBytes, err := valueObject.NewByte(fileInfo.Size())
	if err != nil {
		return archiveFile, errors.New("ArchiveFileSizeBytesParseError")
	}

	downloadUrl, _ := valueObject.NewUrl(
		"https://" + serverHostname.String() + "/api/v1/container/image/archive/" +
			accountId.String() + "/" + imageId.String() + "/",
	)

	rawCreatedAt := fileInfo.ModTime()
	createdAt := valueObject.NewUnixTimeWithGoTime(rawCreatedAt)

	return entity.NewContainerImageArchive(
		imageId, accountId, archiveFilePath, sizeBytes, &downloadUrl, nil, createdAt,
	), nil
}

func (repo *ContainerImageQueryRepo) ReadArchives(
	requestDto dto.ReadContainerImageArchivesRequest,
) (responseDto dto.ReadContainerImageArchivesResponse, err error) {
	archiveFiles := []entity.ContainerImageArchive{}

	archiveFilesBaseDirectoryStr := infraEnvs.UserDataDirectory
	findPathFlagValue := "*/archives/*"
	if requestDto.ArchivesDirectory != nil {
		archiveFilesBaseDirectoryStr = requestDto.ArchivesDirectory.String()
		findPathFlagValue = "*"
	}

	findResult, err := infraHelper.RunCmd(
		"find", archiveFilesBaseDirectoryStr,
		"-type", "f",
		"-path", findPathFlagValue,
		"-maxdepth", "3",
		"-regex", `.*\.\(`+strings.Join(valueObject.ValidCompressionFormats, `\|`)+`\)$`,
	)
	if err != nil {
		return responseDto, errors.New("FindArchiveFilesError: " + err.Error())
	}

	rawArchiveFilesPaths := strings.Split(findResult, "\n")
	if len(rawArchiveFilesPaths) == 0 {
		return responseDto, nil
	}

	serverHostname, err := infraHelper.ReadServerHostname()
	if err != nil {
		return responseDto, errors.New("InvalidServerHostname: " + err.Error())
	}

	for _, rawArchiveFilePath := range rawArchiveFilesPaths {
		if rawArchiveFilePath == "" {
			continue
		}

		archiveFilePath, err := valueObject.NewUnixFilePath(rawArchiveFilePath)
		if err != nil {
			slog.Debug(err.Error(), slog.String("path", rawArchiveFilePath))
			continue
		}

		archiveFile, err := repo.archiveFileFactory(archiveFilePath, serverHostname)
		if err != nil {
			slog.Debug(err.Error(), slog.String("path", rawArchiveFilePath))
			continue
		}
		archiveFiles = append(archiveFiles, archiveFile)
	}

	return dto.ReadContainerImageArchivesResponse{
		Pagination: requestDto.Pagination,
		Archives:   archiveFiles,
	}, nil
}

func (repo *ContainerImageQueryRepo) ReadArchive(
	readDto dto.ReadContainerImageArchive,
) (archiveFile entity.ContainerImageArchive, err error) {
	accountEntity, err := repo.accountQueryRepo.ReadFirst(dto.ReadAccountsRequest{
		AccountId: &readDto.AccountId,
	})
	if err != nil {
		return archiveFile, err
	}

	archiveDirStr := accountEntity.HomeDirectory.String() + "/archives"
	rawArchiveFilePath, err := infraHelper.RunCmdAsUser(
		readDto.AccountId,
		"find", archiveDirStr, "-type", "f", "-name", "*"+readDto.ImageId.String()+"*",
	)
	if err != nil {
		return archiveFile, errors.New("FindArchiveFileError: " + err.Error())
	}
	if len(rawArchiveFilePath) == 0 {
		return archiveFile, errors.New("ArchiveFileNotFound")
	}

	rawArchiveFilePathLines := strings.Split(rawArchiveFilePath, "\n")
	if len(rawArchiveFilePathLines) == 0 {
		return archiveFile, errors.New("ArchiveFileNotFound")
	}

	archiveFilePath, err := valueObject.NewUnixFilePath(rawArchiveFilePathLines[0])
	if err != nil {
		return archiveFile, err
	}

	serverHostname, err := infraHelper.ReadServerHostname()
	if err != nil {
		return archiveFile, errors.New("InvalidServerHostname: " + err.Error())
	}

	return repo.archiveFileFactory(archiveFilePath, serverHostname)
}
