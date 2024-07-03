package ui

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/infra/db"
	"github.com/speedianet/control/src/presentation/api"
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
var frontFiles embed.FS

func UiFs() http.Handler {
	frontFileFs, err := fs.Sub(frontFiles, "dist")
	if err != nil {
		panic(err)
	}

	return http.FileServer(http.FS(frontFileFs))
}

func (router *Router) rootRoute() {
	router.baseRoute.GET("/*", echo.WrapHandler(
		http.StripPrefix("/_", UiFs())),
	)
}

func (router *Router) RegisterRoutes() {
	router.rootRoute()

	router.baseRoute.RouteNotFound("/*", func(c echo.Context) error {
		urlPath := c.Request().URL.Path
		isApi := strings.HasPrefix(urlPath, api.ApiBasePath)
		if isApi {
			return c.NoContent(http.StatusNotFound)
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/_/")
	})
}
