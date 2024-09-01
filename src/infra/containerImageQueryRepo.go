package infra

import (
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	infraHelper "github.com/speedianet/control/src/infra/helper"
)

type ContainerImageQueryRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewContainerImageQueryRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerImageQueryRepo {
	return &ContainerImageQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *ContainerImageQueryRepo) containerImageFactory(
	accountId valueObject.AccountId,
	rawContainerImage map[string]interface{},
) (containerImage entity.ContainerImage, err error) {
	rawImageId, assertOk := rawContainerImage["Id"].(string)
	if !assertOk {
		return containerImage, errors.New("InvalidContainerImageId")
	}
	if len(rawImageId) > 12 {
		rawImageId = rawImageId[:12]
	}
	imageId, err := valueObject.NewContainerImageId(rawImageId)
	if err != nil {
		return containerImage, err
	}

	rawImageNames, assertOk := rawContainerImage["NamesHistory"].([]interface{})
	if !assertOk {
		return containerImage, errors.New("InvalidContainerImageNames")
	}
	if len(rawImageNames) == 0 {
		return containerImage, errors.New("EmptyContainerImageNames")
	}

	imageAddressStr, assertOk := rawImageNames[0].(string)
	if !assertOk {
		return containerImage, errors.New("InvalidContainerImageAddress")
	}
	imageAddress, err := valueObject.NewContainerImageAddress(imageAddressStr)
	if err != nil {
		return containerImage, err
	}

	rawImageDigest, assertOk := rawContainerImage["Digest"].(string)
	if !assertOk {
		return containerImage, errors.New("InvalidContainerImageDigest")
	}
	rawImageDigest = strings.TrimPrefix(rawImageDigest, "sha256:")
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

	rawConfig, assertOk := rawContainerImage["Config"].(map[string]interface{})
	if !assertOk {
		return containerImage, errors.New("InvalidContainerImageConfig")
	}

	rawPortBindings, assertOk := rawConfig["ExposedPorts"].(map[string]interface{})
	if !assertOk {
		return containerImage, errors.New("InvalidContainerImagePortBindings")
	}
	portBindings := []valueObject.PortBinding{}
	for rawPortBinding := range rawPortBindings {
		rawPortBinding = strings.ReplaceAll(rawPortBinding, "/tcp", "")
		parsedPortBindings, err := valueObject.NewPortBindingFromString(rawPortBinding)
		if err != nil {
			return containerImage, err
		}

		portBindings = append(portBindings, parsedPortBindings...)
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

	rawEntrypointSlice, assertOk := rawConfig["Entrypoint"].([]interface{})
	if !assertOk {
		return containerImage, errors.New("InvalidContainerImageEntrypoint")
	}
	rawEntrypoint := ""
	for _, rawEntrypointItem := range rawEntrypointSlice {
		rawEntrypointPart, assertOk := rawEntrypointItem.(string)
		if !assertOk {
			continue
		}
		rawEntrypoint += rawEntrypointPart + " "
	}
	var entrypointPtr *valueObject.ContainerEntrypoint
	if rawEntrypoint != "" {
		entrypoint, err := valueObject.NewContainerEntrypoint(rawEntrypoint)
		if err != nil {
			return containerImage, err
		}
		entrypointPtr = &entrypoint
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
		imageId, accountId, imageAddress, imageHash, isa, sizeBytes,
		portBindings, envs, entrypointPtr, createdAt,
	), nil
}

func (repo *ContainerImageQueryRepo) Read() ([]entity.ContainerImage, error) {
	containerImages := []entity.ContainerImage{}

	accountsList, err := NewAccountQueryRepo(repo.persistentDbSvc).Read()
	if err != nil {
		return containerImages, err
	}

	for _, account := range accountsList {
		rawContainerImagesIdsStr, err := infraHelper.RunCmdAsUser(
			account.Id, "podman", "images", "--format", "{{.Id}}",
		)
		if err != nil {
			slog.Debug(
				"PodmanListImagesIdError",
				slog.String("accountId", account.Id.String()),
				slog.Any("error", err),
			)
			continue
		}

		rawContainerImagesIds := strings.Split(rawContainerImagesIdsStr, "\n")
		if len(rawContainerImagesIds) == 0 {
			continue
		}

		accountIdStr := account.Id.String()
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

			containerImage, err := repo.ReadById(account.Id, imageId)
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

func (repo *ContainerImageQueryRepo) ReadArchiveFiles() (
	[]entity.ContainerImageArchiveFile, error,
) {
	return []entity.ContainerImageArchiveFile{}, errors.New("NotImplemented")
}

func (repo *ContainerImageQueryRepo) ReadArchiveFile(
	readDto dto.ReadContainerImageArchiveFile,
) (archiveFile entity.ContainerImageArchiveFile, err error) {
	return archiveFile, errors.New("NotImplemented")
}
