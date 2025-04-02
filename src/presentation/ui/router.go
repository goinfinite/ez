package ui

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
	"github.com/goinfinite/ez/src/infra/db"
	"github.com/goinfinite/ez/src/presentation/api"
	"github.com/goinfinite/ez/src/presentation/ui/presenter"
	"github.com/labstack/echo/v4"
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

func (router *Router) accountsRoutes() {
	accountsGroup := router.baseRoute.Group("/accounts")

	accountsPresenter := presenter.NewAccountsPresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	accountsGroup.GET("/", accountsPresenter.Handler)
}

func (router *Router) backupRoutes() {
	backupGroup := router.baseRoute.Group("/backup")

	backupPresenter := presenter.NewBackupPresenter(
		router.persistentDbSvc, router.transientDbSvc, router.trailDbSvc,
	)
	backupGroup.GET("/", backupPresenter.Handler)
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

func (router *Router) loginRoutes() {
	loginGroup := router.baseRoute.Group("/login")

	loginPresenter := presenter.NewLoginPresenter()
	loginGroup.GET("/", loginPresenter.Handler)
}

func (router *Router) mappingsRoutes() {
	mappingsGroup := router.baseRoute.Group("/mappings")

	mappingsPresenter := presenter.NewMappingsPresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	mappingsGroup.GET("/", mappingsPresenter.Handler)
}

func (router *Router) overviewRoutes() {
	overviewGroup := router.baseRoute.Group("/overview")

	overviewPresenter := presenter.NewOverviewPresenter(
		router.persistentDbSvc, router.transientDbSvc, router.trailDbSvc,
	)
	overviewGroup.GET("/", overviewPresenter.Handler)
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

func (router *Router) RegisterRoutes() {
	router.assetsRoute()
	router.accountsRoutes()
	router.backupRoutes()
	router.containerRoutes()
	router.loginRoutes()
	router.mappingsRoutes()
	router.overviewRoutes()

	if isDevMode, _ := voHelper.InterfaceToBool(os.Getenv("DEV_MODE")); isDevMode {
		router.devRoutes()
	}

	router.fragmentRoutes()

	router.baseRoute.RouteNotFound("/*", func(c echo.Context) error {
		urlPath := c.Request().URL.Path
		isApi := strings.HasPrefix(urlPath, api.ApiBasePath)
		if isApi {
			return c.NoContent(http.StatusNotFound)
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/overview/")
	})
}
