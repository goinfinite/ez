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
)

type AuthController struct {
	persistentDbSvc *db.PersistentDatabaseService
	transientDbSbc  *db.TransientDatabaseService
}

func NewAuthController(
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
) *AuthController {
	return &AuthController{
		persistentDbSvc: persistentDbSvc,
		transientDbSbc:  transientDbSvc,
	}
}

func (controller *AuthController) Login(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	requiredParams := []string{"username", "password"}
	apiHelper.CheckMissingParams(requestBody, requiredParams)

	username, err := valueObject.NewUsername(requestBody["username"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	password, err := valueObject.NewPassword(requestBody["password"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	ipAddress, err := valueObject.NewIpAddress(c.RealIP())
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	loginDto := dto.NewLogin(username, password, &ipAddress)

	authQueryRepo := infra.NewAuthQueryRepo(controller.persistentDbSvc)
	authCmdRepo := infra.AuthCmdRepo{}
	accountQueryRepo := infra.NewAccountQueryRepo(controller.persistentDbSvc)

	accessToken, err := useCase.GenerateSessionToken(
		authQueryRepo, authCmdRepo, accountQueryRepo, loginDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusUnauthorized, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, accessToken)
}
