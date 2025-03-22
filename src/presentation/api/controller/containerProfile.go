package apiController

import (
	"net/http"

	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	apiHelper "github.com/goinfinite/ez/src/presentation/api/helper"
	"github.com/goinfinite/ez/src/presentation/service"
	"github.com/labstack/echo/v4"
)

type ContainerProfileController struct {
	containerProfileService *service.ContainerProfileService
}

func NewContainerProfileController(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *ContainerProfileController {
	return &ContainerProfileController{
		containerProfileService: service.NewContainerProfileService(
			persistentDbSvc, trailDbSvc,
		),
	}
}

// ReadContainerProfiles	 godoc
// @Summary      ReadContainerProfiles
// @Description  List container profiles.
// @Tags         containerProfile
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.ContainerProfile
// @Router       /v1/container/profile/ [get]
func (controller *ContainerProfileController) Read(c echo.Context) error {
	return apiHelper.ServiceResponseWrapper(c, controller.containerProfileService.Read())
}

func parseContainerSpecs(
	rawSpecs map[string]interface{},
) (specs valueObject.ContainerSpecs, err error) {
	defaultSpecs := valueObject.NewContainerSpecsWithDefaultValues()

	millicores := defaultSpecs.Millicores
	if rawSpecs["cpuCores"] != nil {
		millicores, err = valueObject.NewCpuCores(rawSpecs["cpuCores"])
		if err != nil {
			return specs, err
		}
	}

	if rawSpecs["millicores"] != nil {
		millicores, err = valueObject.NewMillicores(rawSpecs["millicores"])
		if err != nil {
			return specs, err
		}
	}

	memoryBytes := defaultSpecs.MemoryBytes
	if rawSpecs["memoryMebibytes"] != nil {
		memoryBytes, err = valueObject.NewMebibyte(rawSpecs["memoryMebibytes"])
		if err != nil {
			return specs, err
		}
	}

	if rawSpecs["memoryGibibytes"] != nil {
		memoryBytes, err = valueObject.NewGibibyte(rawSpecs["memoryGibibytes"])
		if err != nil {
			return specs, err
		}
	}

	if rawSpecs["memoryBytes"] != nil {
		memoryBytes, err = valueObject.NewByte(rawSpecs["memoryBytes"])
		if err != nil {
			return specs, err
		}
	}

	storagePerformanceUnits := defaultSpecs.StoragePerformanceUnits
	if rawSpecs["storagePerformanceUnits"] != nil {
		storagePerformanceUnits, err = valueObject.NewStoragePerformanceUnits(
			rawSpecs["storagePerformanceUnits"],
		)
		if err != nil {
			return specs, err
		}
	}

	return valueObject.NewContainerSpecs(
		millicores, memoryBytes, storagePerformanceUnits,
	), nil
}

// CreateContainerProfile	 godoc
// @Summary      CreateNewContainerProfile
// @Description  Create a new container profile.
// @Tags         containerProfile
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createContainerProfileDto 	  body    dto.CreateContainerProfile  true  "Only 'name' and 'baseSpecs' are required. Human-readable fields ('cpuCores', 'memoryMebibytes' etc) will be converted to their technical counterpart ('millicores' etc) automatically."
// @Success      201 {object} object{} "ContainerProfileCreated"
// @Router       /v1/container/profile/ [post]
func (controller *ContainerProfileController) Create(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	if requestBody["baseSpecs"] != nil {
		baseSpecsMap, assertOk := requestBody["baseSpecs"].(map[string]interface{})
		if !assertOk {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "InvalidBaseSpecsStructure",
			)
		}

		baseSpecs, err := parseContainerSpecs(baseSpecsMap)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}

		requestBody["baseSpecs"] = baseSpecs
	}

	if requestBody["maxSpecs"] != nil {
		maxSpecsMap, assertOk := requestBody["maxSpecs"].(map[string]interface{})
		if !assertOk {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "InvalidMaxSpecsStructure",
			)
		}

		maxSpecs, err := parseContainerSpecs(maxSpecsMap)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}

		requestBody["maxSpecs"] = maxSpecs
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.containerProfileService.Create(requestBody),
	)
}

// UpdateContainerProfile godoc
// @Summary      UpdateContainerProfile
// @Description  Update a container profile.
// @Tags         containerProfile
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateContainerProfileDto 	  body dto.UpdateContainerProfile  true  "Only 'id' is required. Human-readable fields ('cpuCores', 'memoryMebibytes' etc) will be converted to their technical counterpart ('millicores' etc) automatically."
// @Success      200 {object} object{} "ContainerProfileUpdated"
// @Router       /v1/container/profile/ [put]
func (controller *ContainerProfileController) Update(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	if requestBody["baseSpecs"] != nil {
		baseSpecsMap, assertOk := requestBody["baseSpecs"].(map[string]interface{})
		if !assertOk {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "InvalidBaseSpecsStructure",
			)
		}

		baseSpecs, err := parseContainerSpecs(baseSpecsMap)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}

		requestBody["baseSpecs"] = baseSpecs
	}

	if requestBody["maxSpecs"] != nil {
		maxSpecsMap, assertOk := requestBody["maxSpecs"].(map[string]interface{})
		if !assertOk {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "InvalidMaxSpecsStructure",
			)
		}

		maxSpecs, err := parseContainerSpecs(maxSpecsMap)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}

		requestBody["maxSpecs"] = maxSpecs
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.containerProfileService.Update(requestBody),
	)
}

// DeleteContainerProfile godoc
// @Summary      DeleteContainerProfile
// @Description  Delete a container profile.
// @Tags         containerProfile
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        accountId 	  path   string  true  "AccountId"
// @Param        profileId 	  path   string  true  "ProfileId"
// @Success      200 {object} object{} "ContainerProfileDeleted"
// @Router       /v1/container/profile/{accountId}/{profileId}/ [delete]
func (controller *ContainerProfileController) Delete(c echo.Context) error {
	requestBody := map[string]interface{}{
		"accountId":         c.Param("accountId"),
		"profileId":         c.Param("profileId"),
		"operatorAccountId": c.Get("accountId"),
		"operatorIpAddress": c.RealIP(),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.containerProfileService.Delete(requestBody),
	)
}
