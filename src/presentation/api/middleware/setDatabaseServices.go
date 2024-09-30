package apiMiddleware

import (
	"github.com/goinfinite/ez/src/infra/db"
	"github.com/labstack/echo/v4"
)

func SetDatabaseServices(
	persistentDbSvc *db.PersistentDatabaseService,
	transientDbSvc *db.TransientDatabaseService,
	trailDbSvc *db.TrailDatabaseService,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("persistentDbSvc", persistentDbSvc)
			c.Set("transientDbSvc", transientDbSvc)
			c.Set("trailDbSvc", trailDbSvc)
			return next(c)
		}
	}
}
