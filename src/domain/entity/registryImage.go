package entity

import "github.com/speedianet/control/src/domain/valueObject"

type RegistryImage struct {
	Name          valueObject.RegistryImageName            `json:"name"`
	Description   *valueObject.RegistryImageDescription    `json:"description"`
	RegistryName  valueObject.RegistryName                 `json:"registryName"`
	PublisherName valueObject.RegistryPublisherName        `json:"publisherName"`
	ImageAddress  valueObject.ContainerImageAddress        `json:"imageAddress"`
	Isas          []valueObject.InstructionSetArchitecture `json:"isas"`
	PullCount     uint64                                   `json:"pullCount"`
	StarCount     *uint64                                  `json:"starCount"`
	LogoUrl       *valueObject.Url                         `json:"logoUrl"`
	CreatedAt     *valueObject.UnixTime                    `json:"createdAt"`
	UpdatedAt     *valueObject.UnixTime                    `json:"updatedAt"`
}

func NewRegistryImage(
	name valueObject.RegistryImageName,
	description *valueObject.RegistryImageDescription,
	registryName valueObject.RegistryName,
	publisherName valueObject.RegistryPublisherName,
	imageAddress valueObject.ContainerImageAddress,
	isas []valueObject.InstructionSetArchitecture,
	pullCount uint64,
	starCount *uint64,
	logoUrl *valueObject.Url,
	createdAt *valueObject.UnixTime,
	updatedAt *valueObject.UnixTime,
) RegistryImage {
	return RegistryImage{
		Name:          name,
		Description:   description,
		RegistryName:  registryName,
		PublisherName: publisherName,
		ImageAddress:  imageAddress,
		Isas:          isas,
		PullCount:     pullCount,
		StarCount:     starCount,
		LogoUrl:       logoUrl,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}
