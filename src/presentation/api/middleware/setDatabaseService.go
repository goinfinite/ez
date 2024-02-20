package apiMiddleware

import (
	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/infra/db"
)

func SetPersistentDatabaseService(persistentDbSvc *db.PersistentDatabaseService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("persistentDbSvc", persistentDbSvc)
			return next(c)
		}
	}
}
