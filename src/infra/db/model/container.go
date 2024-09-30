package dbModel

import (
	"log/slog"
	"strings"
	"time"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type Container struct {
	ID            string `gorm:"primarykey"`
	AccountID     uint64 `gorm:"not null"`
	Hostname      string `gorm:"not null"`
	Status        bool   `gorm:"not null"`
	ImageId       string `gorm:"not null"`
	ImageAddress  string `gorm:"not null"`
	ImageHash     string `gorm:"not null"`
	PortBindings  []ContainerPortBinding
	RestartPolicy string `gorm:"not null"`
	RestartCount  uint64
	Entrypoint    *string
	ProfileID     uint64 `gorm:"not null"`
	Envs          *string
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
	StartedAt     *time.Time
	StoppedAt     *time.Time
}

func (Container) TableName() string {
	return "containers"
}

func (Container) JoinEnvs(envs []valueObject.ContainerEnv) string {
	envsStr := ""
	for _, env := range envs {
		envsStr += env.String() + ";"
	}

	return strings.TrimSuffix(envsStr, ";")
}

func (Container) SplitEnvs(envsStr string) []valueObject.ContainerEnv {
	rawEnvsList := strings.Split(envsStr, ";")
	var envs []valueObject.ContainerEnv
	for _, rawEnv := range rawEnvsList {
		env, err := valueObject.NewContainerEnv(rawEnv)
		if err != nil {
			slog.Debug(
				"CreateContainerEnvError",
				slog.String("rawEnv", rawEnv), slog.Any("error", err),
			)
			continue
		}
		envs = append(envs, env)
	}

	return envs
}

func (model Container) ToEntity() (containerEntity entity.Container, err error) {
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

	imageId, err := valueObject.NewContainerImageId(model.ImageId)
	if err != nil {
		return containerEntity, err
	}

	imageAddress, err := valueObject.NewContainerImageAddress(model.ImageAddress)
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

	createdAt := valueObject.NewUnixTimeWithGoTime(model.CreatedAt)
	updatedAt := valueObject.NewUnixTimeWithGoTime(model.UpdatedAt)

	var startedAtPtr *valueObject.UnixTime
	if model.StartedAt != nil {
		startedAt := valueObject.NewUnixTimeWithGoTime(*model.StartedAt)
		startedAtPtr = &startedAt
	}

	var stoppedAtPtr *valueObject.UnixTime
	if model.StoppedAt != nil {
		stoppedAt := valueObject.NewUnixTimeWithGoTime(*model.StoppedAt)
		stoppedAtPtr = &stoppedAt
	}

	profileId, err := valueObject.NewContainerProfileId(model.ProfileID)
	if err != nil {
		return containerEntity, err
	}

	envs := []valueObject.ContainerEnv{}
	if model.Envs != nil {
		envs = model.SplitEnvs(*model.Envs)
	}

	return entity.NewContainer(
		id, accountId, hostname, model.Status, imageId, imageAddress, imageHash,
		portBindings, restartPolicy, model.RestartCount, entryPointPtr, profileId,
		envs, createdAt, updatedAt, startedAtPtr, stoppedAtPtr,
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
		startedAt := time.Unix(entity.StartedAt.Read(), 0)
		startedAtPtr = &startedAt
	}

	var envsPtr *string
	if len(entity.Envs) > 0 {
		envs := Container{}.JoinEnvs(entity.Envs)
		envsPtr = &envs
	}

	return Container{
		ID:            entity.Id.String(),
		AccountID:     entity.AccountId.Uint64(),
		Hostname:      entity.Hostname.String(),
		Status:        entity.Status,
		ImageId:       entity.ImageId.String(),
		ImageAddress:  entity.ImageAddress.String(),
		ImageHash:     entity.ImageHash.String(),
		PortBindings:  portBindingModels,
		RestartPolicy: entity.RestartPolicy.String(),
		RestartCount:  entity.RestartCount,
		Entrypoint:    entrypointStrPtr,
		CreatedAt:     time.Unix(entity.CreatedAt.Read(), 0),
		StartedAt:     startedAtPtr,
		ProfileID:     entity.ProfileId.Uint64(),
		Envs:          envsPtr,
	}
}
