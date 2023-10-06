package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/sfm/src/domain/dto"
	"github.com/speedianet/sfm/src/domain/useCase"
	"github.com/speedianet/sfm/src/domain/valueObject"
	"github.com/speedianet/sfm/src/infra"
	apiHelper "github.com/speedianet/sfm/src/presentation/api/helper"
	"gorm.io/gorm"
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

	authQueryRepo := infra.NewAuthQueryRepo(c.Get("dbSvc").(*gorm.DB))
	authCmdRepo := infra.AuthCmdRepo{}
	accQueryRepo := infra.NewAccQueryRepo(c.Get("dbSvc").(*gorm.DB))

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
