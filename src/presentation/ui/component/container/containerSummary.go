package componentContainer

import (
	"encoding/json"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
)

type ContainerSummary struct {
	ContainerId      valueObject.ContainerId           `json:"containerId"`
	Hostname         valueObject.Fqdn                  `json:"hostname"`
	ImageAddress     valueObject.ContainerImageAddress `json:"imageAddress"`
	AccountId        valueObject.AccountId             `json:"accountId"`
	AccountUsername  valueObject.Username              `json:"accountUsername"`
	ProfileId        valueObject.ContainerProfileId    `json:"profileId"`
	ProfileName      valueObject.ContainerProfileName  `json:"profileName"`
	ProfileBaseSpecs valueObject.ContainerSpecs        `json:"profileBaseSpecs"`
}

func NewContainerSummary(
	containerId valueObject.ContainerId,
	hostname valueObject.Fqdn,
	imageAddress valueObject.ContainerImageAddress,
	accountId valueObject.AccountId,
	accountUsername valueObject.Username,
	profileId valueObject.ContainerProfileId,
	profileName valueObject.ContainerProfileName,
	profileBaseSpecs valueObject.ContainerSpecs,
) ContainerSummary {
	return ContainerSummary{
		ContainerId:      containerId,
		Hostname:         hostname,
		ImageAddress:     imageAddress,
		AccountId:        accountId,
		AccountUsername:  accountUsername,
		ProfileId:        profileId,
		ProfileName:      profileName,
		ProfileBaseSpecs: profileBaseSpecs,
	}
}

func (summary ContainerSummary) JsonSerialize() string {
	jsonBytes, _ := json.Marshal(summary)
	return string(jsonBytes)
}

func NewContainerSummaries(
	containerEntities []entity.Container,
	profileEntities []entity.ContainerProfile,
	accountEntities []entity.Account,
) []ContainerSummary {
	containerIdEntityMap := map[valueObject.ContainerId]entity.Container{}
	for _, containerEntity := range containerEntities {
		containerIdEntityMap[containerEntity.Id] = containerEntity
	}

	profileIdEntityMap := map[valueObject.ContainerProfileId]entity.ContainerProfile{}
	for _, profileEntity := range profileEntities {
		profileIdEntityMap[profileEntity.Id] = profileEntity
	}

	accountIdEntityMap := map[valueObject.AccountId]entity.Account{}
	for _, accountEntity := range accountEntities {
		accountIdEntityMap[accountEntity.Id] = accountEntity
	}

	return NewContainerSummariesWithMaps(
		containerIdEntityMap, profileIdEntityMap, accountIdEntityMap,
	)
}

func NewContainerSummariesWithMaps(
	containerIdEntityMap map[valueObject.ContainerId]entity.Container,
	profileIdEntityMap map[valueObject.ContainerProfileId]entity.ContainerProfile,
	accountIdEntityMap map[valueObject.AccountId]entity.Account,
) []ContainerSummary {
	containerSummaries := []ContainerSummary{}

	for _, containerEntity := range containerIdEntityMap {
		profileEntity, exists := profileIdEntityMap[containerEntity.ProfileId]
		if !exists {
			continue
		}

		accountEntity, exists := accountIdEntityMap[containerEntity.AccountId]
		if !exists {
			continue
		}

		containerSummary := NewContainerSummary(
			containerEntity.Id, containerEntity.Hostname, containerEntity.ImageAddress,
			containerEntity.AccountId, accountEntity.Username, containerEntity.ProfileId,
			profileEntity.Name, profileEntity.BaseSpecs,
		)

		containerSummaries = append(containerSummaries, containerSummary)
	}

	return containerSummaries
}
