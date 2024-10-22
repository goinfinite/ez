package entity

import "github.com/goinfinite/ez/src/domain/valueObject"

type RegistryImage struct {
	Name              valueObject.RegistryImageName            `json:"name"`
	PublisherName     valueObject.RegistryPublisherName        `json:"publisherName"`
	RegistryName      valueObject.RegistryName                 `json:"registryName"`
	ImageAddress      valueObject.ContainerImageAddress        `json:"imageAddress"`
	Description       *valueObject.RegistryImageDescription    `json:"description"`
	Isas              []valueObject.InstructionSetArchitecture `json:"isas"`
	PullCount         uint64                                   `json:"pullCount"`
	StarCount         *uint64                                  `json:"starCount"`
	LogoUrl           *valueObject.Url                         `json:"logoUrl"`
	CreatedAt         *valueObject.UnixTime                    `json:"createdAt"`
	UpdatedAt         *valueObject.UnixTime                    `json:"updatedAt"`
	UpdatedAtRelative *valueObject.RelativeTime                `json:"updatedAtRelative"`
}

func NewRegistryImage(
	name valueObject.RegistryImageName,
	publisherName valueObject.RegistryPublisherName,
	registryName valueObject.RegistryName,
	imageAddress valueObject.ContainerImageAddress,
	description *valueObject.RegistryImageDescription,
	isas []valueObject.InstructionSetArchitecture,
	pullCount uint64,
	starCount *uint64,
	logoUrl *valueObject.Url,
	createdAt, updatedAt *valueObject.UnixTime,
	updatedAtRelative *valueObject.RelativeTime,
) RegistryImage {
	return RegistryImage{
		Name:              name,
		PublisherName:     publisherName,
		RegistryName:      registryName,
		ImageAddress:      imageAddress,
		Description:       description,
		Isas:              isas,
		PullCount:         pullCount,
		StarCount:         starCount,
		LogoUrl:           logoUrl,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		UpdatedAtRelative: updatedAtRelative,
	}
}
