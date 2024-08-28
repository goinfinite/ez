package entity

import "github.com/speedianet/control/src/domain/valueObject"

type ContainerImage struct {
	Id           valueObject.ContainerImageId           `json:"id"`
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
