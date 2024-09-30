package infra

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/valueObject"
	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
	"github.com/goinfinite/ez/src/infra/db"
)

type ContainerRegistryQueryRepo struct {
	persistentDbSvc *db.PersistentDatabaseService
}

func NewContainerRegistryQueryRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerRegistryQueryRepo {
	return &ContainerRegistryQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *ContainerRegistryQueryRepo) dockerHubImageFactory(
	imageMap map[string]interface{},
) (entity.RegistryImage, error) {
	var registryImage entity.RegistryImage

	rawImageName, assertOk := imageMap["name"].(string)
	if !assertOk {
		return registryImage, errors.New("ParseImageNameError")
	}
	imageName, err := valueObject.NewRegistryImageName(rawImageName)
	if err != nil {
		return registryImage, errors.New(err.Error() + ": " + rawImageName)
	}

	publisherNameStr := "docker"
	imageNameHasPublisherName := strings.Contains(imageName.String(), "/")
	if imageNameHasPublisherName {
		imageNameParts := strings.Split(imageName.String(), "/")
		publisherNameStr = imageNameParts[0]
	}

	if imageMap["publisher"] != nil {
		rawPublisherDetails, assertOk := imageMap["publisher"].(map[string]interface{})
		if !assertOk {
			return registryImage, errors.New("ParsePublisherDetailsError")
		}
		rawPublisherName, assertOk := rawPublisherDetails["name"].(string)
		if assertOk && !strings.Contains(rawPublisherName, " ") {
			publisherNameStr = rawPublisherName
		}
	}

	publisherName, err := valueObject.NewRegistryPublisherName(publisherNameStr)
	if err != nil {
		return registryImage, err
	}

	registryName, _ := valueObject.NewRegistryName("DockerHub")

	rawImageAddress, assertOk := imageMap["slug"].(string)
	if !assertOk {
		return registryImage, errors.New("ParseImageAddressError")
	}
	if rawImageAddress == "" {
		rawImageAddress = imageName.String()
	}
	imageAddress, err := valueObject.NewContainerImageAddress(rawImageAddress)
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

		// TODO: support arm, armv7 and arm64 in the future.
		switch rawIsaName {
		case "amd64", "x86-64":
			rawIsaName = "amd64"
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

	starCount, err := voHelper.InterfaceToUint64(imageMap["star_count"])
	if err != nil {
		return registryImage, err
	}

	var logoUrlPtr *valueObject.Url
	if imageMap["logo_url"] != nil {
		logoMap, assertOk := imageMap["logo_url"].(map[string]interface{})
		if assertOk && logoMap["large"] != nil {
			rawLogoUrl, assertOk := logoMap["large"].(string)
			if assertOk && rawLogoUrl != "" {
				logoUrl, err := valueObject.NewUrl(rawLogoUrl)
				if err != nil {
					return registryImage, err
				}
				logoUrlPtr = &logoUrl
			}
		}
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
		createdAt := valueObject.NewUnixTimeWithGoTime(createdAtUnix)
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
		updatedAt := valueObject.NewUnixTimeWithGoTime(updatedAtUnix)
		updatedAtPtr = &updatedAt
	}

	return entity.NewRegistryImage(
		imageName,
		publisherName,
		registryName,
		imageAddress,
		descriptionPtr,
		isas,
		pullCount,
		&starCount,
		logoUrlPtr,
		createdAtPtr,
		updatedAtPtr,
	), nil
}

func (repo *ContainerRegistryQueryRepo) queryJsonApi(
	apiUrl string,
) ([]byte, error) {
	var responseBody []byte

	httpRequest, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return responseBody, errors.New("HttpRequestError: " + err.Error())
	}
	httpRequest.Header.Set("Search-Version", "v3")

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return responseBody, errors.New("HttpResponseError: " + err.Error())
	}
	defer httpResponse.Body.Close()

	return io.ReadAll(httpResponse.Body)
}

func (repo *ContainerRegistryQueryRepo) getDockerHubImages(
	imageName *valueObject.RegistryImageName,
) ([]entity.RegistryImage, error) {
	registryImages := []entity.RegistryImage{}

	imageNameStr := "speedianet/os"
	if imageName != nil {
		imageNameStr = imageName.String()
	}

	hubApiBase := "https://hub.docker.com/api/content/v1/products/search?q=" +
		imageNameStr +
		"&page=1&page_size=10"
	apiUrls := []string{
		hubApiBase + "&image_filter=store%2Cofficial%2Copen_source",
		hubApiBase + "&source=community",
	}

	rawImagesMetadata := []interface{}{}
	for _, apiUrl := range apiUrls {
		apiResponse, err := repo.queryJsonApi(apiUrl)
		if err != nil {
			log.Printf("DockerHubApiError: %v", err)
			continue
		}

		var parsedResponse map[string]interface{}
		err = json.Unmarshal(apiResponse, &parsedResponse)
		if err != nil {
			log.Printf("ParseDockerHubResponseError: %v", err)
			continue
		}

		summariesMap, assertOk := parsedResponse["summaries"].([]interface{})
		if !assertOk {
			continue
		}

		rawImagesMetadata = append(rawImagesMetadata, summariesMap...)
	}

	for _, image := range rawImagesMetadata {
		imageMap, assertOk := image.(map[string]interface{})
		if !assertOk {
			log.Printf("ParseDockerHubImageError: %v", image)
			continue
		}

		rawType, assertOk := imageMap["type"].(string)
		if !assertOk || rawType != "image" {
			continue
		}

		registryImage, err := repo.dockerHubImageFactory(imageMap)
		if err != nil {
			log.Printf("DockerHubImageFactoryError: %v | %v", err, imageMap)
			continue
		}

		registryImages = append(registryImages, registryImage)
	}

	sort.SliceStable(registryImages, func(i, j int) bool {
		return registryImages[i].PullCount > registryImages[j].PullCount
	})

	return registryImages, nil
}

func (repo *ContainerRegistryQueryRepo) ReadImages(
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

func (repo *ContainerRegistryQueryRepo) getTaggedImageFromDockerHub(
	imageAddress valueObject.ContainerImageAddress,
) (registryTaggedImage entity.RegistryTaggedImage, err error) {
	orgName, err := imageAddress.GetOrgName()
	if err != nil {
		return registryTaggedImage, err
	}

	imageName, err := imageAddress.GetImageName()
	if err != nil {
		return registryTaggedImage, err
	}

	imageTag, err := imageAddress.GetImageTag()
	if err != nil {
		return registryTaggedImage, err
	}

	hubApi := "https://hub.docker.com/v2/namespaces/" +
		orgName.String() + "/repositories/" + imageName.String() +
		"/tags/" + imageTag.String() + "/images"

	apiResponse, err := repo.queryJsonApi(hubApi)
	if err != nil {
		return registryTaggedImage, err
	}

	var parsedResponse []interface{}
	err = json.Unmarshal(apiResponse, &parsedResponse)
	if err != nil {
		return registryTaggedImage, errors.New("ParseResponseBodyError: " + err.Error())
	}

	if len(parsedResponse) == 0 {
		return registryTaggedImage, errors.New("ImageNotFound")
	}

	desiredImageMap := map[string]interface{}{}
	for _, image := range parsedResponse {
		imageMap, assertOk := image.(map[string]interface{})
		if !assertOk {
			slog.Debug("ParseDockerHubImageError", slog.Any("image", image))
			continue
		}

		rawArchitecture, assertOk := imageMap["architecture"].(string)
		if !assertOk {
			continue
		}

		if rawArchitecture != "amd64" {
			slog.Debug(
				"SkippingUnsupportedArchitecture",
				slog.String("isa", rawArchitecture),
				slog.String("imageName", imageName.String()),
				slog.String("imageTag", imageTag.String()),
			)
			continue
		}

		desiredImageMap = imageMap
	}

	if desiredImageMap == nil {
		return registryTaggedImage, errors.New("ImageNotFound")
	}

	rawImageSize, assertOk := desiredImageMap["size"].(float64)
	if !assertOk {
		return registryTaggedImage, errors.New("ParseImageSizeError")
	}
	sizeBytes, err := valueObject.NewByte(rawImageSize)
	if err != nil {
		return registryTaggedImage, err
	}

	rawImageHash, assertOk := desiredImageMap["digest"].(string)
	if !assertOk {
		return registryTaggedImage, errors.New("ParseImageHashError")
	}
	if rawImageHash == "" {
		return registryTaggedImage, errors.New("ImageHashEmpty")
	}
	rawImageHash = strings.ReplaceAll(rawImageHash, "sha256:", "")

	imageHash, err := valueObject.NewHash(rawImageHash)
	if err != nil {
		return registryTaggedImage, err
	}

	rawUpdatedAt, assertOk := desiredImageMap["last_pushed"].(string)
	if !assertOk {
		fakeUpdatedAt := time.Now().UTC().Add(-time.Hour * 24 * 365 * 5)
		rawUpdatedAt = fakeUpdatedAt.Format(time.RFC3339)
	}
	updatedAtUnix, err := time.Parse(time.RFC3339, rawUpdatedAt)
	if err != nil {
		return registryTaggedImage, err
	}
	updatedAt := valueObject.NewUnixTimeWithGoTime(updatedAtUnix)

	portBindings := []valueObject.PortBinding{}
	portBindingsRegex := `\d{1,5}(\/\w{1,4})?`

	imageLayers, assertOk := desiredImageMap["layers"].([]interface{})
	if !assertOk {
		return registryTaggedImage, errors.New("ParseImageLayersError")
	}

	for _, layer := range imageLayers {
		layerMap, assertOk := layer.(map[string]interface{})
		if !assertOk {
			continue
		}

		rawInstruction, assertOk := layerMap["instruction"].(string)
		if !assertOk {
			continue
		}

		rawInstruction = strings.TrimSpace(rawInstruction)
		if !strings.HasPrefix(rawInstruction, "EXPOSE") {
			continue
		}

		bindingsRegex, err := regexp.Compile(portBindingsRegex)
		if err != nil {
			continue
		}
		rawPortBindings := bindingsRegex.FindAllString(rawInstruction, -1)

		for _, rawPortBinding := range rawPortBindings {
			rawPortBinding = strings.ReplaceAll(rawPortBinding, "/tcp", "")
			portBinding, err := valueObject.NewPortBindingFromString(rawPortBinding)
			if err != nil {
				continue
			}

			portBindings = append(portBindings, portBinding...)
		}
	}

	registryName, _ := valueObject.NewRegistryName("DockerHub")
	isa, _ := valueObject.NewInstructionSetArchitecture("amd64")

	return entity.NewRegistryTaggedImage(
		imageTag, imageName, orgName, registryName, imageAddress, imageHash, isa,
		sizeBytes, portBindings, updatedAt,
	), nil
}

func (repo *ContainerRegistryQueryRepo) ReadTaggedImage(
	imageAddress valueObject.ContainerImageAddress,
) (taggedImage entity.RegistryTaggedImage, err error) {
	registryName, err := imageAddress.GetFqdn()
	if err != nil {
		return taggedImage, err
	}

	switch registryName.String() {
	case "docker.io", "registry-1.docker.io":
		return repo.getTaggedImageFromDockerHub(imageAddress)
	default:
		return taggedImage, errors.New("UnknownRegistry")
	}
}
