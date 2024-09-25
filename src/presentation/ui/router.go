package ui

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
	"github.com/speedianet/control/src/infra/db"
	"github.com/speedianet/control/src/presentation/api"
	"github.com/speedianet/control/src/presentation/ui/presenter"
	"golang.org/x/net/websocket"
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
	profilePresenter := presenter.NewContainerProfilePresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	containerProfileGroup.GET("/", profilePresenter.Handler)

	imagePresenter := presenter.NewContainerImagePresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	containerImageGroup := containerGroup.Group("/image")
	containerImageGroup.GET("/", imagePresenter.Handler)
}

func (router *Router) devRoutes() {
	devGroup := router.baseRoute.Group("/dev")
	devGroup.GET("/hot-reload", func(c echo.Context) error {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()
			for {
				err := websocket.Message.Send(ws, "WS Hot Reload Activated!")
				if err != nil {
					break
				}

				msgReceived := ""
				err = websocket.Message.Receive(ws, &msgReceived)
				if err != nil {
					break
				}
			}
		}).ServeHTTP(c.Response(), c.Request())
		return nil
	})
}

func (router *Router) fragmentRoutes() {
	fragmentGroup := router.baseRoute.Group("/fragment")

	footerPresenter := presenter.NewFooterPresenter(
		router.persistentDbSvc, router.transientDbSvc, router.trailDbSvc,
	)
	fragmentGroup.GET("/footer", footerPresenter.Handler)
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

	if isDevMode, _ := voHelper.InterfaceToBool(os.Getenv("DEV_MODE")); isDevMode {
		router.devRoutes()
	}

	router.fragmentRoutes()
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
