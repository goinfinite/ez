package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/useCase"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra"
	"github.com/speedianet/control/src/infra/db"
	apiHelper "github.com/speedianet/control/src/presentation/api/helper"
	serviceHelper "github.com/speedianet/control/src/presentation/service/helper"
)

type AuthController struct {
	persistentDbSvc *db.PersistentDatabaseService
	trailDbSvc      *db.TrailDatabaseService
}

func NewAuthController(
	persistentDbSvc *db.PersistentDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *AuthController {
	return &AuthController{
		persistentDbSvc: persistentDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

// CreateSessionTokenWithCredentials	 godoc
// @Summary      CreateSessionTokenWithCredentials
// @Description  Create a new session token with the provided credentials.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        createSessionToken body dto.CreateSessionToken true "CreateSessionToken"
// @Success      200 {object} entity.AccessToken
// @Failure      401 {object} string
// @Router       /v1/auth/login/ [post]
func (controller *AuthController) Login(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	requiredParams := []string{"username", "password"}
	err = serviceHelper.RequiredParamsInspector(requestBody, requiredParams)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	username, err := valueObject.NewUsername(requestBody["username"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	password, err := valueObject.NewPassword(requestBody["password"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	operatorIpAddress, err := valueObject.NewIpAddress(requestBody["ipAddress"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	createDto := dto.NewCreateSessionToken(username, password, operatorIpAddress)

	authQueryRepo := infra.NewAuthQueryRepo(controller.persistentDbSvc)
	authCmdRepo := infra.AuthCmdRepo{}
	accountQueryRepo := infra.NewAccountQueryRepo(controller.persistentDbSvc)
	activityRecordQueryRepo := infra.NewActivityRecordQueryRepo(controller.trailDbSvc)
	activityRecordCmdRepo := infra.NewActivityRecordCmdRepo(controller.trailDbSvc)

	accessToken, err := useCase.CreateSessionToken(
		authQueryRepo, authCmdRepo, accountQueryRepo,
		activityRecordQueryRepo, activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusUnauthorized, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, accessToken)
}
