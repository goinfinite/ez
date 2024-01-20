package infra

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
	"github.com/speedianet/control/src/infra/db"
)

type ContainerRegistryQueryRepo struct {
	dbSvc *db.DatabaseService
}

func NewContainerRegistryQueryRepo(dbSvc *db.DatabaseService) *ContainerRegistryQueryRepo {
	return &ContainerRegistryQueryRepo{dbSvc: dbSvc}
}

func (repo ContainerRegistryQueryRepo) dockerHubImageFactory(
	imageMap map[string]interface{},
) (entity.RegistryImage, error) {
	var registryImage entity.RegistryImage

	rawImageName, assertOk := imageMap["name"].(string)
	if !assertOk {
		return registryImage, errors.New("ParseImageNameError")
	}
	imageName, err := valueObject.NewRegistryImageName(rawImageName)
	if err != nil {
		return registryImage, err
	}

	var descriptionPtr *valueObject.RegistryImageDescription
	rawDescription, assertOk := imageMap["short_description"].(string)
	if !assertOk {
		return registryImage, errors.New("ParseImageDescriptionError")
	}
	if rawDescription != "" {
		description, err := valueObject.NewRegistryImageDescription(rawDescription)
		if err != nil {
			return registryImage, err
		}
		descriptionPtr = &description
	}

	registryName, _ := valueObject.NewRegistryName("DockerHub")

	rawPublisherDetails, assertOk := imageMap["publisher"].(map[string]interface{})
	if !assertOk {
		return registryImage, errors.New("ParsePublisherDetailsError")
	}
	rawPublisherName, assertOk := rawPublisherDetails["name"].(string)
	if !assertOk {
		return registryImage, errors.New("ParsePublisherNameError")
	}
	publisherName, err := valueObject.NewRegistryPublisherName(rawPublisherName)
	if err != nil {
		return registryImage, err
	}

	rawImageAddress, assertOk := imageMap["slug"].(string)
	if !assertOk {
		return registryImage, errors.New("ParseImageAddressError")
	}
	imageAddress, err := valueObject.NewContainerImageAddress(rawImageAddress)
	if err != nil {
		return registryImage, err
	}

	rawIsas, assertOk := imageMap["architectures"].([]interface{})
	if !assertOk {
		return registryImage, errors.New("ParseIsasError")
	}

	isas := []valueObject.InstructionSetArchitecture{}
	for _, rawIsa := range rawIsas {
		rawIsaMap, assertOk := rawIsa.(map[string]interface{})
		if !assertOk {
			return registryImage, errors.New("ParseIsaError")
		}

		rawIsaName, assertOk := rawIsaMap["name"].(string)
		if !assertOk {
			return registryImage, errors.New("ParseIsaNameError")
		}

		switch rawIsaName {
		case "386":
			rawIsaName = "i386"
		case "amd64", "x86-64":
			rawIsaName = "amd64"
		case "arm", "armv7":
			rawIsaName = "arm"
		case "arm64", "aarch64":
			rawIsaName = "arm64"
		case "riscv64":
			rawIsaName = "riscv64"
		default:
			continue
		}

		isaName, err := valueObject.NewInstructionSetArchitecture(rawIsaName)
		if err != nil {
			log.Printf("UnknownIsaName: %v", rawIsaName)
			continue
		}

		isas = append(isas, isaName)
	}

	pullCount := uint64(0)
	if imageMap["pull_count"] != nil {
		rawPullCount, assertOk := imageMap["pull_count"].(string)
		if !assertOk {
			return registryImage, errors.New("ParsePullCountError")
		}
		pullCountInt, err := voHelper.ExpandNumericAbbreviation(rawPullCount)
		if err != nil {
			return registryImage, err
		}
		pullCount = uint64(pullCountInt)
	}

	starCount, err := voHelper.InterfaceToUint(imageMap["star_count"])
	if err != nil {
		return registryImage, err
	}

	var logoUrlPtr *valueObject.Url
	logoMap, assertOk := imageMap["logo_url"].(map[string]interface{})
	if !assertOk {
		return registryImage, errors.New("ParseLogoUrlError")
	}
	rawLogoUrl, assertOk := logoMap["large"].(string)
	if !assertOk {
		return registryImage, errors.New("ParseLogoUrlError")
	}
	if rawLogoUrl != "" {
		logoUrl, err := valueObject.NewUrl(rawLogoUrl)
		if err != nil {
			return registryImage, err
		}
		logoUrlPtr = &logoUrl
	}

	var createdAtPtr *valueObject.UnixTime
	rawCreatedAt, assertOk := imageMap["created_at"].(string)
	if !assertOk {
		return registryImage, errors.New("ParseCreatedAtError")
	}
	if rawCreatedAt != "" {
		createdAtUnix, err := time.Parse(time.RFC3339, rawCreatedAt)
		if err != nil {
			return registryImage, err
		}
		createdAt := valueObject.UnixTime(createdAtUnix.Unix())
		createdAtPtr = &createdAt
	}

	var updatedAtPtr *valueObject.UnixTime
	rawUpdatedAt, assertOk := imageMap["updated_at"].(string)
	if !assertOk {
		return registryImage, errors.New("ParseUpdatedAtError")
	}
	if rawUpdatedAt != "" {
		updatedAtUnix, err := time.Parse(time.RFC3339, rawUpdatedAt)
		if err != nil {
			return registryImage, err
		}
		updatedAt := valueObject.UnixTime(updatedAtUnix.Unix())
		updatedAtPtr = &updatedAt
	}

	return entity.NewRegistryImage(
		imageName,
		descriptionPtr,
		registryName,
		publisherName,
		imageAddress,
		isas,
		pullCount,
		&starCount,
		logoUrlPtr,
		createdAtPtr,
		updatedAtPtr,
	), nil
}

func (repo ContainerRegistryQueryRepo) queryJsonApi(
	apiUrl string,
) (map[string]interface{}, error) {
	var parsedResponse map[string]interface{}

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	httpResponse, err := httpClient.Get(apiUrl)
	if err != nil {
		return parsedResponse, errors.New("HttpRequestError: " + err.Error())
	}
	defer httpResponse.Body.Close()

	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return parsedResponse, errors.New("HttpResponseError: " + err.Error())
	}

	err = json.Unmarshal(responseBody, &parsedResponse)
	if err != nil {
		return parsedResponse, errors.New("HttpResponseError: " + err.Error())
	}

	return parsedResponse, nil
}

func (repo ContainerRegistryQueryRepo) getDockerHubImages(
	imageName *valueObject.RegistryImageName,
) ([]entity.RegistryImage, error) {
	registryImages := []entity.RegistryImage{}

	imageNameStr := "speedianet/os"
	if imageName != nil {
		imageNameStr = imageName.String()
	}

	hubApiBase := "https://hub.docker.com/api/content/v1/products/search?q=" + imageNameStr
	apiUrls := []string{
		hubApiBase + "&image_filter=store%2Cofficial%2Copen_source&page=1&page_size=10",
		hubApiBase + "&source=community&page=1&page_size=10",
	}

	rawImagesMetadata := []interface{}{}
	for _, apiUrl := range apiUrls {
		parsedResponse, err := repo.queryJsonApi(apiUrl)
		if err != nil {
			return registryImages, err
		}

		summariesMap, assertOk := parsedResponse["summaries"].([]interface{})
		if !assertOk {
			return registryImages, nil
		}

		rawImagesMetadata = append(rawImagesMetadata, summariesMap...)
	}

	for _, image := range rawImagesMetadata {
		imageMap, assertOk := image.(map[string]interface{})
		if !assertOk {
			log.Printf("ParseDockerHubImageError: %v", image)
			continue
		}

		registryImage, err := repo.dockerHubImageFactory(imageMap)
		if err != nil {
			log.Printf("DockerHubImageFactoryError: %v | %v", err, imageMap)
			continue
		}

		registryImages = append(registryImages, registryImage)
	}

	return registryImages, nil
}

func (repo ContainerRegistryQueryRepo) GetImages(
	imageName *valueObject.RegistryImageName,
) ([]entity.RegistryImage, error) {
	registryImages := []entity.RegistryImage{}

	dockerHubImages, err := repo.getDockerHubImages(imageName)
	if err != nil {
		return registryImages, err
	}

	registryImages = append(registryImages, dockerHubImages...)

	return registryImages, nil
}
