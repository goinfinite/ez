package api

import (
	_ "embed"
	"net/http"

	"github.com/labstack/echo/v4"
	apiController "github.com/speedianet/sfm/src/presentation/api/controller"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//go:embed docs/swagger.json
var swaggerJson []byte

func swaggerRoute(baseRoute *echo.Group) {
	swaggerGroup := baseRoute.Group("/swagger")

	swaggerGroup.GET("/swagger.json", func(c echo.Context) error {
		return c.Blob(http.StatusOK, echo.MIMEApplicationJSON, swaggerJson)
	})

	swaggerUrl := echoSwagger.URL("swagger.json")
	swaggerGroup.GET("/*", echoSwagger.EchoWrapHandler(swaggerUrl))
}

func authRoutes(baseRoute *echo.Group) {
	authGroup := baseRoute.Group("/auth")
	authGroup.POST("/login/", apiController.AuthLoginController)
}

func accountRoutes(baseRoute *echo.Group) {
	accountGroup := baseRoute.Group("/account")
	accountGroup.GET("/", apiController.GetAccountsController)
	accountGroup.POST("/", apiController.AddAccountController)
	accountGroup.PUT("/", apiController.UpdateAccountController)
	go apiController.AutoUpdateAccountsQuotaUsageController()
}

func containerRoutes(baseRoute *echo.Group) {
	containerGroup := baseRoute.Group("/container")
	containerGroup.GET("/", apiController.GetContainersController)
	containerGroup.POST("/", apiController.AddContainerController)
	containerGroup.PUT("/", apiController.UpdateContainerController)
	containerGroup.DELETE(
		"/:accId/:containerId",
		apiController.DeleteContainerController,
	)
}

func o11yRoutes(baseRoute *echo.Group) {
	o11yGroup := baseRoute.Group("/o11y")
	o11yGroup.GET("/overview/", apiController.O11yOverviewController)
}

func registerApiRoutes(baseRoute *echo.Group) {
	swaggerRoute(baseRoute)
	authRoutes(baseRoute)
	accountRoutes(baseRoute)
	containerRoutes(baseRoute)
	o11yRoutes(baseRoute)
}
