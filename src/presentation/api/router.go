package api

import (
	_ "embed"

	"github.com/goinfinite/ez/src/infra/db"
	apiController "github.com/goinfinite/ez/src/presentation/api/controller"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/goinfinite/ez/src/presentation/api/docs"
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
	go accountController.AutoRefreshAccountQuotas()
}

func (router *Router) backupRoutes() {
	backupGroup := router.baseRoute.Group("/v1/backup")
	backupController := apiController.NewBackupController(
		router.persistentDbSvc, router.trailDbSvc,
	)

	destinationGroup := backupGroup.Group("/destination")
	destinationGroup.GET("/", backupController.ReadDestination)
	destinationGroup.POST("/", backupController.CreateDestination)
	destinationGroup.PUT("/", backupController.UpdateDestination)
	destinationGroup.DELETE("/:accountId/:destinationId/", backupController.DeleteDestination)

	jobGroup := backupGroup.Group("/job")
	jobGroup.GET("/", backupController.ReadJob)
	jobGroup.POST("/", backupController.CreateJob)
	jobGroup.PUT("/", backupController.UpdateJob)
	jobGroup.DELETE("/:accountId/:jobId/", backupController.DeleteJob)
	jobGroup.POST("/:accountId/:jobId/run/", backupController.RunJob)

	taskGroup := backupGroup.Group("/task")
	taskGroup.GET("/", backupController.ReadTask)
	taskGroup.DELETE("/:taskId/", backupController.DeleteTask)
	taskArchiveGroup := taskGroup.Group("/archive")
	taskArchiveGroup.GET("/", backupController.ReadTaskArchive)
	taskArchiveGroup.POST("/", backupController.CreateTaskArchive)
}

func (router *Router) containerRoutes() {
	containerGroup := router.baseRoute.Group("/v1/container")
	containerController := apiController.NewContainerController(
		router.persistentDbSvc, router.trailDbSvc,
	)
	containerGroup.GET("/", containerController.Read)
	containerGroup.GET("/metrics/", containerController.Read)
	containerGroup.GET(
		"/session/:accountId/:containerId/", containerController.CreateContainerSessionToken,
	)
	containerGroup.POST("/", containerController.Create)
	containerGroup.PUT("/", containerController.Update)
	containerGroup.DELETE(
		"/:accountId/:containerId/", containerController.Delete,
	)

	containerProfileGroup := containerGroup.Group("/profile")
	containerProfileController := apiController.NewContainerProfileController(
		router.persistentDbSvc, router.trailDbSvc,
	)
	containerProfileGroup.GET("/", containerProfileController.Read)
	containerProfileGroup.POST("/", containerProfileController.Create)
	containerProfileGroup.PUT("/", containerProfileController.Update)
	containerProfileGroup.DELETE(
		"/:accountId/:profileId/", containerProfileController.Delete,
	)

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

func (router *Router) mappingRoutes() {
	mappingGroup := router.baseRoute.Group("/v1/mapping")

	mappingController := apiController.NewMappingController(
		router.persistentDbSvc, router.trailDbSvc,
	)
	mappingGroup.GET("/", mappingController.Read)
	mappingGroup.POST("/", mappingController.Create)
	mappingGroup.DELETE("/:mappingId/", mappingController.Delete)
	mappingGroup.DELETE("/:accountId/:mappingId/", mappingController.Delete)
	mappingGroup.POST("/target/", mappingController.CreateTarget)
	mappingGroup.DELETE(
		"/:mappingId/target/:targetId/", mappingController.DeleteTarget,
	)
	mappingGroup.DELETE(
		"/:accountId/:mappingId/target/:targetId/", mappingController.DeleteTarget,
	)
}

func (router *Router) marketplaceRoutes() {
	marketplaceGroup := router.baseRoute.Group("/v1/marketplace")

	marketplaceController := apiController.NewMarketplaceController()
	marketplaceGroup.GET("/", marketplaceController.Read)

	go marketplaceController.RefreshMarketplace()
}

func (router *Router) o11yRoutes() {
	o11yGroup := router.baseRoute.Group("/v1/o11y")
	o11yController := apiController.NewO11yController(router.transientDbSvc)

	o11yGroup.GET("/overview/", o11yController.ReadOverview)
}

func (router *Router) scheduledTaskRoutes() {
	scheduledTaskGroup := router.baseRoute.Group("/v1/scheduled-task")

	scheduledTaskController := apiController.NewScheduledTaskController(router.persistentDbSvc)
	scheduledTaskGroup.GET("/", scheduledTaskController.Read)
	scheduledTaskGroup.PUT("/", scheduledTaskController.Update)
	go scheduledTaskController.Run()
}

func (router *Router) RegisterRoutes() {
	router.swaggerRoute()
	router.authRoutes()
	router.accountRoutes()
	router.backupRoutes()
	router.containerRoutes()
	router.mappingRoutes()
	router.marketplaceRoutes()
	router.o11yRoutes()
	router.scheduledTaskRoutes()
}
