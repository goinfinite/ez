package apiMiddleware

import (
	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/infra/db"
)

func SetPersistentDatabaseService(persistDbSvc *db.PersistentDatabaseService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("persistDbSvc", persistDbSvc)
			return next(c)
		}
	}
}
