package apiController

import (
	"net/http"

	"github.com/goinfinite/fleet/src/domain/dto"
	"github.com/goinfinite/fleet/src/domain/useCase"
	"github.com/goinfinite/fleet/src/domain/valueObject"
	"github.com/goinfinite/fleet/src/infra"
	"github.com/goinfinite/fleet/src/infra/db"
	apiHelper "github.com/goinfinite/fleet/src/presentation/api/helper"
	"github.com/labstack/echo/v4"
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
// @Router       /auth/login/ [post]
func AuthLoginController(c echo.Context) error {
	requiredParams := []string{"username", "password"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	loginDto := dto.NewLogin(
		valueObject.NewUsernamePanic(requestBody["username"].(string)),
		valueObject.NewPasswordPanic(requestBody["password"].(string)),
	)

	dbSvc := c.Get("dbSvc").(*db.DatabaseService)
	authQueryRepo := infra.NewAuthQueryRepo(dbSvc)
	authCmdRepo := infra.AuthCmdRepo{}
	accQueryRepo := infra.NewAccQueryRepo(dbSvc)

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
