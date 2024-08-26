package entity

import "github.com/speedianet/control/src/domain/valueObject"

type ContainerImage struct {
	ImageAddress       valueObject.ContainerImageAddress        `json:"imageAddress"`
	ImageHash          valueObject.Hash                         `json:"imageHash"`
	Isa                valueObject.InstructionSetArchitecture   `json:"isa"`
	SizeBytes          valueObject.Byte                         `json:"sizeBytes"`
	SnapshotAttributes *valueObject.ContainerSnapshotAttributes `json:"snapshotAttributes"`
	CreatedAt          valueObject.UnixTime                     `json:"createdAt"`
}

func NewContainerImage(
	imageAddress valueObject.ContainerImageAddress,
	imageHash valueObject.Hash,
	isa valueObject.InstructionSetArchitecture,
	sizeBytes valueObject.Byte,
	snapshotAttributes *valueObject.ContainerSnapshotAttributes,
	createdAt valueObject.UnixTime,
) ContainerImage {
	return ContainerImage{
		ImageAddress:       imageAddress,
		ImageHash:          imageHash,
		Isa:                isa,
		SizeBytes:          sizeBytes,
		SnapshotAttributes: snapshotAttributes,
		CreatedAt:          createdAt,
	}
}
