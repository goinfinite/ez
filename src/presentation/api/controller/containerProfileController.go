package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	apiHelper "github.com/speedianet/control/src/presentation/api/helper"
)

// GetContainerProfiles	 godoc
// @Summary      GetContainerProfiles
// @Description  List container profiles.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.ContainerProfile
// @Router       /v1/container/profile/ [get]
func GetContainerProfilesController(c echo.Context) error {
	persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)
	queryRepo := infra.NewContainerProfileQueryRepo(persistentDbSvc)
	profilesList, err := useCase.GetContainerProfiles(queryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, profilesList)
}

func parseContainerSpecs(specs map[string]interface{}) valueObject.ContainerSpecs {
	cpuCores := valueObject.NewCpuCoresCountPanic(specs["cpuCores"])
	memoryBytes := valueObject.NewBytePanic(specs["memoryBytes"])

	return valueObject.NewContainerSpecs(cpuCores, memoryBytes)
}

// CreateContainerProfile	 godoc
// @Summary      AddNewContainerProfile
// @Description  Add a new container profile.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createContainerProfileDto 	  body    dto.CreateContainerProfile  true  "NewContainerProfile (Only name and baseSpecs are required.)"
// @Success      201 {object} object{} "ContainerProfileCreated"
// @Router       /v1/container/profile/ [post]
func CreateContainerProfileController(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	requiredParams := []string{"name", "baseSpecs"}
	apiHelper.CheckMissingParams(requestBody, requiredParams)

	name := valueObject.NewContainerProfileNamePanic(requestBody["name"].(string))

	baseSpecsMap, assertOk := requestBody["baseSpecs"].(map[string]interface{})
	if !assertOk {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "InvalidBaseSpecs")
	}
	baseSpecs := parseContainerSpecs(baseSpecsMap)

	var maxSpecsPtr *valueObject.ContainerSpecs
	if requestBody["maxSpecs"] != nil {
		maxSpecsMap, assertOk := requestBody["maxSpecs"].(map[string]interface{})
		if !assertOk {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "InvalidMaxSpecs")
		}
		maxSpecs := parseContainerSpecs(maxSpecsMap)
		maxSpecsPtr = &maxSpecs
	}

	var scalingPolicyPtr *valueObject.ScalingPolicy
	if requestBody["scalingPolicy"] != nil {
		scalingPolicy := valueObject.NewScalingPolicyPanic(
			requestBody["scalingPolicy"].(string),
		)
		scalingPolicyPtr = &scalingPolicy
	}

	var scalingThresholdPtr *uint64
	if requestBody["scalingThreshold"] != nil {
		scalingThreshold, err := voHelper.InterfaceToUint(requestBody["scalingThreshold"])
		if err != nil {
			return apiHelper.ResponseWrapper(
				c,
				http.StatusBadRequest,
				"InvalidScalingThreshold",
			)
		}
		scalingThresholdPtr = &scalingThreshold
	}

	var scalingMaxDurationSecsPtr *uint64
	if requestBody["scalingMaxDurationSecs"] != nil {
		scalingMaxDurationSecs, err := voHelper.InterfaceToUint(
			requestBody["scalingMaxDurationSecs"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(
				c,
				http.StatusBadRequest,
				"InvalidScalingMaxDurationSecs",
			)
		}
		scalingMaxDurationSecsPtr = &scalingMaxDurationSecs
	}

	var scalingIntervalSecsPtr *uint64
	if requestBody["scalingIntervalSecs"] != nil {
		scalingIntervalSecs, err := voHelper.InterfaceToUint(
			requestBody["scalingIntervalSecs"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(
				c,
				http.StatusBadRequest,
				"InvalidScalingIntervalSecs",
			)
		}
		scalingIntervalSecsPtr = &scalingIntervalSecs
	}

	var hostMinCapacityPercentPtr *valueObject.HostMinCapacity
	if requestBody["hostMinCapacityPercent"] != nil {
		hostMinCapacityPercent := valueObject.NewHostMinCapacityPanic(
			requestBody["hostMinCapacityPercent"],
		)
		hostMinCapacityPercentPtr = &hostMinCapacityPercent
	}

	dto := dto.NewCreateContainerProfile(
		name,
		baseSpecs,
		maxSpecsPtr,
		scalingPolicyPtr,
		scalingThresholdPtr,
		scalingMaxDurationSecsPtr,
		scalingIntervalSecsPtr,
		hostMinCapacityPercentPtr,
	)

	persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)
	containerProfileCmdRepo := infra.NewContainerProfileCmdRepo(persistentDbSvc)

	err = useCase.CreateContainerProfile(
		containerProfileCmdRepo,
		dto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "ContainerProfileCreated")
}

// UpdateContainerProfile godoc
// @Summary      UpdateContainerProfile
// @Description  Update a container profile.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateContainerProfileDto 	  body dto.UpdateContainerProfile  true  "UpdateContainerProfile (Only id is required.)"
// @Success      200 {object} object{} "ContainerProfileUpdated"
// @Router       /v1/container/profile/ [put]
func UpdateContainerProfileController(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	requiredParams := []string{"id"}
	apiHelper.CheckMissingParams(requestBody, requiredParams)

	id := valueObject.NewContainerProfileIdPanic(requestBody["id"])

	var namePtr *valueObject.ContainerProfileName
	if requestBody["name"] != nil {
		name := valueObject.NewContainerProfileNamePanic(requestBody["name"].(string))
		namePtr = &name
	}

	var baseSpecsPtr *valueObject.ContainerSpecs
	if requestBody["baseSpecs"] != nil {
		baseSpecsMap, assertOk := requestBody["baseSpecs"].(map[string]interface{})
		if !assertOk {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "InvalidBaseSpecs")
		}
		baseSpecs := parseContainerSpecs(baseSpecsMap)
		baseSpecsPtr = &baseSpecs
	}

	var maxSpecsPtr *valueObject.ContainerSpecs
	if requestBody["maxSpecs"] != nil {
		maxSpecsMap, assertOk := requestBody["maxSpecs"].(map[string]interface{})
		if !assertOk {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "InvalidMaxSpecs")
		}
		maxSpecs := parseContainerSpecs(maxSpecsMap)
		maxSpecsPtr = &maxSpecs
	}

	var scalingPolicyPtr *valueObject.ScalingPolicy
	if requestBody["scalingPolicy"] != nil {
		scalingPolicy := valueObject.NewScalingPolicyPanic(
			requestBody["scalingPolicy"].(string),
		)
		scalingPolicyPtr = &scalingPolicy
	}

	var scalingThresholdPtr *uint64
	if requestBody["scalingThreshold"] != nil {
		scalingThreshold, err := voHelper.InterfaceToUint(requestBody["scalingThreshold"])
		if err != nil {
			return apiHelper.ResponseWrapper(
				c,
				http.StatusBadRequest,
				"InvalidScalingThreshold",
			)
		}
		scalingThresholdPtr = &scalingThreshold
	}

	var scalingMaxDurationSecsPtr *uint64
	if requestBody["scalingMaxDurationSecs"] != nil {
		scalingMaxDurationSecs, err := voHelper.InterfaceToUint(
			requestBody["scalingMaxDurationSecs"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(
				c,
				http.StatusBadRequest,
				"InvalidScalingMaxDurationSecs",
			)
		}
		scalingMaxDurationSecsPtr = &scalingMaxDurationSecs
	}

	var scalingIntervalSecsPtr *uint64
	if requestBody["scalingIntervalSecs"] != nil {
		scalingIntervalSecs, err := voHelper.InterfaceToUint(
			requestBody["scalingIntervalSecs"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(
				c,
				http.StatusBadRequest,
				"InvalidScalingIntervalSecs",
			)
		}
		scalingIntervalSecsPtr = &scalingIntervalSecs
	}

	var hostMinCapacityPercentPtr *valueObject.HostMinCapacity
	if requestBody["hostMinCapacityPercent"] != nil {
		hostMinCapacityPercent := valueObject.NewHostMinCapacityPanic(
			requestBody["hostMinCapacityPercent"],
		)
		hostMinCapacityPercentPtr = &hostMinCapacityPercent
	}

	dto := dto.NewUpdateContainerProfile(
		id,
		namePtr,
		baseSpecsPtr,
		maxSpecsPtr,
		scalingPolicyPtr,
		scalingThresholdPtr,
		scalingMaxDurationSecsPtr,
		scalingIntervalSecsPtr,
		hostMinCapacityPercentPtr,
	)

	persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)
	containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(persistentDbSvc)
	containerProfileCmdRepo := infra.NewContainerProfileCmdRepo(persistentDbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(persistentDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(persistentDbSvc)

	err = useCase.UpdateContainerProfile(
		containerProfileQueryRepo,
		containerProfileCmdRepo,
		containerQueryRepo,
		containerCmdRepo,
		dto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "ContainerProfileUpdated")
}

// DeleteContainerProfile godoc
// @Summary      DeleteContainerProfile
// @Description  Delete a container profile.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        profileId 	  path   string  true  "ProfileId"
// @Success      200 {object} object{} "ContainerProfileDeleted"
// @Router       /v1/container/profile/{profileId}/ [delete]
func DeleteContainerProfileController(c echo.Context) error {
	profileId := valueObject.NewContainerProfileIdPanic(c.Param("profileId"))

	persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)
	containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(persistentDbSvc)
	containerProfileCmdRepo := infra.NewContainerProfileCmdRepo(persistentDbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(persistentDbSvc)
	containerCmdRepo := infra.NewContainerCmdRepo(persistentDbSvc)

	err := useCase.DeleteContainerProfile(
		containerProfileQueryRepo,
		containerProfileCmdRepo,
		containerQueryRepo,
		containerCmdRepo,
		profileId,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "ContainerProfileDeleted")
}
