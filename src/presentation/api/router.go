package api

import (
	_ "embed"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/infra/db"
	apiController "github.com/speedianet/control/src/presentation/api/controller"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/speedianet/control/src/presentation/api/docs"
)

func swaggerRoute(baseRoute *echo.Group) {
	swaggerGroup := baseRoute.Group("/swagger")
	swaggerGroup.GET("/*", echoSwagger.WrapHandler)
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
	containerGroup.GET("/metrics/", apiController.GetContainersWithMetricsController)
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

func mappingsRoutes(baseRoute *echo.Group) {
	mappingsGroup := baseRoute.Group("/mappings")
	mappingsGroup.GET("/", apiController.GetMappingsController)
	mappingsGroup.POST("/", apiController.AddMappingController)
	mappingsGroup.DELETE("/:mappingId/", apiController.DeleteMappingController)
	mappingsGroup.POST("/target/", apiController.AddMappingTargetController)
	mappingsGroup.DELETE(
		"/:mappingId/target/:targetId/",
		apiController.DeleteMappingTargetController,
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
