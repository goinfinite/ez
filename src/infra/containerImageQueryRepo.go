package infra

import (
	"encoding/json"
	"errors"
	"log/slog"
	"strings"

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

	rawImageNames, assertOk := rawContainerImage["Names"].([]interface{})
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

	isa, err := valueObject.NewInstructionSetArchitecture("amd64")
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

	rawCreatedAt, assertOk := rawContainerImage["Created"].(float64)
	if !assertOk {
		return containerImage, errors.New("InvalidContainerImageCreatedAt")
	}
	createdAt, err := valueObject.NewUnixTime(rawCreatedAt)
	if err != nil {
		return containerImage, err
	}

	return entity.NewContainerImage(
		imageId, imageAddress, imageHash, isa, sizeBytes, nil, createdAt,
	), nil
}

func (repo *ContainerImageQueryRepo) Read() ([]entity.ContainerImage, error) {
	containerImages := []entity.ContainerImage{}

	accountsList, err := NewAccountQueryRepo(repo.persistentDbSvc).Read()
	if err != nil {
		return containerImages, err
	}

	for _, account := range accountsList {
		containerImagesStr, err := infraHelper.RunCmdAsUser(
			account.Id, "podman", "images", "--format", "json",
		)
		if err != nil {
			slog.Debug(
				"PodmanListImagesError",
				slog.String("accountId", account.Id.String()),
				slog.Any("error", err),
			)
			continue
		}

		rawContainerImages := []map[string]interface{}{}
		err = json.Unmarshal([]byte(containerImagesStr), &rawContainerImages)
		if err != nil {
			slog.Debug(
				"UnmarshalContainerImagesError",
				slog.String("accountId", account.Id.String()),
				slog.Any("error", err),
			)
		}

		for _, rawContainerImage := range rawContainerImages {
			containerImage, err := repo.containerImageFactory(rawContainerImage)
			if err != nil {
				slog.Debug(
					"ContainerImageFactoryError",
					slog.String("accountId", account.Id.String()),
					slog.Any("error", err),
				)
				continue
			}
			containerImages = append(containerImages, containerImage)
		}
	}

	return containerImages, nil
}
