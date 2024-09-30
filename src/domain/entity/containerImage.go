package entity

import (
	"encoding/json"

	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ContainerImage struct {
	Id           valueObject.ContainerImageId           `json:"id"`
	AccountId    valueObject.AccountId                  `json:"accountId"`
	ImageAddress valueObject.ContainerImageAddress      `json:"imageAddress"`
	ImageHash    valueObject.Hash                       `json:"imageHash"`
	Isa          valueObject.InstructionSetArchitecture `json:"isa"`
	SizeBytes    valueObject.Byte                       `json:"sizeBytes"`
	PortBindings []valueObject.PortBinding              `json:"portBindings"`
	Envs         []valueObject.ContainerEnv             `json:"envs"`
	Entrypoint   *valueObject.ContainerEntrypoint       `json:"entrypoint"`
	CreatedAt    valueObject.UnixTime                   `json:"createdAt"`
}

func NewContainerImage(
	id valueObject.ContainerImageId,
	accountId valueObject.AccountId,
	imageAddress valueObject.ContainerImageAddress,
	imageHash valueObject.Hash,
	isa valueObject.InstructionSetArchitecture,
	sizeBytes valueObject.Byte,
	portBindings []valueObject.PortBinding,
	envs []valueObject.ContainerEnv,
	entrypoint *valueObject.ContainerEntrypoint,
	createdAt valueObject.UnixTime,
) ContainerImage {
	return ContainerImage{
		Id:           id,
		AccountId:    accountId,
		ImageAddress: imageAddress,
		ImageHash:    imageHash,
		Isa:          isa,
		SizeBytes:    sizeBytes,
		PortBindings: portBindings,
		Envs:         envs,
		Entrypoint:   entrypoint,
		CreatedAt:    createdAt,
	}
}

func (entity ContainerImage) JsonSerialize() string {
	jsonBytes, _ := json.Marshal(entity)
	return string(jsonBytes)
}
