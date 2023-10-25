package apiMiddleware

import (
	"net/http"

	"github.com/goinfinite/fleet/src/infra/db"
	"github.com/labstack/echo/v4"
)

func DatabaseInit() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			dbSvc, err := db.NewDatabaseService()
			if err != nil {
				return echo.NewHTTPError(
					http.StatusInternalServerError,
					map[string]interface{}{
						"status": http.StatusInternalServerError,
						"body":   "DatabaseConnectionError",
					})
			}

			c.Set("dbSvc", dbSvc)

			return next(c)
		}
	}
}
