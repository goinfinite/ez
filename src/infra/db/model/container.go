package dbModel

import (
	"errors"
	"strings"
	"time"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type Container struct {
	ID            string `gorm:"primarykey"`
	AccountID     uint   `gorm:"not null"`
	Hostname      string `gorm:"not null"`
	Status        bool   `gorm:"not null"`
	ImageAddr     string `gorm:"not null"`
	ImageHash     string `gorm:"not null"`
	PortBindings  []ContainerPortBinding
	RestartPolicy string `gorm:"not null"`
	RestartCount  uint
	Entrypoint    *string
	ProfileID     uint `gorm:"not null"`
	Envs          *string
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
	StartedAt     *time.Time
	StoppedAt     *time.Time
}

func (Container) TableName() string {
	return "containers"
}

func (model Container) ToEntity() (entity.Container, error) {
	var containerEntity entity.Container

	id, err := valueObject.NewContainerId(model.ID)
	if err != nil {
		return containerEntity, err
	}

	accountId, err := valueObject.NewAccountId(model.AccountID)
	if err != nil {
		return containerEntity, err
	}

	hostname, err := valueObject.NewFqdn(model.Hostname)
	if err != nil {
		return containerEntity, err
	}

	imageAddr, err := valueObject.NewContainerImageAddress(model.ImageAddr)
	if err != nil {
		return containerEntity, err
	}

	imageHash, err := valueObject.NewHash(model.ImageHash)
	if err != nil {
		return containerEntity, err
	}

	var portBindings []valueObject.PortBinding
	if len(model.PortBindings) > 0 {
		for _, portBindingModel := range model.PortBindings {
			portBinding, err := portBindingModel.ToValueObject()
			if err != nil {
				return containerEntity, err
			}
			portBindings = append(portBindings, portBinding)
		}
	}

	restartPolicy, err := valueObject.NewContainerRestartPolicy(model.RestartPolicy)
	if err != nil {
		return containerEntity, err
	}

	var entryPointPtr *valueObject.ContainerEntrypoint
	if model.Entrypoint != nil {
		entryPoint, err := valueObject.NewContainerEntrypoint(*model.Entrypoint)
		if err != nil {
			return containerEntity, err
		}
		entryPointPtr = &entryPoint
	}

	createdAt := valueObject.UnixTime(model.CreatedAt.Unix())
	updatedAt := valueObject.UnixTime(model.UpdatedAt.Unix())

	var startedAtPtr *valueObject.UnixTime
	if model.StartedAt != nil {
		startedAt := valueObject.UnixTime(model.StartedAt.Unix())
		startedAtPtr = &startedAt
	}

	var stoppedAtPtr *valueObject.UnixTime
	if model.StoppedAt != nil {
		stoppedAt := valueObject.UnixTime(model.StoppedAt.Unix())
		stoppedAtPtr = &stoppedAt
	}

	profileId, err := valueObject.NewContainerProfileId(model.ProfileID)
	if err != nil {
		return containerEntity, err
	}

	envs := []valueObject.ContainerEnv{}
	if model.Envs != nil {
		envsParts := strings.Split(*model.Envs, ";")
		if len(envsParts) == 0 {
			return containerEntity, errors.New("InvalidEnvs")
		}

		for _, envPart := range envsParts {
			env, err := valueObject.NewContainerEnv(envPart)
			if err != nil {
				return containerEntity, err
			}
			envs = append(envs, env)
		}
	}

	return entity.NewContainer(
		id,
		accountId,
		hostname,
		model.Status,
		imageAddr,
		imageHash,
		portBindings,
		restartPolicy,
		uint64(model.RestartCount),
		entryPointPtr,
		profileId,
		envs,
		createdAt,
		updatedAt,
		startedAtPtr,
		stoppedAtPtr,
	), nil
}

func (Container) ToModel(entity entity.Container) Container {
	portBindingModels := []ContainerPortBinding{}
	for _, portBinding := range entity.PortBindings {
		portBindingModel := ContainerPortBinding{}.ToModel(
			entity.Id,
			portBinding,
		)
		portBindingModels = append(portBindingModels, portBindingModel)
	}

	var entrypointStrPtr *string
	if entity.Entrypoint != nil {
		entrypointStr := entity.Entrypoint.String()
		entrypointStrPtr = &entrypointStr
	}

	var startedAtPtr *time.Time
	if entity.StartedAt != nil {
		startedAt := time.Unix(entity.StartedAt.Get(), 0)
		startedAtPtr = &startedAt
	}

	var envsPtr *string
	if len(entity.Envs) > 0 {
		envs := ""
		for _, env := range entity.Envs {
			envs += env.String() + ";"
		}
		envsPtr = &envs
	}

	return Container{
		ID:            entity.Id.String(),
		AccountID:     uint(entity.AccountId.Get()),
		Hostname:      entity.Hostname.String(),
		Status:        entity.Status,
		ImageAddr:     entity.ImageAddr.String(),
		ImageHash:     entity.ImageHash.String(),
		PortBindings:  portBindingModels,
		RestartPolicy: entity.RestartPolicy.String(),
		RestartCount:  uint(entity.RestartCount),
		Entrypoint:    entrypointStrPtr,
		CreatedAt:     time.Unix(entity.CreatedAt.Get(), 0),
		StartedAt:     startedAtPtr,
		ProfileID:     uint(entity.ProfileId.Get()),
		Envs:          envsPtr,
	}
}
