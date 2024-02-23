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

// AuthLogin godoc
// @Summary      GenerateJwtWithCredentials
// @Description  Generate JWT with credentials
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginDto 	  body    dto.Login  true  "Login"
// @Success      200 {object} entity.AccessToken
// @Failure      401 {object} string
// @Router       /v1/auth/login/ [post]
func AuthLoginController(c echo.Context) error {
	requiredParams := []string{"username", "password"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	loginDto := dto.NewLogin(
		valueObject.NewUsernamePanic(requestBody["username"].(string)),
		valueObject.NewPasswordPanic(requestBody["password"].(string)),
	)

	persistentDbSvc := c.Get("persistentDbSvc").(*db.PersistentDatabaseService)
	authQueryRepo := infra.NewAuthQueryRepo(persistentDbSvc)
	authCmdRepo := infra.AuthCmdRepo{}
	accQueryRepo := infra.NewAccQueryRepo(persistentDbSvc)

	ipAddress := valueObject.NewIpAddressPanic(c.RealIP())

	accessToken, err := useCase.GetSessionToken(
		authQueryRepo,
		authCmdRepo,
		accQueryRepo,
		loginDto,
		ipAddress,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusUnauthorized, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, accessToken)
}
