package api

import (
	_ "embed"
	"net/http"

	"github.com/goinfinite/fleet/src/infra/db"
	apiController "github.com/goinfinite/fleet/src/presentation/api/controller"
	"github.com/labstack/echo/v4"
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

func accountRoutes(baseRoute *echo.Group, dbSvc *db.DatabaseService) {
	accountGroup := baseRoute.Group("/account")
	accountGroup.GET("/", apiController.GetAccountsController)
	accountGroup.POST("/", apiController.AddAccountController)
	accountGroup.PUT("/", apiController.UpdateAccountController)
	accountGroup.DELETE("/:accountId/", apiController.DeleteAccountController)
	go apiController.AutoUpdateAccountsQuotaUsageController(dbSvc)
}

func containerRoutes(baseRoute *echo.Group) {
	containerGroup := baseRoute.Group("/container")
	containerGroup.GET("/", apiController.GetContainersController)
	containerGroup.POST("/", apiController.AddContainerController)
	containerGroup.PUT("/", apiController.UpdateContainerController)
	containerGroup.DELETE(
		"/:accountId/:containerId/",
		apiController.DeleteContainerController,
	)

	containerProfileGroup := containerGroup.Group("/profile")
	containerProfileGroup.GET("/", apiController.GetContainerProfilesController)
	containerProfileGroup.POST("/", apiController.AddContainerProfileController)
	containerProfileGroup.PUT("/", apiController.UpdateContainerProfileController)
	containerProfileGroup.DELETE(
		"/:profileId/",
		apiController.DeleteContainerProfileController,
	)
}

func o11yRoutes(baseRoute *echo.Group) {
	o11yGroup := baseRoute.Group("/o11y")
	o11yGroup.GET("/overview/", apiController.O11yOverviewController)
}

func registerApiRoutes(baseRoute *echo.Group, dbSvc *db.DatabaseService) {
	swaggerRoute(baseRoute)
	authRoutes(baseRoute)
	accountRoutes(baseRoute, dbSvc)
	containerRoutes(baseRoute)
	o11yRoutes(baseRoute)
}
