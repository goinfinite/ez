package apiController

import (
	"net/http"

	"github.com/goinfinite/fleet/src/domain/dto"
	"github.com/goinfinite/fleet/src/domain/useCase"
	"github.com/goinfinite/fleet/src/domain/valueObject"
	voHelper "github.com/goinfinite/fleet/src/domain/valueObject/helper"
	"github.com/goinfinite/fleet/src/infra"
	"github.com/goinfinite/fleet/src/infra/db"
	apiHelper "github.com/goinfinite/fleet/src/presentation/api/helper"
	"github.com/labstack/echo/v4"
)

// GetContainerProfiles	 godoc
// @Summary      GetContainerProfiles
// @Description  List container profiles.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.ContainerProfile
// @Router       /container/profile/ [get]
func GetContainerProfilesController(c echo.Context) error {
	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	queryRepo := infra.NewContainerProfileQueryRepo(dbSvc)
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

// AddContainerProfile	 godoc
// @Summary      AddNewContainerProfile
// @Description  Add a new container profile.
// @Tags         container
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        addContainerProfileDto 	  body    dto.AddContainerProfile  true  "NewContainerProfile (Only name and baseSpecs are required.)"
// @Success      201 {object} object{} "ContainerProfileCreated"
// @Router       /container/profile/ [post]
func AddContainerProfileController(c echo.Context) error {
	requiredParams := []string{"name", "baseSpecs"}
	requestBody, _ := apiHelper.GetRequestBody(c)

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

	dto := dto.NewAddContainerProfile(
		name,
		baseSpecs,
		maxSpecsPtr,
		scalingPolicyPtr,
		scalingThresholdPtr,
		scalingMaxDurationSecsPtr,
		scalingIntervalSecsPtr,
		hostMinCapacityPercentPtr,
	)

	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	containerProfileCmdRepo := infra.NewContainerProfileCmdRepo(dbSvc)

	err := useCase.AddContainerProfile(
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
// @Router       /container/profile/ [put]
func UpdateContainerProfileController(c echo.Context) error {
	requiredParams := []string{"id"}
	requestBody, _ := apiHelper.GetRequestBody(c)

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

	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(dbSvc)
	containerProfileCmdRepo := infra.NewContainerProfileCmdRepo(dbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)
	containerCmdRepo := infra.ContainerCmdRepo{}

	err := useCase.UpdateContainerProfile(
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
// @Router       /container/profile/{profileId}/ [delete]
func DeleteContainerProfileController(c echo.Context) error {
	profileId := valueObject.NewContainerProfileIdPanic(c.Param("profileId"))

	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	containerProfileQueryRepo := infra.NewContainerProfileQueryRepo(dbSvc)
	containerProfileCmdRepo := infra.NewContainerProfileCmdRepo(dbSvc)
	containerQueryRepo := infra.NewContainerQueryRepo(dbSvc)
	containerCmdRepo := infra.ContainerCmdRepo{}

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
