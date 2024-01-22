package entity

import "github.com/speedianet/control/src/domain/valueObject"

type RegistryTaggedImage struct {
	TagName       valueObject.RegistryImageTag           `json:"tagName"`
	ImageName     valueObject.RegistryImageName          `json:"imageName"`
	PublisherName valueObject.RegistryPublisherName      `json:"publisherName"`
	RegistryName  valueObject.RegistryName               `json:"registryName"`
	ImageAddress  valueObject.ContainerImageAddress      `json:"imageAddress"`
	ImageHash     valueObject.Hash                       `json:"imageHash"`
	Isa           valueObject.InstructionSetArchitecture `json:"isa"`
	SizeBytes     valueObject.Byte                       `json:"sizeBytes"`
	PortBindings  []valueObject.PortBinding              `json:"portBindings"`
	UpdatedAt     valueObject.UnixTime                   `json:"updatedAt"`
}

func NewRegistryTaggedImage(
	tagName valueObject.RegistryImageTag,
	imageName valueObject.RegistryImageName,
	publisherName valueObject.RegistryPublisherName,
	registryName valueObject.RegistryName,
	imageAddress valueObject.ContainerImageAddress,
	imageHash valueObject.Hash,
	isa valueObject.InstructionSetArchitecture,
	sizeBytes valueObject.Byte,
	portBindings []valueObject.PortBinding,
	updatedAt valueObject.UnixTime,
) RegistryTaggedImage {
	return RegistryTaggedImage{
		TagName:       tagName,
		ImageName:     imageName,
		PublisherName: publisherName,
		RegistryName:  registryName,
		ImageAddress:  imageAddress,
		ImageHash:     imageHash,
		Isa:           isa,
		SizeBytes:     sizeBytes,
		PortBindings:  portBindings,
		UpdatedAt:     updatedAt,
	}
}
