package ui

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/infra/db"
	"github.com/speedianet/control/src/presentation/api"
	"github.com/speedianet/control/src/presentation/ui/presenter"
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

//go:embed dist/*
var previousDashFiles embed.FS

//go:embed assets/*
var assetsFiles embed.FS

func (router *Router) assetsRoute() {
	assetsFs, err := fs.Sub(assetsFiles, "assets")
	if err != nil {
		slog.Error("ReadAssetsFilesError", slog.Any("error", err))
		os.Exit(1)
	}
	assetsFileServer := http.FileServer(http.FS(assetsFs))

	router.baseRoute.GET(
		"/assets/*",
		echo.WrapHandler(http.StripPrefix("/assets/", assetsFileServer)),
	)
}

func (router *Router) containerRoutes() {
	containerGroup := router.baseRoute.Group("/container")

	containerProfileGroup := containerGroup.Group("/profile")
	profilePresenter := presenter.NewContainerProfilePresenter(router.persistentDbSvc)
	containerProfileGroup.GET("/", profilePresenter.Handler)
}

func (router *Router) previousDashboardRoute() {
	dashFilesFs, err := fs.Sub(previousDashFiles, "dist")
	if err != nil {
		slog.Error("ReadPreviousDashFilesError", slog.Any("error", err))
		os.Exit(1)
	}
	dashFileServer := http.FileServer(http.FS(dashFilesFs))

	previousDashGroup := router.baseRoute.Group("/_")
	previousDashGroup.GET(
		"/*",
		echo.WrapHandler(http.StripPrefix("/_", dashFileServer)),
	)
}

func (router *Router) RegisterRoutes() {
	router.assetsRoute()
	router.containerRoutes()
	router.previousDashboardRoute()

	router.baseRoute.RouteNotFound("/*", func(c echo.Context) error {
		urlPath := c.Request().URL.Path
		isApi := strings.HasPrefix(urlPath, api.ApiBasePath)
		if isApi {
			return c.NoContent(http.StatusNotFound)
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/_/")
	})
}
