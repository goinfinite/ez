package api

import (
	_ "embed"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/infra/db"
	apiController "github.com/speedianet/control/src/presentation/api/controller"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/speedianet/control/src/presentation/api/docs"
)

type Router struct {
	baseRoute       *echo.Group
	persistentDbSvc *db.PersistentDatabaseService
	transientDbSvc  *db.TransientDatabaseService
}

func NewRouter(
	baseRoute *echo.Group,
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
) *Router {
	return &Router{
		baseRoute:       baseRoute,
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
	}
}

func (router *Router) swaggerRoute() {
	swaggerGroup := router.baseRoute.Group("/swagger")
	swaggerGroup.GET("/*", echoSwagger.WrapHandler)
}

func (router *Router) authRoutes() {
	authGroup := router.baseRoute.Group("/auth")
	authGroup.POST("/login/", apiController.AuthLoginController)
}

func (router *Router) accountRoutes() {
	accountGroup := router.baseRoute.Group("/account")
	accountGroup.GET("/", apiController.GetAccountsController)
	accountGroup.POST("/", apiController.AddAccountController)
	accountGroup.PUT("/", apiController.UpdateAccountController)
	accountGroup.DELETE("/:accountId/", apiController.DeleteAccountController)
	go apiController.AutoUpdateAccountsQuotaUsageController(
		router.persistentDbSvc,
	)
}

func (router *Router) containerRoutes() {
	containerGroup := router.baseRoute.Group("/container")
	containerController := apiController.NewContainerController(router.persistentDbSvc)
	containerGroup.GET("/", containerController.GetContainers)
	containerGroup.GET("/metrics/", containerController.GetContainersWithMetrics)
	containerGroup.POST("/", containerController.CreateContainer)
	containerGroup.PUT("/", containerController.UpdateContainer)
	containerGroup.DELETE(
		"/:accountId/:containerId/",
		containerController.DeleteContainer,
	)

	containerProfileGroup := containerGroup.Group("/profile")
	containerProfileGroup.GET("/", apiController.GetContainerProfilesController)
	containerProfileGroup.POST("/", apiController.AddContainerProfileController)
	containerProfileGroup.PUT("/", apiController.UpdateContainerProfileController)
	containerProfileGroup.DELETE(
		"/:profileId/",
		apiController.DeleteContainerProfileController,
	)

	containerRegistryGroup := containerGroup.Group("/registry")
	containerRegistryGroup.GET("/image/", apiController.GetContainerRegistryImagesController)
	containerRegistryGroup.GET(
		"/image/tagged/",
		apiController.GetContainerRegistryTaggedImageController,
	)
}

func (router *Router) licenseRoutes() {
	licenseGroup := router.baseRoute.Group("/license")
	licenseGroup.GET("/", apiController.GetLicenseInfoController)
	go apiController.AutoLicenseValidationController(
		router.persistentDbSvc,
		router.transientDbSvc,
	)
}

func (router *Router) mappingRoutes() {
	mappingGroup := router.baseRoute.Group("/mapping")
	mappingGroup.GET("/", apiController.GetMappingsController)
	mappingGroup.POST("/", apiController.AddMappingController)
	mappingGroup.DELETE("/:mappingId/", apiController.DeleteMappingController)
	mappingGroup.POST("/target/", apiController.AddMappingTargetController)
	mappingGroup.DELETE(
		"/:mappingId/target/:targetId/",
		apiController.DeleteMappingTargetController,
	)
}

func (router *Router) o11yRoutes() {
	o11yGroup := router.baseRoute.Group("/o11y")
	o11yGroup.GET("/overview/", apiController.O11yOverviewController)
}

func (router *Router) RegisterRoutes() {
	router.swaggerRoute()
	router.authRoutes()
	router.accountRoutes()
	router.containerRoutes()
	router.licenseRoutes()
	router.mappingRoutes()
	router.o11yRoutes()
}
