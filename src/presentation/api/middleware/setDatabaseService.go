package apiMiddleware

import (
	"github.com/goinfinite/fleet/src/infra/db"
	"github.com/labstack/echo/v4"
)

func SetDatabaseService(dbSvc *db.DatabaseService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("dbSvc", dbSvc)
			return next(c)
		}
	}
}
