package apiMiddleware

import (
	"github.com/labstack/echo/v4"
	"github.com/speedianet/control/src/infra/db"
)

func SetDatabaseService(dbSvc *db.DatabaseService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("dbSvc", dbSvc)
			return next(c)
		}
	}
}
