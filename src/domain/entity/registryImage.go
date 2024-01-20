package entity

import "github.com/speedianet/control/src/domain/valueObject"

type RegistryImage struct {
	Name          valueObject.RegistryImageName            `json:"name"`
	Description   valueObject.RegistryImageDescription     `json:"description"`
	RegistryName  valueObject.RegistryName                 `json:"registryName"`
	PublisherName valueObject.RegistryPublisherName        `json:"publisherName"`
	ImageAddress  valueObject.ContainerImageAddress        `json:"imageAddress"`
	Isa           []valueObject.InstructionSetArchitecture `json:"isa"`
	IsVerified    bool                                     `json:"isVerified"`
	PullCount     uint64                                   `json:"pullCount"`
	StarCount     *uint64                                  `json:"starCount"`
	LogoUrl       *valueObject.Url                         `json:"logoUrl"`
	CreatedAt     *valueObject.UnixTime                    `json:"createdAt"`
	UpdatedAt     *valueObject.UnixTime                    `json:"updatedAt"`
}
