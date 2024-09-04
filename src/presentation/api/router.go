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
	trailDbSvc      *db.TrailDatabaseService
}

func NewRouter(
	baseRoute *echo.Group,
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) *Router {
	return &Router{
		baseRoute:       baseRoute,
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

func (router *Router) swaggerRoute() {
	swaggerGroup := router.baseRoute.Group("/swagger")
	swaggerGroup.GET("/*", echoSwagger.WrapHandler)
}

func (router *Router) authRoutes() {
	authGroup := router.baseRoute.Group("/v1/auth")

	authController := apiController.NewAuthController(
		router.persistentDbSvc, router.trailDbSvc,
	)
	authGroup.POST("/login/", authController.Login)
}

func (router *Router) accountRoutes() {
	accountGroup := router.baseRoute.Group("/v1/account")

	accountController := apiController.NewAccountController(
		router.persistentDbSvc, router.trailDbSvc,
	)
	accountGroup.GET("/", accountController.Read)
	accountGroup.POST("/", accountController.Create)
	accountGroup.PUT("/", accountController.Update)
	accountGroup.DELETE("/:accountId/", accountController.Delete)
	go accountController.AutoUpdateAccountsQuotaUsage()
}

func (router *Router) containerRoutes() {
	containerGroup := router.baseRoute.Group("/v1/container")
	containerController := apiController.NewContainerController(
		router.persistentDbSvc, router.trailDbSvc,
	)
	containerGroup.GET("/", containerController.Read)
	containerGroup.GET("/metrics/", containerController.ReadWithMetrics)
	containerGroup.GET(
		"/session/:accountId/:containerId/", containerController.CreateContainerSessionToken,
	)
	containerGroup.POST("/", containerController.Create)
	containerGroup.PUT("/", containerController.Update)
	containerGroup.DELETE(
		"/:accountId/:containerId/",
		containerController.Delete,
	)

	containerProfileGroup := containerGroup.Group("/profile")
	containerProfileController := apiController.NewContainerProfileController(
		router.persistentDbSvc, router.trailDbSvc,
	)
	containerProfileGroup.GET("/", containerProfileController.Read)
	containerProfileGroup.POST("/", containerProfileController.Create)
	containerProfileGroup.PUT("/", containerProfileController.Update)
	containerProfileGroup.DELETE("/:profileId/", containerProfileController.Delete)

	containerRegistryGroup := containerGroup.Group("/registry")
	containerRegistryGroup.GET("/image/", apiController.GetContainerRegistryImagesController)
	containerRegistryGroup.GET(
		"/image/tagged/",
		apiController.GetContainerRegistryTaggedImageController,
	)

	containerImageGroup := containerGroup.Group("/image")
	containerImageController := apiController.NewContainerImageController(
		router.persistentDbSvc, router.trailDbSvc,
	)
	containerImageGroup.GET("/", containerImageController.Read)
	containerImageGroup.DELETE("/:accountId/:imageId/", containerImageController.Delete)
	containerImageGroup.POST("/snapshot/", containerImageController.CreateSnapshot)

	containerImageArchiveGroup := containerImageGroup.Group("/archive")
	containerImageArchiveGroup.GET("/", containerImageController.ReadArchiveFiles)
	containerImageArchiveGroup.GET(
		"/:accountId/:imageId/", containerImageController.ReadArchiveFile,
	)
	containerImageArchiveGroup.POST("/", containerImageController.CreateArchiveFile)
	containerImageArchiveGroup.POST(
		"/import/", containerImageController.ImportArchiveFile,
	)
	containerImageArchiveGroup.DELETE(
		"/:accountId/:imageId/", containerImageController.DeleteArchiveFile,
	)
}

func (router *Router) licenseRoutes() {
	licenseGroup := router.baseRoute.Group("/v1/license")
	licenseGroup.GET("/", apiController.ReadLicenseInfoController)
	go apiController.AutoLicenseValidationController(
		router.persistentDbSvc,
		router.transientDbSvc,
	)
}

func (router *Router) mappingRoutes() {
	mappingGroup := router.baseRoute.Group("/v1/mapping")

	mappingController := apiController.NewMappingController(
		router.persistentDbSvc, router.trailDbSvc,
	)
	mappingGroup.GET("/", mappingController.Read)
	mappingGroup.POST("/", mappingController.Create)
	mappingGroup.DELETE("/:mappingId/", mappingController.Delete)
	mappingGroup.POST("/target/", mappingController.CreateTarget)
	mappingGroup.DELETE(
		"/:mappingId/target/:targetId/",
		mappingController.DeleteTarget,
	)
}

func (router *Router) o11yRoutes() {
	o11yGroup := router.baseRoute.Group("/v1/o11y")
	o11yGroup.GET("/overview/", apiController.O11yOverviewController)
}

func (router *Router) scheduledTaskRoutes() {
	scheduledTaskGroup := router.baseRoute.Group("/v1/task")

	scheduledTaskController := apiController.NewScheduledTaskController(router.persistentDbSvc)
	scheduledTaskGroup.GET("/", scheduledTaskController.Read)
	scheduledTaskGroup.PUT("/", scheduledTaskController.Update)
	go scheduledTaskController.Run()
}

func (router *Router) RegisterRoutes() {
	router.swaggerRoute()
	router.authRoutes()
	router.accountRoutes()
	router.containerRoutes()
	router.licenseRoutes()
	router.mappingRoutes()
	router.o11yRoutes()
	router.scheduledTaskRoutes()
}
