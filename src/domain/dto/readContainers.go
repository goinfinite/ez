package dto

import (
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type ReadContainersRequest struct {
	Pagination               Pagination                          `json:"pagination"`
	ContainerId              []valueObject.ContainerId           `json:"containerId,omitempty"`
	ContainerAccountId       []valueObject.AccountId             `json:"containerAccountId,omitempty"`
	ExceptContainerId        []valueObject.ContainerId           `json:"exceptContainerId,omitempty"`
	ExceptContainerAccountId []valueObject.AccountId             `json:"exceptContainerAccountId,omitempty"`
	ContainerHostname        *valueObject.Fqdn                   `json:"containerHostname,omitempty"`
	ContainerStatus          *bool                               `json:"containerStatus,omitempty"`
	ContainerImageId         *valueObject.ContainerImageId       `json:"containerImageId,omitempty"`
	ContainerImageAddress    *valueObject.ContainerImageAddress  `json:"containerImageAddress,omitempty"`
	ContainerImageHash       *valueObject.Hash                   `json:"containerImageHash,omitempty"`
	ContainerPortBindings    []valueObject.PortBinding           `json:"containerPortBindings,omitempty"`
	ContainerRestartPolicy   *valueObject.ContainerRestartPolicy `json:"containerRestartPolicy,omitempty"`
	ContainerProfileId       *valueObject.ContainerProfileId     `json:"containerProfileId,omitempty"`
	ContainerEnv             []valueObject.ContainerEnv          `json:"containerEnv,omitempty"`
	CreatedBeforeAt          *valueObject.UnixTime               `json:"createdBeforeAt,omitempty"`
	CreatedAfterAt           *valueObject.UnixTime               `json:"createdAfterAt,omitempty"`
	StartedBeforeAt          *valueObject.UnixTime               `json:"startedBeforeAt,omitempty"`
	StartedAfterAt           *valueObject.UnixTime               `json:"startedAfterAt,omitempty"`
	StoppedBeforeAt          *valueObject.UnixTime               `json:"stoppedBeforeAt,omitempty"`
	StoppedAfterAt           *valueObject.UnixTime               `json:"stoppedAfterAt,omitempty"`
	WithMetrics              *bool                               `json:"withMetrics,omitempty"`
}

type ReadContainersResponse struct {
	Pagination            Pagination             `json:"pagination"`
	Containers            []entity.Container     `json:"containers"`
	ContainersWithMetrics []ContainerWithMetrics `json:"containersWithMetrics"`
}
